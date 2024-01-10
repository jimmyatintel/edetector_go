package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"math"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

var ruleMatchWorkingPath = "ruleMatchWorking"
var ruleMatchUnstage = "ruleMatchUnstage"
var pathWorkingPath = "pathWorking"

func init() {
	file.CheckDir(ruleMatchWorkingPath)
	file.CheckDir(ruleMatchUnstage)
	file.CheckDir(pathWorkingPath)
}

func ReadyYaraRule(p packet.Packet, conn net.Conn, dataRight chan net.Conn) (task.TaskResult, error) {
	srcPath := filepath.Join("static", "yaraRule")
	dstPath := filepath.Join("static", "yaraRule.zip")
	logger.Info("ReadyYaraRule: " + p.GetRkey() + "::" + p.GetMessage())
	// zip the file
	err := file.ZipDirectory(srcPath, dstPath)
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
		query.Failed_task(p.GetRkey(), "StartYaraRule", 6)
		return
	}
	start := 0
	for {
		conn := <-dataRight
		if start >= fileLen {
			logger.Info("ServerSend GiveYaraRuleEnd: " + p.GetRkey())
			err = clientsearchsend.SendDataTCPtoClient(p, task.GIVE_YARA_RULE_END, []byte{}, conn)
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
	srcPath := filepath.Join(pathWorkingPath, key)
	dstPath := filepath.Join(pathWorkingPath, key+".zip")
	// q := fmt.Sprintf(`{
	// 	"query": {
	// 		"bool": {
	// 			"must": [
	// 				{ "term": { "agent": "%s" } }
	// 			]
	// 		}
	// 	}
	// }`, key)
	// hitsArray := elastic.SearchRequest(config.Viper.GetString("ELASTIC_PREFIX")+"_explorer", q)
	// err := file.CreateFile(srcPath)
	// if err != nil {
	// 	logger.Error("Create file error: " + err.Error())
	// 	return err
	// }
	// for _, hit := range hitsArray {
	// 	hitMap, ok := hit.(map[string]interface{})
	// 	if !ok {
	// 		logger.Error("Assert hitMap error")
	// 		return err
	// 	}
	// 	path := hitMap["_source"].(map[string]interface{})["path"].(string)
	// 	path = strings.ReplaceAll(path, "\\\\", "\\")
	// 	path = strings.ReplaceAll(path, "////", "//")
	// 	pathByte := []byte(path + "\n")
	// 	err := file.WriteFile(srcPath, pathByte)
	// 	if err != nil {
	// 		logger.Error("Write file error: " + err.Error())
	// 		return err
	// 	}
	// }

	file.ZipFile(srcPath, dstPath)
	fileInfo, err := os.Stat(dstPath)
	if err != nil {
		logger.Error("Get file info error: " + err.Error())
		return err
	}
	fileLen := int(fileInfo.Size())
	logger.Info("ServerSend GivePathInfo: " + key + "::" + strconv.Itoa(fileLen))
	err = clientsearchsend.SendTCPtoClient(p, task.GIVE_PATH_INFO, strconv.Itoa(fileLen), conn)
	if err != nil {
		logger.Error("Send GivePathInfo error: " + err.Error())
		return err
	}
	err = GivePath(p, fileLen, dstPath, dataRight)
	if err != nil {
		return err
	}
	return nil
}

func GivePath(p packet.Packet, fileLen int, path string, dataRight chan net.Conn) error {
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Read file error: " + err.Error())
		return err
	}
	start := 0
	for {
		conn := <-dataRight
		if start >= fileLen {
			logger.Info("ServerSend GivePathEnd: " + p.GetRkey())
			err = clientsearchsend.SendDataTCPtoClient(p, task.GIVE_PATH_END, []byte{}, conn)
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
