package taskchannel

import (
	packet "edetector_go/internal/packet"
)

var Task_worker_channel map[string](chan packet.Packet)
var Task_detect_channel map[string](chan packet.Packet)