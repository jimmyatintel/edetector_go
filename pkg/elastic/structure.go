package elastic

import (
	"encoding/json"
)

type Request_data interface {
	Elastical() ([]byte, error)
}

type MainSource struct {
	UUID      string `json:"uuid"`
	Index     string `json:"index"`
	Agent     string `json:"agent"`
	AgentIP   string `json:"agentIP"`
	AgentName string `json:"agentName"`
	ItemMain  string `json:"item_main"`
	DateMain  int    `json:"date_main"`
	TypeMain  string `json:"type_main"`
	EtcMain   string `json:"etc_main"`
}

func (s *MainSource) Elastical() ([]byte, error) {
	return json.Marshal(s)
}

type FinishSignal struct {
	Agent    string `json:"agent"`
	TaskType string `json:"taskType"`
}

func (s *FinishSignal) Elastical() ([]byte, error) {
	return json.Marshal(s)
}
