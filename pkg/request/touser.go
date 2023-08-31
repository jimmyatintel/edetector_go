package request

import (
	"bytes"
	"context"
	"edetector_go/config"
	"edetector_go/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Create an HTTP request
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	ip := config.Viper.GetString("WS_HOST")
	port := config.Viper.GetString("WS_PORT")
	path := fmt.Sprintf("http://%s:%s/updateTask", ip, port)
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(payload))
	if err != nil {
		logger.Error("Error creating HTTP request: ", zap.Any("error", err.Error()))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the HTTP request
	response, err := client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				logger.Error("Request timed out: ", zap.Any("error", err.Error()))
				return
			}
		default:
			logger.Error("Error sending HTTP request: ", zap.Any("error", err.Error()))
			return
		}
	}
	defer response.Body.Close()
	// Check the response status code
	if response.StatusCode != http.StatusOK {
		logger.Error("Request failed with status code:", zap.Any("error", response.StatusCode))
		return
	}
}
