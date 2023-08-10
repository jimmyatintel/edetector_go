package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/memory"
	"edetector_go/internal/packet"
	risklevel "edetector_go/internal/risklevel"
	"edetector_go/internal/task"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// type ProcessDetectJson struct {
// 	0 PID         int    `json:"pid"`
// 	1 Parent_PID  int    `json:"parent_pid"`
// 	2 ProcessName string `json:"process_name"`
// 	3 ProcessTime int    `json:"process_time"`
// 	4 ParentName  string `json:"parent_name"`
// 	5 ParentTime  int    `json:"parent_time"`
// }

// func (n ProcessDetectJson) Elastical() ([]byte, error) {
// 	return json.Marshal(n)
// }

func GiveProcessHistory(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveProcessHistory: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_PROCESS_HISTORY_INFO,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveProcessHistoryData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveProcessHistoryData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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

func GiveProcessHistoryEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveProcessHistoryEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))

	// send to elasticsearch
	lines := strings.Split(p.GetMessage(), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		//! tmp version
		original := strings.Split(line, "|")
		int_date, err := strconv.Atoi(original[3])
		if err != nil {
			logger.Debug("Invalid date: ", zap.Any("message", original[3]))
			original[3] = "0"
			int_date = 0
		}
		line = original[2] + "@|@" + original[3] + "@|@detecting@|@cmd@|@md5@|@path@|@" + original[1] + "@|@" + original[4] + "@|@parentPath@|@sign@|@" + original[0] + "@|@-12345@|@0,0@|@0@|@0,0@|@null@|@0@|@0,0@|@detect"
		values := strings.Split(line, "@|@")
		//! tmp version
		uuid := uuid.NewString()
		m_tmp := memory.Memory{}
		_, err = elasticquery.StringToStruct(uuid, p.GetRkey(), line, &m_tmp)
		if err != nil {
			logger.Error("Error converting to struct: ", zap.Any("error", err.Error()))
		}
		m_tmp.RiskLevel, err = risklevel.Getriskscore(m_tmp)
		if err != nil {
			logger.Error("Error converting to struct: ", zap.Any("error", err.Error()))
		}
		line = strings.ReplaceAll(line, "-12345", strconv.Itoa(m_tmp.RiskLevel))
		err = elasticquery.SendToMainElastic(uuid, "ed_memory", p.GetRkey(), values[0], int_date, "memory", strconv.Itoa(m_tmp.RiskLevel), "ed_high")
		if err != nil {
			logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
		}
		err = elasticquery.SendToDetailsElastic(uuid, "ed_memory", p.GetRkey(), line, &m_tmp, "ed_high")
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
