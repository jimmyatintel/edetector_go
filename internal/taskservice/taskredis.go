package taskservice

import "edetector_go/pkg/redis"

func loadfromredis(key string) string {
	return redis.RedisGetString(key)
}
