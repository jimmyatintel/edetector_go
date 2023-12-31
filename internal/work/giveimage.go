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
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var imageWorkingPath = "imageWorking"
var imageUstagePath = "imageUnstage"

func init() {
	file.CheckDir(imageWorkingPath)
	file.CheckDir(imageUstagePath)
}

func ReadyImage(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("ReadyImage: " + key + "::" + p.GetMessage())
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveImageInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveImageInfo: " + key + "::" + p.GetMessage())
	redis.RedisSet(key+"-ImageProgress", 0)
	go updateImageProgress(key)
	// init image info
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ImageTotal", total)
	redis.RedisSet(key+"-ImageCount", 0)
	// create or truncate the zip file
	path := filepath.Join(imageWorkingPath, p.GetRkey())
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

func GiveImage(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveImage: " + key)
	// write file
	path := filepath.Join(imageWorkingPath, key)
	err := file.WriteFile(path, p)
	if err != nil {
		return task.FAIL, err
	}
	// update progress
	redis.RedisSet_AddInteger((key + "-ImageCount"), 1)
	progress := getProgressByCount(redis.RedisGetInt(key+"-ImageCount"), redis.RedisGetInt(key+"-ImageTotal"), 65436, 100)
	redis.RedisSet(key+"-ImageProgress", progress)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveImageEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveImageEnd: " + key + "::" + p.GetMessage())
	workPath := filepath.Join(imageWorkingPath, key)
	extension := ".tar.gz"
	f, err := os.Open(workPath)
	if err != nil {
		return task.FAIL, err
	}
	defer f.Close()
	var firstByte [1]byte
	_, err = f.Read(firstByte[:])
	if err != nil {
		return task.FAIL, err
	}
	if firstByte[0] == 'P' {
		extension = ".zip"
	}
	unstagePath := filepath.Join(imageUstagePath, (key + extension))
	// truncate data
	err = file.TruncateFile(workPath, redis.RedisGetInt(key+"-ImageTotal"))
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
	query.Finish_task(key, "StartGetImage")
	return task.SUCCESS, nil
}

func updateImageProgress(key string) {
	for {
		result, err := query.Load_stored_task("nil", key, 2, "StartGetImage")
		if err != nil {
			logger.Error("Get handling tasks failed: " + err.Error())
			return
		}
		if len(result) == 0 {
			return
		}
		query.Update_progress(redis.RedisGetInt(key+"-ImageProgress"), key, "StartGetImage")
		time.Sleep(time.Duration(config.Viper.GetInt("UPDATE_INTERVAL")) * time.Second)
	}
}
