package parsedb

import (
	"database/sql"
	"edetector_go/pkg/logger"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var currentDir = ""
var unstagePath = "../dbUnstage"

type JSONData []byte

func (j JSONData) Elastical() ([]byte, error) {
	return j, nil
}

func init() {
	curDir, err := os.Getwd()
	if err != nil {
		logger.Error("Error getting current dir:", zap.Any("error", err.Error()))
	}
	currentDir = curDir
}

func Main() {
	dir := filepath.Join(currentDir, unstagePath)
	// for { // infinite loop
	dbFiles, err := getDBFiles(dir)
	if err != nil {
		logger.Error("Error getting database files: ", zap.Any("error", err.Error()))
		return
	}
	// loop all db files
	for _, dbFile := range dbFiles {
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			logger.Error("Error opening database file: ", zap.Any("error", dbFile+" err: "+err.Error()))
			continue
		}
		logger.Info("Open db file: ", zap.Any("message", dbFile))
		// tableNames, err := getTableNames(db)
		// if err != nil {
		// 	logger.Error("Error getting table names: ", zap.Any("error", err.Error()))
		// 	return
		// }
		tableName := "ARPCache"
		// loop all tables in the db file
		// for _, tableName := range tableNames {
		rows, err := db.Query("SELECT * FROM " + tableName)
		if err != nil {
			logger.Error("Error executing query: ", zap.Any("error", err.Error()))
			return
		}
		logger.Info("Executing query: ", zap.Any("message", tableName))
		rawData, err := rowsToString(rows)
		if err != nil {
			logger.Error("Error converting to JSON: ", zap.Any("error", err.Error()))
			return
		}
		fmt.Println(rawData)
		// Data := changeARPCacheToJson(rawData)
		// logger.Info("Data: ", zap.Any("message", Data))
		// slice := strings.Split((strings.Split(dbFile, ".db")[0]), "/")
		// agentID := slice[len(slice)-1]
		// template := elasticquery.New_source(agentID, "ARPCache")
		// elasticquery.Send_to_elastic("ed_process_history", template, Data)
		rows.Close()
		// }
		db.Close()
	}
	// }
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

func rowsToString(rows *sql.Rows) (string, error) {
	var builder strings.Builder
	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("error getting column names: %v", err)
	}

	// Create a slice to store the values of each row
	values := make([]interface{}, len(columns))
	rowData := make([]string, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}
	// Iterate through the rows
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			return "", fmt.Errorf("error scanning row values: %v", err)
		}
		// Build the row data string with | separators
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
		// Append the line to the result with \n separator
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error after iterating through rows: %v", err)
	}
	return builder.String(), nil
}

// func changeARPCacheToJson(msg string) []elasticquery.Request_data {
// 	lines := strings.Split(msg, "\n")
// 	var dataSlice []elasticquery.Request_data
// 	for _, line := range lines {
// 		if len(line) == 0 {
// 			continue
// 		}
// 		data := ARPCache{}
// 		work.To_json(line, &data)
// 		dataSlice = append(dataSlice, elasticquery.Request_data(data))
// 	}
// 	logger.Info("ChangeProcessToJson", zap.Any("message", dataSlice))
// 	return dataSlice
// }

// func getTableNames(db *sql.DB) ([]string, error) {
// 	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var tableNames []string
// 	for rows.Next() {
// 		var tableName string
// 		if err := rows.Scan(&tableName); err != nil {

// 			return nil, err
// 		}
// 		tableNames = append(tableNames, tableName)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return tableNames, nil
// }
