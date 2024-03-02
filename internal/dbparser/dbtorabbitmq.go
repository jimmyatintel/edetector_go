package dbparser

import (
	"database/sql"
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func sendCollectToRabbitMQ(db *sql.DB, tableName string, agent string) error {
	taskID := query.Load_task_id(agent, "StartCollect", 2)
	logger.Debug("Handling table (" + agent + "): " + tableName)
	rows, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		logger.Error("Error getting rows (" + agent + "): " + err.Error())
		return err
	}
	defer rows.Close()
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
				values[i] = string(v)
			default:
				values[i] = fmt.Sprintf("%v", v)
			}
			if values[i] == "" || values[i] == " " || values[i] == "<nil>" {
				values[i] = "0"
			}
		}
		var err error
		index := config.Viper.GetString("ELASTIC_PREFIX") + "_" + strings.ToLower(tableName) //! developing
		switch tableName {
		case "AppResourceUsageMonitor":
			err = toRabbitMQ(index, agent, values, values[1], values[19], "software", values[14], &AppResourceUsageMonitor{}, taskID)
		case "ARPCache":
			err = toRabbitMQ(index, agent, values, values[1], "0", "volatile", values[2], &ARPCache{}, taskID)
		case "BaseService":
			values[14] = toBoolean(values[14])
			values[15] = toBoolean(values[15])
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[5], &BaseService{}, taskID)
		case "ChromeBookmarks":
			err = toRabbitMQ(index, agent, values, values[4], values[6], "website_bookmark", values[3], &ChromeBookmarks{}, taskID)
		case "ChromeCache":
			values[8] = RFCToTimestamp(values[8])
			values[9] = RFCToTimestamp(values[9])
			values[10] = RFCToTimestamp(values[10])
			err = toRabbitMQ(index, agent, values, values[1], values[8], "cookie_cache", values[2], &ChromeCache{}, taskID)
		case "ChromeDownload":
			values[11] = toBoolean(values[11])
			err = toRabbitMQ(index, agent, values, values[0], values[6], "website_bookmark", values[3], &ChromeDownload{}, taskID)
		case "ChromeHistory":
			err = toRabbitMQ(index, agent, values, values[0], values[2], "website_bookmark", values[1], &ChromeHistory{}, taskID)
		case "ChromeKeywordSearch":
			err = toRabbitMQ(index, agent, values, values[0], "0", "website_bookmark", "", &ChromeKeywordSearch{}, taskID)
		case "ChromeLogin":
			err = toRabbitMQ(index, agent, values, values[0], values[6], "website_bookmark", values[3], &ChromeLogin{}, taskID)
		case "DNSInfo":
			err = toRabbitMQ(index, agent, values, values[9], "0", "software", values[6], &DNSInfo{}, taskID)
		case "EdgeBookmarks":
			err = toRabbitMQ(index, agent, values, values[3], values[7], "website_bookmark", values[4], &EdgeBookmarks{}, taskID)
		case "EdgeCache":
			values[8] = RFCToTimestamp(values[8])
			values[9] = RFCToTimestamp(values[9])
			values[10] = RFCToTimestamp(values[10])
			err = toRabbitMQ(index, agent, values, values[1], values[10], "cookie_cache", values[2], &EdgeCache{}, taskID)
		case "EdgeCookies":
			err = toRabbitMQ(index, agent, values, values[3], values[7], "cookie_cache", values[2], &EdgeCookies{}, taskID)
		case "EdgeHistory":
			err = toRabbitMQ(index, agent, values, values[1], values[5], "website_bookmark", values[2], &EdgeHistory{}, taskID)
		case "EdgeLogin":
			err = toRabbitMQ(index, agent, values, values[1], values[7], "website_bookmark", values[4], &EdgeLogin{}, taskID)
		case "EventApplication":
			err = toRabbitMQ(index, agent, values, values[3], values[9], "software", values[17], &EventApplication{}, taskID)
		case "EventSecurity":
			err = toRabbitMQ(index, agent, values, values[3], values[9], "usb", values[17], &EventSecurity{}, taskID)
		case "EventSystem":
			err = toRabbitMQ(index, agent, values, values[3], values[9], "usb", values[17], &EventSystem{}, taskID)
		case "FirefoxBookmarks":
			err = toRabbitMQ(index, agent, values, values[8], values[5], "website_bookmark", values[3], &FirefoxBookmarks{}, taskID)
		case "FirefoxCache":
			err = toRabbitMQ(index, agent, values, values[1], values[8], "cookie_cache", values[2], &FirefoxCache{}, taskID)
		case "FirefoxCookies":
			err = toRabbitMQ(index, agent, values, values[1], values[5], "cookie_cache", values[3], &FirefoxCookies{}, taskID)
		case "FirefoxHistory":
			err = toRabbitMQ(index, agent, values, values[0], values[9], "website_bookmark", values[1], &FirefoxHistory{}, taskID)
		case "IEHistory":
			err = toRabbitMQ(index, agent, values, values[0], values[4], "website_bookmark", values[1], &IEHistory{}, taskID)
		case "InstalledSoftware":
			values[3] = DigitToTimestamp(values[3])
			err = toRabbitMQ(index, agent, values, values[0], values[17], "network_record", values[6], &InstalledSoftware{}, taskID)
		case "JumpList":
			err = toRabbitMQ(index, agent, values, values[0], values[5], "software", values[1], &JumpList{}, taskID)
		case "MUICache":
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[1], &MUICache{}, taskID)
		case "Network":
			err = toRabbitMQ(index, agent, values, values[1], "0", "volatile", values[4], &Network{}, taskID)
		case "NetworkDataUsageMonitor":
			err = toRabbitMQ(index, agent, values, values[1], values[10], "software", values[5], &NetworkDataUsageMonitor{}, taskID)
		case "NetworkResources":
			err = toRabbitMQ(index, agent, values, values[0], "0", "network_record", values[8], &NetworkResources{}, taskID)
		case "OpenedFiles":
			err = toRabbitMQ(index, agent, values, values[1], "0", "volatile", values[0], &OpenedFiles{}, taskID)
		case "Prefetch":
			err = toRabbitMQ(index, agent, values, values[1], values[2], "software", values[3], &Prefetch{}, taskID)
		case "Process":
			err = toRabbitMQ(index, agent, values, values[1], values[3], "volatile", values[4], &Process{}, taskID)
		case "Service":
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[5], &Service{}, taskID)
		case "Shortcuts":
			values[7] = toBoolean(values[7])
			err = toRabbitMQ(index, agent, values, values[0], values[10], "document", values[2], &Shortcuts{}, taskID)
		case "StartRun":
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[1], &StartRun{}, taskID)
		case "TaskSchedule":
			err = toRabbitMQ(index, agent, values, values[0], values[3], "software", values[1], &TaskSchedule{}, taskID)
		case "USBdevices":
			err = toRabbitMQ(index, agent, values, values[1], values[14], "usb", values[10], &USBdevices{}, taskID)
		case "UserAssist":
			err = toRabbitMQ(index, agent, values, values[0], values[5], "software", values[2], &UserAssist{}, taskID)
		case "UserProfiles":
			values[3] = toBoolean(values[3])
			err = toRabbitMQ(index, agent, values, values[0], values[6], "document", values[2], &UserProfiles{}, taskID)
		case "WindowsActivity":
			err = toRabbitMQ(index, agent, values, values[1], values[15], "document", values[3], &WindowsActivity{}, taskID)
		case "Wireless":
			err = toRabbitMQ(index, agent, values, values[0], values[8], "network_record", values[1], &Wireless{}, taskID)
		case "Email":
			err = toRabbitMQ(index, agent, values, values[9], values[5], "emails", values[3], &Email{}, taskID)
		case "EmailPath":
			err = toRabbitMQ(index, agent, values, values[1], "0", "emails", "0", &EmailPath{}, taskID)
		case "FirefoxLogin":
			err = toRabbitMQ(index, agent, values, values[0], values[4], "website_bookmark", values[1], &FirefoxLogin{}, taskID)
		case "IECache":
			err = toRabbitMQ(index, agent, values, values[0], values[3], "cookie_cache", values[1], &IECache{}, taskID)
		case "IELogin":
			err = toRabbitMQ(index, agent, values, values[1], values[4], "website_bookmark", values[2], &IELogin{}, taskID)
		case "NetAdapters":
			err = toRabbitMQ(index, agent, values, values[0], values[11], "software", values[3], &Netadapters{}, taskID)
		case "RecentFile":
			err = toRabbitMQ(index, agent, values, values[2], "0", "document", values[0], &RecentFile{}, taskID)
		case "Shellbags":
			err = toRabbitMQ(index, agent, values, values[0], values[6], "document", values[1], &Shellbags{}, taskID)
		case "SystemInfo":
			err = toRabbitMQ(index, agent, values, values[13], "0", "network_record", values[1], &SystemInfo{}, taskID)
		default:
			logger.Error("Unknown table name (" + agent + "): " + tableName)
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func toRabbitMQ(index string, agent string, values []string, item string, date string, ttype string, etc string, st elastic.Request_data, taskID string) error {
	ip, name, err := query.GetMachineIPandName(agent)
	if err != nil {
		return err
	}
	uuid := uuid.NewString()
	err = rabbitmq.ToRabbitMQ_Details(index, st, values, uuid, agent, ip, name, item, date, ttype, etc, "ed_low", "StartCollect", taskID)
	if err != nil {
		return err
	}
	return nil
}

func toBoolean(b string) string {
	if b == "Yes" {
		return "1"
	} else {
		return "0"
	}
}
