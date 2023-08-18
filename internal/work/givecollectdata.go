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
	"path/filepath"
	"strconv"
	"sync"

	"net"

	"go.uber.org/zap"
)

var collectMu *sync.Mutex
var collectTotalMap = make(map[string]int)
var collectCountMap = make(map[string]int)
var collectProgressMap = make(map[string]int)
var dbWorkingPath = "dbWorking"
var dbUstagePath = "dbUnstage"

func init() {
	collectMu = &sync.Mutex{}
	file.CheckDir(dbWorkingPath)
	file.CheckDir(dbUstagePath)
}

func GiveCollectInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	key := p.GetRkey()
	logger.Info("GiveCollectInfo: ", zap.Any("message", key+", Msg: "+p.GetMessage()))
	collectMu.Lock()
	collectProgressMap[key] = 0
	collectMu.Unlock()
	go updateCollectProgress(key)
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

func GiveCollectProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectProgress: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	// update progress
	progress, err := getProgressByMsg(p.GetMessage(), 20)
	if err != nil {
		return task.FAIL, err
	}
	collectMu.Lock()
	collectProgressMap[p.GetRkey()] = progress
	collectMu.Unlock()

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
	logger.Info("GiveCollectDataInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	// init collect info
	len, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	collectCountMap[p.GetRkey()] = 0
	collectTotalMap[p.GetRkey()] = len

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
	collectCountMap[key] += 1
	collectMu.Lock()
	collectProgressMap[key] = getProgressByCount(collectCountMap[key], collectTotalMap[key], 80)
	collectMu.Unlock()
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
	err := file.TruncateFile(srcPath, collectTotalMap[key])
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
	return task.SUCCESS, nil
}

func updateCollectProgress(key string) {
	for {
		collectMu.Lock()
		if collectProgressMap[key] >= 100 {
			break
		}
		rowsAffected := query.Update_progress(collectProgressMap[key], key, "StartGetDrive")
		collectMu.Unlock()
		if rowsAffected != 0 {
			go taskservice.RequestToUser(key)
		}
	}
}
