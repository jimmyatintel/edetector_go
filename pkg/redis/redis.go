package redis

import (
	"context"
	"edetector_go/config"
	"edetector_go/pkg/logger"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

	oldValue, err := RedisGetString(key)
	if err != nil {
		return err
	}

	newValue := oldValue + value
	return RedisClient.Set(context.Background(), key, newValue, 0).Err()
}

func RedisSet_AddInteger(key string, value int) error {
	if !checkflag() {
		return nil
	}

	oldValue, err := RedisGetInt(key)
	if err != nil {
		return err
	}

	newValue := oldValue + value

	return RedisClient.Set(context.Background(), key, newValue, 0).Err()
}

func RedisGet(key string) (string, error) {
	if !checkflag() {
		return "", nil
	}
	return RedisClient.Get(context.Background(), key).Result()
}

func RedisGetString(key string) (string, error) {
	if !checkflag() {
		return "", errors.New("Redis is disabled")
	}
	val, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func RedisGetInt(key string) (int, error) {
	if !checkflag() {
		return 0, errors.New("Redis is disabled")
	}

	val, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	val_int, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return val_int, nil
}

func RedisDelete(keys ...string) error {
	if !checkflag() {
		return nil
	}
	return RedisClient.Del(context.Background(), keys...).Err()
}

func GetKeysMatchingPattern(pattern string) []string {
	if !checkflag() {
		return nil
	}
	keys, err := RedisClient.Keys(context.Background(), pattern).Result()
	if err != nil {
		logger.Error("Error getting keys from redis: " + err.Error())
		return []string{}
	}
	return keys
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
			logger.Error("Error getting value from redis: " + err.Error())
		} else {
			values[key] = value
		}
	}

	return values
}

func UpdateDumpProgress(taskId string, progress int) {
	oldInfo, err := RedisGetString(taskId + "-LoadDumpTask")
	if err != nil {
		logger.Error("Error getting progress from redis: " + err.Error())
		return
	}

	newInfo := strings.Split(oldInfo, "|")
	if len(newInfo) == 5 {
		newInfo[4] = strconv.Itoa(int(progress))
	} else {
		logger.Error("The format of dump task info in redis is inccorect: " + oldInfo)
	}

	// set newInfo to redis
	err = RedisSet(taskId+"-LoadDumpTask", strings.Join(newInfo, "|"))
	if err != nil {
		logger.Error("Error setting progress to redis: " + err.Error())
		return
	}
}
