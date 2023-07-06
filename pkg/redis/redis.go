package redis

import (
	"context"
	"edetector_go/config"
	"edetector_go/internal/fflag"
	"fmt"

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
		Addr:     "0.0.0.0:" + config.Viper.GetString("REDIS_PORT"),
		Password: config.Viper.GetString("REDIS_PASSWORD"),
		DB:       config.Viper.GetInt("REDIS_DB"),
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to redis")
		return nil
	}
	return RedisClient
}

func Redis_close() {
	if !checkflag() {
		return
	}
	RedisClient.Close()
}

func Redis_set(key string, value string) {
	if !checkflag() {
		return
	}
	RedisClient.Set(context.Background(), key, value, 0)
}

func Redis_get(key string) string {
	if !checkflag() {
		return ""
	}
	val, _ := RedisClient.Get(context.Background(), key).Result()
	return val
}
