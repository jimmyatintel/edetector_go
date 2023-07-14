package taskservice

import (
	"bytes"
	"edetector_go/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Request struct {
	DeviceId string `json:"deviceId"`
}

func RequestToUser(id string) {
	request := Request{
		DeviceId: id,
	}
	// Marshal payload into JSON
	payload, err := json.Marshal(request)
	if err != nil {
		logger.Error("Error marshaling JSON:", zap.Any("error", err.Error()))
		return
	}
	// Create an HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://192.168.200.161:5000/agent", bytes.NewBuffer(payload))
	if err != nil {
		logger.Error("Error creating HTTP request: ", zap.Any("error", err.Error()))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the HTTP request
	response, err := client.Do(req)
	if err != nil {
		logger.Error("Error sending HTTP request: ", zap.Any("error", err.Error()))
		return
	}
	defer response.Body.Close()
	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status code:", response.StatusCode)
		return
	}
}
