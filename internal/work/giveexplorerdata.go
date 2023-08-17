package work

import (
	"archive/zip"
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/file"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"errors"
	"io"
	"os"

	"path/filepath"
	"strconv"

	"net"
	"strings"

	"go.uber.org/zap"
)

// var driveMu sync.Mutex
var ExplorerTotalMap = make(map[string]int)
var diskMap = make(map[string]string)

// var explorerCountMap = make(map[string]int)
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
	if len(parts) == 4 {
		total, err := strconv.Atoi(parts[1])
		if err != nil {
			return task.FAIL, err
		}
		ExplorerTotalMap[key] = total
		diskMap[key] = parts[2]
		// create or truncate the db file
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

	// // update progress
	// parts := strings.Split(p.GetMessage(), "|")
	// count, err := strconv.Atoi(parts[0])
	// if err != nil {
	// 	return task.FAIL, err
	// }
	// explorerCountMap[key] = count
	// driveMu.Lock()
	// driveProgressMap[key] = int(((float64(driveCountMap[key]) / float64(driveTotalMap[key])) + (float64(explorerCountMap[key]) / float64(ExplorerTotalMap[key]) / float64(driveTotalMap[key]))) * 100)
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
	path := filepath.Join(fileWorkingPath, (key + "-" + diskMap[key] + ".zip"))
	err := file.TruncateFile(path, ExplorerTotalMap[key])
	if err != nil {
		return task.FAIL, err
	}
	err = unzipFile(key)
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

func unzipFile(key string) error {
	// open the zip file for reading
	path := filepath.Join(fileWorkingPath, (key + "-" + diskMap[key] + ".zip"))
	reader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	// extract the files from the zip archive
	for _, file := range reader.File {
		if !file.FileInfo().IsDir() {
			destFile, err := os.Create(filepath.Join(fileUnstagePath, key+"-"+diskMap[key]+".txt"))
			if err != nil {
				return err
			}
			srcFile, err := file.Open()
			if err != nil {
				return err
			}
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
			destFile.Close()
			srcFile.Close()
		} else {
			err = errors.New("the zip file contains a directory")
			return err
		}
	}
	reader.Close()
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
