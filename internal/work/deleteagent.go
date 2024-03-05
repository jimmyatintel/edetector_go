package work

import (
	"edetector_go/pkg/logger"
	mq "edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
)

func DeleteAgentData(key string) {
	mq.DeleteAgent(key)
	err := redis.RedisDelete(key)
	if err != nil {
		logger.Error("Error deleting key from redis: " + err.Error())
	}
}
