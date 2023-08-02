package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	elasticquery "edetector_go/pkg/elastic/query"
	"edetector_go/pkg/logger"
	"encoding/json"
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// type ProcessOverJson struct {
//0 	PID               int    `json:"pid"`
//1 	Mode              string `json:"mode"`
//2 	ProcessCTime      int    `json:"process_c_time"`
//3 	ProcessTime       string `json:"process_time"`
//4 	ProcessName       string `json:"process_name"`
//5 	ProcessPath       string `json:"process_path"`
//6 	ProcessHash       string `json:"process_hash"`
//7 	Parent_PID        int    `json:"parent_pid"`
//8 	ParentCTime       int    `json:"parent_C_time"`
//9 	ParentPath        string `json:"parent_path"`
//10 	InjectedHash      string `json:"injected_hash"`
//11 	StartRun          int    `json:"start_run"`
//12 	HideAttribute     int    `json:"hide_attribute"`
//13 	HideProcess       int    `json:"hide_process"`
//14 	SignerSubjectName string `json:"signer_subject_name"`
//15 	Injection         string `json:"injection"`
//16 	DllStr            string `json:"dll_str"`
//17 	InlineStr         string `json:"inline_str"`
//18 	NetStr            string `json:"net_str"`
// }

type Memory struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	ProcessName       string `json:"processName"`
	ProcessCreateTime int    `json:"processCreateTime"`
	ProcessConnectIP  string `json:"processConnectIP"`
	DynamicCommand    string `json:"dynamicCommand"`
	ProcessMD5        string `json:"processMD5"`
	ProcessPath       string `json:"processPath"`
	ParentProcessId   int    `json:"parentProcessId"`
	ParentProcessName string `json:"parentProcessName"`
	ParentProcessPath string `json:"parentProcessPath"`
	DigitalSign       string `json:"digitalSign"`
	ProcessId         int    `json:"processId"`
	RiskLevel         int    `json:"riskLevel"`
	InjectActive      string `json:"injectActive"`
	ProcessBeInjected int    `json:"processBeInjected"`
	Boot              string `json:"boot"`
	Hook              string `json:"hook"`
	ImportOtherDLL    bool   `json:"importOtherDLL"`
	Hide              string `json:"hide"`
	Mode              string `json:"mode"`
}

func (n Memory) Elastical() ([]byte, error) {
	return json.Marshal(n)
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

	// send to elasticsearch
	lines := strings.Split(p.GetMessage(), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		//! tmp version
		original := strings.Split(line, "|")
		line = original[4] + "@|@" + original[2] + "@|@detecting@|@cmd@|@" + original[6] + "@|@" + original[5] + "@|@" + original[7] + "@|@parentName@|@" + original[9] + "@|@" + original[14] + "@|@" + original[0] + "@|@1@|@0,0@|@0@|@0,0@|@null@|@0@|@" + original[13] + "," + original[12] + "@|@detect"
		values := strings.Split(line, "@|@")
		//! tmp version
		uuid := uuid.NewString()
		int_date, err := strconv.Atoi(values[1])
		if err != nil {
			logger.Error("Invalid date: ", zap.Any("message", values[1]))
			int_date = 0
		}
		logger.Info("memory info", zap.Any("message", line))
		elasticquery.SendToMainElastic(uuid, "ed_de_memory", p.GetRkey(), values[0], int_date, "memory", values[12], "ed_high")
		elasticquery.SendToDetailsElastic(uuid, "ed_de_memory", p.GetRkey(), line, &Memory{}, "ed_high")
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
