package work

import (
	"bytes"
	"edetector_go/internal/C_AES"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"

	"net"
)

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
		task.GIVE_DETECT_NETWORK:      GiveDetectNetwork,
		task.GIVE_DETECT_PROCESS_FRAG: GiveDetectProcessFrag,
		task.GIVE_DETECT_PROCESS:      GiveDetectProcess,

		// new scan
		task.READY_SCAN:          ReadyScan,
		task.GIVE_SCAN_INFO:      GiveScanInfo,
		task.GIVE_SCAN_PROGRESS:  GiveScanProgress,
		task.GIVE_SCAN_DATA_INFO: GiveScanDataInfo,
		task.GIVE_SCAN:           GiveScan,
		task.GIVE_SCAN_END:       GiveScanEnd,

		// drive
		task.GIVE_DRIVE_INFO:        GiveDriveInfo,
		task.EXPLORER:               Explorer,
		task.GIVE_EXPLORER_PROGRESS: GiveExplorerProgress,
		task.GIVE_EXPLORER_INFO:     GiveExplorerInfo,
		task.GIVE_EXPLORER_DATA:     GiveExplorerData,
		task.GIVE_EXPLORER_END:      GiveExplorerEnd,
		task.GIVE_EXPLORER_ERROR:    GiveExplorerError,

		// collection
		task.GIVE_COLLECT_PROGRESS:   GiveCollectProgress,
		task.GIVE_COLLECT_DATA_INFO:  GiveCollectDataInfo,
		task.GIVE_COLLECT_DATA:       GiveCollectData,
		task.GIVE_COLLECT_DATA_END:   GiveCollectDataEnd,
		task.GIVE_COLLECT_DATA_ERROR: GiveCollectDataError,

		// image
		task.READY_IMAGE:     ReadyImage,
		task.GIVE_IMAGE_INFO: GiveImageInfo,
		task.GIVE_IMAGE:      GiveImage,
		task.GIVE_IMAGE_END:  GiveImageEnd,

		// rule match
		task.GIVE_RULE_MATCH_INFO: GiveRuleMatchInfo,
		task.GIVE_RULE_MATCH:      GiveRuleMatch,
		task.GIVE_RULE_MATCH_END:  GiveRuleMatchEnd,
		
		// terminate
		task.FINISH_TERMINATE: FinishTerminate,
	}
}

func getDataPacketContent(p packet.Packet) []byte {
	dp := packet.CheckIsData(p)
	decrypt_buf := bytes.Repeat([]byte{0}, len(dp.Raw_data))
	C_AES.Decryptbuffer(dp.Raw_data, len(dp.Raw_data), decrypt_buf)
	decrypt_buf = decrypt_buf[100:]
	return decrypt_buf
}