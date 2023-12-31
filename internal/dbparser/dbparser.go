package dbparser

import (
	"database/sql"
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"edetector_go/pkg/redis"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var dbUnstagePath = "dbUnstage"
var dbStagedPath = "dbStaged"

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
}

func Main(version string) {
	parser_init()
	logger.Info("Welcome to edetector dbparser: " + version)
outerloop:
	for {
		dbFile, agent, _ := file.GetOldestFile(dbUnstagePath, ".db")
		elastic.DeleteByQueryRequest("agent", agent, "StartCollect")
		time.Sleep(3 * time.Second) // wait for fully copy
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			logger.Error("Error opening database file: " + err.Error())
			query.Failed_task(agent, "StartCollect")
			clearParser(db, dbFile, agent)
			continue
		}
		logger.Info("Open db file: " + dbFile)
		tableNames, err := getTableNames(db)
		if err != nil {
			logger.Error("Error getting table names: " + err.Error())
			query.Failed_task(agent, "StartCollect")
			clearParser(db, dbFile, agent)
			continue
		}
		for _, tableName := range tableNames {
			if terminateCollect(db, dbFile, agent) {
				continue outerloop
			}
			err = sendCollectToRabbitMQ(db, tableName, agent)
			if err != nil {
				logger.Error("Error sending to rabbitMQ: " + err.Error())
				query.Failed_task(agent, "StartCollect")
				break
			}
		}
		if terminateCollect(db, dbFile, agent) {
			continue
		}
		clearParser(db, dbFile, agent)
		query.Finish_task(agent, "StartCollect")
		logger.Info("DB parser task finished: " + agent)
	}
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

func terminateCollect(db *sql.DB, dbFile string, agent string) bool {
	var flag = false
	if redis.RedisExists(agent+"-terminateCollect") && redis.RedisGetInt(agent+"-terminateCollect") == 1 {
		logger.Info("Terminate collect: " + agent)
		flag = true
		elastic.DeleteByQueryRequest("agent", agent, "StartCollect")
		clearParser(db, dbFile, agent)
	}
	return flag
}

func clearParser(db *sql.DB, dbFile string, agent string) {
	redis.RedisSet(agent+"-terminateCollect", 0)
	db.Close()
	err := file.MoveFile(dbFile, filepath.Join(dbStagedPath, agent+".db"))
	if err != nil {
		logger.Error("Error moving file: " + err.Error())
	}
}
