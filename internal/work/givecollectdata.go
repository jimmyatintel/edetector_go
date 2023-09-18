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
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"net"
)

var dbWorkingPath = "dbWorking"
var dbUstagePath = "dbUnstage"
var collectFirstPart float64
var collectSecondPart float64

func init() {
	file.CheckDir(dbWorkingPath)
	file.CheckDir(dbUstagePath)
}

func GiveCollectProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveCollectProgress: " + key + "|" + p.GetMessage())
	// update progress
	if strings.Split(p.GetMessage(), "/")[0] == "1" {
		collectFirstPart = float64(config.Viper.GetInt("COLLECT_FIRST_PART"))
		collectSecondPart = 100 - collectFirstPart
		go updateCollectProgress(key)
	}
	progress, err := getProgressByMsg(p.GetMessage(), collectFirstPart)
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-CollectProgress", progress)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectDataInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveCollectDataInfo: " + key + "|" + p.GetMessage())
	// init collect info
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-CollectTotal", total)
	redis.RedisSet(key+"-CollectCount", 0)
	// create or truncate the zip file
	path := filepath.Join(dbWorkingPath, (p.GetRkey() + ".zip"))
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

func GiveCollectData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveCollectData: " + key + "|" + p.GetMessage())
	// write file
	path := filepath.Join(dbWorkingPath, (key + ".zip"))
	err := file.WriteFile(path, p)
	if err != nil {
		return task.FAIL, err
	}
	// update progress
	redis.RedisSet_AddInteger((key + "-CollectCount"), 1)
	progress := int(collectFirstPart) + getProgressByCount(redis.RedisGetInt(key+"-CollectCount"), redis.RedisGetInt(key+"-CollectTotal"), 65436, collectSecondPart)
	redis.RedisSet(key+"-CollectProgress", progress)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectDataEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveCollectDataEnd: " + key + "|" + p.GetMessage())

	srcPath := filepath.Join(dbWorkingPath, (key + ".zip"))
	workPath := filepath.Join(dbWorkingPath, key+".db")
	unstagePath := filepath.Join(dbUstagePath, (key + ".db"))
	// unzip data
	err := file.UnzipFile(srcPath, workPath, redis.RedisGetInt(key+"-CollectTotal"))
	if err != nil {
		return task.FAIL, err
	}
	// move to Unstage
	err = file.MoveFile(workPath, unstagePath)
	if err != nil {
		return task.FAIL, err
	}
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectDataError(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataError: " + p.GetRkey() + "|" + p.GetMessage())
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.FAIL, errors.New(p.GetMessage())
}

func updateCollectProgress(key string) {
	for {
		if redis.RedisGetInt(key+"-CollectProgress") >= 100 {
			break
		}
		query.Update_progress(redis.RedisGetInt(key+"-CollectProgress"), key, "StartCollect")
		time.Sleep(time.Duration(config.Viper.GetInt("UPDATE_INTERVAL")) * time.Second)
	}
}
