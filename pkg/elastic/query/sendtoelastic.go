package elasticquery

import (
	"edetector_go/internal/rbconnector"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/rabbitmq"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func SendToMainElastic(uuid string, index string, agent string, item string, date int, ttype string, etc string) error {
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
	var msg = rbconnector.Message{
		Index: "ed_main",
		Data:  string(request),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	logger.Info("ed_main", zap.Any("message", string(msgBytes)))
	if index == "ed_memory" {
		err = rabbitmq.Publish("ed_mid", msgBytes)
	} else {
		err = rabbitmq.Publish("ed_mid", msgBytes)
	}
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
	var msg = rbconnector.Message{
		Index: index,
		Data:  string(request),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	logger.Info("ed_memory", zap.Any("message", string(msgBytes)))
	if index == "ed_memory" {
		err = rabbitmq.Publish("ed_mid", msgBytes)
	} else {
		err = rabbitmq.Publish("ed_mid", msgBytes)
	}
	if err != nil {
		return err
	}
	return nil
}

func stringToStruct(uuid string, agentID string, mes string, data Request_data) (Request_data, error) {
	v := reflect.Indirect(reflect.ValueOf(data))
	line := uuid + "||" + agentID + "||" + mes
	values := strings.Split(line, "||")
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
