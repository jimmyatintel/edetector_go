package work

import (
	"edetector_go/config"
	"edetector_go/internal/channelmap"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"errors"
	"strconv"
	"time"

	"path/filepath"

	"net"
	"strings"
)

var fileWorkingPath = "fileWorking"
var fileUnstagePath = "fileUnstage"
var explorerFirstPart float64
var explorerSecondPart float64

func init() {
	file.CheckDir(fileWorkingPath)
	file.CheckDir(fileUnstagePath)
}

func Explorer(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("Explorer: " + key + "|" + p.GetMessage())
	explorerFirstPart = float64(config.Viper.GetInt("EXPLORER_FIRST_PART"))
	explorerSecondPart = 100 - explorerFirstPart
	parts := strings.Split(p.GetMessage(), "|")
	redis.RedisSet(key+"-Disk", parts[0])
	// create or truncate the zip file
	path := filepath.Join(fileWorkingPath, (key + "-" + parts[0] + ".zip"))
	err := file.CreateFile(path)
	if err != nil {
		return task.FAIL, err
	}
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveExplorerProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveExplorerProgress: " + key + "|" + p.GetMessage())
	// update progress
	progress, err := getProgressByMsg(p.GetMessage(), explorerFirstPart)
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ExplorerProgress", progress)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveExplorerInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveExplorerInfo: " + key + "|" + p.GetMessage())
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ExplorerTotal", total)
	redis.RedisSet(key+"-ExplorerCount", 0)
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveExplorerData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveExplorerData: " + key + "|" + p.GetMessage())
	// write file
	path := filepath.Join(fileWorkingPath, (key + "-" + redis.RedisGetString(key+"-Disk") + ".zip"))
	err := file.WriteFile(path, p)
	if err != nil {
		return task.FAIL, err
	}

	// update progress
	redis.RedisSet_AddInteger((key + "-ExplorerCount"), 1)
	progress := int(explorerFirstPart) + getProgressByCount(redis.RedisGetInt(key+"-ExplorerCount"), redis.RedisGetInt(key+"-ExplorerTotal"), 65426, explorerSecondPart)
	redis.RedisSet(key+"-ExplorerProgress", progress)

	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveExplorerEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveExplorerEnd: " + key + "|" + p.GetMessage())

	filename := key + "-" + redis.RedisGetString(key+"-Disk")
	srcPath := filepath.Join(fileWorkingPath, (filename + ".zip"))
	workPath := filepath.Join(fileWorkingPath, filename+".txt")
	unstagePath := filepath.Join(fileUnstagePath, (filename + ".txt"))
	// unzip data
	err := file.UnzipFile(srcPath, workPath, redis.RedisGetInt(key+"-ExplorerTotal"))
	if err != nil {
		return task.FAIL, err
	}
	// move to Unstage
	err = file.MoveFile(workPath, unstagePath)
	if err != nil {
		return task.FAIL, err
	}
	inject_chan, err := channelmap.GetDiskChannel(key)
	if err != nil {
		return task.FAIL, err
	}
	<-inject_chan
	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveExplorerError(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Error("GiveExplorerError: " + p.GetRkey() + "|" + p.GetMessage())
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.FAIL, errors.New(p.GetMessage())
}

func updateDriveProgress(key string) {
	for {
		driveProgress := int((float64(redis.RedisGetInt(key+"-DriveCount"))/float64(redis.RedisGetInt(key+"-DriveTotal")))*100 + float64(redis.RedisGetInt(key+"-ExplorerProgress"))/float64(redis.RedisGetInt(key+"-DriveTotal")))
		if driveProgress >= 100 {
			break
		}
		query.Update_progress(driveProgress, key, "StartGetDrive")
		time.Sleep(time.Duration(config.Viper.GetInt("UPDATE_INTERVAL")) * time.Second)
	}
}
