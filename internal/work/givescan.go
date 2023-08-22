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
	"edetector_go/pkg/redis"
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var scanFirstPart float64 = 50
var scanSecondPart float64 = 100 - scanFirstPart

// new scan
func GiveScanInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveScanInfo: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ScanTotal", total)
	redis.RedisSet(key+"-ScanCount", 0)
	redis.RedisSet(key+"-ScanProgress", 0)
	redis.RedisSet(key+"-ScanMsg", "")
	go updateScanProgress(key)
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
	progress, err := getProgressByMsg(p.GetMessage(), scanFirstPart)
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ScanProgress", progress)
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
	key := p.GetRkey()
	logger.Debug("GiveScanFragment: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	redis.RedisSet_AddString(key+"-ScanMsg", p.GetMessage())
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
	redis.RedisSet_AddString(key+"-ScanMsg", p.GetMessage())
	// send to elasticsearch
	lines := strings.Split(redis.RedisGetString(key+"-ScanMsg"), "\n")
	redis.RedisSet(key+"-ScanMsg", "")
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
	redis.RedisSet_AddInteger((key + "-ScanCount"), 1)
	progress := int(scanFirstPart) + getProgressByCount(redis.RedisGetInt(key+"-ScanCount"), redis.RedisGetInt(key+"-ScanTotal"), 1, scanSecondPart)
	redis.RedisSet(key+"-ScanProgress", progress)
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
		if redis.RedisGetInt(key+"-ScanProgress") >= 100 {
			break
		}
		rowsAffected := query.Update_progress(redis.RedisGetInt(key+"-ScanProgress"), key, "StartScan")
		if rowsAffected != 0 {
			go taskservice.RequestToUser(key)
		}
	}
}
