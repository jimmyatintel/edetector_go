package work

import (
	"edetector_go/internal/packet"
	"edetector_go/internal/task"

	"net"
)

var WorkMap map[task.TaskType]func(packet.Packet, net.Conn) (task.TaskResult, error)

type ProcessJson struct {
	ProcessName       string `json:"processName"`
	ProcessConnectIP  string `json:"processConnectIP"`
	DynamicCommand    string `json:"dynamicCommand"`
	ProcessMD5        string `json:"processMD5"`
	ProcessPath       string `json:"processPath"`
	ParentProcessId   int    `json:"parentProcessId"`
	ParentProcessName string `json:"parentProcessName"`
	ParentProcessPath string `json:"parentProcessPath"`
	DigitalSign       string `json:"digitalSign"`
	ImportOtherDLL    bool   `json:"importOtherDLL"`
	ProcessId         int    `json:"processId"`
	RiskLevel         int    `json:"riskLevel"`
	InjectActive      int    `json:"injectActive"`
	ProcessBeInjected string `json:"processBeInjected"`
	Boot              int    `json:"boot"`
	Hook              string `json:"hook"`
	Hide              string `json:"hide"`
	Mode              string `json:"mode"`
}

func init() {
	WorkMap = map[task.TaskType]func(packet.Packet, net.Conn) (task.TaskResult, error){
		// handshake
		task.GIVE_INFO:              GiveInfo,
		task.GIVE_DETECT_INFO_FIRST: GiveDetectInfoFirst,
		task.GIVE_DETECT_INFO:       GiveDetectInfo,

		// check connect & ack
		task.CHECK_CONNECT: CheckConnect,

		// process history
		task.GIVE_PROCESS_HISTORY:      GiveProcessHistory,
		task.GIVE_PROCESS_HISTORY_DATA: GiveProcessHistoryData,
		task.GIVE_PROCESS_HISTORY_END:  GiveProcessHistoryEnd,

		// process risk
		task.GIVE_DETECT_PROCESS_RISK: GiveDetectProcessRisk,
		task.GIVE_DETECT_PROCESS_OVER: GiveDetectProcessOver,
		task.GIVE_DETECT_PROCESS_END:  GiveDetectProcessEnd,

		// process info
		task.GIVE_PROCESS_INFORMATION: GiveProcessInformation,
		task.GIVE_PROCESS_INFO_DATA:   GiveProcessInfoData,
		task.GIVE_PROCESS_INFO_END:    GiveProcessInfoEnd,

		// network history
		task.GIVE_NETWORK_HISTORY:      GiveNetworkHistory,
		task.GIVE_NETWORK_HISTORY_DATA: GiveNetworkHistoryData,
		task.GIVE_NETWORK_HISTORY_END:  GiveNetworkHistoryEnd,

		// drive
		task.GIVE_DRIVE_INFO:     GiveDriveInfo,
		task.EXPLORER:            Explorer,
		task.GIVE_EXPLORER_DATA:  GiveExplorerData,
		task.GIVE_EXPLORER_END:   GiveExplorerEnd,
		task.GIVE_EXPLORER_ERROR: GiveExplorerError,

		// collection
		task.IMPORT_STARTUP:          ImportStartup,
		task.COLLECT_INFO:            CollectInfo,
		task.GIVE_COLLECT_PROGRESS:   GiveCollectProgress,
		task.GIVE_COLLECT_DATA_INFO:  GiveCollectDataInfo,
		task.GIVE_COLLECT_DATA:       GiveCollectData,
		task.GIVE_COLLECT_DATA_END:   GiveCollectDataEnd,
		task.GIVE_COLLECT_DATA_ERROR: GiveCollectDataError,
		// scan
		// task.GET_PROCESS_INFO:            GetProcessInfo,
		task.GET_SCAN_INFO_DATA:    GetScanInfoData,
		task.PROCESS:               Process,
		task.GIVE_PROCESS_DATA:     GiveProcessData,
		task.GIVE_PROCESS_DATA_END: GiveProcessDataEnd,
		task.GIVE_SCAN_PROGRESS:    GiveScanProgress,
		task.GIVE_SCAN_DATA:        GiveScanData,
		task.GIVE_SCAN_DATA_INFO:   GiveScanDataInfo,
		task.GIVE_SCAN_DATA_OVER:   GiveScanDataOver,
		task.GIVE_SCAN_DATA_END:    GiveScanDataEnd,
	}
}
