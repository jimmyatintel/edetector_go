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
	REJECT_AGENT           TaskType = "RejectAgent"
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
	GET_IMAGE            TaskType = "GetImage"
	READY_IMAGE          TaskType = "ReadyImage"
	GIVE_IMAGE_PATH_INFO TaskType = "GiveImagePathInfo"
	GIVE_IMAGE_PATH      TaskType = "GiveImagePath"
	GIVE_IMAGE_PATH_END  TaskType = "GiveImagePathEnd"
	GIVE_IMAGE_PROGRESS  TaskType = "GiveImageProgress"
	GIVE_IMAGE_INFO      TaskType = "GiveImageInfo"
	GIVE_IMAGE           TaskType = "GiveImage"
	GIVE_IMAGE_END       TaskType = "GiveImageEnd"

	// update
	UPDATE_AGENT       TaskType = "UpdateAgent"
	READY_UPDATE_AGENT TaskType = "ReadyUpdateAgent"
	GIVE_UPDATE_INFO   TaskType = "GiveUpdateInfo"
	GIVE_UPDATE        TaskType = "GiveUpdate"
	GIVE_UPDATE_END    TaskType = "GiveUpdateEnd"

	// remove
	REMOVE_AGENT TaskType = "RemoveAgent"

	// yara rule
	YARA_RULE            TaskType = "YaraRule"
	READY_YARA_RULE      TaskType = "ReadyYaraRule"
	GIVE_YARA_RULE_INFO  TaskType = "GiveYaraRuleInfo"
	GIVE_YARA_RULE       TaskType = "GiveYaraRule"
	GIVE_YARA_RULE_END   TaskType = "GiveYaraRuleEnd"
	GIVE_PATH_INFO       TaskType = "GivePathInfo"
	GIVE_PATH            TaskType = "GivePath"
	GIVE_PATH_END        TaskType = "GivePathEnd"
	GIVE_RULE_MATCH_INFO TaskType = "GiveRuleMatchInfo"
	GIVE_RULE_MATCH      TaskType = "GiveRuleMatch"
	GIVE_RULE_MATCH_END  TaskType = "GiveRuleMatchEnd"
	GIVE_YARA_PROGRESS   TaskType = "GiveYaraProgress"

	// terminate
	TERMINATE_ALL    TaskType = "TerminateAll"
	FINISH_TERMINATE TaskType = "FinishTerminate"

	// memory tree
	GET_MEMORY_TREE           TaskType = "GetMemoryTree"
	GIVE_MEMORY_TREE_INFO     TaskType = "GiveMemoryTreeInfo"
	GIVE_MEMORY_TREE          TaskType = "GiveMemoryTree"
	GIVE_MEMORY_TREE_END      TaskType = "GiveMemoryTreeEnd"
	GIVE_MEMORY_TREE_PROGRESS TaskType = "GiveMemoryTreeProgress"

	// dump dll
	GET_DUMP_DLL       TaskType = "GetDumpDll"
	GIVE_DUMP_DLL_INFO TaskType = "GiveDumpDllInfo"
	GIVE_DUMP_DLL_DATA TaskType = "GiveDumpDllData"
	GIVE_DUMP_DLL_END  TaskType = "GiveDumpDllEnd"

	// dump process
	GET_DUMP_PROCESS       TaskType = "GetDumpProcess"
	GIVE_DUMP_PROCESS_INFO TaskType = "GiveDumpProcessInfo"
	GIVE_DUMP_PROCESS_DATA TaskType = "GiveDumpProcessData"
	GIVE_DUMP_PROCESS_END  TaskType = "GiveDumpProcessEnd"

	// load dll
	GET_LOAD_DLL       TaskType = "GetLoadDll"
	GIVE_LOAD_DLL_DATA TaskType = "GiveLoadDllData"
	GIVE_LOAD_DLL_END  TaskType = "GiveLoadDllEnd"

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
	START_REMOVE            UserTaskType = "StartRemove"
	START_YARA_RULE         UserTaskType = "StartYaraRule"
	START_MEMORY_TREE       UserTaskType = "StartMemoryTree"
	START_LOAD_DLL          UserTaskType = "StartLoadDll"
	START_DUMP_DLL          UserTaskType = "StartDumpDll"
	START_DUMP_PROCESS      UserTaskType = "StartDumpProcess"
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
	READY_IMAGE,
	GIVE_IMAGE_PROGRESS,
	GIVE_IMAGE_INFO,
	GIVE_IMAGE,
	GIVE_IMAGE_END,
	GIVE_IMAGE_PATH_INFO,
	GIVE_IMAGE_PATH,
	GIVE_IMAGE_PATH_END,
	UPDATE_AGENT,
	READY_UPDATE_AGENT,
	GIVE_UPDATE_INFO,
	GIVE_UPDATE,
	GIVE_UPDATE_END,
	REMOVE_AGENT,
	YARA_RULE,
	READY_YARA_RULE,
	GIVE_YARA_RULE_INFO,
	GIVE_YARA_RULE,
	GIVE_YARA_RULE_END,
	GIVE_PATH_INFO,
	GIVE_PATH,
	GIVE_PATH_END,
	GIVE_RULE_MATCH_INFO,
	GIVE_RULE_MATCH,
	GIVE_RULE_MATCH_END,
	GIVE_YARA_PROGRESS,
	TERMINATE_ALL,
	FINISH_TERMINATE,
	GET_MEMORY_TREE,
	GIVE_MEMORY_TREE_INFO,
	GIVE_MEMORY_TREE,
	GIVE_MEMORY_TREE_END,
	GIVE_MEMORY_TREE_PROGRESS,
	GET_DUMP_DLL,
	GIVE_DUMP_DLL_INFO,
	GIVE_DUMP_DLL_DATA,
	GIVE_DUMP_DLL_END,
	GET_DUMP_PROCESS,
	GIVE_DUMP_PROCESS_INFO,
	GIVE_DUMP_PROCESS_DATA,
	GIVE_DUMP_PROCESS_END,
	GET_LOAD_DLL,
	GIVE_LOAD_DLL_DATA,
	GIVE_LOAD_DLL_END,
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
	START_REMOVE,
	START_MEMORY_TREE,
	START_LOAD_DLL,
	START_DUMP_DLL,
	START_DUMP_PROCESS,
	START_YARA_RULE,
	TERMINATE,
	USER_UNDEFINE,
}
