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

var collectFirstPart float64
var collectSecondPart float64

func GiveCollectProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveCollectProgress: " + key + "::" + p.GetMessage())

	// update progress
	if strings.Split(p.GetMessage(), "/")[0] == "1" {
		collectFirstPart = float64(config.Viper.GetInt("COLLECT_FIRST_PART"))
		collectSecondPart = float64(config.Viper.GetInt("COLLECT_SECOND_PART")) - collectFirstPart
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

func CollectReady(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("CollectReady: " + p.GetRkey() + "::" + p.GetMessage())

	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveCollectDataInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveCollectDataInfo: " + key + "::" + p.GetMessage())

	// init collect info
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-CollectTotal", total)
	redis.RedisSet(key+"-CollectCount", 0)

	// create or truncate the zip file
	path := filepath.Join(dbWorkingPath, p.GetRkey())
	if err := file.CreateFile(path); err != nil {
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
	logger.Debug("GiveCollectData: " + key)

	// write file
	path := filepath.Join(dbWorkingPath, key)
	content := getDataPacketContent(p)
	if err := file.WriteFile(path, content); err != nil {
		return task.FAIL, err
	}

	// update progress
	redis.RedisSet_AddInteger((key + "-CollectCount"), 1)
	collectCount, err := redis.RedisGetInt(key + "-CollectCount")
	if err != nil {
		return task.FAIL, err
	}
	collectTotal, err := redis.RedisGetInt(key + "-CollectTotal")
	if err != nil {
		return task.FAIL, err
	}

	progress := int(collectFirstPart) + getProgressByCount(collectCount, collectTotal, 65436, collectSecondPart)
	err = redis.RedisSet(key+"-CollectProgress", progress)
	if err != nil {
		return task.FAIL, err
	}

	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}

	return task.SUCCESS, nil
}

func GiveCollectDataEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveCollectDataEnd: " + key + "::" + p.GetMessage())

	progress := int(collectFirstPart) + int(collectSecondPart)
	srcPath := filepath.Join(dbWorkingPath, key)
	workPath := filepath.Join(dbWorkingPath, key+".db")
	unstagePath := filepath.Join(dbUstagePath, (key + ".db"))
	err := redis.RedisSet(key+"-CollectProgress", progress)
	if err != nil {
		return task.FAIL, err
	}
	collectTotal, err := redis.RedisGetInt(key + "-CollectTotal")
	if err != nil {
		return task.FAIL, err
	}

	// unzip data
	err = file.DecompressFile(srcPath, workPath, collectTotal)
	if err != nil {
		return task.FAIL, err
	}

	// move to Unstage
	if err := file.MoveFile(workPath, unstagePath); err != nil {
		return task.FAIL, err
	}

	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}

	// delete redis key: CollectTotal, CollectCount, CollectProgress
	err = redis.RedisDelete(key+"-CollectTotal", key+"-CollectCount", key+"-CollectProgress")
	if err != nil {
		logger.Error("Delete redis key failed: " + err.Error())
	}

	return task.SUCCESS, nil
}

func GiveCollectDataError(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataError: " + p.GetRkey() + "::" + p.GetMessage())

	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}

	return task.FAIL, errors.New(p.GetMessage())
}

func updateCollectProgress(key string) {
	for {
		result, err := query.Load_stored_task("nil", key, 2, "StartCollect")
		if err != nil {
			logger.Error("Get handling tasks failed: " + err.Error())
			return
		}

		if len(result) == 0 {
			return
		}

		collectProgress, err := redis.RedisGetInt(key + "-CollectProgress")
		if err != nil {
			logger.Error("Get collect progress failed: " + err.Error())
			return
		}

		err = query.Update_progress(collectProgress, key, "StartCollect")
		if err != nil {
			logger.Error("Update progress failed: " + err.Error())
			return
		}
		time.Sleep(time.Duration(config.Viper.GetInt("UPDATE_INTERVAL")) * time.Second)
	}
}
