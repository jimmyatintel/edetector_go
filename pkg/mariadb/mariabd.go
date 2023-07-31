package mariadb

import (
	"database/sql"
	"edetector_go/pkg/logger"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

var DB *sql.DB

func Connect_init() error {
	var err error
	dbUser := os.Getenv("MARIADB_USER")
	dbPass := os.Getenv("MARIADB_PASSWORD")
	dbHost := os.Getenv("MARIADB_HOST")
	dbPort := os.Getenv("MARIADB_PORT")
	dbName := os.Getenv("MARIADB_DATABASE")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	logger.Info("connectionString", zap.Any("message", connectionString))
	DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	return nil
}
