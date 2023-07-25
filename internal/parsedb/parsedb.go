package parsedb

import (
	"database/sql"
	"edetector_go/config"
	"edetector_go/internal/fflag"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/rabbitmq"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var currentDir = ""
var unstagePath = "../../dbUnstage"

func init() {
	curDir, err := os.Getwd()
	if err != nil {
		logger.Error("Error getting current dir:", zap.Any("error", err.Error()))
	}
	currentDir = curDir

	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		fmt.Println("Error loading feature flag")
		return
	}
	vp := config.LoadConfig()
	if vp == nil {
		fmt.Println("Error loading config file")
		return
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("DB_LOG_FILE"))
		fmt.Println("logger is enabled please check all out info in log file: ", config.Viper.GetString("DB_LOG_FILE"))
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("rabbit_enable"); enable && err == nil {
		rabbitmq.Rabbit_init()
		fmt.Println("rabbit is enabled.")
	}
}

func Main() {
	dir := filepath.Join(currentDir, unstagePath)
	// for {
	dbFiles, err := getDBFiles(dir)
	if err != nil {
		logger.Error("Error getting database files: ", zap.Any("error", err.Error()))
		return
	}
	// loop all db files
	for _, dbFile := range dbFiles {
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			logger.Error("Error opening database file: ", zap.Any("error", err.Error()))
			continue
		}
		logger.Info("Open db file: ", zap.Any("message", dbFile))
		// tableNames, err := getTableNames(db)
		// if err != nil {
		// 	logger.Error("Error getting table names: ", zap.Any("error", err.Error()))
		// 	return
		// }
		var tableNames []string
		tableNames = append(tableNames, "ARPCache")
		tableNames = append(tableNames, "ChromeDownload")
		// loop all tables in the db file
		for _, tableName := range tableNames {
			rows, err := db.Query("SELECT * FROM " + tableName)
			if err != nil {
				logger.Error("Error getting rows: ", zap.Any("error", err.Error()))
				return
			}
			logger.Info("Handling table: ", zap.Any("message", tableName))
			strData, err := rowsToString(rows)
			if err != nil {
				logger.Error("Error converting to string: ", zap.Any("error", err.Error()))
				return
			}
			err = sendCollectToElastic(dbFile, strData, tableName)
			if err != nil {
				logger.Error("Error sending to elastic: ", zap.Any("error", err.Error()))
				return
			}
			rows.Close()
		}
		db.Close()
		// }
	}
}

func getDBFiles(dir string) ([]string, error) {
	var dbFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".db" {
			dbFiles = append(dbFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dbFiles, nil
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

func rowsToString(rows *sql.Rows) (string, error) {
	var builder strings.Builder
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	values := make([]interface{}, len(columns))
	rowData := make([]string, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			return "", err
		}
		for i, val := range values {
			switch v := (*val.(*interface{})).(type) {
			case int, int64, float64:
				rowData[i] = fmt.Sprintf("%v", v)
			case []byte:
				rowData[i] = string(v)
			default:
				rowData[i] = fmt.Sprintf("%v", v)
			}
		}
		line := strings.Join(rowData, "|")
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func sendCollectToElastic(dbFile string, rawData string, tableName string) error {
	path := strings.Split(strings.Split(dbFile, ".db")[0], "/")
	agent := path[len(path)-1]
	lines := strings.Split(rawData, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		values := strings.Split(line, "|")
		err := errors.New("")
		details := "ed_" + strings.ToLower(tableName)
		switch tableName {
		case "AppResourceUsageMonitor":
			err = toElastic(details, agent, line, values[1], values[19], "software", values[14], &AppResourceUsageMonitor{})
		case "ARPCache":
			err = toElastic(details, agent, line, values[1], "-1", "volatile", values[2], &ARPCache{})
		case "ChromeDownload":
			err = toElastic(details, agent, line, values[0], values[6], "website_bookmark", values[3], &ChromeDownload{})
		default:
			logger.Error("Unknown table name: ", zap.Any("message", tableName))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func toElastic(details string, agent string, line string, item string, date string, ttype string, etc string, st elasticquery.Request_data) error {
	uuid := uuid.NewString()
	err := elasticquery.SendToMainElastic(uuid, "ed_main", agent, item, date, ttype, etc)
	if err != nil {
		return err
	}
	err = elasticquery.SendToDetailsElastic(uuid, details, agent, line, st)
	if err != nil {
		return err
	}
	return nil
}
