package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/redis"
	"net"

	"go.uber.org/zap"
)

func GiveDetectProcessFrag(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveDetectProcessFrag: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	redis.RedisSet_AddString(key+"-DetectMsg", p.GetMessage())
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

func GiveDetectProcess(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	// key := p.GetRkey()
	// logger.Info("GiveDetectProcess: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	// redis.RedisSet_AddString(key+"-DetectMsg", p.GetMessage())
	// // send to elasticsearch
	// lines := strings.Split(redis.RedisGetString(key+"-DetectMsg"), "\n")
	// redis.RedisSet(key+"-DetectMsg", "")
	// // send to elasticsearch
	// for _, line := range lines {
	// 	if len(line) == 0 {
	// 		continue
	// 	}
	// 	line = strings.ReplaceAll(line, "|", "@|@")
	// 	values := strings.Split(line, "@|@")
	// 	int_date, err := strconv.Atoi(values[1])
	// 	if err != nil {
	// 		logger.Error("Invalid date: ", zap.Any("message", values[1]))
	// 		int_date = 0
	// 	}
	// 	query := fmt.Sprintf(`{
	// 		"query": {
	// 			"bool": {
	// 				"must": [
	// 					{ "term": { "agent": "%s" } },
	// 					{ "term": { "processId": %s } },
	// 					{ "term": { "processCreateTime": %s } },
	// 					{ "term": { "processConnectIP": "true" } },
	// 					{ "term": { "mode": "detect" } }
	// 				]
	// 			}
	// 		}
	// 	}`, p.GetRkey(), values[9], values[1])
	// 	doc := elastic.SearchRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", query)
	// 	var network string
	// 	if doc == "" {
	// 		network = "detecting"
	// 	} else {
	// 		network = "true"
	// 		logger.Debug("Update information of the detect process: ", zap.Any("message", values[9]+" "+values[1]))
	// 	}
	// 	line = line + "@|@" + network + "@|@riskLevel@|@detect"
	// 	uuid := uuid.NewString()
	// 	m_tmp := memory.Memory{}
	// 	// _, err = elasticquery.StringToStruct(uuid, p.GetRkey(), line, &m_tmp, "0", 0, "0", "0")
	// 	if err != nil {
	// 		logger.Error("Error converting to struct: ", zap.Any("error", err.Error()))
	// 	}
	// 	m_tmp.RiskLevel, err = risklevel.Getriskscore(m_tmp)
	// 	if err != nil {
	// 		logger.Error("Error getting risk level: ", zap.Any("error", err.Error()))
	// 	}
	// 	line = strings.ReplaceAll(line, "riskLevel", strconv.Itoa(m_tmp.RiskLevel))
	// 	// err = elasticquery.SendToMainElastic(uuid, config.Viper.GetString("ELASTIC_PREFIX")+"_memory", p.GetRkey(), values[0], int_date, "memory", strconv.Itoa(m_tmp.RiskLevel), "ed_mid")
	// 	if err != nil {
	// 		logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
	// 	}
	// 	// err = elasticquery.SendToDetailsElastic(uuid, config.Viper.GetString("ELASTIC_PREFIX")+"_memory", p.GetRkey(), line, &m_tmp, "ed_mid", values[0], int_date, "memory", strconv.Itoa(m_tmp.RiskLevel))
	// 	if err != nil {
	// 		logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
	// 	}
	// }
	// var send_packet = packet.WorkPacket{
	// 	MacAddress: p.GetMacAddress(),
	// 	IpAddress:  p.GetipAddress(),
	// 	Work:       task.DATA_RIGHT,
	// 	Message:    "",
	// }
	// err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	// if err != nil {
	// 	return task.FAIL, err
	// }
	return task.SUCCESS, nil
}
