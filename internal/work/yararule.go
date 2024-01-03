package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"math"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var srcPath = filepath.Join("static", "yaraRule")
var dstPath = filepath.Join("static", "yaraRule.zip")
var ruleMatchWorkingPath = "ruleMatchWorking"
var ruleMatchUnstaged = "ruleMatchUnstaged"

func init() {
	file.CheckDir(ruleMatchWorkingPath)
	file.CheckDir(ruleMatchUnstaged)
}

func ReadyYaraRule(p packet.Packet, conn net.Conn, dataRight chan net.Conn) (task.TaskResult, error) {
	logger.Info("ReadyYaraRule: " + p.GetRkey() + "::" + p.GetMessage())
	// zip the file
	err := file.ZipFile(srcPath, dstPath)
	if err != nil {
		return task.FAIL, err
	}
	fileInfo, err := os.Stat(dstPath)
	if err != nil {
		return task.FAIL, err
	}
	fileLen := int(fileInfo.Size())
	logger.Info("ServerSend GiveYaraRuleInfo: " + p.GetRkey() + "::" + strconv.Itoa(fileLen))
	err = clientsearchsend.SendTCPtoClient(p, task.GIVE_YARA_RULE_INFO, strconv.Itoa(fileLen), conn)
	if err != nil {
		return task.FAIL, err
	}
	go GiveYaraRule(p, fileLen, dstPath, dataRight)
	return task.SUCCESS, nil
}

func GiveYaraRule(p packet.Packet, fileLen int, path string, dataRight chan net.Conn) {
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Read file error: " + err.Error())
	}
	start := 0
	for {
		conn := <-dataRight
		if start >= fileLen {
			logger.Info("ServerSend GiveYaraRuleEnd: " + p.GetRkey())
			err = clientsearchsend.SendDataTCPtoClient(p, task.GIVE_YARA_RULE_END, []byte{}, conn)
			if err != nil {
				logger.Error("Send GiveYaraRuleEnd error: " + err.Error())
			}
			query.Finish_task(p.GetRkey(), "StartYaraRule")
			break
		}
		end := int(math.Min(float64(fileLen), float64(start+65436)))
		data := content[start:end]
		logger.Info("ServerSend GiveYaraRule: " + p.GetRkey())
		err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_YARA_RULE, data, conn)
		if err != nil {
			logger.Error("Send GiveYaraRule error: " + err.Error())
		}
		start += 65436
	}
}

func GiveRuleMatchInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveRuleMatchInfo: " + p.GetRkey())
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-RuleMatchTotal", total)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveRuleMatch(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveRuleMatch: " + key)
	// write file
	path := filepath.Join(ruleMatchWorkingPath, key)
	content := getDataPacketContent(p)
	err := file.WriteFile(path, content)
	if err != nil {
		return task.FAIL, err
	}
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveRuleMatchEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveRuleMatchEnd: " + key)
	srcPath := filepath.Join(ruleMatchWorkingPath, key)
	workPath := filepath.Join(ruleMatchWorkingPath, key+".txt")
	dstPath := filepath.Join(ruleMatchUnstaged, key+".txt")
	err := file.DecompressionFile(srcPath, workPath, redis.RedisGetInt(key+"-RuleMatchTotal"))
	if err != nil {
		return task.FAIL, err
	}
	err = file.MoveFile(workPath, dstPath)
	if err != nil {
		return task.FAIL, err
	}
	err = parseRuleMatch(dstPath, p.GetRkey())
	if err != nil {
		return task.FAIL, err
	}
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func parseRuleMatch(path string, key string) error {
	logger.Info("ParseRuleMatch: " + key)
	// send to elasticsearch
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		values := strings.Split(line, "|")
		if len(values) == 2 {
			rules := strings.Split(values[0], ";")
			hits := len(rules)
			updateRuleMatch(key, values[0], values[1], hits)
		} else if line != "" {
			logger.Error("Error format: " + line)
		}
	}
	return nil
}

func updateRuleMatch(key string, rule string, path string, hits int) {
	logger.Info("UpdateRuleMatch: " + key + "|" + rule + "|" + path + "|" + strconv.Itoa(hits))
}
