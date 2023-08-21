package redis

import (
	"context"
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/logger"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func checkflag() bool {
	if enable, err := fflag.FFLAG.FeatureEnabled("redis_enable"); enable && err == nil {
		return true
	}
	return false
}
func Redis_init() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Viper.GetString("REDIS_HOST") + ":" + config.Viper.GetString("REDIS_PORT"),
		Password: config.Viper.GetString("REDIS_PASSWORD"),
		DB:       config.Viper.GetInt("REDIS_DB"),
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("Error connecting to redis")
		return nil
	}
	logger.Info("redis is enabled")
	return RedisClient
}

func RedisClose() {
	if !checkflag() {
		return
	}
	RedisClient.Close()
}

func RedisSet(key string, value interface{}) error {
	if !checkflag() {
		return nil
	}
	return RedisClient.Set(context.Background(), key, value, 0).Err()
}

func RedisSetAdd(key string, value interface{}) error {
	if !checkflag() {
		return nil
	}
	return RedisClient.Set(context.Background(), key, value, 0).Err()
}

func RedisGetString(key string) string {
	if !checkflag() {
		return ""
	}
	val, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		logger.Error("Error getting value from redis " + err.Error())
		return ""
	}
	return val
}

func RedisGetInt(key string) int {
	if !checkflag() {
		return 0
	}
	val, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		logger.Error("Error getting value from redis " + err.Error())
		return 0
	}
	val_int, err := strconv.Atoi(val)
	if err != nil {
		logger.Error("Error converting to integer: " + err.Error())
		return 0
	}
	return val_int
}
