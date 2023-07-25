package elasticquery

import (
	"encoding/json"
)

type Request_data interface {
	Elastical() ([]byte, error)
}

type mainSource struct {
	UUID  string `json:"uuid"`
	Index string `json:"index"`
	Agent string `json:"agent"`
	Item  string `json:"item"`
	Date  string `json:"date"`
	Type  string `json:"type"`
	Etc   string `json:"etc"`
}

func (s *mainSource) Elastical() ([]byte, error) {
	return json.Marshal(s)
}
