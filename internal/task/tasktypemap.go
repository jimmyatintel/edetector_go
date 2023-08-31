package task

var TaskTypeMap map[TaskType]string
var UserTaskTypeMap map[UserTaskType]string

func init() {
	TaskTypeMap = map[TaskType]string{
		GIVE_SCAN_INFO:     "StartScan",
		GIVE_SCAN_PROGRESS: "StartScan",
		GIVE_SCAN_FRAGMENT: "StartScan",
		GIVE_SCAN:          "StartScan",
		GIVE_SCAN_END:      "StartScan",

		GIVE_DRIVE_INFO:        "StartGetDrive",
		EXPLORER:               "StartGetDrive",
		GIVE_EXPLORER_PROGRESS: "StartGetDrive",
		GIVE_EXPLORER_INFO:     "StartGetDrive",
		GIVE_EXPLORER_DATA:     "StartGetDrive",
		GIVE_EXPLORER_END:      "StartGetDrive",
		GIVE_EXPLORER_ERROR:    "StartGetDrive",

		GIVE_COLLECT_PROGRESS:   "StartCollect",
		GIVE_COLLECT_DATA_INFO:  "StartCollect",
		GIVE_COLLECT_DATA:       "StartCollect",
		GIVE_COLLECT_DATA_END:   "StartCollect",
		GIVE_COLLECT_DATA_ERROR: "StartCollect",
	}

	UserTaskTypeMap = map[UserTaskType]string{
		CHANGE_DETECT_MODE: "ChangeDetectMode",
		START_SCAN:         "StartScan",
		START_GET_DRIVE:    "StartGetDrive",
		START_COLLECTION:   "StartCollect",
		TERMINATE:          "Terminate",
	}
}
