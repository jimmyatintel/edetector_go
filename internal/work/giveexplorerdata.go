package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"os"
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

var DetailsMap = make(map[string](string))
var Finished = make(chan string, 1000)

func Explorer(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("Explorer: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	parts := strings.Split(p.GetMessage(), "|")
	total, err := strconv.Atoi(parts[0])
	if err != nil {
		return task.FAIL, err
	}

	ExplorerTotalMap[p.GetRkey()] = total
	msg := parts[1] + "|" + parts[2] + "|" + parts[3] + "|" + parts[4]
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.TRANSPORT_EXPLORER,
		Message:    msg,
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
	lastNewlineInd := strings.LastIndex(p.GetMessage(), "\n")
	var realData string
	if lastNewlineInd >= 0 {
		realData = p.GetMessage()[:lastNewlineInd+1]
	} else {
		logger.Error("Invalid GiveExplorerData")
	}

	DetailsMap[key] += realData
	// write file
	file, err := os.OpenFile("explorer.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return task.FAIL, err
	}
	_, err = file.Seek(0, 2)
	if err != nil {
		return task.FAIL, err
	}
	messageBytes := []byte(realData)
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
	Finished <- p.GetRkey()

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
