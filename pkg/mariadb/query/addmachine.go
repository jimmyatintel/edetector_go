package query

import (
	"edetector_go/internal/client/clientinfo"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/redis"
	"edetector_go/pkg/mariadb/method"
	"encoding/json"
)

type ClientOnlineStatus struct {
	Status int
	Time   string
}

func (c *ClientOnlineStatus) Marshal() string {
	json, err := json.Marshal(c)
	if err != nil {
		logger.Error("Error in json marshal" + err.Error())
	}
	return string(json)
}

func Checkindex(KeyNum string, ip string, mac string) {
	res, err := mariadb.DB.Query("SELECT EXISTS(SELECT * FROM client WHERE client_id = ?)", KeyNum)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer res.Close()
	var check int
	for res.Next() {
		err := res.Scan(&check)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
	if check == 0 {
		method.Exec(
			"INSERT INTO client (client_id, ip, mac) VALUE (?,?,?) ON DUPLICATE KEY UPDATE client_id = ?",
			KeyNum, ip, mac, KeyNum,
		)
	}
}

func Addmachine(ClientInfo clientinfo.ClientInfo) {
	_, err := method.Exec(
		"INSERT INTO client_info (client_id, sysinfo, osinfo, computername, username, fileversion, boottime) VALUE (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE client_id = ?",
		ClientInfo.KeyNum, ClientInfo.SysInfo, ClientInfo.OsInfo, ClientInfo.ComputerName, ClientInfo.UserName, ClientInfo.FileVersion, ClientInfo.BootTime, ClientInfo.KeyNum,
	)
	if err != nil {
		logger.Error(err.Error())
	}
	_, err = method.Exec(
		"INSERT INTO client_task_status (client_id, scan_schedule, scan_finish_time, collect_schedule, collect_finish_time, file_schedule, file_finish_time, image_finish_time) VALUE (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE client_id = ?",
		ClientInfo.KeyNum, "null", "null", "null", "null", "null", "null", "null", ClientInfo.KeyNum,
	)
	if err != nil {
		logger.Error(err.Error())
	}
	onlineStatusInfo := ClientOnlineStatus {
		Status: 1,
		Time:   "20230706110846",
	}
	redis.Redis_set(ClientInfo.KeyNum, onlineStatusInfo.Marshal())
}
