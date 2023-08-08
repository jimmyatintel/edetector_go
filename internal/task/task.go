package task

import (
// fflag "edetector_go/internal/fflag"
// logger "edetector_go/pkg/logger"
// "go.uber.org/zap"
)

func GetTaskType(work string) TaskType {
	for _, v := range Worklist {
		if v == TaskType(work) {
			return v
		}
	}
	return UNDEFINE
}

func GetUserTaskType(work string) UserTaskType {
	for _, v := range UserWorklist {
		if v == UserTaskType(work) {
			return v
		}
	}
	return USER_UNDEFINE
}