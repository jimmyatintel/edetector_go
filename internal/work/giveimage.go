package work

import (
	"bytes"
	"edetector_go/config"
	C_AES "edetector_go/internal/C_AES"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"net"
	"path/filepath"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var imageWorkingPath = "imageWorking"
var imageUstagePath = "imageUnstage"

func init() {
	file.CheckDir(imageWorkingPath)
	file.CheckDir(imageUstagePath)
}

func GiveImageInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveImageInfo: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	go updateImageProgress(key)
	// init image info
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ImageTotal", total)
	redis.RedisSet(key+"-ImageCount", 0)
	// create or truncate the zip file
	path := filepath.Join(imageWorkingPath, (p.GetRkey() + ".zip"))
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
	logger.Debug("GiveImage: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	// write file
	dp := packet.CheckIsData(p)
	decrypt_buf := bytes.Repeat([]byte{0}, len(dp.Raw_data))
	C_AES.Decryptbuffer(dp.Raw_data, len(dp.Raw_data), decrypt_buf)
	decrypt_buf = decrypt_buf[100:]
	path := filepath.Join(imageWorkingPath, (key + ".zip"))
	err := file.WriteFile(path, decrypt_buf)
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
	logger.Info("GiveImageEnd: ", zap.Any("message", key+", Msg: "+p.GetMessage()))

	workPath := filepath.Join(imageWorkingPath, key+".zip")
	unstagePath := filepath.Join(imageUstagePath, (key + ".zip"))
	// truncate data
	err := file.TruncateFile(workPath, redis.RedisGetInt(key+"-ImageTotal"))
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

func updateImageProgress(key string) {
	for {
		if redis.RedisGetInt(key+"-ImageProgress") >= 100 {
			break
		}
		query.Update_progress(redis.RedisGetInt(key+"-ImageProgress"), key, "Image")
		time.Sleep(time.Duration(config.Viper.GetInt("UPDATE_INTERVAL")) * time.Second)
	}
}
