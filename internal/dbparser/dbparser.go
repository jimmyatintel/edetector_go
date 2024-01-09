package dbparser

import (
	"context"
	"database/sql"
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"edetector_go/pkg/redis"
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var dbUnstagePath = "dbUnstage"
var dbStagedPath = "dbStaged"
var limit int
var count int
var cancelMap = map[string]context.CancelFunc{}

func parser_init() {
	file.CheckDir(dbUnstagePath)
	file.CheckDir(dbStagedPath)
	// fflag.Get_fflag()
	// if fflag.FFLAG == nil {
	// 	logger.Panic("Error loading feature flag")
	// 	panic("Error loading feature flag")
	// }
	vp, err := config.LoadConfig()
	if vp == nil {
		logger.Panic("Error loading config file: " + err.Error())
		panic(err)
	}
	if true {
		logger.InitLogger(config.Viper.GetString("PARSER_LOG_FILE"), "dbparser", "DBPARSR")
		logger.Info("logger is enabled please check all out info in log file: " + config.Viper.GetString("PARSER_LOG_FILE"))
	}
	connString, err := mariadb.Connect_init()
	if err != nil {
		logger.Panic("Error connecting to mariadb: " + err.Error())
		panic(err)
	} else {
		logger.Info("Mariadb connectionString: " + connString)
	}
	if true {
		if db := redis.Redis_init(); db == nil {
			logger.Panic("Error connecting to redis")
			panic(err)
		}
	}
	if true {
		rabbitmq.Rabbit_init()
		logger.Info("rabbit is enabled.")
	}
	if true {
		elastic.Elastic_init()
		logger.Info("elastic is enabled.")
	}
	limit = config.Viper.GetInt("PARSER_BUILDER_LIMIT")
}

func Main(version string) {
	parser_init()
	logger.Info("Welcome to edetector dbparser: " + version)
	count = 0
	go terminateCollect()
	for {
		if count < limit {
			dbFile, agent, _ := file.GetOldestFile(dbUnstagePath, ".db")
			count++
			ctx, cancel := context.WithCancel(context.Background())
			cancelMap[agent] = cancel
			go dbParser(ctx, dbFile, agent)
		}
		time.Sleep(10 * time.Second)
	}
}

func terminateCollect() {
	terminating := 5
	for {
		time.Sleep(10 * time.Second)
		handlingTasks, err := query.Load_stored_task("nil", "nil", terminating, "StartCollect")
		if err != nil {
			logger.Error("Error loading stored task: " + err.Error())
			continue
		}
		for _, t := range handlingTasks {
			logger.Info("Received terminate collect: " + t[1])
			if cancelMap[t[1]] != nil {
				cancelMap[t[1]]()
			}
			query.Terminated_task(t[1], "StartCollect", terminating)
		}
	}
}

func dbParser(ctx context.Context, dbFile string, agent string) {
	deleteQuery := fmt.Sprintf(`
	{
		"query": {
			"term": {
				"agent": "%s"
			}
		}
	}
	`, agent)
	elastic.DeleteByQueryRequest(deleteQuery, "StartCollect")
	time.Sleep(3 * time.Second) // wait for fully copy
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		logger.Error("Error opening database file (" + agent + "): " + err.Error())
		query.Failed_task(agent, "StartCollect", 6)
		clearParser(db, dbFile, agent)
		return
	}
	logger.Info("Open db file: " + dbFile)
	tableNames, err := getTableNames(db)
	if err != nil {
		logger.Error("Error getting table names (" + agent + "): " + err.Error())
		query.Failed_task(agent, "StartCollect", 6)
		clearParser(db, dbFile, agent)
		return
	}
	for _, tableName := range tableNames {
		select {
		case <-ctx.Done():
			logger.Info("Terminate collect: " + agent)
			clearParser(db, dbFile, agent)
			return
		default:
			err = sendCollectToRabbitMQ(db, tableName, agent)
			if err != nil {
				logger.Error("Error sending to rabbitMQ (" + agent + "): " + err.Error())
				query.Failed_task(agent, "StartCollect", 6)
				clearParser(db, dbFile, agent)
				return
			}
		}
	}
	clearParser(db, dbFile, agent)
	err = rabbitmq.ToRabbitMQ_FinishSignal(agent, "StartCollect", "ed_low")
	if err != nil {
		logger.Error("Error sending finish signal to rabbitMQ (" + agent + "): " + err.Error())
		query.Failed_task(agent, "StartCollect", 6)
		return
	}
	logger.Info("DB parser task finished: " + agent)
}

func getTableNames(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, tableName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tableNames, nil
}

func clearParser(db *sql.DB, dbFile string, agent string) {
	count--
	cancelMap[agent] = nil
	db.Close()
	err := file.MoveFile(dbFile, filepath.Join(dbStagedPath, agent+".db"))
	if err != nil {
		logger.Error("Error moving file (" + agent + "): " + err.Error())
	} else {
		logger.Info("Move db file to staged: " + dbFile)
	}
}
