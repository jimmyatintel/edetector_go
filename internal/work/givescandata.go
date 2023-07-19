package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"math"
	"net"
	"strconv"

	"go.uber.org/zap"
)

type ProcessScanJson struct {
	PID         int    `json:"pid"`
	Parent_PID  int    `json:"parent_pid"`
	ProcessName string `json:"process_name"`
	ProcessTime int    `json:"process_time"`
	ParentName  string `json:"parent_name"`
	ParentTime  int    `json:"parent_time"`
	FilePath    string `json:"file_path"`
	UserName    string `json:"user_name"`
	IsPacked    bool   `json:"is_packeted"`
	CommandLine string `json:"command_line"`
	IsHide      bool   `json:"is_hide"`
}

type ScanInfoJson struct {
	PID        int    `json:"pid"`
	FineName   string `json:"file_name"`
	FilePath   string `json:"file_path"`
	FileHash   string `json:"file_hash"`
	Isinjected int    `json:"is_injected"`
	Mode       string `json:"mode"`
	Count      int    `json:"count"`
	AllCount   int    `json:"all_count"`
	IsStartRun bool   `json:"is_start_run"`
}

type ScanOverJson struct {
	PID               int    `json:"pid"`
	Mode              string `json:"mode"`
	ProcessTime       int    `json:"process_time"`
	DetectTime        string `json:"detect_time"`
	ProcessName       string `json:"process_name"`
	ProcessPath       string `json:"process_path"`
	ProcessHash       string `json:"process_hash"`
	Parent_PID        int    `json:"parent_pid"`
	ParentTime        int    `json:"parent_time"`
	ParentPath        string `json:"parent_path"`
	InjectedHash      string `json:"injected_hash"`
	StartRun          int    `json:"start_run"`
	HideAttribute     int    `json:"hide_attribute"`
	HideProcess       int    `json:"hide_process"`
	SignerSubjectName string `json:"signer_subject_name"`
	IsInjection       string `json:"is_injection"`
	IsOtherdll        bool   `json:"is_other_dll"`
	IsInlineHook      string `json:"is_inline_hook"`
	IsNetwork         int    `json:"is_network"`
}

type ScanJson struct {
}

func Process(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("Process: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_SCAN_INFO_DATA,
		Message:    "Ring0Process",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GetScanInfoData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GetScanInfoData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_PROCESS_INFO,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveProcessData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveProcessData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveProcessDataEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveProcessDataEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveScanProgress: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	progress, err_ := strconv.Atoi(p.GetMessage())
	if err_ != nil {
		return task.FAIL, err
	}
	progress = int(math.Min(float64(progress*2), 99))
	query.Update_progress(progress, p.GetRkey(), "StartScan")
	go taskservice.RequestToUser(p.GetRkey())
	return task.SUCCESS, nil
}

func GiveScanDataInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveScanDataInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveScanData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanDataOver(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveScanDataOver: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveScanDataEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveScanDataEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "null",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	taskservice.Finish_task(p.GetRkey(), "StartScan")
	return task.SUCCESS, nil
}
