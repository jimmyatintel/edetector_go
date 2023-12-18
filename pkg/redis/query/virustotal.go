package query

import (
	"edetector_go/pkg/redis"
)

func SetVTCache(ip string, result string) error {
	err := redis.RedisSet(ip, result)
	if err != nil {
		return err
	}
	return nil
}

func GetVTCache(ip string) (string, error) {
	result, err := redis.RedisGet(ip)
	if err != nil {
		if err.Error() == "redis: nil" {
			return "null", nil
		}
		return "", err
	}
	return result, nil
}
