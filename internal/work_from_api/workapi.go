package workfromapi

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	// "edetector_go/internal/work_functions"
)

var WorkapiMap map[task.UserTaskType]func(packet.UserPacket) (task.TaskResult, error)

func init() {
	WorkapiMap = map[task.UserTaskType]func(packet.UserPacket) (task.TaskResult, error){
		task.CHANGE_DETECT_MODE: ChangeDetectMode,
		task.START_SCAN:         StartScan,
		task.START_GET_DRIVE:    StartGetDrive,
		task.START_COLLECTION:   StartCollect,
		task.START_GET_IMAGE:    StartGetImage,
		task.START_UPDATE:       StartUpdate,
		task.START_REMOVE:       StartRemove,
		task.START_MEMORY_TREE:  StartMemoryTree,
		task.START_LOAD_DLL:     StartLoadDll,
		task.START_DUMP_DLL:     StartDumpDll,
		task.START_DUMP_PROCESS: StartDumpProcess,
		task.START_YARA_RULE:    StartYaraRule,
		task.TERMINATE:          Terminate,
	}
}
