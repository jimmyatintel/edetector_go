package work

import "encoding/json"

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
	Malicious         int    `json:"malicious"`
	VirusTotal        int    `json:"virusTotal"`
}

type Collect_Memory struct {
	Memory        Memory `json:"memory"`
	UUID          string `json:"uuid"`
	Agent         string `json:"agent"`
	AgentIP       string `json:"agentIP"`
	AgentName     string `json:"agentName"`
	ItemMain      string `json:"item_main"`
	DateMain      int    `json:"date_main"`
	TypeMain      string `json:"type_main"`
	EtcMain       string `json:"etc_main"`
	Task_id       string `json:"task_id"`
	Category      string `json:"category"`
	TaskTimestamp int    `json:"task_timestamp"`
}

func (n Memory) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

func (n Collect_Memory) Elastical() ([]byte, error) {
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
	Malicious         int    `json:"malicious"`
	VirusTotal        int    `json:"virusTotal"`
}

type Collect_MemoryNetwork struct {
	MemoryNetwork MemoryNetwork `json:"memory_network"`
	UUID          string        `json:"uuid"`
	Agent         string        `json:"agent"`
	AgentIP       string        `json:"agentIP"`
	AgentName     string        `json:"agentName"`
	ItemMain      string        `json:"item_main"`
	DateMain      int           `json:"date_main"`
	TypeMain      string        `json:"type_main"`
	EtcMain       string        `json:"etc_main"`
	Task_id       string        `json:"task_id"`
	Category      string        `json:"category"`
	TaskTimestamp int           `json:"task_timestamp"`
}

func (n MemoryNetwork) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

func (n Collect_MemoryNetwork) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

type MemoryTree struct {
	ProcessId               int      `json:"processId"`
	ParentProcessId         int      `json:"parentProcessId"`
	ProcessName             string   `json:"processName"`
	ProcessCreateTime       int      `json:"processCreateTime"`
	ParentProcessName       string   `json:"parentProcessName"`
	ParentProcessCreateTime int      `json:"parentProcessCreateTime"`
	ProcessPath             string   `json:"processPath"`
	UserName                string   `json:"userName"`
	IsPacked                bool     `json:"isPacked"`
	DynamicCommand          string   `json:"dynamicCommand"`
	IsHide                  bool     `json:"isHide"`
	IsRoot                  bool     `json:"isRoot"`
	Child                   []string `json:"child"`
}

type Collect_MemoryTree struct {
	MemoryTree    MemoryTree `json:"memory_tree"`
	UUID          string     `json:"uuid"`
	Agent         string     `json:"agent"`
	AgentIP       string     `json:"agentIP"`
	AgentName     string     `json:"agentName"`
	ItemMain      string     `json:"item_main"`
	DateMain      int        `json:"date_main"`
	TypeMain      string     `json:"type_main"`
	EtcMain       string     `json:"etc_main"`
	Task_id       string     `json:"task_id"`
	Category      string     `json:"category"`
	TaskTimestamp int        `json:"task_timestamp"`
}

func (n MemoryTree) Elastical() ([]byte, error) {
	return json.Marshal(n)
}

func (n Collect_MemoryTree) Elastical() ([]byte, error) {
	return json.Marshal(n)
}
