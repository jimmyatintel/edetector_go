package mariadb

import (
	"database/sql"
	"edetector_go/config"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect_init() error {
	var err error
	dbUser := config.Viper.GetString("MARIADB_USER")
	dbPass := config.Viper.GetString("MARIADB_PASSWORD")
	dbHost := config.Viper.GetString("MARIADB_HOST")
	dbPort := config.Viper.GetInt("MARIADB_PORT")
	dbName := config.Viper.GetString("MARIADB_DATABASE")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	fmt.Println("connectionString", connectionString)
	DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	return nil
}
