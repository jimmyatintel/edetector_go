package dbparser

import (
	"database/sql"
	"edetector_go/config"
	"edetector_go/internal/taskservice"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/fflag"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/rabbitmq"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var dbUnstagePath = "dbUnstage"
var dbStagedPath = "dbStaged"

func parser_init() {
	file.CheckDir(dbUnstagePath)
	file.CheckDir(dbStagedPath)

	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Error("Error loading feature flag")
		return
	}
	vp, err := config.LoadConfig()
	if vp == nil {
		logger.Error("Error loading config file", zap.Any("error", err.Error()))
		return
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("PARSER_LOG_FILE"), "dbparser", "DBPARSR")
		logger.Info("logger is enabled please check all out info in log file: ", zap.Any("message", config.Viper.GetString("PARSER_LOG_FILE")))
	}
	if err := mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("rabbit_enable"); enable && err == nil {
		rabbitmq.Rabbit_init()
		logger.Info("rabbit is enabled.")
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("elastic_enable"); enable && err == nil {
		elastic.Elastic_init()
		logger.Info("elastic is enabled.")
	}
}

func Main(version string) {
	parser_init()
	logger.Info("Welcome to edetector dbparser: ", zap.Any("version", version))
	for {
		dbFile, agent := file.GetOldestFile(dbUnstagePath, ".db")
		elastic.DeleteByQueryRequest("agent", agent, "StartCollect")
		time.Sleep(3 * time.Second) // wait for fully copy
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			logger.Error("Error opening database file: ", zap.Any("error", err.Error()))
			err = file.MoveFile(dbFile, filepath.Join(dbStagedPath, agent+".db"))
			if err != nil {
				logger.Error("Error moving file: ", zap.Any("error", err.Error()))
			}
			continue
		}
		logger.Info("Open db file: ", zap.Any("message", dbFile))
		tableNames, err := getTableNames(db)
		if err != nil {
			logger.Error("Error getting table names: ", zap.Any("error", err.Error()))
			err = file.MoveFile(dbFile, filepath.Join(dbStagedPath, agent+".db"))
			if err != nil {
				logger.Error("Error moving file: ", zap.Any("error", err.Error()))
			}
			continue
		}
		for _, tableName := range tableNames {
			err = sendCollectToRabbitMQ(db, tableName, agent)
			if err != nil {
				logger.Error("Error sending to elastic: ", zap.Any("error", err.Error()))
				continue
			}
		}
		db.Close()
		err = file.MoveFile(dbFile, filepath.Join(dbStagedPath, agent+".db"))
		if err != nil {
			logger.Error("Error moving file: ", zap.Any("error", err.Error()))
		}
		taskservice.Finish_task(agent, "StartCollect")
		logger.Info("Task finished: ", zap.Any("message", agent))
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
