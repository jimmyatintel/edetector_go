package main

import (
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"fmt"

	"go.uber.org/zap"
)

func main() {
	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Error("Error loading feature flag")
		return
	}
	vp, err := config.LoadConfig()
	if vp == nil {
		logger.Error("Error loading config file", zap.Any("error", err.Error()))
		return
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("WORKER_LOG_FILE"), "server", "SERVER")
		logger.Info("logger is enabled please check all out info in log file: ", zap.Any("message", config.Viper.GetString("WORKER_LOG_FILE")))
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("redis_enable"); enable && err == nil {
		if db := redis.Redis_init(); db == nil {
			logger.Error("Error connecting to redis")
		}
	}
	key := "569a2191ae414802a5a72bc0b8e0bd1e"
	// fmt.Println(redis.RedisGetInt(key+"-DriveTotal"))
	driveProgress := int((float64(redis.RedisGetInt(key+"-DriveCount"))/float64(redis.RedisGetInt(key+"-DriveTotal")))*100 + float64(redis.RedisGetInt(key+"-ExplorerProgress"))/float64(redis.RedisGetInt(key+"-DriveTotal")))
	fmt.Println(driveProgress)
}
