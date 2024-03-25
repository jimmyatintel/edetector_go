package work

import (
	"edetector_go/config"
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
	"time"
)

var ruleMatchWorkingPath = "ruleMatchWorking"
var ruleMatchUnstage = "ruleMatchUnstage"
var pathStaged = "pathStaged"
var yaraRulePath = filepath.Join("static", "yaraRule")

func init() {
	file.ClearDirContent(ruleMatchWorkingPath)
	file.CheckDir(ruleMatchUnstage)
	file.CheckDir(yaraRulePath)
}

func ReadyYaraRule(p packet.Packet, conn net.Conn, dataRight chan net.Conn) (task.TaskResult, error) {
	logger.Info("ReadyYaraRule: " + p.GetRkey() + "::" + p.GetMessage())
	path := filepath.Join(yaraRulePath, "yara.zip")
	content, err := os.ReadFile(path)
	if err != nil {
		return task.FAIL, err
	}
	fileLen := len(content)
	logger.Info("ServerSend GiveYaraRuleInfo: " + p.GetRkey() + "::" + strconv.Itoa(fileLen))
	err = clientsearchsend.SendTCPtoClient(p, task.GIVE_YARA_RULE_INFO, strconv.Itoa(fileLen), conn)
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(p.GetRkey()+"-YaraProgress", 0)
	go updateYaraRuleProgress(p.GetRkey())
	go GiveYaraRule(p, fileLen, content, dataRight)
	return task.SUCCESS, nil
}

func GiveYaraRule(p packet.Packet, fileLen int, content []byte, dataRight chan net.Conn) {
	start := 0
	for {
		conn := <-dataRight
		if start >= fileLen {
			logger.Info("ServerSend GiveYaraRuleEnd: " + p.GetRkey())
			err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_YARA_RULE_END, []byte{}, conn)
			if err != nil {
				logger.Error("Send GiveYaraRuleEnd error: " + err.Error())
				query.Failed_task(p.GetRkey(), "StartYaraRule", 6)
				return
			}
			<-dataRight
			err = GivePathInfo(p, p.GetRkey(), dataRight, conn)
			if err != nil {
				query.Failed_task(p.GetRkey(), "StartYaraRule", 6)
				return
			}
			break
		}
		end := int(math.Min(float64(fileLen), float64(start+65436)))
		data := content[start:end]
		logger.Info("ServerSend GiveYaraRule: " + p.GetRkey())
		err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_YARA_RULE, data, conn)
		if err != nil {
			logger.Error("Send GiveYaraRule error: " + err.Error())
			query.Failed_task(p.GetRkey(), "StartYaraRule", 6)
			return
		}
		start += 65436
	}
}

func GivePathInfo(p packet.Packet, key string, dataRight chan net.Conn, conn net.Conn) error {
	path := filepath.Join(pathStaged, key+".zip")
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Read file error: " + err.Error())
		return err
	}
	fileLen := len(content)
	logger.Info("ServerSend GivePathInfo: " + key + "::" + strconv.Itoa(fileLen))
	err = clientsearchsend.SendTCPtoClient(p, task.GIVE_PATH_INFO, strconv.Itoa(fileLen), conn)
	if err != nil {
		logger.Error("Send GivePathInfo error: " + err.Error())
		return err
	}
	err = GivePath(p, fileLen, content, dataRight)
	if err != nil {
		return err
	}
	return nil
}

func GivePath(p packet.Packet, fileLen int, content []byte, dataRight chan net.Conn) error {
	start := 0
	for {
		conn := <-dataRight
		if start >= fileLen {
			logger.Info("ServerSend GivePathEnd: " + p.GetRkey())
			err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_PATH_END, []byte{}, conn)
			if err != nil {
				logger.Error("Send GivePathEnd error: " + err.Error())
				return err
			}
			<-dataRight
			break
		}
		end := int(math.Min(float64(fileLen), float64(start+65436)))
		data := content[start:end]
		logger.Info("ServerSend GivePath: " + p.GetRkey())
		err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_PATH, data, conn)
		if err != nil {
			logger.Error("Send GivePath error: " + err.Error())
			return err
		}
		start += 65436
	}
	return nil
}

func GiveYaraProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveYaraProgress: " + key + "::" + p.GetMessage())
	progress, err := getProgressByMsg(p.GetMessage(), 95)
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-YaraProgress", progress)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func updateYaraRuleProgress(key string) {
	for {
		result, err := query.Load_stored_task("nil", key, 2, "StartYaraRule")
		if err != nil {
			logger.Error("Get handling tasks failed: " + err.Error())
			return
		}
		if len(result) == 0 {
			return
		}
		query.Update_progress(redis.RedisGetInt(key+"-YaraProgress"), key, "StartYaraRule")
		time.Sleep(time.Duration(config.Viper.GetInt("UPDATE_INTERVAL")) * time.Second)
	}
}
