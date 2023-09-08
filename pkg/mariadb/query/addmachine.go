package query

import (
	"edetector_go/internal/client/clientinfo"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/mariadb/method"
	"edetector_go/pkg/redis"

	"go.uber.org/zap"
)

func Checkindex(KeyNum string, ip string, mac string) {
	res, err := mariadb.DB.Query("SELECT EXISTS(SELECT * FROM client WHERE client_id = ?)", KeyNum)
	if err != nil {
		logger.Error("error check index query: " + err.Error())
		return
	}
	defer res.Close()
	var check int
	for res.Next() {
		err := res.Scan(&check)
		if err != nil {
			logger.Error("error scan: " + err.Error())
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
	// client_info table
	_, err := method.Exec(
		"INSERT INTO client_info (client_id, sysinfo, osinfo, computername, username, fileversion, boottime) VALUE (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE client_id = ?",
		ClientInfo.KeyNum, ClientInfo.SysInfo, ClientInfo.OsInfo, ClientInfo.ComputerName, ClientInfo.UserName, ClientInfo.FileVersion, ClientInfo.BootTime, ClientInfo.KeyNum,
	)
	if err != nil {
		logger.Error("error add client_info: " + err.Error())
	}
	// client_task_status table
	_, err = method.Exec(
		"INSERT INTO client_task_status (client_id, scan_schedule, scan_finish_time, collect_schedule, collect_finish_time, file_schedule, file_finish_time, image_finish_time) VALUE (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE client_id = ?",
		ClientInfo.KeyNum, nil, nil, nil, nil, nil, nil, nil, ClientInfo.KeyNum,
	)
	if err != nil {
		logger.Error("error add client_task_status: " + err.Error())
	}
	_, err = method.Exec(
		"INSERT INTO client_setting (client_id, networkreport, processreport) VALUE (?,0,0) ON DUPLICATE KEY UPDATE client_id = ?",
		ClientInfo.KeyNum, ClientInfo.KeyNum,
	)
	if err != nil {
		logger.Error("error add client_setting: " + err.Error())
	}
	err = redis.Offline(ClientInfo.KeyNum)
	if err != nil {
		logger.Error("Update Offline failed:", zap.Any("error", err.Error()))
	}
}
