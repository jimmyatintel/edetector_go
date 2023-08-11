package memory

import "encoding/json"

type Memory struct {
	UUID              string `json:"uuid"`
	Agent             string `json:"agent"`
	AgentIP            string `json:"agentIP"`
	AgentName          string `json:"agentName"`
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
	Hook              string `json:"hook"`
	ImportOtherDLL    bool   `json:"importOtherDLL"`
	ProcessConnectIP  string `json:"processConnectIP"`
	RiskLevel         int    `json:"riskLevel"`
	Mode              string `json:"mode"`
}

func (n Memory) Elastical() ([]byte, error) {
	return json.Marshal(n)
}