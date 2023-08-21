package work

import (
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/file"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	"strconv"

	"path/filepath"

	"net"
	"strings"

	"go.uber.org/zap"
)

var fileWorkingPath = "fileWorking"
var fileUnstagePath = "fileUnstage"

func init() {
	file.CheckDir(fileWorkingPath)
	file.CheckDir(fileUnstagePath)
}

func Explorer(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("Explorer: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	parts := strings.Split(p.GetMessage(), "|")
	redis.RedisSet(key+"-Disk", parts[0])
	// create or truncate the zip file
	path := filepath.Join(fileWorkingPath, (key + "-" + parts[0] + ".zip"))
	err := file.CreateFile(path)
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

func GiveExplorerProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveExplorerProgress: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	// update progress
	progress, err := getProgressByMsg(p.GetMessage(), 75)
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ExplorerProgress", progress)
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

func GiveExplorerInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveExplorerInfo: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	redis.RedisSet(key+"-ExplorerTotal", total)
	redis.RedisSet(key+"-ExplorerCount", 0)
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

func GiveExplorerData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveExplorerData: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	// write file
	dp := packet.CheckIsData(p)
	decrypt_buf := bytes.Repeat([]byte{0}, len(dp.Raw_data))
	C_AES.Decryptbuffer(dp.Raw_data, len(dp.Raw_data), decrypt_buf)
	decrypt_buf = decrypt_buf[100:]
	path := filepath.Join(fileWorkingPath, (key + "-" + redis.RedisGetString(key+"-Disk") + ".zip"))
	err := file.WriteFile(path, decrypt_buf)
	if err != nil {
		return task.FAIL, err
	}

	// update progress
	redis.RedisSet_AddInteger((key + "-ExplorerCount"), 1)
	progress := getProgressByCount(redis.RedisGetInt(key+"-ExplorerCount"), redis.RedisGetInt(key+"-ExplorerTotal"), 25)
	redis.RedisSet(key+"-ExplorerProgress", progress)

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

func GiveExplorerEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveExplorerEnd: ", zap.Any("message", key+", Msg: "+p.GetMessage()))

	filename := key + "-" + redis.RedisGetString(key+"-Disk")
	srcPath := filepath.Join(fileWorkingPath, (filename + ".zip"))
	workPath := filepath.Join(fileWorkingPath, filename+".txt")
	unstagePath := filepath.Join(fileUnstagePath, (filename + ".txt"))
	// truncate data
	err := file.TruncateFile(srcPath, redis.RedisGetInt(key+"-ExplorerTotal"))
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
	<-user_explorer[key]
	return task.SUCCESS, nil
}

func GiveExplorerError(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveExplorerError: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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
	return task.SUCCESS, nil
}

func updateDriveProgress(key string) {
	for {
		driveProgress := int((float64(redis.RedisGetInt(key+"-DriveCount"))/float64(redis.RedisGetInt(key+"-DriveTotal")))*100 + float64(redis.RedisGetInt(key+"-ExplorerProgress"))/float64(redis.RedisGetInt(key+"-DriveTotal")))
		if driveProgress >= 100 {
			break
		}
		rowsAffected := query.Update_progress(driveProgress, key, "StartGetDrive")
		if rowsAffected != 0 {
			go taskservice.RequestToUser(key)
		}
	}
}
