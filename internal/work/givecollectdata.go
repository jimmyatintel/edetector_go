package work

import (
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/mariadb/query"
	"errors"
	"math"
	"os"
	"strconv"
	"strings"

	// elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"net"

	// "encoding/json"
	// "fmt"
	// "strings"

	"go.uber.org/zap"
)

var dataLenMap map[string]int
var fileName = "db.db"

func ImportStartup(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("ImportStartup: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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
	logger.Debug("CollectInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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
	logger.Info("GiveCollectProgress: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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
	parts := strings.Split(p.GetMessage(), "/")
	if len(parts) != 2 {
		return task.FAIL, errors.New("invalid progress format")
	}
	numerator, err := strconv.Atoi(parts[0])
	if err != nil {
		return task.FAIL, err
	}
	denominator, err := strconv.Atoi(parts[1])
	if err != nil {
		return task.FAIL, err
	}
	progress := int(math.Min((float64(numerator) / float64(denominator) * 100), 99))
	query.Update_progress(progress, p.GetRkey(), "StartCollect")
	go taskservice.RequestToUser(p.GetRkey())
	return task.SUCCESS, nil
}

func GiveCollectDataInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	len, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	dataLenMap[p.GetRkey()] = len
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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
	logger.Debug("GiveCollectData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	dp := packet.CheckIsData(p)
	decrypt_buf := bytes.Repeat([]byte{0}, len(dp.Raw_data))
	C_AES.Decryptbuffer(dp.Raw_data, len(dp.Raw_data), decrypt_buf)
	decrypt_buf = decrypt_buf[100:]

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
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
	return task.SUCCESS, nil
}

func GiveCollectDataEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	data, err := os.ReadFile(fileName)
	if err != nil {
		return task.FAIL, err
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return task.FAIL, err
	}
	realLen := fileInfo.Size()
	if int(realLen) < dataLenMap[p.GetRkey()] {
		return task.FAIL, errors.New("incomplete data")
	}
	err = os.WriteFile(fileName, data[:dataLenMap[p.GetRkey()]], 0644)
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
