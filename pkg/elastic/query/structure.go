package elasticquery

import (
	"encoding/json"
)

type Request_data interface {
	Elastical() ([]byte, error)
}

type mainSource struct {
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

func (s *mainSource) Elastical() ([]byte, error) {
	return json.Marshal(s)
}
