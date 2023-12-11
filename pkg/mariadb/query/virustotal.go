package query

import (
	"database/sql"
	"edetector_go/pkg/mariadb"
)

func StoreVTKey(key string) error {
	query := "UPDATE virustotal SET api_key = ? WHERE id = 1"
	_, err := mariadb.DB.Exec(query, key)
	if err != nil {
		return err
	}
	return nil
}

func LoadVTKey() (string, error) {
	var key string
	query := "SELECT api_key FROM virustotal WHERE id = 1"
	err := mariadb.DB.QueryRow(query).Scan(&key)
	if err != nil {
		if err == sql.ErrNoRows {
			return "null", nil
		}
		return "", err
	}
	return key, nil
}
