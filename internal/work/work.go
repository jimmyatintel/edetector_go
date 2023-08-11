package work

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"

	"net"
)

var elasticPrefix = "james_"

var WorkMap map[task.TaskType]func(packet.Packet, net.Conn) (task.TaskResult, error)

func init() {
	WorkMap = map[task.TaskType]func(packet.Packet, net.Conn) (task.TaskResult, error){
		// handshake
		task.GIVE_INFO:              GiveInfo,
		task.GIVE_DETECT_INFO_FIRST: GiveDetectInfoFirst,
		task.GIVE_DETECT_INFO:       GiveDetectInfo,

		// check connect & ack
		task.CHECK_CONNECT: CheckConnect,

		// new detect
		task.GIVE_DETECT_NETWORK: GiveDetectNetwork,
		task.GIVE_DETECT_PROCESS: GiveDetectProcess,

		// new scan
		task.GIVE_SCAN_INFO: GiveScanInfo,
		task.GIVE_SCAN:      GiveScan,
		task.GIVE_SCAN_END:  GiveScanEnd,

		// drive
		task.GIVE_DRIVE_INFO:     GiveDriveInfo,
		task.EXPLORER:            Explorer,
		task.GIVE_EXPLORER_DATA:  GiveExplorerData,
		task.GIVE_EXPLORER_END:   GiveExplorerEnd,
		task.GIVE_EXPLORER_ERROR: GiveExplorerError,

		// collection
		task.IMPORT_STARTUP:          ImportStartup,
		task.COLLECT_INFO:            CollectInfo,
		task.GIVE_COLLECT_PROGRESS:   GiveCollectProgress,
		task.GIVE_COLLECT_DATA_INFO:  GiveCollectDataInfo,
		task.GIVE_COLLECT_DATA:       GiveCollectData,
		task.GIVE_COLLECT_DATA_END:   GiveCollectDataEnd,
		task.GIVE_COLLECT_DATA_ERROR: GiveCollectDataError,
	}
}
