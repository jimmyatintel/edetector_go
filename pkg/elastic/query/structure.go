package elasticquery

import (
	"encoding/json"
)

var elasticPrefix = "ed_"

type Request_data interface {
	Elastical() ([]byte, error)
}

type mainSource struct {
	UUID      string `json:"uuid"`
	Index     string `json:"index"`
	Agent     string `json:"agent"`
	AgentIP   string `json:"agentIP"`
	AgentName string `json:"agentName"`
	Item      string `json:"item"`
	Date      int    `json:"date"`
	Type      string `json:"type"`
	Etc       string `json:"etc"`
}

func (s *mainSource) Elastical() ([]byte, error) {
	return json.Marshal(s)
}
