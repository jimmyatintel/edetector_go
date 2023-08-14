package work

import (
	"archive/zip"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/internal/treebuilder"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"net"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var driveMu sync.Mutex
var ExplorerTotalMap = make(map[string]int)
var explorerCountMap = make(map[string]int)
var driveProgressMap = make(map[string]int)
var treeWorkingPath = "treeWorking"
var treeUnstagePath = "treeUnstage"

func Explorer(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("Explorer: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	path := filepath.Join(treeWorkingPath, (p.GetRkey() + ".zip"))
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return task.FAIL, err
	}
	file.Close()
	parts := strings.Split(p.GetMessage(), "|")
	total, err := strconv.Atoi(parts[0])
	if err != nil {
		return task.FAIL, err
	}
	ExplorerTotalMap[p.GetRkey()] = total
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
	logger.Debug("GiveExplorerData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	key := p.GetRkey()
	// lastNewlineInd := strings.LastIndex(p.GetMessage(), "\n")
	// var realData string
	// if lastNewlineInd >= 0 {
	// 	realData = p.GetMessage()[:lastNewlineInd+1]
	// } else {
	// 	logger.Error("Invalid GiveExplorerData")
	// }
	// treebuilder.DetailsMap[key] += realData

	// write file
	path := filepath.Join(treeWorkingPath, (p.GetRkey() + ".zip"))
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return task.FAIL, err
	}
	_, err = file.Seek(0, 2)
	if err != nil {
		return task.FAIL, err
	}
	messageBytes := []byte(p.GetMessage())
	_, err = file.Write(messageBytes)
	if err != nil {
		return task.FAIL, err
	}
	file.Close()

	// update progress
	parts := strings.Split(p.GetMessage(), "|")
	count, err := strconv.Atoi(parts[0])
	if err != nil {
		return task.FAIL, err
	}
	explorerCountMap[key] = count
	driveMu.Lock()
	driveProgressMap[key] = int(((float64(driveCountMap[key]) / float64(driveTotalMap[key])) + (float64(explorerCountMap[key]) / float64(ExplorerTotalMap[key]) / float64(driveTotalMap[key]))) * 100)
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
	logger.Info("GiveExplorerEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))

	// open the zip file for reading
	path := filepath.Join(treeWorkingPath, (p.GetRkey() + ".zip"))
	reader, err := zip.OpenReader(path)
	if err != nil {
		logger.Error("Error opening zip file: ", zap.Any("message", err.Error()))
		return task.FAIL, err
	}
	defer reader.Close()
	// Extract the files from the zip archive
	for _, file := range reader.File {
		destPath := filepath.Join(treeUnstagePath, file.Name)
		if !file.FileInfo().IsDir() {
			// Create the file
			destFile, err := os.Create(destPath)
			if err != nil {
				logger.Error("Error creating file:", zap.Any("error", err.Error()))
				continue
			}
			defer destFile.Close()
			// Open the file inside the zip archive
			srcFile, err := file.Open()
			if err != nil {
				logger.Error("Error opening file inside zip:", zap.Any("error", err.Error()))
				continue
			}
			defer srcFile.Close()
			// Copy the contents from the source to the destination file
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				logger.Error("Error copying file contents:", zap.Any("error", err.Error()))
				continue
			}
		} else {
			logger.Error("the zip file contains directory")
		}
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
	treebuilder.Finished <- p.GetRkey()
	<-user_explorer[p.GetRkey()]
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
	for {
		driveMu.Lock()
		if driveProgressMap[clientid] >= 100 {
			break
		}
		rowsAffected := query.Update_progress(driveProgressMap[clientid], clientid, "StartGetDrive")
		driveMu.Unlock()
		if rowsAffected != 0 {
			go taskservice.RequestToUser(clientid)
		}

	}
}
