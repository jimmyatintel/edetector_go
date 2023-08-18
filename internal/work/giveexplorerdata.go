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
	"strconv"
	"sync"

	"path/filepath"

	"net"
	"strings"

	"go.uber.org/zap"
)

var driveMu *sync.Mutex
var explorerTotalMap = make(map[string]int)
var explorerCountMap = make(map[string]int)
var explorerProgressMap = make(map[string]int)
var diskMap = make(map[string]string)
var fileWorkingPath = "fileWorking"
var fileUnstagePath = "fileUnstage"

func init() {
	driveMu = &sync.Mutex{}
	file.CheckDir(fileWorkingPath)
	file.CheckDir(fileUnstagePath)
}

func Explorer(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("Explorer: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	parts := strings.Split(p.GetMessage(), "|")
	diskMap[key] = parts[0]
	// create or truncate the zip file
	path := filepath.Join(fileWorkingPath, (key + "-" + diskMap[key] + ".zip"))
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
	logger.Info("GiveExplorerProgress: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	// update progress
	progress, err := getProgressByMsg(p.GetMessage(), 50)
	if err != nil {
		return task.FAIL, err
	}
	driveMu.Lock()
	explorerProgressMap[p.GetRkey()] = progress
	driveMu.Unlock()

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
	logger.Info("GiveExplorerInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	total, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	explorerCountMap[p.GetRkey()] = 0
	explorerTotalMap[p.GetRkey()] = total
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
	path := filepath.Join(fileWorkingPath, (key + "-" + diskMap[key] + ".txt"))
	err := file.WriteFile(path, decrypt_buf)
	if err != nil {
		return task.FAIL, err
	}

	// update progress
	explorerCountMap[key] += 1
	driveMu.Lock()
	explorerProgressMap[key] = getProgressByCount(explorerCountMap[key], explorerTotalMap[key], 50)
	driveMu.Unlock()

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

	filename := key + "-" + diskMap[key]
	srcPath := filepath.Join(fileWorkingPath, (filename + ".zip"))
	workPath := filepath.Join(fileWorkingPath, filename+".txt")
	unstagePath := filepath.Join(fileUnstagePath, (filename + ".txt"))
	// truncate data
	err := file.TruncateFile(srcPath, explorerTotalMap[key])
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
		driveMu.Lock()
		driveProgress := int((float64(driveCountMap[key])/float64(driveTotalMap[key]))*100 + float64(explorerProgressMap[key])/float64(driveTotalMap[key]))
		driveMu.Unlock()
		if driveProgress >= 100 {
			break
		}
		rowsAffected := query.Update_progress(driveProgress, key, "StartGetDrive")
		if rowsAffected != 0 {
			go taskservice.RequestToUser(key)
		}
	}
}
