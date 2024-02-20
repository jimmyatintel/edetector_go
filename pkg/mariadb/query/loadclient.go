package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
)
// return a list of client
func Load_all_client() []string{
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