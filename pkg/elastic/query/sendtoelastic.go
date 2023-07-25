package elasticquery

import (
	"edetector_go/pkg/elastic"
	"reflect"
	"strconv"
	"strings"
)

func SendToMainElastic(uuid string, index string, agent string, item string, date string, ttype string, etc string) error {
	template := mainSource{
		UUID:  uuid,
		Index: index,
		Agent: agent,
		Item:  item,
		Date:  date,
		Type:  ttype,
		Etc:   etc,
	}
	request, err := template.Elastical()
	if err != nil {
		return err
	}
	err = elastic.IndexRequest(index, string(request))
	if err != nil {
		return err
	}
	return nil
}

func SendToDetailsElastic(uuid string, index string, agentID string, mes string, data Request_data) error {
	template, err := stringToStruct(uuid, agentID, mes, data)
	if err != nil {
		return err
	}
	request, err := template.Elastical()
	if err != nil {
		return err
	}
	err = elastic.IndexRequest(index, string(request))
	if err != nil {
		return err
	}
	return nil
}

func stringToStruct(uuid string, agentID string, mes string, data Request_data) (Request_data, error) {
	v := reflect.Indirect(reflect.ValueOf(data))
	line := uuid + "|" + agentID + "|" + mes
	values := strings.Split(line, "|")
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Int:
			value, err := strconv.Atoi(values[i])
			if err != nil {
				break
			}
			field.Set(reflect.ValueOf(value))
		case reflect.String:
			field.Set(reflect.ValueOf(values[i]))
		case reflect.Bool:
			value, err := strconv.ParseBool(values[i])
			if err != nil {
				break
			}
			field.Set(reflect.ValueOf(value))
		}
	}
	return data, nil
}
