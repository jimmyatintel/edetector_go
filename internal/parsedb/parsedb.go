package parsedb

import (
	"database/sql"
	"edetector_go/config"
	"edetector_go/internal/fflag"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/rabbitmq"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var currentDir string
var unstagePath string
var stagedPath string

func parser_init() {
	curDir, err := os.Getwd()
	if err != nil {
		logger.Error("Error getting current dir:", zap.Any("error", err.Error()))
	}
	currentDir = curDir
	unstagePath = filepath.Join(currentDir, "../../dbUnstage")
	stagedPath = filepath.Join(currentDir, "../../dbStaged")
	CheckDir(unstagePath)
	CheckDir(stagedPath)

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
	logger.Info("Check & Create DB dir")
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("PARSER_LOG_FILE"))
		logger.Info("logger is enabled please check all out info in log file: ", zap.Any("message", config.Viper.GetString("PARSER_LOG_FILE")))
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("rabbit_enable"); enable && err == nil {
		rabbitmq.Rabbit_init()
		logger.Info("rabbit is enabled.")
	}
}

func CheckDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			logger.Error("error creating working dir:", zap.Any("error", err.Error()))
		}
		logger.Info("create dir:", zap.Any("message", path))
	}
}

func Main() {
	parser_init()
	// for {
	dbFiles, err := getDBFiles(unstagePath)
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
		// var tableNames []string
		// tableNames = append(tableNames, "ARPCache")
		if err != nil {
			logger.Error("Error getting table names: ", zap.Any("error", err.Error()))
			continue
		}
		// loop all tables in the db file
		for _, tableName := range tableNames {
			rows, err := db.Query("SELECT * FROM " + tableName)
			if err != nil {
				logger.Error("Error getting rows: ", zap.Any("error", err.Error()))
				continue
			}
			logger.Info("Handling table: ", zap.Any("message", tableName))
			strData, err := rowsToString(rows, tableName)
			if err != nil {
				logger.Error("Error converting to string: ", zap.Any("error", err.Error()))
				continue
			}
			err = sendCollectToElastic(dbFile, strData, tableName)
			if err != nil {
				logger.Error("Error sending to elastic: ", zap.Any("error", err.Error()))
				continue
			}
			rows.Close()
		}
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

func rowsToString(rows *sql.Rows, tablename string) (string, error) {
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
		line := strings.Join(rowData, "||")
		line = strings.ReplaceAll(line, "<nil>", "0")
		builder.WriteString(line)
		builder.WriteString("#newline#")
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func sendCollectToElastic(dbFile string, rawData string, tableName string) error {
	if tableName == "sqlite_sequence" {
		return nil
	}
	path := strings.Split(strings.Split(dbFile, ".db")[0], "/")
	agent := path[len(path)-1]
	lines := strings.Split(rawData, "#newline#")
outerLoop:
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		values := strings.Split(line, "||")
		var err error
		details := "ed_" + strings.ToLower(tableName)
		switch tableName {
		case "AppResourceUsageMonitor":
			err = toElastic(details, agent, line, values[1], values[19], "software", values[14], &AppResourceUsageMonitor{})
		case "ARPCache":
			err = toElastic(details, agent, line, values[1], "-1", "volatile", values[2], &ARPCache{})
		case "BaseService":
			err = toElastic(details, agent, line, values[0], "-1", "software", values[5], &BaseService{})
		case "ChromeBookmarks":
			err = toElastic(details, agent, line, values[4], values[6], "website_bookmark", values[3], &ChromeBookmarks{})
		case "ChromeCache":
			err = toElastic(details, agent, line, values[1], values[8], "cookie_cache", values[2], &ChromeCache{})
		case "ChromeDownload":
			err = toElastic(details, agent, line, values[0], values[6], "website_bookmark", values[3], &ChromeDownload{})
		case "ChromeHistory":
			err = toElastic(details, agent, line, values[0], values[2], "website_bookmark", values[1], &ChromeHistory{})
		case "ChromeKeywordSearch":
			err = toElastic(details, agent, line, values[0], "-1", "website_bookmark", "", &ChromeKeywordSearch{})
		case "ChromeLogin":
			err = toElastic(details, agent, line, values[0], values[6], "website_bookmark", values[3], &ChromeLogin{})
		case "DNSInfo":
			err = toElastic(details, agent, line, values[9], "-1", "software", values[6], &DNSInfo{})
		case "EdgeBookmarks":
			err = toElastic(details, agent, line, values[3], values[7], "website_bookmark", values[4], &EdgeBookmarks{})
		case "EdgeCache":
			err = toElastic(details, agent, line, values[1], values[10], "cookie_cache", values[2], &EdgeCache{})
		case "EdgeCookies":
			err = toElastic(details, agent, line, values[3], values[7], "cookie_cache", values[2], &EdgeCookies{})
		case "EdgeHistory":
			err = toElastic(details, agent, line, values[1], values[5], "website_bookmark", values[2], &EdgeHistory{})
		case "EdgeLogin":
			err = toElastic(details, agent, line, values[1], values[7], "website_bookmark", values[4], &EdgeLogin{})
		case "EventApplication":
			err = toElastic(details, agent, line, values[3], values[9], "software", values[17], &EventApplication{})
		case "EventSecurity":
			err = toElastic(details, agent, line, values[3], values[9], "usb", values[17], &EventSecurity{})
		case "EventSystem":
			err = toElastic(details, agent, line, values[3], values[9], "usb", values[17], &EventSystem{})
		case "FirefoxBookmarks":
			err = toElastic(details, agent, line, values[8], values[5], "website_bookmark", values[3], &FirefoxBookmarks{})
		case "FirefoxCache":
			err = toElastic(details, agent, line, values[1], values[8], "cookie_cache", values[2], &FirefoxCache{})
		case "FirefoxCookies":
			err = toElastic(details, agent, line, values[1], values[5], "cookie_cache", values[3], &FirefoxCookies{})
		case "FirefoxHistory":
			err = toElastic(details, agent, line, values[0], values[9], "website_bookmark", values[1], &FirefoxHistory{})
		case "IEHistory":
			err = toElastic(details, agent, line, values[0], values[4], "website_bookmark", values[1], &IEHistory{})
		case "InstalledSoftware":
			err = toElastic(details, agent, line, values[0], values[17], "network_record", values[6], &InstalledSoftware{})
		case "JumpList":
			err = toElastic(details, agent, line, values[0], values[5], "software", values[1], &JumpList{})
		case "MUICache":
			err = toElastic(details, agent, line, values[0], "-1", "software", values[1], &MUICache{})
		case "Network":
			err = toElastic(details, agent, line, values[1], "-1", "volatile", values[4], &Network{})
		case "NetworkDataUsageMonitor":
			err = toElastic(details, agent, line, values[1], values[10], "software", values[5], &NetworkDataUsageMonitor{})
		case "NetworkResources":
			err = toElastic(details, agent, line, values[0], "-1", "network_record", values[8], &NetworkResources{})
		case "OpenedFiles":
			err = toElastic(details, agent, line, values[1], "-1", "volatile", values[0], &OpenedFiles{})
		case "Prefetch":
			err = toElastic(details, agent, line, values[1], values[2], "software", values[3], &Prefetch{})
		case "Process":
			err = toElastic(details, agent, line, values[1], values[3], "volatile", values[4], &Process{})
		case "Service":
			err = toElastic(details, agent, line, values[0], "-1", "software", values[5], &Service{})
		case "Shortcuts":
			err = toElastic(details, agent, line, values[0], values[10], "document", values[2], &Shortcuts{})
		case "StartRun":
			err = toElastic(details, agent, line, values[0], "-1", "software", values[1], &StartRun{})
		case "TaskSchedule":
			err = toElastic(details, agent, line, values[0], values[3], "software", values[1], &TaskSchedule{})
		case "USBdevices":
			err = toElastic(details, agent, line, values[1], values[14], "usb", values[10], &USBdevices{})
		case "UserAssist":
			err = toElastic(details, agent, line, values[0], values[5], "software", values[2], &UserAssist{})
		case "UserProfiles":
			err = toElastic(details, agent, line, values[0], values[6], "document", values[2], &UserProfiles{})
		case "WindowsActivity":
			err = toElastic(details, agent, line, values[1], values[15], "document", values[3], &WindowsActivity{})
		default:
			logger.Error("Unknown table name: ", zap.Any("message", tableName))
			break outerLoop
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func toElastic(details string, agent string, line string, item string, date string, ttype string, etc string, st elasticquery.Request_data) error {
	uuid := uuid.NewString()
	int_date, err := strconv.Atoi(date)
	if err != nil {
		// logger.Error("Invalid date: ", zap.Any("message", date))
		int_date = 0
	}
	err = elasticquery.SendToMainElastic(uuid, details, agent, item, int_date, ttype, etc)
	if err != nil {
		return err
	}
	err = elasticquery.SendToDetailsElastic(uuid, details, agent, line, st)
	if err != nil {
		return err
	}
	return nil
}
