package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"

	"go.uber.org/zap"
)

func Update_progress(progress int, clientid string, tasktype string) int {
	qu := "update task set progress = ? where client_id = ? and type = ? and status = 2"
	result, err := mariadb.DB.Exec(qu, progress, clientid, tasktype)
	if err != nil {
		logger.Error("update failed: ", zap.Any("error", err.Error()))
		return 0
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("get rows failed: ", zap.Any("error", err.Error()))
		return 0
	}
	return int(rowsAffected)
}
