package dbparser

import (
	"database/sql"
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/internal/file"
	"edetector_go/internal/taskservice"
	"edetector_go/pkg/elastic"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var dbUnstagePath = "dbUnstage"
var dbStagedPath = "dbStaged"

func init() {
	file.CheckDir(dbUnstagePath)
	file.CheckDir(dbStagedPath)

	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Error("Error loading feature flag")
		return
	}
	vp := config.LoadConfig()
	if vp == nil {
		logger.Error("Error loading config file")
		return
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("PARSER_LOG_FILE"))
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
		err := elastic.SetElkClient()
		if err != nil {
			logger.Error("Error connecting to elastic: " + err.Error())
		}
		logger.Info("elastic is enabled.")
	}
}

func Main() {
	for {
		dbFile, agent := file.GetOldestFile(dbUnstagePath, ".db")
		elastic.DeleteByQueryRequest("agent", agent, "StartCollect")
		time.Sleep(3 * time.Second) // wait for fully copy
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			logger.Error("Error opening database file: ", zap.Any("error", err.Error()))
			continue
		}
		logger.Info("Open db file: ", zap.Any("message", dbFile))
		tableNames, err := getTableNames(db)
		if err != nil {
			logger.Error("Error getting table names: ", zap.Any("error", err.Error()))
			continue
		}
		// loop all tables in the db file
		for _, tableName := range tableNames {
			logger.Info("Handling table: ", zap.Any("message", tableName))
			rows, err := db.Query("SELECT * FROM " + tableName)
			if err != nil {
				logger.Error("Error getting rows: ", zap.Any("error", err.Error()))
				continue
			}
			err = sendCollectToElastic(rows, tableName, agent)
			if err != nil {
				logger.Error("Error sending to elastic: ", zap.Any("error", err.Error()))
				continue
			}
			rows.Close()
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

func sendCollectToElastic(rows *sql.Rows, tableName string, agent string) error {
	if tableName == "sqlite_sequence" {
		return nil
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	colValues := make([]interface{}, len(columns))
	for i := range colValues {
		colValues[i] = new(interface{})
	}
	for rows.Next() {
		err = rows.Scan(colValues...)
		if err != nil {
			return err
		}
		values := make([]string, len(columns))
		for i, val := range colValues {
			switch v := (*val.(*interface{})).(type) {
			case []byte:
				logger.Info("byte type")
				values[i] = string(v)
			default:
				values[i] = fmt.Sprintf("%v", v)
			}
			if values[i] == "" || values[i] == "<nil>" {
				logger.Info("empty", zap.Any("message", tableName+" "+strconv.Itoa(i)+" "+values[i]))
				values[i] = "0"
			}
			logger.Info("data", zap.Any("message", tableName+" "+strconv.Itoa(i)+" "+values[i]))
		}
		var err error
		index := config.Viper.GetString("ELASTIC_PREFIX") + "_" + strings.ToLower(tableName) //! developing
		switch tableName {
		case "AppResourceUsageMonitor":
			err = toElastic(index, agent, values, values[1], values[19], "software", values[14], &AppResourceUsageMonitor{})
		case "ARPCache":
			err = toElastic(index, agent, values, values[1], "0", "volatile", values[2], &ARPCache{})
		case "BaseService":
			err = toElastic(index, agent, values, values[0], "0", "software", values[5], &BaseService{})
		case "ChromeBookmarks":
			err = toElastic(index, agent, values, values[4], values[6], "website_bookmark", values[3], &ChromeBookmarks{})
		case "ChromeCache":
			err = toElastic(index, agent, values, values[1], values[8], "cookie_cache", values[2], &ChromeCache{})
		case "ChromeDownload":
			err = toElastic(index, agent, values, values[0], values[6], "website_bookmark", values[3], &ChromeDownload{})
		case "ChromeHistory":
			err = toElastic(index, agent, values, values[0], values[2], "website_bookmark", values[1], &ChromeHistory{})
		case "ChromeKeywordSearch":
			err = toElastic(index, agent, values, values[0], "0", "website_bookmark", "", &ChromeKeywordSearch{})
		case "ChromeLogin":
			err = toElastic(index, agent, values, values[0], values[6], "website_bookmark", values[3], &ChromeLogin{})
		case "DNSInfo":
			err = toElastic(index, agent, values, values[9], "0", "software", values[6], &DNSInfo{})
		case "EdgeBookmarks":
			err = toElastic(index, agent, values, values[3], values[7], "website_bookmark", values[4], &EdgeBookmarks{})
		case "EdgeCache":
			err = toElastic(index, agent, values, values[1], values[10], "cookie_cache", values[2], &EdgeCache{})
		case "EdgeCookies":
			err = toElastic(index, agent, values, values[3], values[7], "cookie_cache", values[2], &EdgeCookies{})
		case "EdgeHistory":
			err = toElastic(index, agent, values, values[1], values[5], "website_bookmark", values[2], &EdgeHistory{})
		case "EdgeLogin":
			err = toElastic(index, agent, values, values[1], values[7], "website_bookmark", values[4], &EdgeLogin{})
		case "EventApplication":
			err = toElastic(index, agent, values, values[3], values[9], "software", values[17], &EventApplication{})
		case "EventSecurity":
			err = toElastic(index, agent, values, values[3], values[9], "usb", values[17], &EventSecurity{})
		case "EventSystem":
			err = toElastic(index, agent, values, values[3], values[9], "usb", values[17], &EventSystem{})
		case "FirefoxBookmarks":
			err = toElastic(index, agent, values, values[8], values[5], "website_bookmark", values[3], &FirefoxBookmarks{})
		case "FirefoxCache":
			err = toElastic(index, agent, values, values[1], values[8], "cookie_cache", values[2], &FirefoxCache{})
		case "FirefoxCookies":
			err = toElastic(index, agent, values, values[1], values[5], "cookie_cache", values[3], &FirefoxCookies{})
		case "FirefoxHistory":
			err = toElastic(index, agent, values, values[0], values[9], "website_bookmark", values[1], &FirefoxHistory{})
		case "IEHistory":
			err = toElastic(index, agent, values, values[0], values[4], "website_bookmark", values[1], &IEHistory{})
		case "InstalledSoftware":
			err = toElastic(index, agent, values, values[0], values[17], "network_record", values[6], &InstalledSoftware{})
		case "JumpList":
			err = toElastic(index, agent, values, values[0], values[5], "software", values[1], &JumpList{})
		case "MUICache":
			err = toElastic(index, agent, values, values[0], "0", "software", values[1], &MUICache{})
		case "Network":
			err = toElastic(index, agent, values, values[1], "0", "volatile", values[4], &Network{})
		case "NetworkDataUsageMonitor":
			err = toElastic(index, agent, values, values[1], values[10], "software", values[5], &NetworkDataUsageMonitor{})
		case "NetworkResources":
			err = toElastic(index, agent, values, values[0], "0", "network_record", values[8], &NetworkResources{})
		case "OpenedFiles":
			err = toElastic(index, agent, values, values[1], "0", "volatile", values[0], &OpenedFiles{})
		case "Prefetch":
			err = toElastic(index, agent, values, values[1], values[2], "software", values[3], &Prefetch{})
		case "Process":
			err = toElastic(index, agent, values, values[1], values[3], "volatile", values[4], &Process{})
		case "Service":
			err = toElastic(index, agent, values, values[0], "0", "software", values[5], &Service{})
		case "Shortcuts":
			err = toElastic(index, agent, values, values[0], values[10], "document", values[2], &Shortcuts{})
		case "StartRun":
			err = toElastic(index, agent, values, values[0], "0", "software", values[1], &StartRun{})
		case "TaskSchedule":
			err = toElastic(index, agent, values, values[0], values[3], "software", values[1], &TaskSchedule{})
		case "USBdevices":
			err = toElastic(index, agent, values, values[1], values[14], "usb", values[10], &USBdevices{})
		case "UserAssist":
			err = toElastic(index, agent, values, values[0], values[5], "software", values[2], &UserAssist{})
		case "UserProfiles":
			err = toElastic(index, agent, values, values[0], values[6], "document", values[2], &UserProfiles{})
		case "WindowsActivity":
			err = toElastic(index, agent, values, values[1], values[15], "document", values[3], &WindowsActivity{})
		default:
			logger.Error("Unknown table name: ", zap.Any("message", tableName))
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func toElastic(index string, agent string, values []string, item string, date string, ttype string, etc string, st elasticquery.Request_data) error {
	ip, name := query.GetMachineIPandName(agent)
	uuid := uuid.NewString()
	int_date, err := strconv.Atoi(date)
	if err != nil {
		logger.Error("Invalid date: ", zap.Any("message", date))
		int_date = 0
		date = "0"
	}
	err = elasticquery.SendToMainElastic(index, uuid, agent, ip, name, item, int_date, ttype, etc, "ed_high")
	if err != nil {
		return err
	}
	err = elasticquery.SendToDetailsElastic(index, st, values, uuid, agent, ip, name, item, date, ttype, etc, "ed_high")
	if err != nil {
		return err
	}
	return nil
}
