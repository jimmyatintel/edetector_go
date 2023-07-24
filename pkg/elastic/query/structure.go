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
	Item  string
	Date  string
	Type  string
	Etc   string
}

func (s *mainSource) Elastical() ([]byte, error) {
	return json.Marshal(s)
}

func New_main(uuid string, index string, item string, date string, ttype string, etc string) mainSource {
	return mainSource {
		UUID:  uuid,
		Index: index,
		Item:  item,
		Date:  date,
		Type:  ttype,
		Etc:   etc,
	}
}

// type source struct {
// 	Uuid string
// 	Time int64
// 	Type string
// 	Data Request_data
// }

// func (s *source) Elastical() ([]byte, error) {
// 	return json.Marshal(s)
// }

// func New_source(uuid string, Type string) source {
// 	return source{
// 		Uuid: uuid,
// 		Time: time.Now().Unix(),
// 		Type: Type,
// 	}
// }
