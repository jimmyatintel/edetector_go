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
	logger.Info("Handling table: " + tableName)
	rows, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		logger.Error("Error getting rows: " + err.Error())
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
			err = toRabbitMQ(index, agent, values, values[1], values[19], "software", values[14], &AppResourceUsageMonitor{})
		case "ARPCache":
			err = toRabbitMQ(index, agent, values, values[1], "0", "volatile", values[2], &ARPCache{})
		case "BaseService":
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[5], &BaseService{})
		case "ChromeBookmarks":
			err = toRabbitMQ(index, agent, values, values[4], values[6], "website_bookmark", values[3], &ChromeBookmarks{})
		case "ChromeCache":
			RFCToTimestamp(&values)
			err = toRabbitMQ(index, agent, values, values[1], values[8], "cookie_cache", values[2], &ChromeCache{})
		case "ChromeDownload":
			err = toRabbitMQ(index, agent, values, values[0], values[6], "website_bookmark", values[3], &ChromeDownload{})
		case "ChromeHistory":
			err = toRabbitMQ(index, agent, values, values[0], values[2], "website_bookmark", values[1], &ChromeHistory{})
		case "ChromeKeywordSearch":
			err = toRabbitMQ(index, agent, values, values[0], "0", "website_bookmark", "", &ChromeKeywordSearch{})
		case "ChromeLogin":
			err = toRabbitMQ(index, agent, values, values[0], values[6], "website_bookmark", values[3], &ChromeLogin{})
		case "DNSInfo":
			err = toRabbitMQ(index, agent, values, values[9], "0", "software", values[6], &DNSInfo{})
		case "EdgeBookmarks":
			err = toRabbitMQ(index, agent, values, values[3], values[7], "website_bookmark", values[4], &EdgeBookmarks{})
		case "EdgeCache":
			RFCToTimestamp(&values)
			err = toRabbitMQ(index, agent, values, values[1], values[10], "cookie_cache", values[2], &EdgeCache{})
		case "EdgeCookies":
			err = toRabbitMQ(index, agent, values, values[3], values[7], "cookie_cache", values[2], &EdgeCookies{})
		case "EdgeHistory":
			err = toRabbitMQ(index, agent, values, values[1], values[5], "website_bookmark", values[2], &EdgeHistory{})
		case "EdgeLogin":
			err = toRabbitMQ(index, agent, values, values[1], values[7], "website_bookmark", values[4], &EdgeLogin{})
		case "EventApplication":
			err = toRabbitMQ(index, agent, values, values[3], values[9], "software", values[17], &EventApplication{})
		case "EventSecurity":
			err = toRabbitMQ(index, agent, values, values[3], values[9], "usb", values[17], &EventSecurity{})
		case "EventSystem":
			err = toRabbitMQ(index, agent, values, values[3], values[9], "usb", values[17], &EventSystem{})
		case "FirefoxBookmarks":
			err = toRabbitMQ(index, agent, values, values[8], values[5], "website_bookmark", values[3], &FirefoxBookmarks{})
		case "FirefoxCache":
			err = toRabbitMQ(index, agent, values, values[1], values[8], "cookie_cache", values[2], &FirefoxCache{})
		case "FirefoxCookies":
			err = toRabbitMQ(index, agent, values, values[1], values[5], "cookie_cache", values[3], &FirefoxCookies{})
		case "FirefoxHistory":
			err = toRabbitMQ(index, agent, values, values[0], values[9], "website_bookmark", values[1], &FirefoxHistory{})
		case "IEHistory":
			err = toRabbitMQ(index, agent, values, values[0], values[4], "website_bookmark", values[1], &IEHistory{})
		case "InstalledSoftware":
			DigitToTimestamp(&values)
			err = toRabbitMQ(index, agent, values, values[0], values[17], "network_record", values[6], &InstalledSoftware{})
		case "JumpList":
			err = toRabbitMQ(index, agent, values, values[0], values[5], "software", values[1], &JumpList{})
		case "MUICache":
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[1], &MUICache{})
		case "Network":
			err = toRabbitMQ(index, agent, values, values[1], "0", "volatile", values[4], &Network{})
		case "NetworkDataUsageMonitor":
			err = toRabbitMQ(index, agent, values, values[1], values[10], "software", values[5], &NetworkDataUsageMonitor{})
		case "NetworkResources":
			err = toRabbitMQ(index, agent, values, values[0], "0", "network_record", values[8], &NetworkResources{})
		case "OpenedFiles":
			err = toRabbitMQ(index, agent, values, values[1], "0", "volatile", values[0], &OpenedFiles{})
		case "Prefetch":
			err = toRabbitMQ(index, agent, values, values[1], values[2], "software", values[3], &Prefetch{})
		case "Process":
			err = toRabbitMQ(index, agent, values, values[1], values[3], "volatile", values[4], &Process{})
		case "Service":
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[5], &Service{})
		case "Shortcuts":
			err = toRabbitMQ(index, agent, values, values[0], values[10], "document", values[2], &Shortcuts{})
		case "StartRun":
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[1], &StartRun{})
		case "TaskSchedule":
			err = toRabbitMQ(index, agent, values, values[0], values[3], "software", values[1], &TaskSchedule{})
		case "USBdevices":
			err = toRabbitMQ(index, agent, values, values[1], values[14], "usb", values[10], &USBdevices{})
		case "UserAssist":
			err = toRabbitMQ(index, agent, values, values[0], values[5], "software", values[2], &UserAssist{})
		case "UserProfiles":
			err = toRabbitMQ(index, agent, values, values[0], values[6], "document", values[2], &UserProfiles{})
		case "WindowsActivity":
			err = toRabbitMQ(index, agent, values, values[1], values[15], "document", values[3], &WindowsActivity{})
		default:
			logger.Error("Unknown table name: " + tableName)
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func toRabbitMQ(index string, agent string, values []string, item string, date string, ttype string, etc string, st elastic.Request_data) error {
	ip, name := query.GetMachineIPandName(agent)
	uuid := uuid.NewString()
	err := rabbitmq.ToRabbitMQ_Main(index, uuid, agent, ip, name, item, date, ttype, etc, "ed_low")
	if err != nil {
		return err
	}
	err = rabbitmq.ToRabbitMQ_Details(index, st, values, uuid, agent, ip, name, item, date, ttype, etc, "ed_low")
	if err != nil {
		return err
	}
	return nil
}
