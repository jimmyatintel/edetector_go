package workfromapi

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"

	// "edetector_go/internal/work_functions"
	"net"
)

var WrokapiMap map[task.TaskType]func(packet.Packet, *string, net.Conn) (task.TaskResult, error)

func init() {
	WrokapiMap = map[task.TaskType]func(packet.Packet, *string, net.Conn) (task.TaskResult, error){
		task.CHANGE_DETECT_MODE: ChangeDetectMode,
	}
}
