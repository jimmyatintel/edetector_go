package work

import (
	"bytes"
	"edetector_go/internal/C_AES"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	mq "edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"

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
		task.GIVE_IMAGE_PROGRESS: GiveImageProgress,
		task.GIVE_IMAGE_INFO:     GiveImageInfo,
		task.GIVE_IMAGE:          GiveImage,
		task.GIVE_IMAGE_END:      GiveImageEnd,

		// rule match
		task.GIVE_RULE_MATCH_INFO: GiveRuleMatchInfo,
		task.GIVE_RULE_MATCH:      GiveRuleMatch,
		task.GIVE_RULE_MATCH_END:  GiveRuleMatchEnd,
		task.GIVE_YARA_PROGRESS:   GiveYaraProgress,

		// terminate
		task.FINISH_TERMINATE: FinishTerminate,

		// memory tree
		task.GIVE_MEMORY_TREE_INFO: GiveMemoryTreeInfo,
		task.GIVE_MEMORY_TREE:      GiveMemoryTree,
		task.GIVE_MEMORY_TREE_END:  GiveMemoryTreeEnd,
		// task.GIVE_MEMORY_TREE_PROGRESS: GiveMemoryTreeProgress,

		// dump dll
		task.GIVE_DUMP_DLL_INFO: GiveDumpDllInfo,
		task.GIVE_DUMP_DLL_DATA: GiveDumpDllData,
		task.GIVE_DUMP_DLL_END:  GiveDumpDllEnd,

		// dump process
		task.GIVE_DUMP_PROCESS_INFO: GiveDumpProcessInfo,
		task.GIVE_DUMP_PROCESS_DATA: GiveDumpProcessData,
		task.GIVE_DUMP_PROCESS_END:  GiveDumpProcessEnd,

		// load dll
		task.GIVE_LOAD_DLL_DATA: GiveLoadDllData,
		task.GIVE_LOAD_DLL_END:  GiveLoadDllEnd,
	}
}

func getTaskMsg(key string, ttype string) string {
	taskID := mq.Load_task_id(key, ttype, 2)
	content := []byte(redis.RedisGetString(taskID))
	NewPacket := new(packet.TaskPacket)
	err := NewPacket.NewPacket(content)
	if err != nil {
		logger.Error("Error getting task msg: " + err.Error())
		return "Unknown"
	}
	return NewPacket.GetMessage()
}

func getDataPacketContent(p packet.Packet) []byte {
	dp := packet.CheckIsData(p)
	decrypt_buf := bytes.Repeat([]byte{0}, len(dp.Raw_data))
	C_AES.Decryptbuffer(dp.Raw_data, len(dp.Raw_data), decrypt_buf)
	decrypt_buf = decrypt_buf[100:]
	return decrypt_buf
}
