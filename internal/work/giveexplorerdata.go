package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	// elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"net"
	"strings"

	// "encoding/json"
	// "fmt"

	"go.uber.org/zap"
)

type ExplorerJson struct {
	Ind                 int `json:"ind"`
	FileName            string `json:"file_name"`
	Parent_Ind          int `json:"parent_ind"`
	IsDeleted           bool `json:"isDeleted"`
	IsDirectory         bool `json:"isDirectory"`
	CreateTime          int `json:"create_time"`
	WriteTime           int `json:"write_time"`
	AccessTime          int `json:"access_time"`
	EntryModifiedTime   int `json:"entry_modified_time"`
	Datalen             int `json:"data_len"`
}

func Explorer(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("Explorer: ", zap.Any("message", p.GetMessage()))
	parts := strings.Split(p.GetMessage(), "|")
	msg := parts[1] + "|" + parts[2] + "|" + parts[3] + "|" + parts[4]
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.TRANSPORT_EXPLORER,
		Message:    msg,
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveExplorerData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveExplorerData: ", zap.Any("message", p.GetMessage()))
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

func GiveExplorerEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveExplorerEnd: ", zap.Any("message", p.GetMessage()))
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
	taskservice.Finish_task(p.GetRkey(), "StartGetDrive")
	return task.SUCCESS, nil
}

func GiveExplorerError(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveExplorerError: ", zap.Any("message", p.GetMessage()))
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