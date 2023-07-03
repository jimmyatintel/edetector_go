package elasticquery

import (
	"encoding/json"
	"time"
)

type Request_data interface {
	Elastical() ([]byte, error)
}

type source struct {
	Uuid string
	Time int64
	Type string
	Data Request_data
}

func (s *source) Elastical() ([]byte, error) {
	return json.Marshal(s)
}

func New_source(uuid string, Type string) source {
	return source{
		Uuid: uuid,
		Time: time.Now().Unix(),
		Type: Type,
	}
}
