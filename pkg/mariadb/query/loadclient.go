package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"strings"
)

// return a list of client
func Load_all_client() []string {
	var clients []string
	res, err := mariadb.DB.Query("SELECT client_id FROM client")
	if err != nil {
		logger.Error("Error loading client: " + err.Error())
		return clients
	}
	defer res.Close()
	for res.Next() {
		var client string
		err = res.Scan(&client)
		if err != nil {
			logger.Error("Error loading client: " + err.Error())
			return clients
		}
		clients = append(clients, client)
	}
	return clients
}

func Get_client_os(client string) string {
	os := ""
	res, err := mariadb.DB.Query("SELECT osinfo FROM client_info WHERE client_id = ?", client)
	if err != nil {
		logger.Error("Error loading client os: " + err.Error())
		return os
	}
	defer res.Close()
	for res.Next() {
		err = res.Scan(&os)
		if err != nil {
			logger.Error("Error loading client os: " + err.Error())
			return os
		}
	}
	if strings.Contains(os, "Windows") {
		return "windows"
	} else if strings.Contains(os, "Ubuntu") {
		return "ubuntu"
	}
	return os
}
