package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/mariadb/method"
	"strconv"
	"strings"
)

func First_detect_info(KeyNum string, message string) string {
	res, err := mariadb.DB.Query("SELECT EXISTS(SELECT * FROM client_setting WHERE client_id = ?)", KeyNum)
	if err != nil {
		logger.Error("Error first detect info query: " + err.Error())
		return ""
	}
	defer res.Close()
	var check int
	for res.Next() {
		err := res.Scan(&check)
		if err != nil {
			logger.Error("Error scan:" + err.Error())
			return ""
		}
	}
	if check == 0 {
		data_splited := strings.Split(message, "|")
		if len(data_splited) < 2 {
			logger.Error("Invalid GiveDetectInfoFirst format")
		}
		_, err = method.Exec(
			"INSERT INTO client_setting (client_id, networkreport, processreport) VALUE (?,?,?) ON DUPLICATE KEY UPDATE networkreport = VALUES(networkreport), processreport = VALUES(processreport);",
			KeyNum, data_splited[1], data_splited[0], KeyNum,
		)
		if err != nil {
			logger.Error("Error insert client_setting: " + err.Error())
		}
		return message
	} else {
		res2, err := mariadb.DB.Query("SELECT networkreport, processreport FROM client_setting WHERE client_id=?", KeyNum)
		if err != nil {
			logger.Error("Error select client_setting: " + err.Error())
			return ""
		}
		defer res2.Close()
		var network int
		var process int
		for res2.Next() {
			err := res2.Scan(&network, &process)
			if err != nil {
				logger.Error("Error scan: " + err.Error())
				return ""
			}
		}
		return strconv.Itoa(process) + "|" + strconv.Itoa(network)
	}
}
