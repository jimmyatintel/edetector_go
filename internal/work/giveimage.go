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
var imageFilePath = "ImageFile"

func init() {
	file.ClearDirContent(imageWorkingPath)
	file.CheckDir(imageFilePath)
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
	content := getDataPacketContent(p)
	err := file.WriteFile(path, content)
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
	srcPath := filepath.Join(imageWorkingPath, key)
	// truncate data
	err := file.TruncateFile(srcPath, redis.RedisGetInt(key+"-ImageTotal"))
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
