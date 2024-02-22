package query

import (
	"edetector_go/pkg/mariadb"
)

func Load_white_list() ([][]string, error) {
	qu := "SELECT filename, md5, signature, path FROM white_list"
	var result [][]string
	res, err := mariadb.DB.Query(qu)
	if err != nil {
		return result, err
	}
	defer res.Close()
	l, _ := res.Columns()
	for res.Next() {
		tmp := make([]string, len(l))
		err := res.Scan(&tmp[0], &tmp[1], &tmp[2], &tmp[3])
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}
	// if len(result) == 0 {
	// 	logger.Info("White list is empty")
	// }
	return result, nil
}

func Load_black_list() ([][]string, error) {
	qu := "SELECT filename, md5, signature, path FROM black_list"
	var result [][]string
	res, err := mariadb.DB.Query(qu)
	if err != nil {
		return result, err
	}
	defer res.Close()
	l, _ := res.Columns()
	for res.Next() {
		tmp := make([]string, len(l))
		err := res.Scan(&tmp[0], &tmp[1], &tmp[2], &tmp[3])
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}
	// if len(result) == 0 {
	// 	logger.Info("Black list is empty")
	// }
	return result, nil
}

func Load_hack_list() ([][]string, error) {
	qu := "SELECT process_name, cmd, path, adding_point FROM hack_list"
	var result [][]string
	res, err := mariadb.DB.Query(qu)
	if err != nil {
		return result, err
	}
	defer res.Close()
	l, _ := res.Columns()
	for res.Next() {
		tmp := make([]string, len(l))
		err := res.Scan(&tmp[0], &tmp[1], &tmp[2], &tmp[3])
		if err != nil {
			return result, err
		}
		result = append(result, tmp)
	}
	// if len(result) == 0 {
	// 	logger.Info("Hack list is empty")
	// }
	return result, nil
}
