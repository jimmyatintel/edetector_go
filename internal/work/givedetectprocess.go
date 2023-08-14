package work

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/memory"
	packet "edetector_go/internal/packet"
	"edetector_go/internal/risklevel"
	task "edetector_go/internal/task"
	"edetector_go/pkg/elastic"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GiveDetectProcess(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectProcess: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	// send to elasticsearch
	lines := strings.Split(p.GetMessage(), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		line = strings.ReplaceAll(line, "|", "@|@")
		values := strings.Split(line, "@|@")
		int_date, err := strconv.Atoi(values[1])
		if err != nil {
			logger.Error("Invalid date: ", zap.Any("message", values[1]))
			int_date = 0
		}
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
		var network string
		if doc == "" {
			network = "detecting"
		} else {
			network = "true"
			logger.Debug("Update information of the detect process: ", zap.Any("message", values[9]+" "+values[1]))
		}
		lastElement := strings.LastIndex(line, "@|@")
		line = line[:lastElement] + network + "@|@riskLevel@|@detect"
		uuid := uuid.NewString()
		m_tmp := memory.Memory{}
		_, err = elasticquery.StringToStruct(uuid, p.GetRkey(), line, &m_tmp)
		if err != nil {
			logger.Error("Error converting to struct: ", zap.Any("error", err.Error()))
		}
		m_tmp.RiskLevel, err = risklevel.Getriskscore(m_tmp)
		if err != nil {
			logger.Error("Error getting risk level: ", zap.Any("error", err.Error()))
		}
		line = strings.ReplaceAll(line, "riskLevel", strconv.Itoa(m_tmp.RiskLevel))
		err = elasticquery.SendToMainElastic(uuid, config.Viper.GetString("ELASTIC_PREFIX")+"_memory", p.GetRkey(), values[0], int_date, "memory", strconv.Itoa(m_tmp.RiskLevel), "ed_mid")
		if err != nil {
			logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
		}
		err = elasticquery.SendToDetailsElastic(uuid, config.Viper.GetString("ELASTIC_PREFIX")+"_memory", p.GetRkey(), line, &m_tmp, "ed_mid")
		if err != nil {
			logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
		}
	}
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}
