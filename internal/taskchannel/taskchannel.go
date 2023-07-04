package taskchannel

import (
    packet "edetector_go/internal/packet"
)

var Task_channel map[string](chan packet.WorkPacket)