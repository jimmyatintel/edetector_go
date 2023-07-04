package workfromapi

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"

	// "edetector_go/internal/work_functions"
	"net"
)

var WorkapiMap map[task.UserTaskType]func(packet.UserPacket, *string, net.Conn) (task.TaskResult, error)

func init() {
	WorkapiMap = map[task.UserTaskType]func(packet.UserPacket, *string, net.Conn) (task.TaskResult, error){
		task.CHANGE_DETECT_MODE: ChangeDetectMode,
	}
}
