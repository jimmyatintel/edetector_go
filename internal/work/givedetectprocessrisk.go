package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"net"

	"go.uber.org/zap"
)

type ProcessOverJson struct {
	PID               int    `json:"pid"`
	Mode              string `json:"mode"`
	ProcessCTime      int    `json:"process_c_time"`
	ProcessTime       string `json:"process_time"`
	ProcessName       string `json:"process_name"`
	ProcessPath       string `json:"process_path"`
	ProcessHash       string `json:"process_hash"`
	Parent_PID        int    `json:"parent_pid"`
	ParentCTime       int    `json:"parent_C_time"`
	ParentPath        string `json:"parent_path"`
	InjectedHash      string `json:"injected_hash"`
	StartRun          int    `json:"start_run"`
	HideAttribute     int    `json:"hide_attribute"`
	HideProcess       int    `json:"hide_process"`
	SignerSubjectName string `json:"signer_subject_name"`
	Injection         string `json:"injection"`
	DllStr            string `json:"dll_str"`
	InlineStr         string `json:"inline_str"`
	NetStr            string `json:"net_str"`
}

func GiveDetectProcessRisk(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectProcessRisk: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_DETECT_PROCESS_RISK,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveDetectProcessOver(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveDetectProcessOver: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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

func GiveDetectProcessEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveDetectProcessEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
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
