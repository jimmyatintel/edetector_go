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
		task.TERMINATE:          Terminate,
	}
}
