package work

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/memory"
	"edetector_go/internal/packet"
	risklevel "edetector_go/internal/risklevel"
	"edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Process(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("Process: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_SCAN_INFO_DATA,
		Message:    "Ring0Process",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GetScanInfoData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GetScanInfoData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_PROCESS_INFO,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveProcessData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveProcessData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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

func GiveProcessDataEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveProcessDataEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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

func GiveScanProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveScanProgress: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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
	progress, err_ := strconv.Atoi(p.GetMessage())
	if err_ != nil {
		return task.FAIL, err
	}
	progress = int(math.Min(float64(progress*2), 99))
	query.Update_progress(progress, p.GetRkey(), "StartScan")
	go taskservice.RequestToUser(p.GetRkey())
	return task.SUCCESS, nil
}

func GiveScanDataInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveScanDataInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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

func GiveScanData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveScanData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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

func GiveScanDataOver(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveScanDataOver: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))

	// send to elasticsearch
	lines := strings.Split(p.GetMessage(), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		//! tmp version
		original := strings.Split(line, "|")
		int_date, err := strconv.Atoi(original[2])
		if err != nil {
			logger.Debug("Invalid date: ", zap.Any("message", original[2]))
			original[2] = "0"
			int_date = 0
		}
		network := "false"
		if original[18] != "null" {
			network = "true"
			go scanNetworkElastic(p.GetRkey(), original[18])
		}
		line = original[4] + "@|@" + original[2] + "@|@" + network + "@|@cmd@|@" + original[6] + " @|@" + original[5] + "@|@" + original[7] + "@|@parentName@|@" + original[9] + "@|@" + original[14] + "@|@" + original[0] + "@|@-12345@|@0,0@|@0@|@0,0@|@" + original[17] + "@|@" + original[16] + "@|@" + original[13] + "," + original[12] + "@|@scan"
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
		err = elasticquery.SendToMainElastic(uuid, config.Viper.GetString("ELASTIC_PREFIX")+"memory", p.GetRkey(), values[0], int_date, "memory", strconv.Itoa(m_tmp.RiskLevel), "ed_mid")
		if err != nil {
			logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
		}
		err = elasticquery.SendToDetailsElastic(uuid, config.Viper.GetString("ELASTIC_PREFIX")+"memory", p.GetRkey(), line, &m_tmp, "ed_mid")
		if err != nil {
			logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
		}
	}

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

func GiveScanDataEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
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

func scanNetworkElastic(key string, data string) {
	// networkSet := make(map[string]struct{})
	// lines := strings.Split(data, ";")
	// for _, line := range lines {
	// 	if len(line) == 0 {
	// 		continue
	// 	}
	// 	line = strings.ReplaceAll(line, "|", "@|@")
	// 	uuid := uuid.NewString()
	// 	values := strings.Split(line, "@|@")
	// 	key := values[0] + "," + values[3]
	// 	networkSet[key] = struct{}{}
	// 	err := elasticquery.SendToDetailsElastic(uuid, config.Viper.GetString("ELASTIC_PREFIX")+"memory_network", p.GetRkey(), line, &MemoryNetwork{}, "ed_high")
	// 	if err != nil {
	// 		logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
	// 	}
	// }
	// elasticquery.UpdateNetworkInfo(p.GetRkey(), networkSet)
}
