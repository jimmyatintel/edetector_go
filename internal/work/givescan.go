package work

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"edetector_go/pkg/redis"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var scanWorkingPath = "scanWorking"
var scanUstagePath = "scanUnstage"
var scanFirstPart float64
var scanSecondPart float64
var Re *regexp.Regexp

func init() {
	file.CheckDir(scanWorkingPath)
	file.CheckDir(scanUstagePath)
}

// new scan
func ReadyScan(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("ReadyScan: " + p.GetRkey() + "::" + p.GetMessage())
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveScanInfo: " + key + "::" + p.GetMessage())
	scanFirstPart = float64(config.Viper.GetInt("SCAN_FIRST_PART"))
	scanSecondPart = 100 - scanFirstPart
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ScanTotal", total)
	redis.RedisSet(key+"-ScanCount", 0)
	redis.RedisSet(key+"-ScanProgress", 0)
	go updateScanProgress(key)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveScanProgress: " + key + "::" + p.GetMessage())
	// update progress
	progress, err := getProgressByMsg(p.GetMessage(), scanFirstPart)
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ScanProgress", progress)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanDataInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveScanDataInfo: " + key + "::" + p.GetMessage())
	// init scan info
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ScanTotal", total)
	redis.RedisSet(key+"-ScanCount", 0)
	// create or truncate the zip file
	path := filepath.Join(scanWorkingPath, (p.GetRkey() + ".zip"))
	err = file.CreateFile(path)
	if err != nil {
		return task.FAIL, err
	}
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScan(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveScan: " + key)
	// write file
	path := filepath.Join(scanWorkingPath, (key + ".zip"))
	err := file.WriteFile(path, p)
	if err != nil {
		return task.FAIL, err
	}
	// update progress
	redis.RedisSet_AddInteger((key + "-ScanCount"), 1)
	progress := int(scanFirstPart) + getProgressByCount(redis.RedisGetInt(key+"-ScanCount"), redis.RedisGetInt(key+"-ScanTotal"), 65436, scanSecondPart)
	redis.RedisSet(key+"-ScanProgress", progress)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveScanEnd: " + key + "::" + p.GetMessage())

	srcPath := filepath.Join(scanWorkingPath, (key + ".zip"))
	workPath := filepath.Join(scanWorkingPath, key+".txt")
	unstagePath := filepath.Join(scanUstagePath, (key + ".txt"))
	// unzip data
	err := file.UnzipFile(srcPath, workPath, redis.RedisGetInt(key+"-ScanTotal"))
	if err != nil {
		return task.FAIL, err
	}
	// move to Unstage
	err = file.MoveFile(workPath, unstagePath)
	if err != nil {
		return task.FAIL, err
	}
	err = parseScan(unstagePath, p.GetRkey())
	if err != nil {
		return task.FAIL, err
	}
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	query.Finish_task(key, "StartScan")
	return task.SUCCESS, nil
}

func updateScanProgress(key string) {
	for {
		if redis.RedisGetInt(key+"-ScanProgress") >= 100 {
			break
		}
		query.Update_progress(redis.RedisGetInt(key+"-ScanProgress"), key, "StartScan")
		time.Sleep(time.Duration(config.Viper.GetInt("UPDATE_INTERVAL")) * time.Second)
	}
}

func parseScan(path string, key string) error {
	ip, name := query.GetMachineIPandName(key)
	logger.Info("ParseScan: " + key)
	// send to elasticsearch
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\r", "")
		values := strings.Split(line, "|")
		if len(values) != 17 {
			if len(values) != 1 {
				logger.Warn("Invalid line: " + line)
			}
			continue
		}
		if values[16] == "null" {
			values[16] = "false"
		} else {
			scanNetworkElastic(values[9], values[1], key, values[16], ip, name)
			values[16] = "true"
		}
		values = append(values, "risklevel", "riskscore", "scan")
		uuid := uuid.NewString()
		m_tmp := Memory{}
		_, err := rabbitmq.StringToStruct(&m_tmp, values, uuid, key, "ip", "name", "item", "date", "ttype", "etc")
		if err != nil {
			return err
		}
		values[17], values[18], err = Getriskscore(m_tmp)
		if err != nil {
			return err
		}
		err = rabbitmq.ToRabbitMQ_Main(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", uuid, key, ip, name, values[0], values[1], "memory", values[17], "ed_mid")
		if err != nil {
			return err
		}
		err = rabbitmq.ToRabbitMQ_Details(config.Viper.GetString("ELASTIC_PREFIX")+"_memory", &m_tmp, values, uuid, key, ip, name, values[0], values[1], "memory", values[17], "ed_mid")
		if err != nil {
			return err
		}
	}
	return nil
}

func scanNetworkElastic(pid string, pCreateTime string, key string, data string, ip string, name string) {
	lines := strings.Split(data, ";")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		conns := strings.Split(line, ",")
		if len(conns) != 5 {
			logger.Warn("Invalid line: " + line)
			continue
		}
		actionAndTime := strings.Split(conns[4], ">")
		if len(actionAndTime) != 2 {
			logger.Warn("Invalid line: " + line)
			continue
		}
		var direction string
		if ip == conns[0] { // i am src -> out
			direction = "1"
		} else { // i am dst -> in
			direction = "0"
		}
		line = pid + "|" + pCreateTime + "|" + actionAndTime[1] + "|" + conns[0] + "|" + conns[1] + "|" + conns[2] + "|" + conns[3] + "|" + actionAndTime[0] + "|" + direction
		values := strings.Split(line, "|")
		uuid := uuid.NewString()
		err := rabbitmq.ToRabbitMQ_Details(config.Viper.GetString("ELASTIC_PREFIX")+"_memory_network", &MemoryNetwork{}, values, uuid, key, ip, name, "0", "0", "0", "0", "ed_mid")
		if err != nil {
			logger.Error("Error sending to rabbitMQ (details): " + err.Error())
		}
	}
}
