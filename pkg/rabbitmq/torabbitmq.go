package rabbitmq

import (
	"edetector_go/config"
	"edetector_go/pkg/elastic"
	"edetector_go/pkg/logger"
	"encoding/json"
	"reflect"
	"strconv"
)

type Message struct {
	Index string `json:"index"`
	Data  string `json:"data"`
}

func ToRabbitMQ_Main(index string, uuid string, agentID string, ip string, name string, item string, date string, ttype string, etc string, priority string) error {
	date_int, err := strconv.Atoi(date)
	if err != nil {
		logger.Error("Error converting time: " + err.Error())
		date_int = 0
	}
	template := elastic.MainSource{
		UUID:      uuid,
		Index:     index,
		Agent:     agentID,
		AgentIP:   ip,
		AgentName: name,
		ItemMain:  item,
		DateMain:  date_int,
		TypeMain:  ttype,
		EtcMain:   etc,
	}
	request, err := template.Elastical()
	if err != nil {
		return err
	}
	var msg = Message{
		Index: config.Viper.GetString("ELASTIC_PREFIX") + "_main",
		Data:  string(request),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = Publish(priority, msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func ToRabbitMQ_Details(index string, st elastic.Request_data, values []string, uuid string, agentID string, ip string, name string, item string, date string, ttype string, etc string, priority string) error {
	template, err := StringToStruct(st, values, uuid, agentID, ip, name, item, date, ttype, etc)
	if err != nil {
		return err
	}
	request, err := template.Elastical()
	if err != nil {
		return err
	}
	var msg = Message{
		Index: index,
		Data:  string(request),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = Publish(priority, msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func ToRabbitMQ_Relation(template elastic.Request_data, priority string) error {
	request, err := template.Elastical()
	if err != nil {
		return err
	}
	var msg = Message{
		Index: config.Viper.GetString("ELASTIC_PREFIX") + "_explorer_relation",
		Data:  string(request),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = Publish(priority, msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func StringToStruct(st elastic.Request_data, values []string, uuid string, agentID string, ip string, name string, item string, date string, ttype string, etc string) (elastic.Request_data, error) {
	v := reflect.Indirect(reflect.ValueOf(st))
	values = append(values, uuid, agentID, ip, name, item, date, ttype, etc)
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
	return st, nil
}
