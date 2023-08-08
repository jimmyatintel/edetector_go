package elasticquery

import (
	"edetector_go/internal/rbconnector"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/rabbitmq"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

func SendToMainElastic(uuid string, index string, agent string, item string, date int, ttype string, etc string, priority string) error {
	agentIP := query.GetMachineIP(agent)
	if agentIP == "" {
		logger.Error("Error getting machine ip")
	}
	agentName := query.GetMachineName(agent)
	if agentName == "" {
		logger.Error("Error getting machine name")
	}
	template := mainSource{
		UUID:      uuid,
		Index:     index,
		Agent:     agent,
		AgentIP:   agentIP,
		AgentName: agentName,
		Item:      item,
		Date:      date,
		Type:      ttype,
		Etc:       etc,
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
	err = rabbitmq.Publish(priority, msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func SendToDetailsElastic(uuid string, index string, agentID string, mes string, data Request_data, priority string) error {
	template, err := StringToStruct(uuid, agentID, mes, data)
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
	err = rabbitmq.Publish(priority, msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func SendToRelationElastic(template Request_data, priority string) error {
	request, err := template.Elastical()
	if err != nil {
		return err
	}
	var msg = rbconnector.Message{
		Index: "ed_de_explorer_relation",
		Data:  string(request),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = rabbitmq.Publish(priority, msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func StringToStruct(uuid string, agentID string, mes string, data Request_data) (Request_data, error) {
	agentIP := query.GetMachineIP(agentID)
	if agentIP == "" {
		logger.Error("Error getting machine ip")
	}
	agentName := query.GetMachineName(agentID)
	if agentName == "" {
		logger.Error("Error getting machine name")
	}
	v := reflect.Indirect(reflect.ValueOf(data))
	line := uuid + "@|@" + agentID + "@|@" + agentIP + "@|@" + agentName + "@|@" + mes
	values := strings.Split(line, "@|@")
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
