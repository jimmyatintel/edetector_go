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
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var imageFirstPart float64
var imageSecondPart float64

func GiveImageProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveImageProgress: " + key + "::" + p.GetMessage())
	// update progress
	if strings.Split(p.GetMessage(), "/")[0] == "1" { // the first
		imageFirstPart = float64(config.Viper.GetInt("IMAGE_FIRST_PART"))
		// imageSecondPart = 100 - imageFirstPart
		redis.RedisSet(key+"-ImageProgress", 0)
		err := file.ClearDirContent(filepath.Join(imageWorkingPath, key)) // clear or create imageWorking directory for the agent
		if err != nil {
			return task.FAIL, err
		}
		go updateImageProgress(key)
	}
	progress, err := getProgressByMsg(p.GetMessage(), imageFirstPart)
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ImageProgress", progress)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveImageInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveImageInfo: " + key + "::" + p.GetMessage())
	// init image info
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ImageTotal", total)
	// redis.RedisSet(key+"-ImageCount", 0)
	// create or truncate the zip file
	path := filepath.Join(imageWorkingPath, key+".zip")
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
	path := filepath.Join(imageWorkingPath, key+".zip")
	content := getDataPacketContent(p)
	err := file.WriteFile(path, content)
	if err != nil {
		return task.FAIL, err
	}
	// update progress
	// redis.RedisSet_AddInteger((key + "-ImageCount"), 1)
	// progress := int(imageFirstPart) + getProgressByCount(redis.RedisGetInt(key+"-ImageCount"), redis.RedisGetInt(key+"-ImageTotal"), 65436, imageSecondPart)
	// redis.RedisSet(key+"-ImageProgress", progress)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveImageEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveImageEnd: " + key + "::" + p.GetMessage())
	srcPath := filepath.Join(imageWorkingPath, key+".zip")
	dstPath := filepath.Join(imageWorkingPath, key)
	// truncate data
	err := file.TruncateFile(srcPath, redis.RedisGetInt(key+"-ImageTotal"))
	if err != nil {
		return task.FAIL, err
	}
	// decompress the file and move to the imageWorkingPath/key
	err = file.DecompressionDir(srcPath, dstPath)
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

func FinishImage(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("FinishImage: " + key)
	ip, _, err := query.GetMachineIPandName(key)
	if err != nil {
		return task.FAIL, err
	}
	time := time.Now().Format("2006_0102_150405")
	// clear all the content of the directory
	err = file.ClearDirContent(filepath.Join(imageFilePath, ip))
	if err != nil {
		return task.FAIL, err
	}
	// zip the file and move to ImagePath
	srcPath := filepath.Join(imageWorkingPath, key)
	dstPath := filepath.Join(imageFilePath, ip, (("Obtained_" + time + "_" + getTaskMsg(key, "StartGetImage")) + ".zip"))
	err = file.ZipDir(srcPath, dstPath)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func ImageError(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	return task.FAIL, errors.New("receive ImageError")
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
