package redis

import (
	"context"
	"edetector_go/config"
	"edetector_go/pkg/logger"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func checkflag() bool {
	// if enable, err := fflag.FFLAG.FeatureEnabled("redis_enable"); enable && err == nil {
	return true
	// }
	// return false
}

func Redis_init() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Viper.GetString("REDIS_HOST") + ":" + config.Viper.GetString("REDIS_PORT"),
		Password: config.Viper.GetString("REDIS_PASSWORD"),
		DB:       config.Viper.GetInt("REDIS_DB"),
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Panic("Error connecting to redis")
		panic(err)
	}
	logger.Info("Redis is enabled")
	return RedisClient
}

func RedisClose() {
	if !checkflag() {
		return
	}
	RedisClient.Close()
}

func RedisExists(key string) bool {
	exists, err := RedisClient.Exists(context.Background(), key).Result()
	if err != nil {
		logger.Error("Error checking key existence: " + err.Error())
		return false
	}
	if exists == 1 {
		return true
	} else {
		return false
	}
}

func RedisSet(key string, value interface{}) error {
	if !checkflag() {
		return nil
	}
	return RedisClient.Set(context.Background(), key, value, 0).Err()
}

func RedisSet_AddString(key string, value string) error {
	if !checkflag() {
		return nil
	}
	newValue := RedisGetString(key) + value
	return RedisClient.Set(context.Background(), key, newValue, 0).Err()
}

func RedisSet_AddInteger(key string, value int) error {
	if !checkflag() {
		return nil
	}
	newValue := RedisGetInt(key) + value
	return RedisClient.Set(context.Background(), key, newValue, 0).Err()
}

func RedisGet(key string) (string, error) {
	if !checkflag() {
		return "", nil
	}
	return RedisClient.Get(context.Background(), key).Result()
}

func RedisGetString(key string) string {
	if !checkflag() {
		return ""
	}
	val, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		logger.Error("Error getting value from redis: " + err.Error())
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
		logger.Error("Error getting value from redis: " + err.Error())
		return 0
	}
	val_int, err := strconv.Atoi(val)
	if err != nil {
		logger.Error("Error converting to integer: " + err.Error())
		return 0
	}
	return val_int
}

func RedisDelete(key string) error {
	if !checkflag() {
		return nil
	}
	return RedisClient.Del(context.Background(), key).Err()
}

func GetKeysByLength(length int) []string {
	keys, err := RedisClient.Keys(context.Background(), "*").Result()
	if err != nil {
		fmt.Println("Error getting keys from redis:", err)
		return nil
	}

	var keysWithLength []string
	for _, key := range keys {
		if len(key) == length {
			keysWithLength = append(keysWithLength, key)
		}
	}

	return keysWithLength
}

func GetValuesForKeys(keys []string) map[string]string {
	values := make(map[string]string)

	for _, key := range keys {
		value, err := RedisClient.Get(context.Background(), key).Result()
		if err != nil {
			fmt.Println("Error getting value from redis:", err)
		} else {
			values[key] = value
		}
	}

	return values
}
