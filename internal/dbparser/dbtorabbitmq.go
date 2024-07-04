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

		index := config.Viper.GetString("ELASTIC_PREFIX") + "_" + "collection" //! developing
		category := strings.ToLower(tableName)
		taskData, err := query.Load_stored_task(taskID, "nil", -1, "nil")
		if err != nil {
			return err
		}

		switch tableName {
		case "AppResourceUsageMonitor":
			err = toRabbitMQ(index, agent, values, values[1], values[19], "software", values[14], &Collect_AppResourceUsageMonitor{}, &AppResourceUsageMonitor{}, taskID, category, taskData[0][6])
		case "ARPCache":
			err = toRabbitMQ(index, agent, values, values[1], "0", "volatile", values[2], &Collect_ARPCache{}, &ARPCache{}, taskID, category, taskData[0][6])
		case "BaseService":
			values[14] = toBoolean(values[14])
			values[15] = toBoolean(values[15])
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[5], &Collect_BaseService{}, &BaseService{}, taskID, category, taskData[0][6])
		case "ChromeBookmarks":
			err = toRabbitMQ(index, agent, values, values[4], values[6], "website_bookmark", values[3], &Collect_ChromeBookmarks{}, &ChromeBookmarks{}, taskID, category, taskData[0][6])
		case "ChromeCache":
			values[8] = RFCToTimestamp(values[8])
			values[9] = RFCToTimestamp(values[9])
			values[10] = RFCToTimestamp(values[10])
			err = toRabbitMQ(index, agent, values, values[1], values[8], "cookie_cache", values[2], &Collect_ChromeCache{}, &ChromeCache{}, taskID, category, taskData[0][6])
		case "ChromeDownload":
			values[11] = toBoolean(values[11])
			values[17] = RFCToTimestamp(values[17])
			err = toRabbitMQ(index, agent, values, values[0], values[6], "website_bookmark", values[3], &Collect_ChromeDownload{}, &ChromeDownload{}, taskID, category, taskData[0][6])
		case "ChromeHistory":
			err = toRabbitMQ(index, agent, values, values[0], values[2], "website_bookmark", values[1], &Collect_ChromeHistory{}, &ChromeHistory{}, taskID, category, taskData[0][6])
		case "ChromeKeywordSearch":
			err = toRabbitMQ(index, agent, values, values[0], "0", "website_bookmark", "", &Collect_ChromeKeywordSearch{}, &ChromeKeywordSearch{}, taskID, category, taskData[0][6])
		case "ChromeLogin":
			err = toRabbitMQ(index, agent, values, values[0], values[4], "website_bookmark", values[3], &Collect_ChromeLogin{}, &ChromeLogin{}, taskID, category, taskData[0][6])
		case "DNSInfo":
			err = toRabbitMQ(index, agent, values, values[9], "0", "software", values[6], &Collect_DNSInfo{}, &DNSInfo{}, taskID, category, taskData[0][6])
		case "EdgeBookmarks":
			err = toRabbitMQ(index, agent, values, values[3], values[7], "website_bookmark", values[4], &Collect_EdgeBookmarks{}, &EdgeBookmarks{}, taskID, category, taskData[0][6])
		case "EdgeCache":
			values[8] = RFCToTimestamp(values[8])
			values[9] = RFCToTimestamp(values[9])
			values[10] = RFCToTimestamp(values[10])
			err = toRabbitMQ(index, agent, values, values[1], values[10], "cookie_cache", values[2], &Collect_EdgeCache{}, &EdgeCache{}, taskID, category, taskData[0][6])
		case "EdgeCookies":
			err = toRabbitMQ(index, agent, values, values[3], values[7], "cookie_cache", values[2], &Collect_EdgeCookies{}, &EdgeCookies{}, taskID, category, taskData[0][6])
		case "EdgeHistory":
			err = toRabbitMQ(index, agent, values, values[1], values[5], "website_bookmark", values[2], &Collect_EdgeHistory{}, &EdgeHistory{}, taskID, category, taskData[0][6])
		case "EdgeLogin":
			err = toRabbitMQ(index, agent, values, values[1], values[5], "website_bookmark", values[4], &Collect_EdgeLogin{}, &EdgeLogin{}, taskID, category, taskData[0][6])
		case "EventApplication":
			err = toRabbitMQ(index, agent, values, values[3], values[9], "software", values[17], &Collect_EventApplication{}, &EventApplication{}, taskID, category, taskData[0][6])
		case "EventSecurity":
			err = toRabbitMQ(index, agent, values, values[3], values[9], "usb", values[17], &Collect_EventSecurity{}, &EventSecurity{}, taskID, category, taskData[0][6])
		case "EventSystem":
			err = toRabbitMQ(index, agent, values, values[3], values[9], "usb", values[17], &Collect_EventSystem{}, &EventSystem{}, taskID, category, taskData[0][6])
		case "FirefoxBookmarks":
			err = toRabbitMQ(index, agent, values, values[8], values[5], "website_bookmark", values[3], &Collect_FirefoxBookmarks{}, &FirefoxBookmarks{}, taskID, category, taskData[0][6])
		case "FirefoxCache":
			err = toRabbitMQ(index, agent, values, values[1], values[8], "cookie_cache", values[2], &Collect_FirefoxCache{}, &FirefoxCache{}, taskID, category, taskData[0][6])
		case "FirefoxCookies":
			err = toRabbitMQ(index, agent, values, values[1], values[5], "cookie_cache", values[3], &Collect_FirefoxCookies{}, &FirefoxCookies{}, taskID, category, taskData[0][6])
		case "FirefoxHistory":
			err = toRabbitMQ(index, agent, values, values[0], values[9], "website_bookmark", values[1], &Collect_FirefoxHistory{}, &FirefoxHistory{}, taskID, category, taskData[0][6])
		case "IEHistory":
			err = toRabbitMQ(index, agent, values, values[0], values[4], "website_bookmark", values[1], &Collect_IEHistory{}, &IEHistory{}, taskID, category, taskData[0][6])
		case "InstalledSoftware":
			values[3] = DigitToTimestamp(values[3])
			err = toRabbitMQ(index, agent, values, values[0], values[17], "network_record", values[6], &Collect_InstalledSoftware{}, &InstalledSoftware{}, taskID, category, taskData[0][6])
		case "JumpList":
			err = toRabbitMQ(index, agent, values, values[0], values[5], "software", values[1], &Collect_JumpList{}, &JumpList{}, taskID, category, taskData[0][6])
		case "MUICache":
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[1], &Collect_MUICache{}, &MUICache{}, taskID, category, taskData[0][6])
		case "Network":
			err = toRabbitMQ(index, agent, values, values[1], "0", "volatile", values[4], &Collect_Network{}, &Network{}, taskID, category, taskData[0][6])
		case "NetworkDataUsageMonitor":
			err = toRabbitMQ(index, agent, values, values[1], values[10], "software", values[5], &Collect_NetworkDataUsageMonitor{}, &NetworkDataUsageMonitor{}, taskID, category, taskData[0][6])
		case "NetworkResources":
			err = toRabbitMQ(index, agent, values, values[0], "0", "network_record", values[8], &Collect_NetworkResources{}, &NetworkResources{}, taskID, category, taskData[0][6])
		case "OpenedFiles":
			err = toRabbitMQ(index, agent, values, values[1], "0", "volatile", values[0], &Collect_OpenedFiles{}, &OpenedFiles{}, taskID, category, taskData[0][6])
		case "Prefetch":
			err = toRabbitMQ(index, agent, values, values[1], values[2], "software", values[3], &Collect_Prefetch{}, &Prefetch{}, taskID, category, taskData[0][6])
		case "Process":
			err = toRabbitMQ(index, agent, values, values[1], values[3], "volatile", values[4], &Collect_Process{}, &Process{}, taskID, category, taskData[0][6])
		case "Service":
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[5], &Collect_Service{}, &Service{}, taskID, category, taskData[0][6])
		case "Shortcuts":
			values[7] = toBoolean(values[7])
			err = toRabbitMQ(index, agent, values, values[0], values[10], "document", values[2], &Collect_Shortcuts{}, &Shortcuts{}, taskID, category, taskData[0][6])
		case "StartRun":
			err = toRabbitMQ(index, agent, values, values[0], "0", "software", values[1], &Collect_StartRun{}, &StartRun{}, taskID, category, taskData[0][6])
		case "TaskSchedule":
			err = toRabbitMQ(index, agent, values, values[0], values[3], "software", values[1], &Collect_TaskSchedule{}, &TaskSchedule{}, taskID, category, taskData[0][6])
		case "USBdevices":
			err = toRabbitMQ(index, agent, values, values[1], values[14], "usb", values[10], &Collect_USBdevices{}, &USBdevices{}, taskID, category, taskData[0][6])
		case "UserAssist":
			err = toRabbitMQ(index, agent, values, values[0], values[5], "software", values[2], &Collect_UserAssist{}, &UserAssist{}, taskID, category, taskData[0][6])
		case "UserProfiles":
			values[3] = toBoolean(values[3])
			err = toRabbitMQ(index, agent, values, values[0], values[6], "document", values[2], &Collect_UserProfiles{}, &UserProfiles{}, taskID, category, taskData[0][6])
		case "WindowsActivity":
			err = toRabbitMQ(index, agent, values, values[1], values[15], "document", values[3], &Collect_WindowsActivity{}, &WindowsActivity{}, taskID, category, taskData[0][6])
		case "Wireless":
			err = toRabbitMQ(index, agent, values, values[0], values[8], "network_record", values[1], &Collect_Wireless{}, &Wireless{}, taskID, category, taskData[0][6])
		case "Email":
			values[5] = RFCToTimestamp(values[5])
			values[6] = RFCToTimestamp(values[6])
			err = toRabbitMQ(index, agent, values, values[9], values[5], "emails", values[3], &Collect_Email{}, &Email{}, taskID, category, taskData[0][6])
		case "EmailPath":
			err = toRabbitMQ(index, agent, values, values[1], "0", "emails", "0", &Collect_EmailPath{}, &EmailPath{}, taskID, category, taskData[0][6])
		case "FirefoxLogin":
			err = toRabbitMQ(index, agent, values, values[0], values[4], "website_bookmark", values[1], &Collect_FirefoxLogin{}, &FirefoxLogin{}, taskID, category, taskData[0][6])
		case "IECache":
			err = toRabbitMQ(index, agent, values, values[0], values[3], "cookie_cache", values[1], &Collect_IECache{}, &IECache{}, taskID, category, taskData[0][6])
		case "IELogin":
			err = toRabbitMQ(index, agent, values, values[1], values[4], "website_bookmark", values[2], &Collect_IELogin{}, &IELogin{}, taskID, category, taskData[0][6])
		case "NetAdapters":
			err = toRabbitMQ(index, agent, values, values[0], values[11], "software", values[3], &Collect_Netadapters{}, &Netadapters{}, taskID, category, taskData[0][6])
		case "RecentFile":
			err = toRabbitMQ(index, agent, values, values[2], "0", "document", values[0], &Collect_RecentFile{}, &RecentFile{}, taskID, category, taskData[0][6])
		case "ShellBags":
			err = toRabbitMQ(index, agent, values, values[0], values[6], "document", values[1], &Collect_Shellbags{}, &Shellbags{}, taskID, category, taskData[0][6])
		case "SystemInfo":
			err = toRabbitMQ(index, agent, values, values[13], "0", "network_record", values[1], &Collect_SystemInfo{}, &SystemInfo{}, taskID, category, taskData[0][6])
		case "ChromeCookies":
			err = toRabbitMQ(index, agent, values, values[3], values[7], "cookie_cache", values[2], &Collect_ChromeCookies{}, &ChromeCookies{}, taskID, category, taskData[0][6])
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

func toRabbitMQ(index string, agent string, values []string, item string, date string, ttype string, etc string, st elastic.Request_data, sub_st elastic.Request_data, taskID string, category string, taskTimestamp string) error {
	ip, name, err := query.GetMachineIPandName(agent)
	if err != nil {
		return err
	}
	uuid := uuid.NewString()
	if item == "-1" { // empty table -> not insert
		return nil
	} else if item == "-2" { // failed table -> update mariadb
		query.Add_fail_table(agent, category)
		return nil
	}
	err = rabbitmq.ToRabbitMQ_Details(index, st, sub_st, values, uuid, agent, ip, name, item, date, ttype, etc, "ed_low_collect", "StartCollect", taskID, category, taskTimestamp)
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
