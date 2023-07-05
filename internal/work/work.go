package work

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"

	// "edetector_go/internal/work_functions"
	"net"
)

var WorkMap map[task.TaskType]func(packet.Packet, *string, net.Conn) (task.TaskResult, error)

type Request_data interface {
	Elastical() ([]byte, error)
}

func init() {
	WorkMap = map[task.TaskType]func(packet.Packet, *string, net.Conn) (task.TaskResult, error){
		// handshake
		task.GIVE_INFO:                GiveInfo,
		task.GIVE_DETECT_INFO_FIRST:   GiveDetectInfoFirst,
		task.GIVE_DETECT_INFO:         GiveDetectInfo,

		// check connect & ack
		task.CHECK_CONNECT:            CheckConnect,

		// process history
		task.GIVE_PROCESS_HISTORY:         GiveProcessHistory,
		task.GIVE_PROCESS_HISTORY_DATA:    GiveProcessHistoryData,
		task.GIVE_PROCESS_HISTORY_END:     GiveProcessHistoryEnd,

		// process risk
		task.GIVE_DETECT_PROCESS_RISK:     GiveDetectProcessRisk,
		task.GIVE_DETECT_PROCESS_OVER:     GiveDetectProcessOver,
		task.GIVE_DETECT_PROCESS_END:      GiveDetectProcessEnd,

		// process info
		task.GIVE_PROCESS_INFORMATION:    GiveProcessInformation,
		task.GIVE_PROCESS_INFO_DATA:      GiveProcessInfoData,
		task.GIVE_PROCESS_INFO_END:       GiveProcessInfoEnd,

		// network history
		task.GIVE_NETWORK_HISTORY:        GiveNetworkHistory,
		task.GIVE_NETWORK_HISTORY_DATA:   GiveNetworkHistoryData,
		task.GIVE_NETWORK_HISTORY_END:    GiveNetworkHistoryEnd,

		// drive
		task.GIVE_DRIVE_INFO:             GiveDriveInfo,

		// scan
		// task.GET_PROCESS_INFO:            GetProcessInfo,
		task.GET_SCAN_INFO_DATA:          GetScanInfoData,
		task.PROCESS:                     Process,
		task.GIVE_PROCESS_DATA:           GiveProcessData,
	}
}
