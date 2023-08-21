package work

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/memory"
	"edetector_go/internal/packet"
	"edetector_go/internal/risklevel"
	"edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var scanMu *sync.Mutex
var scanTotalMap = make(map[string]int)
var scanCountMap = make(map[string]int)
var scanProgressMap = make(map[string]int)
var scanMap = make(map[string](string))

func init() {
	scanMu = &sync.Mutex{}
}

// new scan
func GiveScanInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveScanInfo: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	scanTotalMap[key] = total
	scanMu.Lock()
	scanProgressMap[key] = 0
	scanMu.Unlock()
	go updateScanProgress(key)
	scanMap[key] = ""
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err = clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveScanProgress: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	// update progress
	progress, err := getProgressByMsg(p.GetMessage(), 50)
	if err != nil {
		return task.FAIL, err
	}
	scanMu.Lock()
	scanProgressMap[key] = progress
	scanMu.Unlock()
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err = clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanFragment(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveScanFragment: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	scanMap[p.GetRkey()] += p.GetMessage()
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
	key := p.GetRkey()
	logger.Debug("GiveScan: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	scanMap[key] += p.GetMessage()
	// send to elasticsearch
	lines := strings.Split(scanMap[key], "\n")
	scanMap[key] = ""
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
		network := "true"
		if values[16] == "null" {
			network = "false"
		}
		lastElement := strings.LastIndex(line, "@|@")
		line = line[:lastElement] + "@|@" + network + "@|@riskLevel@|@scan"
		uuid := uuid.NewString()
		m_tmp := memory.Memory{}
		_, err = elasticquery.StringToStruct(uuid, p.GetRkey(), line, &m_tmp, "0", 0, "0", "0")
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
		err = elasticquery.SendToDetailsElastic(uuid, config.Viper.GetString("ELASTIC_PREFIX")+"_memory", p.GetRkey(), line, &m_tmp, "ed_mid", values[0], int_date, "memory", strconv.Itoa(m_tmp.RiskLevel))
		if err != nil {
			logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
		}
	}
	// update progress
	scanCountMap[key] += 1
	scanMu.Lock()
	scanProgressMap[key] = int(50 + float64(scanCountMap[key])/(float64(scanTotalMap[key]))*50)
	scanMu.Unlock()
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

func GiveScanEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveScanEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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
	taskservice.Finish_task(p.GetRkey(), "StartScan")
	return task.SUCCESS, nil
}

func updateScanProgress(key string) {
	for {
		scanMu.Lock()
		if scanProgressMap[key] >= 100 {
			break
		}
		rowsAffected := query.Update_progress(scanProgressMap[key], key, "StartScan")
		scanMu.Unlock()
		if rowsAffected != 0 {
			logger.Info("update progress", zap.Any("message", scanProgressMap[key]))
			go taskservice.RequestToUser(key)
		}
	}
}
