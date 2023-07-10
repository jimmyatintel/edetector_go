package task

type TaskType string
type UserTaskType string

const (
	// handshake
	GIVE_INFO              TaskType = "GiveInfo"
	GIVE_DETECT_PORT_INFO  TaskType = "GiveDetectPortInfo"
	OPEN_CHECK_THREAD      TaskType = "OpenCheckthread"
	GIVE_DETECT_INFO_FIRST TaskType = "GiveDetectInfoFirst"
	GIVE_DETECT_INFO       TaskType = "GiveDetectInfo"
	UNDEFINE               TaskType = "Undefine"

	// check connect & ack
	CHECK_CONNECT          TaskType = "CheckConnect"
	DATA_RIGHT                TaskType = "DataRight"

	// process history
	GET_PROCESS_HISTORY_INFO  TaskType = "GetProcessHistoryInfo"
	GIVE_PROCESS_HISTORY      TaskType = "GiveProcessHistory"
	GIVE_PROCESS_HISTORY_DATA TaskType = "GiveProcessHistoryData"
	GIVE_PROCESS_HISTORY_END  TaskType = "GiveProcessHistoryEnd"

	// process risk
	GET_DETECT_PROCESS_RISK  TaskType = "GetDetectProcessRisk"
	GIVE_DETECT_PROCESS_RISK TaskType = "GiveDetectProcessRisk"
	GIVE_DETECT_PROCESS_OVER TaskType = "GiveDetectProcessOver"
	GIVE_DETECT_PROCESS_END  TaskType = "GiveDetectProcessEnd"

	// process info
	GET_PROCESS_INFORMATION  TaskType = "GetProcessInformation"
	GIVE_PROCESS_INFORMATION TaskType = "GiveProcessInformation"
	GIVE_PROCESS_INFO_DATA   TaskType = "GiveProcessInfoData"
	GIVE_PROCESS_INFO_END    TaskType = "GiveProcessInfoEnd"

	// network history
	GET_NETWORK_HISTORY_INFO  TaskType = "GetNetworkHistoryInfo"
	GIVE_NETWORK_HISTORY      TaskType = "GiveNetworkHistory"
	GIVE_NETWORK_HISTORY_DATA TaskType = "GiveNetworkHistoryData"
	GIVE_NETWORK_HISTORY_END  TaskType = "GiveNetworkHistoryEnd"

	// detect
	UPDATE_DETECT_MODE        TaskType = "UpdateDetectMode"

	// scan
	GET_SCAN_INFO_DATA        TaskType = "GetScanInfoData"
	GET_PROCESS_INFO          TaskType = "GetProcessInfo"
	PROCESS                   TaskType = "Process"
	GIVE_PROCESS_DATA         TaskType = "GiveProcessData"
	GIVE_PROCESS_DATA_END     TaskType = "GiveProcessDataEnd"
	SCAN_PROGRESS             TaskType = "ScanProgress"
	GIVE_SCAN_PROGRESS        TaskType = "GiveScanProgress"
	GIVE_SCAN_DATA            TaskType = "GiveScanData"
	GIVE_SCAN_DATA_INFO       TaskType = "GiveScanDataInfo"
	GIVE_SCAN_DATA_OVER       TaskType = "GiveScanDataOver"

	// collection
	GET_COLLECT_INFO_DATA     TaskType = "GetCollectInfoData"
	GIVE_COLLECT_PROGRESS     TaskType = "GiveCollectProgress"
	GIVE_COLLECT_DATA_INFO    TaskType = "GiveCollectDataInfo"
	GIVE_COLLECT_DATA         TaskType = "GiveCollectData"
	GIVE_COLLECT_DATA_END     TaskType = "GiveCollectDataEnd"
	GIVE_COLLECT_DATA_ERROR   TaskType = "GiveCollectDataError"
	
	// drive
	GET_DRIVE                 TaskType = "GetDrive" 
	GIVE_DRIVE_INFO           TaskType = "GiveDriveInfo"
	TRANSPORT_EXPLORER        TaskType = "TransportExplorer"
	GIVE_EXPLORER_DATA        TaskType = "GiveExplorerData"
	GIVE_EXPLORER_END         TaskType = "GiveExplorerEnd"

	// task from API
	CHANGE_DETECT_MODE        UserTaskType = "ChangeDetectMode"
	START_SCAN                UserTaskType = "StartScan"
	START_GET_DRIVE           UserTaskType = "StartGetDrive"
	START_COLLECTION          UserTaskType = "StartCollection"
	USER_TRANSPORT_EXPLORER   UserTaskType = "TransportExplorer"
	USER_UNDEFINE             UserTaskType = "Undefine"
)

var Worklist = []TaskType{GIVE_INFO, GIVE_DETECT_PORT_INFO, OPEN_CHECK_THREAD, GIVE_DETECT_INFO_FIRST, GIVE_DETECT_INFO, UPDATE_DETECT_MODE, CHECK_CONNECT,
	GIVE_PROCESS_HISTORY, GET_PROCESS_INFORMATION, GIVE_DETECT_PROCESS_RISK, GIVE_PROCESS_INFORMATION, GIVE_PROCESS_INFO_DATA, GIVE_PROCESS_INFO_END,
	GIVE_NETWORK_HISTORY, GIVE_NETWORK_HISTORY_DATA, GIVE_NETWORK_HISTORY_END, GIVE_PROCESS_HISTORY_DATA, GIVE_PROCESS_HISTORY_END,
	GIVE_DETECT_PROCESS_OVER, GIVE_DETECT_PROCESS_END, GET_SCAN_INFO_DATA, GIVE_SCAN_PROGRESS, GIVE_SCAN_DATA, GIVE_SCAN_DATA_INFO,
	GIVE_SCAN_DATA_OVER, GET_COLLECT_INFO_DATA, GIVE_COLLECT_PROGRESS, GIVE_COLLECT_DATA_INFO, GIVE_COLLECT_DATA, GIVE_COLLECT_DATA_END,
	GIVE_COLLECT_DATA_ERROR, GET_DRIVE, GIVE_DRIVE_INFO, GET_PROCESS_INFO, PROCESS, GIVE_PROCESS_DATA, GIVE_PROCESS_DATA_END,
	GIVE_EXPLORER_DATA, GIVE_EXPLORER_END}

var UserWorklist = []UserTaskType{CHANGE_DETECT_MODE, START_SCAN, START_GET_DRIVE, START_COLLECTION}