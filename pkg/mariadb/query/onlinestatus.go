package query

import (
	// "edetector_go/pkg/logger"
	// "edetector_go/pkg/mariadb"
	// "edetector_go/pkg/mariadb/method"
)

func Online(KeyNum string, time int64) {
	// _, err := mariadb.DB.Exec("INSERT INTO online_status (client_id, time) VALUE (?,FROM_UNIXTIME(?)) ON DUPLICATE KEY UPDATE client_id = ?", KeyNum, time, KeyNum)
	// if err != nil {
	// 	logger.Error(err.Error())
	// }
}

func Update_time(KeyNum string, time int64) {
	// _, err := method.Exec("UPDATE online_status SET time = FROM_UNIXTIME(?) WHERE client_id = ?", time, KeyNum)
	// if err != nil {
	// 	logger.Error(err.Error())
	// }
}
