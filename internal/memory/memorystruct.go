package memory

import "encoding/json"

type Memory struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	AgentIP           string `json:"agentIP"`
	AgentName         string `json:"agentName"`
	ProcessName       string `json:"processName"`
	ProcessCreateTime int    `json:"processCreateTime"`
	DynamicCommand    string `json:"dynamicCommand"`
	ProcessMD5        string `json:"processMD5"`
	ProcessPath       string `json:"processPath"`
	ParentProcessId   int    `json:"parentProcessId"`
	ParentProcessName string `json:"parentProcessName"`
	ParentProcessPath string `json:"parentProcessPath"`
	DigitalSign       string `json:"digitalSign"`
	ProcessId         int    `json:"processId"`
	InjectActive      string `json:"injectActive"`
	ProcessBeInjected int    `json:"processBeInjected"`
	Boot              string `json:"boot"`
	Hide              string `json:"hide"`
	ImportOtherDLL    string `json:"importOtherDLL"`
	Hook              string `json:"hook"`
	ProcessConnectIP  string `json:"processConnectIP"`
	RiskLevel         int    `json:"riskLevel"`
	Mode              string `json:"mode"`
	ItemMain          string `json:"item_main"`
	DateMain          int    `json:"date_main"`
	TypeMain          string `json:"type_main"`
	EtcMain           string `json:"etc_main"`
}

func (n Memory) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type MemoryNetworkDetect struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	AgentIP           string `json:"agentIP"`
	AgentName         string `json:"agentName"`
	ProcessId         int    `json:"processId"`
	Address           string `json:"address"`
	Timestamp         int    `json:"timestamp"`
	ProcessCreateTime int    `json:"processCreateTime"`
	ConnectionINorOUT bool   `json:"connectionInOrOut"`
	AgentPort         int    `json:"agentPort"`
	ItemMain          string `json:"item_main"`
	DateMain          int    `json:"date_main"`
	TypeMain          string `json:"type_main"`
	EtcMain           string `json:"etc_main"`
}

func (n MemoryNetworkDetect) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type MemoryNetworkScan struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	AgentIP           string `json:"agentIP"`
	AgentName         string `json:"agentName"`
	ProcessId         int    `json:"processId"`
	ProcessCreateTime int    `json:"processCreateTime"`
	SrcAddress        string `json:"srcAddress"`
	DstAddress        string `json:"dstAddress"`
	Action            string `json:"action"`
	Timestamp         int    `json:"timestamp"`
	ItemMain          string `json:"item_main"`
	DateMain          int    `json:"date_main"`
	TypeMain          string `json:"type_main"`
	EtcMain           string `json:"etc_main"`
}

func (n MemoryNetworkScan) Elastical() ([]byte, error) {
	return json.Marshal(n)
}
