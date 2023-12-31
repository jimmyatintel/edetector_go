package work

import "encoding/json"

type MemoryRelation struct {
	Agent  string   `json:"agent"`
	IsRoot bool     `json:"isRoot"`
	Parent string   `json:"parent"`
	Child  []string `json:"child"`
}

func (n MemoryRelation) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type Memory struct {
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
	RiskScore         int    `json:"riskScore"`
	Mode              string `json:"mode"`
	ProcessKey        string `json:"processKey"`
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	AgentIP           string `json:"agentIP"`
	AgentName         string `json:"agentName"`
	ItemMain          string `json:"item_main"`
	DateMain          int    `json:"date_main"`
	TypeMain          string `json:"type_main"`
	EtcMain           string `json:"etc_main"`
}

func (n Memory) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type MemoryNetwork struct {
	ProcessId         int    `json:"processId"`
	ProcessCreateTime int    `json:"processCreateTime"`
	Timestamp         int    `json:"timestamp"`
	SrcAddress        string `json:"srcAddress"`
	SrcPort           int    `json:"srcPort"`
	DstAddress        string `json:"dstAddress"`
	DstPort           int    `json:"dstPort"`
	Action            string `json:"action"`
	ConnectionINorOUT string `json:"connectionInOrOut"`
	Mode              string `json:"mode"`
	AgentPort         int    `json:"agentPort"`
	AgentCountry      string `json:"agentCountry"`
	AgentLongitude    int    `json:"agentLongitude"`
	AgentLatitude     int    `json:"agentLatitude"`
	OtherIP           string `json:"otherIP"`
	OtherPort         int    `json:"otherPort"`
	OtherCountry      string `json:"otherCountry"`
	OtherLongitude    int    `json:"otherLongitude"`
	OtehrLatitude     int    `json:"otherLatitude"`
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	AgentIP           string `json:"agentIP"`
	AgentName         string `json:"agentName"`
	ItemMain          string `json:"item_main"`
	DateMain          int    `json:"date_main"`
	TypeMain          string `json:"type_main"`
	EtcMain           string `json:"etc_main"`
}

func (n MemoryNetwork) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

// type MemoryNetworkScan struct {
// 	ProcessId         int    `json:"processId"`
// 	ProcessCreateTime int    `json:"processCreateTime"`
// 	SrcAddress        string `json:"srcAddress"`
// 	DstAddress        string `json:"dstAddress"`
// 	Action            string `json:"action"`
// 	Timestamp         int    `json:"timestamp"`
// 	UUID              string `json:"uuid"`
// 	Agent             string `json:"agent"`
// 	AgentIP           string `json:"agentIP"`
// 	AgentName         string `json:"agentName"`
// 	ItemMain          string `json:"item_main"`
// 	DateMain          int    `json:"date_main"`
// 	TypeMain          string `json:"type_main"`
// 	EtcMain           string `json:"etc_main"`
// }

// func (n MemoryNetworkScan) Elastical() ([]byte, error) {
// 	return json.Marshal(n)
// }
