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
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func server_init() {
	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Error("Error loading feature flag")
		return
	}
	vp := config.LoadConfig()
	if vp == nil {
		logger.Error("Error loading config file")
		return
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(os.Getenv("WORKER_LOG_FILE"))
		logger.Info("logger is enabled please check all out info in log file: ", zap.Any("message", os.Getenv("WORKER_LOG_FILE")))
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("redis_enable"); enable && err == nil {
		if db := redis.Redis_init(); db == nil {
			logger.Error("Error connecting to redis")
		}
	}
	if err := mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())

	}
	if enable, err := fflag.FFLAG.FeatureEnabled("rabbit_enable"); enable && err == nil {
		rabbitmq.Rabbit_init()
		logger.Info("rabbit is enabled.")
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("elastic_enable"); enable && err == nil {
		err := elastic.SetElkClient()
		if err != nil {
			logger.Error("Error connecting to elastic: " + err.Error())
		}
		logger.Info("elastic is enabled.")
	}
}

func Main() {
	server_init()
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
		err := redis.Offline(client)
		if err != nil {
			logger.Error("Update offline failed:", zap.Any("error", err.Error()))
		}
		logger.Info("offline ", zap.Any("message", client))
		taskservice.RequestToUser(client)
	}
	redis.Redis_close()
}
