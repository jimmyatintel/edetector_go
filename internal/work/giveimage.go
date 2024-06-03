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
	"os"
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
	if strings.Split(p.GetMessage(), "/")[0] == "1" {
		imageFirstPart = float64(config.Viper.GetInt("IMAGE_FIRST_PART"))
		imageSecondPart = 100 - imageFirstPart
		redis.RedisSet(key+"-ImageProgress", 0)
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
	content := getDataPacketContent(p)
	err := file.WriteFile(path, content)
	if err != nil {
		return task.FAIL, err
	}
	// update progress
	redis.RedisSet_AddInteger((key + "-ImageCount"), 1)
	imageCount, err := redis.RedisGetInt(key + "-ImageCount")
	if err != nil {
		return task.FAIL, err
	}
	imageTotal, err := redis.RedisGetInt(key + "-ImageTotal")
	if err != nil {
		return task.FAIL, err
	}
	progress := int(imageFirstPart) + getProgressByCount(imageCount, imageTotal, 65436, imageSecondPart)
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
	srcPath := filepath.Join(imageWorkingPath, key)
	imageTotal, err := redis.RedisGetInt(key + "-ImageTotal")
	if err != nil {
		return task.FAIL, err
	}
	// truncate data
	err = file.TruncateFile(srcPath, imageTotal)
	if err != nil {
		return task.FAIL, err
	}
	// store the ImageFile
	err = storeImageFile(key, srcPath)
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

func ImageError(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("ImageError: " + p.GetRkey() + "::" + p.GetMessage())

	// set task status in mariadb to 0
	rowAffected := query.Update_task_status(p.GetRkey(), "StartGetImage", 2, 0)
	if rowAffected > 0 {
		logger.Info("ImageError: " + p.GetRkey() + "::" + "Update task status to 0 to retry")
	} else {
		logger.Error("ImageError: " + p.GetRkey() + "::" + "Update task status failed")
		return task.FAIL, errors.New("update task status failed")
	}

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
		imageProgress, err := redis.RedisGetInt(key + "-ImageProgress")
		if err != nil {
			logger.Error("Get image progress failed: " + err.Error())
			return
		}
		query.Update_progress(imageProgress, key, "StartGetImage")
		time.Sleep(time.Duration(config.Viper.GetInt("UPDATE_INTERVAL")) * time.Second)
	}
}

func storeImageFile(key string, srcPath string) error {
	extension, err := getExtension(srcPath)
	if err != nil {
		return err
	}
	ip, _, err := query.GetMachineIPandName(key)
	if err != nil {
		return err
	}
	time := time.Now().Format("2006_0102_150405")
	imageType := getTaskMsg(key, "StartGetImage")
	// clear all the content of the directory
	err = file.ClearDirContent(filepath.Join(imageFilePath, ip))
	if err != nil {
		return err
	}
	// move to ImagePath
	dstPath := filepath.Join(imageFilePath, ip, (("Obtained_" + time + "_" + imageType) + extension))
	err = file.MoveFile(srcPath, dstPath)
	if err != nil {
		return err
	}
	return nil
}

func getExtension(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	// check the extension type
	var firstByte [1]byte
	_, err = f.Read(firstByte[:])
	if err != nil {
		return "", err
	}
	extension := ".tar.gz"
	if firstByte[0] == 'P' {
		extension = ".zip"
	}
	return extension, nil
}
