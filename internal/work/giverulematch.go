package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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
	err := file.TruncateFile(srcPath, redis.RedisGetInt(key+"-RuleMatchTotal"))
	if err != nil {
		return task.FAIL, err
	}
	// workPath := filepath.Join(ruleMatchWorkingPath, key+".txt")
	dstPath := filepath.Join(ruleMatchUnstage, key+".txt")
	// err := file.DecompressionFile(srcPath, workPath, redis.RedisGetInt(key+"-RuleMatchTotal"))
	// if err != nil {
	// 	return task.FAIL, err
	// }
	// err = file.MoveFile(workPath, dstPath)
	err = file.MoveFile(srcPath, dstPath)
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
	query.Finish_task(p.GetRkey(), "StartYaraRule")
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
		}
	}
	return nil
}

func updateRuleMatch(key string, rule string, path string, hits int) {
	logger.Info("UpdateRuleMatch: " + key + "|" + rule + "|" + path + "|" + strconv.Itoa(hits))
}