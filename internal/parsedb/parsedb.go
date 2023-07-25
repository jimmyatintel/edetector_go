package parsedb

import (
	"database/sql"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
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
		tableNames, err := getTableNames(db)
		if err != nil {
			logger.Error("Error getting table names: ", zap.Any("error", err.Error()))
			return
		}
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
	agent := path[len(path) - 1]
	lines := strings.Split(rawData, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		values := strings.Split(line, "|")
		tableFunc, ok := dbMap[tableName]
		if !ok {
			err := errors.New("table not found: " + tableName)
			return err
		}
		err := tableFunc(agent, line, values)
		if err != nil {
			return err
		}		
	}
	return nil
}

func AppResourceUsageMonitorTable(agent string, line string, values []string) error {
	uuid := uuid.NewString()
	err := elasticquery.SendToMainElastic(uuid, "ed_main", agent, values[1], "-1", "volatile", values[2])
	if err != nil {
		return err
	}
	err = elasticquery.SendToDetailsElastic(uuid, "ed_arpcache", agent, line, &ARPCache{})
	if err != nil {
		return err
	}
	return nil
}

func ARPCacheTable(agent string, line string, values []string) error {
	uuid := uuid.NewString()
	err := elasticquery.SendToMainElastic(uuid, "ed_main", agent, values[1], "-1", "volatile", values[2])
	if err != nil {
		return err
	}
	err = elasticquery.SendToDetailsElastic(uuid, "ed_arpcache", agent, line, &ARPCache{})
	if err != nil {
		return err
	}
	return nil
}

func BaseServiceTable(agent string, line string, values []string) error {
	uuid := uuid.NewString()
	err := elasticquery.SendToMainElastic(uuid, "ed_main", agent, values[1], "-1", "volatile", values[2])
	if err != nil {
		return err
	}
	err = elasticquery.SendToDetailsElastic(uuid, "ed_arpcache", agent, line, &ARPCache{})
	if err != nil {
		return err
	}
	return nil
}

func ChromeBookmarksTable(agent string, line string, values []string) error {
	uuid := uuid.NewString()
	err := elasticquery.SendToMainElastic(uuid, "ed_main", agent, values[1], "-1", "volatile", values[2])
	if err != nil {
		return err
	}
	err = elasticquery.SendToDetailsElastic(uuid, "ed_arpcache", agent, line, &ARPCache{})
	if err != nil {
		return err
	}
	return nil
}
