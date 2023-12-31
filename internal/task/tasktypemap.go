package task

var TaskTypeMap map[TaskType]string
var UserTaskTypeMap map[UserTaskType]string

func init() {
	TaskTypeMap = map[TaskType]string{
		GIVE_INFO: "Main",

		GIVE_DETECT_NETWORK:      "DetectNetwork",
		GIVE_DETECT_PROCESS_FRAG: "DetectProcess",
		GIVE_DETECT_PROCESS:      "DetectProcess",

		READY_SCAN: "StartScan",

		EXPLORER: "StartGetDrive",

		GIVE_COLLECT_PROGRESS:  "CollectProgress",
		GIVE_COLLECT_DATA_INFO: "StartCollect",

		READY_IMAGE: "StartGetImage",

		READY_UPDATE_AGENT: "StartUpdate",
	}

	UserTaskTypeMap = map[UserTaskType]string{
		CHANGE_DETECT_MODE: "ChangeDetectMode",
		START_SCAN:         "StartScan",
		START_GET_DRIVE:    "StartGetDrive",
		START_COLLECTION:   "StartCollect",
		START_GET_IMAGE:    "StartGetImage",
		START_UPDATE:       "StartUpdate",
		START_REMOVE:       "StartRemove",
		TERMINATE:          "Terminate",
	}
}
