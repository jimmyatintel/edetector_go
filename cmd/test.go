package main

import (
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
)

func init() {
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
	if enable, err := fflag.FFLAG.FeatureEnabled("redis_enable"); enable && err == nil {
		if db := redis.Redis_init(); db == nil {
			logger.Error("Error connecting to redis")
		}
	}
}

func main() {
	redis.RedisSet("test", 1)
	redis.RedisSet_AddString("test", 2)
}
