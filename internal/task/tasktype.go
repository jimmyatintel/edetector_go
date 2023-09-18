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
	CHECK_CONNECT TaskType = "CheckConnect"
	DATA_RIGHT    TaskType = "DataRight"

	// new detect
	UPDATE_DETECT_MODE       TaskType = "UpdateDetectMode"
	GIVE_DETECT_NETWORK      TaskType = "GiveDetectNetwork"
	GIVE_DETECT_PROCESS_FRAG TaskType = "GiveDetectProcessFrag"
	GIVE_DETECT_PROCESS      TaskType = "GiveDetectProcess"

	// new scan
	GET_SCAN            TaskType = "GetScan"
	READY_SCAN          TaskType = "ReadyScan"
	GIVE_SCAN_INFO      TaskType = "GiveScanInfo"
	GIVE_SCAN_PROGRESS  TaskType = "GiveScanProgress"
	GIVE_SCAN_DATA_INFO TaskType = "GiveScanDataInfo"
	GIVE_SCAN           TaskType = "GiveScan"
	GIVE_SCAN_END       TaskType = "GiveScanEnd"

	// new drive
	GET_DRIVE              TaskType = "GetDrive"
	GIVE_DRIVE_INFO        TaskType = "GiveDriveInfo"
	EXPLORER_INFO          TaskType = "ExplorerInfo"
	EXPLORER               TaskType = "Explorer"
	GIVE_EXPLORER_PROGRESS TaskType = "GiveExplorerProgress"
	GIVE_EXPLORER_INFO     TaskType = "GiveExplorerInfo"
	GIVE_EXPLORER_DATA     TaskType = "GiveExplorerData"
	GIVE_EXPLORER_END      TaskType = "GiveExplorerEnd"
	GIVE_EXPLORER_ERROR    TaskType = "GiveExplorerError"

	// new collection
	GET_COLLECT_INFO        TaskType = "GetCollectInfo"
	GIVE_COLLECT_PROGRESS   TaskType = "GiveCollectProgress"
	GIVE_COLLECT_DATA_INFO  TaskType = "GiveCollectDataInfo"
	GIVE_COLLECT_DATA       TaskType = "GiveCollectData"
	GIVE_COLLECT_DATA_END   TaskType = "GiveCollectDataEnd"
	GIVE_COLLECT_DATA_ERROR TaskType = "GiveCollectDataError"

	// image
	GET_IMAGE       TaskType = "GetImage"
	GIVE_IMAGE_INFO TaskType = "GiveImageInfo"
	GIVE_IMAGE      TaskType = "GiveImage"
	GIVE_IMAGE_END  TaskType = "GiveImageEnd"

	// update
	UPDATE_AGENT       TaskType = "UpdateAgent"
	READY_UPDATE_AGENT TaskType = "ReadyUpdateAgent"
	GIVE_UPDATE_INFO   TaskType = "GiveUpdateInfo"
	GIVE_UPDATE        TaskType = "GiveUpdate"
	GIVE_UPDATE_END    TaskType = "GiveUpdateEnd"

	// terminate
	TERMINATE_ALL TaskType = "TerminateAll"

	// task from API
	CHANGE_DETECT_MODE      UserTaskType = "ChangeDetectMode"
	START_SCAN              UserTaskType = "StartScan"
	START_GET_DRIVE         UserTaskType = "StartGetDrive"
	START_GET_EXPLORER      UserTaskType = "StartGetExplorer"
	START_COLLECTION        UserTaskType = "StartCollect"
	USER_TRANSPORT_EXPLORER UserTaskType = "TransportExplorer"
	START_GET_IMAGE         UserTaskType = "StartGetImage"
	START_UPDATE            UserTaskType = "StartUpdate"
	TERMINATE               UserTaskType = "Terminate"
	USER_UNDEFINE           UserTaskType = "Undefine"
)

var Worklist = []TaskType{
	GIVE_INFO,
	GIVE_DETECT_PORT_INFO,
	OPEN_CHECK_THREAD,
	GIVE_DETECT_INFO_FIRST,
	GIVE_DETECT_INFO,
	UNDEFINE,
	CHECK_CONNECT,
	DATA_RIGHT,
	UPDATE_DETECT_MODE,
	GIVE_DETECT_NETWORK,
	GIVE_DETECT_PROCESS_FRAG,
	GIVE_DETECT_PROCESS,
	GET_SCAN,
	READY_SCAN,
	GIVE_SCAN_INFO,
	GIVE_SCAN_PROGRESS,
	GIVE_SCAN_DATA_INFO,
	GIVE_SCAN,
	GIVE_SCAN_END,
	GET_DRIVE,
	GIVE_DRIVE_INFO,
	EXPLORER_INFO,
	EXPLORER,
	GIVE_EXPLORER_PROGRESS,
	GIVE_EXPLORER_INFO,
	GIVE_EXPLORER_DATA,
	GIVE_EXPLORER_END,
	GIVE_EXPLORER_ERROR,
	GET_COLLECT_INFO,
	GIVE_COLLECT_PROGRESS,
	GIVE_COLLECT_DATA_INFO,
	GIVE_COLLECT_DATA,
	GIVE_COLLECT_DATA_END,
	GIVE_COLLECT_DATA_ERROR,
	GET_IMAGE,
	GIVE_IMAGE_INFO,
	GIVE_IMAGE,
	GIVE_IMAGE_END,
	UPDATE_AGENT,
	READY_UPDATE_AGENT,
	GIVE_UPDATE_INFO,
	GIVE_UPDATE,
	GIVE_UPDATE_END,
	TERMINATE_ALL,
}

var UserWorklist = []UserTaskType{
	CHANGE_DETECT_MODE,
	START_SCAN,
	START_GET_DRIVE,
	START_COLLECTION,
	START_GET_EXPLORER,
	USER_TRANSPORT_EXPLORER,
	START_GET_IMAGE,
	START_UPDATE,
	TERMINATE,
	USER_UNDEFINE,
}
