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
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var scanFirstPart float64
var scanSecondPart float64

// new scan
func GiveScanInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveScanInfo: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	scanFirstPart = float64(config.Viper.GetInt("SCAN_FIRST_PART"))
	scanSecondPart = 100 - scanFirstPart
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
	logger.Debug("GiveScanProgress: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
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
	ip, name := query.GetMachineIPandName(key)
	logger.Debug("GiveScan: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	redis.RedisSet_AddString(key+"-ScanMsg", p.GetMessage())
	// send to elasticsearch
	lines := strings.Split(redis.RedisGetString(key+"-ScanMsg"), "\n")
	redis.RedisSet(key+"-ScanMsg", "")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		values := strings.Split(line, "|")
		if values[16] == "null" {
			values[16] = "false"
		} else {
			scanNetworkElastic(values[9], values[1], key, values[16], ip, name)
			values[16] = "true"
		}
		values = append(values, "risklevel", "scan")
		uuid := uuid.NewString()
		m_tmp := memory.Memory{}
		_, err := elasticquery.StringToStruct(&m_tmp, values, uuid, key, "ip", "name", "item", "date", "ttype", "etc")
		if err != nil {
			logger.Error("Error converting to struct: ", zap.Any("error", err.Error()))
		}
		values[17], err = risklevel.Getriskscore(m_tmp)
		if err != nil {
			logger.Error("Error getting risk level: ", zap.Any("error", err.Error()))
		}
		err = elasticquery.SendToMainElastic(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", uuid, key, ip, name, values[0], values[1], "memory", values[17], "ed_mid")
		if err != nil {
			logger.Error("Error sending to main elastic: ", zap.Any("error", err.Error()))
		}
		err = elasticquery.SendToDetailsElastic(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", &m_tmp, values, uuid, key, ip, name, values[0], values[1], "memory", values[17], "ed_mid")
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

func scanNetworkElastic(id string, time string, key string, data string, ip string, name string) {
	lines := strings.Split(data, ";")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		re := regexp.MustCompile(`([^,]+),([^,]+),([^,]+),([^,]+),([^>]+)`)
		line = re.ReplaceAllString(line, "$1:$2|$3:$4|$5|$6")
		line = strings.ReplaceAll(line, ">", "")
		line = id + "|" + time + "|" + line
		logger.Info("scan network", zap.Any("message", line))
		values := strings.Split(line, "|")
		uuid := uuid.NewString()
		err := elasticquery.SendToDetailsElastic(config.Viper.GetString("ELASTIC_PREFIX")+"_memory_network_scan", &memory.MemoryNetworkScan{}, values, uuid, key, ip, name, "0", "0", "0", "0", "ed_mid")
		if err != nil {
			logger.Error("Error sending to details elastic: ", zap.Any("error", err.Error()))
		}
	}
}
