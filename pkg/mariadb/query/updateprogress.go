package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
)

func Update_progress(progress int, clientid string, tasktype string) {
	qu := "update task set progress = ? where client_id = ? and type = ?"
	_, err := mariadb.DB.Exec(qu, progress, clientid, tasktype)
	if err != nil {
		logger.Error(err.Error())
	}
}
