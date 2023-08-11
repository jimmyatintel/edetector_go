package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

// new scan
func GiveScanInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveScanInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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

func GiveScan(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveScan: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	// send to elasticsearch
	// lines := strings.Split(p.GetMessage(), "\n")
	// for _, line := range lines {
	// 	if len(line) == 0 {
	// 		continue
	// 	}
	// 	line = strings.ReplaceAll(line, "|", "@|@")
	// 	values := strings.Split(line, "@|@")
	// 	int_date, err := strconv.Atoi(values[1])
	// 	if err != nil {
	// 		logger.Debug("Invalid date: ", zap.Any("message", values[1]))
	// 		values[1] = "0"
	// 		int_date = 0
	// 	}
	// 	// network := "true"
	// 	// if values[18] == "null" {
	// 	// 	network = "false"
	// 	// }
	// 	uuid := uuid.NewString()
	// 	m_tmp := memory.Memory{}
	// 	_, err = elasticquery.StringToStruct(uuid, p.GetRkey(), line, &m_tmp)
	// 	if err != nil {
	// 		logger.Error("Error converting to struct: ", zap.Any("error", err.Error()))
	// 	}
	// 	m_tmp.RiskLevel, err = risklevel.Getriskscore(m_tmp)
	// 	if err != nil {
	// 		logger.Error("Error converting to struct: ", zap.Any("error", err.Error()))
	// 	}
	// 	line = strings.ReplaceAll(line, "-12345", strconv.Itoa(m_tmp.RiskLevel))
	// 	err = elasticquery.SendToMainElastic(uuid, "ed_de_memory", p.GetRkey(), values[0], int_date, "memory", strconv.Itoa(m_tmp.RiskLevel), "ed_mid")
	// 	if err != nil {
	// 		logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
	// 	}
	// 	err = elasticquery.SendToDetailsElastic(uuid, "ed_de_memory", p.GetRkey(), line, &m_tmp, "ed_mid")
	// 	if err != nil {
	// 		logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
	// 	}
	// }
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveScanDataEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	taskservice.Finish_task(p.GetRkey(), "StartScan")
	return task.SUCCESS, nil
}