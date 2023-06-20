package task

type TaskType string

const (
	GIVE_INFO              TaskType = "GiveInfo"
	OPEN_CHECK_THREAD      TaskType = "OpenCheckthread"
	UNDEFINE               TaskType = "Undefine"
	GIVE_DETECT_INFO_FIRST TaskType = "GiveDetectInfoFirst"
	GIVE_DETECT_INFO       TaskType = "GiveDetectInfo"
	UPDATE_DETECT_MODE     TaskType = "UpdateDetectMode"
	CHECK_CONNECT          TaskType = "CheckConnect"

	GET_PROCESS_HISTORY_INFO  TaskType = "GetProcessHistoryInfo"
	GIVE_PROCESS_HISTORY      TaskType = "GiveProcessHistory"
	GIVE_PROCESS_HISTORY_DATA TaskType = "GiveProcessHistoryData"
	GIVE_PROCESS_HISTORY_END  TaskType = "GiveProcessHistoryEnd"

	GIVE_DETECT_PROCESS_RISK TaskType = "GiveDetectProcessRisk"
	GET_DETECT_PROCESS_RISK  TaskType = "GetDetectProcessRisk"
	GIVE_DETECT_PROCESS_OVER TaskType = "GiveDetectProcessOver"
	GIVE_DETECT_PROCESS_END  TaskType = "GiveDetectProcessEnd"

	GET_PROCESS_INFORMATION  TaskType = "GetProcessInformation"
	GIVE_PROCESS_INFORMATION TaskType = "GiveProcessInformation"
	GIVE_PROCESS_INFO_DATA   TaskType = "GiveProcessInfoData"
	GIVE_PROCESS_INFO_END    TaskType = "GiveProcessInfoEnd"

	GIVE_NETWORK_HISTORY      TaskType = "GiveNetworkHistory"
	GET_NETWORK_HISTORY_INFO  TaskType = "GetNetworkHistoryInfo"
	GIVE_NETWORK_HISTORY_DATA TaskType = "GiveNetworkHistoryData"
	GIVE_NETWORK_HISTORY_END  TaskType = "GiveNetworkHistoryEnd"

	DATA_RIGHT TaskType = "DataRight"

	// Task from API
	CHANGE_DETECT_MODE TaskType = "ChangeDetectMode"
)

var Worklist = []TaskType{GIVE_INFO, OPEN_CHECK_THREAD, GIVE_DETECT_INFO_FIRST, GIVE_DETECT_INFO, UPDATE_DETECT_MODE, CHECK_CONNECT, GIVE_PROCESS_HISTORY,
	GET_PROCESS_INFORMATION, GIVE_DETECT_PROCESS_RISK, GIVE_PROCESS_INFORMATION, GIVE_PROCESS_INFO_DATA, GIVE_PROCESS_INFO_END,
	GIVE_NETWORK_HISTORY, GIVE_NETWORK_HISTORY_DATA, GIVE_NETWORK_HISTORY_END, GIVE_PROCESS_HISTORY_DATA, GIVE_PROCESS_HISTORY_END,
	GIVE_DETECT_PROCESS_OVER, GIVE_DETECT_PROCESS_END, CHANGE_DETECT_MODE}
