package work

import (
	"bytes"
	"edetector_go/config"
	C_AES "edetector_go/internal/C_AES"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/file"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"net"

	"go.uber.org/zap"
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
	logger.Info("GiveCollectProgress: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
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
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err = clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectDataInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveCollectDataInfo: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
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
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err = clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Debug("GiveCollectData: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	// write file
	dp := packet.CheckIsData(p)
	decrypt_buf := bytes.Repeat([]byte{0}, len(dp.Raw_data))
	C_AES.Decryptbuffer(dp.Raw_data, len(dp.Raw_data), decrypt_buf)
	decrypt_buf = decrypt_buf[100:]
	path := filepath.Join(dbWorkingPath, (key + ".zip"))
	err := file.WriteFile(path, decrypt_buf)
	if err != nil {
		return task.FAIL, err
	}
	// update progress
	redis.RedisSet_AddInteger((key + "-CollectCount"), 1)
	progress := int(collectFirstPart) + getProgressByCount(redis.RedisGetInt(key+"-CollectCount"), redis.RedisGetInt(key+"-CollectTotal"), 65436, collectSecondPart)
	redis.RedisSet(key+"-CollectProgress", progress)
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err = clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectDataEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveCollectDataEnd: ", zap.Any("message", key+", Msg: "+p.GetMessage()))

	srcPath := filepath.Join(dbWorkingPath, (key + ".zip"))
	workPath := filepath.Join(dbWorkingPath, key+".db")
	unstagePath := filepath.Join(dbUstagePath, (key + ".db"))
	// truncate data
	err := file.TruncateFile(srcPath, redis.RedisGetInt(key+"-CollectTotal"))
	if err != nil {
		return task.FAIL, err
	}
	// unzip data
	err = file.UnzipFile(srcPath, workPath)
	if err != nil {
		return task.FAIL, err
	}
	// move to Unstage
	err = file.MoveFile(workPath, unstagePath)
	if err != nil {
		return task.FAIL, err
	}

	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err = clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectDataError(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataError: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
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
		rowsAffected := query.Update_progress(redis.RedisGetInt(key+"-CollectProgress"), key, "StartCollect")
		if rowsAffected != 0 {
			go taskservice.RequestToUser(key)
		}
	}
}
