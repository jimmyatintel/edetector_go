package work

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"edetector_go/pkg/redis"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
)

func GiveDetectProcessFrag(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveDetectProcessFrag: " + key + "||" + p.GetMessage())
	redis.RedisSet_AddString(key+"-DetectMsg", p.GetMessage())
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveDetectProcess(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	ip, name := query.GetMachineIPandName(key)
	logger.Info("GiveDetectProcess: " + key + "||" + p.GetMessage())
	redis.RedisSet_AddString(key+"-DetectMsg", p.GetMessage())
	lines := strings.Split(redis.RedisGetString(key+"-DetectMsg"), "\n")
	redis.RedisSet(key+"-DetectMsg", "")
	// send to elasticsearch
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		values := strings.Split(line, "|")
		values = append(values, "network", "risklevel", "detect")
		query := fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{ "term": { "agent": "%s" } },
						{ "term": { "processId": %s } },
						{ "term": { "processCreateTime": %s } },
						{ "term": { "processConnectIP": "true" } },
						{ "term": { "mode": "detect" } }
					]
				}
			}
		}`, p.GetRkey(), values[9], values[1])
		doc := elastic.SearchRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", query)
		if doc == "" {
			values[16] = "detecting"
		} else {
			values[16] = "true"
			logger.Debug("Update information of the detect process: " + values[9] + " " + values[1])
		}
		uuid := uuid.NewString()
		m_tmp := Memory{}
		_, err := rabbitmq.StringToStruct(&m_tmp, values, uuid, key, "ip", "name", "item", "date", "ttype", "etc")
		if err != nil {
			logger.Error("Error converting to struct: " + err.Error())
		}
		values[17], err = Getriskscore(m_tmp)
		if err != nil {
			logger.Error("Error getting risk level: " + err.Error())
		}
		err = rabbitmq.ToRabbitMQ_Main(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", uuid, key, ip, name, values[0], values[1], "memory", values[17], "ed_mid")
		if err != nil {
			logger.Error("Error sending to rabbitMQ (main): " + err.Error())
		}
		err = rabbitmq.ToRabbitMQ_Details(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", &m_tmp, values, uuid, key, ip, name, values[0], values[1], "memory", values[17], "ed_mid")
		if err != nil {
			logger.Error("Error sending to rabbitMQ (details): " + err.Error())
		}
	}
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}
