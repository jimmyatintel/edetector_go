package server

import (
	"context"
	config "edetector_go/config"
	Client "edetector_go/internal/clientsearch"
	fflag "edetector_go/internal/fflag"
	"edetector_go/internal/taskservice"
	"edetector_go/pkg/elastic"
	logger "edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/rabbitmq"

	"edetector_go/pkg/redis"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Main() {
	serverinit()
	Quit := make(chan os.Signal, 1)
	Connection_close := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go Client.Main(ctx, Connection_close)
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	<-Quit
	cancel()
	taskservice.Stop()
	fmt.Println("Server is shutting down...")
	servershutdown()
	select {
	case <-Connection_close:
		logger.Info("Connection closed")
	case <-time.After(5 * time.Second):
		logger.Info("Connection close fail, force shutdown after 5 seconds")
	}
	fmt.Println("Server shutdown complete.")

	defer cancel()
}
func servershutdown() {
	// rabbitmq.Connection_close()
	redis.Redis_close()
}
func serverinit() {
	fflag.Connect_GB()
	if fflag.GB == nil {
		fmt.Println("Error loading feature flag")
		return
	}
	vp := config.LoadConfig()
	if vp == nil {
		fmt.Println("Error loading config file")
		return
	}
	if fflag.GB.Feature("logger_enable").On {
		logger.InitLogger(config.Viper.GetString("WORKER_LOG_FILE"))
		fmt.Println("logger is enabled please check all out info in log file: ", config.Viper.GetString("WORKER_LOG_FILE"))
	}
	if fflag.GB.Feature("redis_enable").On {
		if db := redis.Redis_init(); db == nil {
			logger.Error("Error connecting to redis")
		}
	}

	if err := mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())

	}
	if fflag.GB.Feature("rabbit_enable").On {
		rabbitmq.Rabbit_init()
		fmt.Println("rabbit is enabled.")
	}
	if fflag.GB.Feature("elastic_enable").On {
		err := elastic.SetElkClient()
		if err != nil {
			logger.Error("Error connecting to elastic: " + err.Error())
		}
		fmt.Println("elastic is enabled.")
	}
}
