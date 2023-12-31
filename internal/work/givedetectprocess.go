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
	logger.Info("GiveDetectProcessFrag: " + key + "::" + p.GetMessage())
	redis.RedisSet_AddString(key+"-DetectMsg", p.GetMessage())
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveDetectProcess(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	ip, name, err := query.GetMachineIPandName(key)
	if err != nil {
		return task.FAIL, err
	}
	logger.Info("GiveDetectProcess: " + key + "::" + p.GetMessage())
	redis.RedisSet_AddString(key+"-DetectMsg", p.GetMessage())
	lines := strings.Split(redis.RedisGetString(key+"-DetectMsg"), "\n")
	redis.RedisSet(key+"-DetectMsg", "")
	// send to elasticsearch
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	HandleRelation(lines, key, 16)
	// }()
	for _, line := range lines {
		values := strings.Split(line, "|")
		if len(values) != 16 {
			if len(values) != 1 {
				logger.Warn("Invalid line: " + line)
			}
			continue
		}
		processKey := key + "##" + values[9] + "##" + values[1]
		values = append(values, "network", "risklevel", "riskscore", "detect", processKey)
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
		hitsArray := elastic.SearchRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", query)
		if len(hitsArray) == 0 {
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
		values[17], values[18], err = Getriskscore(m_tmp)
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
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	// wg.Wait()
	return task.SUCCESS, nil
}
