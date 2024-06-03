package work

import (
	"edetector_go/config"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"fmt"
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
	ruleMatchTotal, err := redis.RedisGetInt(key + "-RuleMatchTotal")
	if err != nil {
		return task.FAIL, err
	}
	err = file.TruncateFile(srcPath, ruleMatchTotal)
	if err != nil {
		return task.FAIL, err
	}
	dstPath := filepath.Join(ruleMatchUnstage, key+".txt")
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
	matches := make(map[string]string)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		values := strings.Split(line, "|")
		if len(values) == 2 {
			rule := values[0]
			hitPath := strings.TrimRight(values[1], "\r")
			hitPath = strings.ReplaceAll(hitPath, "\\\\", "_backslash_")
			hitPath = strings.ReplaceAll(hitPath, "\\", "_backslash_")
			hitPath = strings.ReplaceAll(hitPath, "//", "_slash_")
			hitPath = strings.ReplaceAll(hitPath, "/", "_slash_")
			hitPath = strings.ReplaceAll(hitPath, ":", "_colon_")
			// check if matches exist
			if _, ok := matches[hitPath]; !ok {
				matches[hitPath] = rule
			} else {
				matches[hitPath] += ("," + rule)
			}
		}
	}
	for p, rule := range matches {
		count := len(strings.Split(rule, ","))
		UpdateRuleMatch(key, rule, p, count)
	}

	return nil
}

func UpdateRuleMatch(key string, rule string, path string, count int) {
	logger.Info("UpdateRuleMatch: " + key + "|" + rule + "|" + path + "|" + strconv.Itoa(count))
	query := fmt.Sprintf(`
		{
		"script": {
			"source": "ctx._source.explorer.yaraRuleHitCount = params.count ; ctx._source.explorer.yaraRuleHit = params.hit",
			"lang": "painless",
			"params": {
				"count": %d,
				"hit": "%s"
			}
		},
		"query": {
			"bool": {
				"must": [
					{ "term": { "agent": "%s" } },
					{
						"match_phrase": {
						  "explorer.path": "%s"
						}
					}
				]
			}
		}
	}`, count, rule, key, path)
	updateCount, err := elastic.UpdateByQueryRequest(query, config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", 0)
	if err != nil {
		logger.Error("UpdateRuleMatch error: " + err.Error())
		return
	}
	logger.Info("UpdateCount: " + strconv.Itoa(updateCount))
}
