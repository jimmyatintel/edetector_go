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
	// PID         int    `json:"pid"`
	// Mode  int    `json:"parent_pid"`
	// ProcessName string `json:"process_name"`
	// ProcessPath
	// ProcessHash int    `json:"process_time"`
	// DetectTime
	// ProcessTime
	// Parent_PID
	// ParentTime
	// ParentPath
	// UnknownHash
	// Injected
	// StartRun
	// HideAttribute
	// HideProcess
	// SignerSubjectName
	// IsInjectionOther
	// IsInjectionPE
	// IsNetwork
	// IsOtherdll
	// IsInlineHook
}

func GiveDetectProcessRisk(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveDetectProcessRisk: ", zap.Any("message", p.GetMessage()))
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
	logger.Debug("GiveDetectProcessOver: ", zap.Any("message", p.GetMessage()))
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
	logger.Debug("GiveDetectProcessEnd: ", zap.Any("message", p.GetMessage()))
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
