package redis

import (
	"context"
	"edetector_go/internal/fflag"
	"edetector_go/pkg/logger"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisClient *redis.Client

func checkflag() bool {
	if enable, err := fflag.FFLAG.FeatureEnabled("redis_enable"); enable && err == nil {
		return true
	}
	return false
}

func Redis_init() *redis.Client {
	redis_db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logger.Error("REDIS_DB is not set", zap.Any("error", err))
		redis_db = 0
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redis_db,
	})
	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("Error connecting to redis")
		return nil
	}
	logger.Info("redis is enabled")
	return RedisClient
}

func Redis_close() {
	if !checkflag() {
		return
	}
	RedisClient.Close()
}

func Redis_set(key string, value string) error {
	if !checkflag() {
		return nil
	}
	return RedisClient.Set(context.Background(), key, value, 0).Err()
}

func Redis_get(key string) string {
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
