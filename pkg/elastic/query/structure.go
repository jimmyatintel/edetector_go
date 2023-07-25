package elasticquery

import (
	"encoding/json"
)

type Request_data interface {
	Elastical() ([]byte, error)
}

type mainSource struct {
	UUID  string `json:"UUID"`
	Index string `json:"Index"`
	Agent string `json:"Agent"`
	Item  string `json:"Item"`
	Date  string `json:"Date"`
	Type  string `json:"Type"`
	Etc   string `json:"Etc"`
}

func (s *mainSource) Elastical() ([]byte, error) {
	return json.Marshal(s)
}
