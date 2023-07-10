package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
)

func GetMachineIP(KeyNum string) (ip string) {
	res, err := mariadb.DB.Query("SELECT ip FROM client WHERE client_id = ?", KeyNum)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer res.Close()
	for res.Next() {
		err := res.Scan(&ip)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
	return ip
}

func GetMachineMAC(KeyNum string) (mac string) {
	res, err := mariadb.DB.Query("SELECT mac FROM client WHERE client_id = ?", KeyNum)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer res.Close()
	for res.Next() {
		err := res.Scan(&mac)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
	return mac
}

