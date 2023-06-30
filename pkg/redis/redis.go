package redis

import (
	"context"
	"edetector_go/config"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Redis_init() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Viper.GetString("REDIS_HOST") + ":" + config.Viper.GetString("REDIS_PORT"),
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
	RedisClient.Close()
}

func Redis_set(key string, value string) {
	RedisClient.Set(context.Background(), key, value, 0)
}

func Redis_get(key string) string {
	val, _ := RedisClient.Get(context.Background(), key).Result()
	return val
}
