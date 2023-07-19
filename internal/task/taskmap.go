package task

var TaskTypeMap map[TaskType]string
var UserTaskTypeMap map[UserTaskType]string

func init() {
	TaskTypeMap = map[TaskType]string{
		PROCESS:               "StartScan",
		GET_SCAN_INFO_DATA:    "StartScan",
		GIVE_PROCESS_DATA:     "StartScan",
		GIVE_PROCESS_DATA_END: "StartScan",
		GIVE_SCAN_PROGRESS:    "StartScan",
		GIVE_SCAN_DATA_INFO:   "StartScan",
		GIVE_SCAN_DATA:        "StartScan",
		GIVE_SCAN_DATA_OVER:   "StartScan",
		GIVE_SCAN_DATA_END:    "StartScan",

		IMPORT_STARTUP:          "StartCollect",
		COLLECT_INFO:            "StartCollect",
		GIVE_COLLECT_PROGRESS:   "StartCollect",
		GIVE_COLLECT_DATA_INFO:  "StartCollect",
		GIVE_COLLECT_DATA:       "StartCollect",
		GIVE_COLLECT_DATA_END:   "StartCollect",
		GIVE_COLLECT_DATA_ERROR: "StartCollect",

		GIVE_DRIVE_INFO:     "StartGetDrive",
		EXPLORER:            "StartGetDrive",
		GIVE_EXPLORER_DATA:  "StartGetDrive",
		GIVE_EXPLORER_END:   "StartGetDrive",
		GIVE_EXPLORER_ERROR: "StartGetDrive",
	}

	UserTaskTypeMap = map[UserTaskType]string{
		CHANGE_DETECT_MODE: "ChangeDetectMode",
		START_SCAN:         "StartScan",
		START_GET_DRIVE:    "StartGetDrive",
		START_COLLECTION:   "StartCollect",
	}
}
