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
)

var explorerFirstPart float64
var explorerSecondPart float64

func Explorer(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("Explorer: " + key + "::" + p.GetMessage())
	explorerFirstPart = float64(config.Viper.GetInt("EXPLORER_FIRST_PART"))
	explorerSecondPart = float64(config.Viper.GetInt("EXPLORER_SECOND_PART")) - explorerFirstPart
	redis.RedisSet(key+"-Disk", p.GetMessage())
	// create or truncate the zip file
	path := filepath.Join(fileWorkingPath, (key + "." + p.GetMessage()))
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
	logger.Debug("GiveExplorerProgress: " + key + "::" + p.GetMessage())
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
	logger.Info("GiveExplorerInfo: " + key + "::" + p.GetMessage())
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
	logger.Debug("GiveExplorerData: " + key)
	// write file
	path := filepath.Join(fileWorkingPath, (key + "." + redis.RedisGetString(key+"-Disk")))
	content := getDataPacketContent(p)
	err := file.WriteFile(path, content)
	if err != nil {
		return task.FAIL, err
	}
	// update progress
	redis.RedisSet_AddInteger((key + "-ExplorerCount"), 1)

	explorerCount, err := redis.RedisGetInt(key + "-ExplorerCount")
	if err != nil {
		return task.FAIL, err
	}
	explorerTotal, err := redis.RedisGetInt(key + "-ExplorerTotal")
	if err != nil {
		return task.FAIL, err
	}

	progress := int(explorerFirstPart) + getProgressByCount(explorerCount, explorerTotal, 65426, explorerSecondPart)
	redis.RedisSet(key+"-ExplorerProgress", progress)

	err = clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveExplorerEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveExplorerEnd: " + key + "::" + p.GetMessage())

	filename := key + "." + redis.RedisGetString(key+"-Disk")
	srcPath := filepath.Join(fileWorkingPath, filename)
	workPath := filepath.Join(fileWorkingPath, filename+".txt")
	unstagePath := filepath.Join(fileUnstagePath, (filename + ".txt"))
	explorerTotal, err := redis.RedisGetInt(key + "-ExplorerTotal")
	if err != nil {
		return task.FAIL, err
	}
	// unzip data
	err = file.DecompressionFile(srcPath, workPath, explorerTotal)
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
	logger.Error("GiveExplorerError: " + p.GetRkey() + "::" + p.GetMessage())
	err := clientsearchsend.SendTCPtoClient(p, task.DATA_RIGHT, "", conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.FAIL, errors.New(p.GetMessage())
}

func updateDriveProgress(key string) {
	for {
		result, err := query.Load_stored_task("nil", key, 2, "StartGetDrive")
		if err != nil {
			logger.Error("Get handling tasks failed: " + err.Error())
			return
		}
		if len(result) == 0 {
			return
		}
		driveCount, err := redis.RedisGetInt(key + "-DriveCount")
		if err != nil {
			logger.Error("Get drive count failed: " + err.Error())
			return
		}
		driveTotal, err := redis.RedisGetInt(key + "-DriveTotal")
		if err != nil {
			logger.Error("Get drive total failed: " + err.Error())
			return
		}
		explorerProgress, err := redis.RedisGetInt(key + "-ExplorerProgress")
		if err != nil {
			logger.Error("Get explorer progress failed: " + err.Error())
			return
		}

		driveProgress := int((float64(driveCount)/float64(driveTotal))*100 + float64(explorerProgress)/float64(driveTotal))
		query.Update_progress(driveProgress, key, "StartGetDrive")
		time.Sleep(time.Duration(config.Viper.GetInt("UPDATE_INTERVAL")) * time.Second)
	}
}
