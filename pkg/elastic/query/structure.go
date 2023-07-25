package elasticquery

import (
	"encoding/json"
)

type Request_data interface {
	Elastical() ([]byte, error)
}

type mainSource struct {
	UUID  string
	Index string
	Agent string
	Item  string
	Date  string
	Type  string
	Etc   string
}

func (s *mainSource) Elastical() ([]byte, error) {
	return json.Marshal(s)
}