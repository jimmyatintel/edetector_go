package workfromapi

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"net"
)

func ChangeDetectMode(p packet.Packet, Key *string, conn net.Conn) (task.TaskResult, error) {
	return task.SUCCESS, nil
}
