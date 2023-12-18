package task

// fflag "edetector_go/pkg/fflag"
// logger "edetector_go/pkg/logger"
// "go.uber.org/zap"

func GetTaskType(work string) TaskType {
	for _, v := range Worklist {
		if v == TaskType(work) {
			return v
		}
	}
	return UNDEFINE
}

// GetUserTaskType returns the UserTaskType based on the given work string.
// It iterates through the UserWorklist and checks if any UserTaskType matches the work string.
// If a match is found, it returns the corresponding UserTaskType.
// If no match is found, it returns USER_UNDEFINE.
func GetUserTaskType(work string) UserTaskType {
	for _, v := range UserWorklist {
		if v == UserTaskType(work) {
			return v
		}
	}
	return USER_UNDEFINE
}
