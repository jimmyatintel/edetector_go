package query

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
)

func Load_key_image(ttype string) ([][]string, error) {
	qu := "SELECT apptype, path, keyword FROM key_image where type = " + ttype
	var result [][]string
	res, err := mariadb.DB.Query(qu)
	if err != nil {
		return result, err
	}
	defer res.Close()
	l, _ := res.Columns()
	for res.Next() {
		tmp := make([]string, len(l))
		err := res.Scan(&tmp[0], &tmp[1], &tmp[2])
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}
	if len(result) == 0 {
		logger.Info("Key image type " + ttype + " is empty")
	}
	return result, nil
}
