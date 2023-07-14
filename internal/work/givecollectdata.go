package work

import (
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"errors"
	"fmt"
	"os"
	"strconv"

	// elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"io/ioutil"
	"net"

	// "encoding/json"
	// "fmt"
	// "strings"

	"go.uber.org/zap"
)

var DataLen int
var FileName = "db.db"
var count = 0

func ImportStartup(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("ImportStartup: ", zap.Any("message", p.GetMessage()))
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

func CollectInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("CollectInfo: ", zap.Any("message", p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_COLLECT_INFO_DATA,
		Message:    "10",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectProgress: ", zap.Any("message", p.GetMessage()))
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

func GiveCollectDataInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataInfo: ", zap.Any("message", p.GetMessage()))
	len, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	DataLen = len
	file, err := os.OpenFile(FileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return task.FAIL, err
	}
	file.Close()
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
	logger.Debug("GiveCollectData: ", zap.Any("message", p.GetMessage()))
	dp := packet.CheckIsData(p)
	decrypt_buf := bytes.Repeat([]byte{0}, len(dp.Raw_data))
	C_AES.Decryptbuffer(dp.Raw_data, len(dp.Raw_data), decrypt_buf)
	decrypt_buf = decrypt_buf[100:]

	file, err := os.OpenFile(FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return task.FAIL, err
	}
	_, err = file.Seek(0, 2)
	if err != nil {
		return task.FAIL, err
	}
	_, err = file.Write(decrypt_buf)
	if err != nil {
		return task.FAIL, err
	}
	file.Close()
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
	count = count + 1
	return task.SUCCESS, nil
}

func GiveCollectDataEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataEnd: ", zap.Any("message", p.GetMessage()))
	data, err := ioutil.ReadFile(FileName)
	if err != nil {
		return task.FAIL, err
	}
	fileInfo, err := os.Stat(FileName)
	if err != nil {
		return task.FAIL, err
	}
	realLen := fileInfo.Size()
	logger.Info("File Size: " + fmt.Sprint(realLen) + " Count: " + fmt.Sprint(count))
	if int(realLen) < DataLen {
		return task.FAIL, errors.New("incomplete data")
	}
	err = ioutil.WriteFile(FileName, data[:DataLen], 0644)
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
	taskservice.Finish_task(p.GetRkey(), "StartCollect")
	return task.SUCCESS, nil
}

func GiveCollectDataError(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataError: ", zap.Any("message", p.GetMessage()))
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
