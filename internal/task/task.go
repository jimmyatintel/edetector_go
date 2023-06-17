package task

import (
// fflag "edetector_go/internal/fflag"
// logger "edetector_go/pkg/logger"
// "go.uber.org/zap"
)

func init() {
}

func GetTaskType(work string) TaskType {
	for _, v := range Worklist {
		if v == TaskType(work) {
			return v
		}
	}
	return UNDEFINE
}
