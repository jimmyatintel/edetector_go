package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
)

func DeleteAgent(KeyNum string) {
	res, err := mariadb.DB.Query("DELETE FROM client WHERE client_id = ?", KeyNum)
	if err != nil {
		logger.Error("Error delete agent: " + err.Error())
		return
	}
	defer res.Close()
}
