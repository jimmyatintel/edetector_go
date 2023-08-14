package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

func GiveDetectProcess(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectProcess: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	// send to elasticsearch
	// lines := strings.Split(p.GetMessage(), "\n")
	// for _, line := range lines {
	// 	if len(line) == 0 {
	// 		continue
	// 	}
	// 	//! tmp version
	// 	original := strings.Split(line, "|")
	// 	int_date, err := strconv.Atoi(original[2])
	// 	if err != nil {
	// 		logger.Debug("Invalid date: ", zap.Any("message", original[2]))
	// 		original[2] = "0"
	// 		int_date = 0
	// 	}
	// 	line = original[4] + "@|@" + original[2] + "@|@detecting@|@cmd@|@" + original[6] + "@|@" + original[5] + "@|@" + original[7] + "@|@parentName@|@" + original[9] + "@|@" + original[14] + "@|@" + original[0] + "@|@-12345@|@0,0@|@0@|@0,0@|@null@|@0@|@" + original[13] + "," + original[12] + "@|@detect"
	// 	values := strings.Split(line, "@|@")
	// 	//! tmp version
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
	// 	// logger.Info("Risk level", zap.Any("message", strconv.Itoa(m_tmp.RiskLevel)))
	// 	err = elasticquery.SendToMainElastic(uuid, "ed_de_memory", p.GetRkey(), values[0], int_date, "memory", strconv.Itoa(m_tmp.RiskLevel), "ed_high")
	// 	if err != nil {
	// 		logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
	// 	}
	// 	err = elasticquery.SendToDetailsElastic(uuid, "ed_de_memory", p.GetRkey(), line, &m_tmp, "ed_high")
	// 	if err != nil {
	// 		logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
	// 	}
	// }
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
