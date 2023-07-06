package taskchannel

import (
)

var Task_worker_channel map[string](chan []byte)
var Task_detect_channel map[string](chan []byte)