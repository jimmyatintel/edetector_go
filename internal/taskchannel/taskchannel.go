package taskchannel

import (
	packet "edetector_go/internal/packet"
	"sync"
)

var TaskMu *sync.Mutex
var TaskWorkerChannel map[string](*chan packet.Packet)

func init() {
	TaskMu = &sync.Mutex{}
}