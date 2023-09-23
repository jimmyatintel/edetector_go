package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
)

func GetMachineIPandName(KeyNum string) (string, string) {
	agentIP := GetMachineIP(KeyNum)
	if agentIP == "" {
		logger.Error("Error getting machine ip")
	}
	agentName := GetMachineName(KeyNum)
	if agentName == "" {
		logger.Error("Error getting machine name")
	}
	return agentIP, agentName
}

func GetMachineIP(KeyNum string) (ip string) {
	res, err := mariadb.DB.Query("SELECT ip FROM client WHERE client_id = ?", KeyNum)
	if err != nil {
		logger.Error("Error getting machine ip" + err.Error())
		return ""
	}
	defer res.Close()
	for res.Next() {
		err := res.Scan(&ip)
		if err != nil {
			logger.Error("Error scan: " + err.Error())
			return ""
		}
	}
	return ip
}

func GetMachineMAC(KeyNum string) (mac string) {
	res, err := mariadb.DB.Query("SELECT mac FROM client WHERE client_id = ?", KeyNum)
	if err != nil {
		logger.Error("Error getting machine mac" + err.Error())
		return ""
	}
	defer res.Close()
	for res.Next() {
		err := res.Scan(&mac)
		if err != nil {
			logger.Error("Error scan: " + err.Error())
			return ""
		}
	}
	return mac
}

func GetMachineName(KeyNum string) (name string) {
	res, err := mariadb.DB.Query("SELECT computername FROM client_info WHERE client_id = ?", KeyNum)
	if err != nil {
		logger.Error("Error getting machine name" + err.Error())
		return ""
	}
	defer res.Close()
	for res.Next() {
		err := res.Scan(&name)
		if err != nil {
			logger.Error("Error scan: " + err.Error())
			return ""
		}
	}
	return name
}
