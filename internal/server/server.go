package server

import (
	"context"
	config "edetector_go/config"
	Client "edetector_go/internal/clientsearch"
	"edetector_go/pkg/elastic"
	fflag "edetector_go/pkg/fflag"
	logger "edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/rabbitmq"

	"edetector_go/pkg/redis"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func server_init() {
	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Panic("Error loading feature flag")
		panic("Error loading feature flag")
	}
	vp, err := config.LoadConfig()
	if vp == nil {
		logger.Panic("Error loading config file: " + err.Error())
		panic(err)
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("WORKER_LOG_FILE"), "server", "SERVER")
		logger.Info("Logger is enabled please check all out info in log file: " + config.Viper.GetString("WORKER_LOG_FILE"))
	}
	connString, err := mariadb.Connect_init()
	if err != nil {
		logger.Panic("Error connecting to mariadb: " + err.Error())
		panic(err)
	} else {
		logger.Info("Mariadb connectionString: " + connString)
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("redis_enable"); enable && err == nil {
		if db := redis.Redis_init(); db == nil {
			logger.Panic("Error connecting to redis")
			panic(err)
		}
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("rabbit_enable"); enable && err == nil {
		rabbitmq.Rabbit_init()
		logger.Info("Rabbit is enabled.")
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("elastic_enable"); enable && err == nil {
		elastic.Elastic_init()
		logger.Info("Elastic is enabled.")
	}
}

func Main(version string) {
	server_init()
	logger.Info("Welcome to edetector main server: " + version)
	Quit := make(chan os.Signal, 1)
	Connection_close := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go Client.Main(ctx, Connection_close)
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	<-Quit
	cancel()
	// taskservice.Stop()
	logger.Info("Server is shutting down...")
	servershutdown()
	select {
	case <-Connection_close:
		logger.Info("Connection closed")
	case <-time.After(5 * time.Second):
		logger.Info("Connection close fail, force shutdown after 5 seconds")
	}
	logger.Info("Server shutdown complete")
	defer cancel()
}

func servershutdown() {
	// rabbitmq.Connection_close()
	for _, client := range Client.Clientlist {
		redis.Offline(client)
	}
	redis.RedisClose()
}
