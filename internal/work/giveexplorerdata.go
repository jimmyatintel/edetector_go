package work

import (
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/file"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"errors"

	"path/filepath"
	"strconv"

	"net"
	"strings"

	"go.uber.org/zap"
)

// var driveMu sync.Mutex
var explorerTotalMap = make(map[string]int)
var explorerCountMap = make(map[string]int)
var driveProgressMap = make(map[string]int)
var diskMap = make(map[string]string)
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
	if len(parts) >= 2 {
		total, err := strconv.Atoi(parts[0])
		if err != nil {
			return task.FAIL, err
		}
		explorerTotalMap[key] = total
		diskMap[key] = parts[1]
		// create or truncate the zip file
		path := filepath.Join(fileWorkingPath, (key + "-" + diskMap[key] + ".zip"))
		err = file.CreateFile(path)
		if err != nil {
			return task.FAIL, err
		}
	} else {
		err := errors.New("invalid msg format")
		return task.FAIL, err
	}
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
	// parts := strings.Split(p.GetMessage(), "|")
	// count, err := strconv.Atoi(parts[0])
	// if err != nil {
	// 	return task.FAIL, err
	// }
	// explorerCountMap[key] = count
	// driveMu.Lock()
	// driveProgressMap[key] = int(((float64(driveCountMap[key]) / float64(driveTotalMap[key])) + (float64(explorerCountMap[key]) / float64(explorerTotalMap[key]) / float64(driveTotalMap[key]))) * 100)
	// driveMu.Unlock()

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

func driveProgress(clientid string) {
	// for {
	// 	driveMu.Lock()
	// 	if driveProgressMap[clientid] >= 100 {
	// 		break
	// 	}
	// 	rowsAffected := query.Update_progress(driveProgressMap[clientid], clientid, "StartGetDrive")
	// 	driveMu.Unlock()
	// 	if rowsAffected != 0 {
	// 		go taskservice.RequestToUser(clientid)
	// 	}

	// }
}
