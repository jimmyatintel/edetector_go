package clientinfo

import (
	logger "edetector_go/pkg/logger"
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strings"
)

type ClientInfo struct {
	SysInfo      string
	OsInfo       string
	ComputerName string
	UserName     string
	FileVersion  string
	BootTime     string
	KeyNum       string
}

func (c *ClientInfo) Load_data(data string) {
	var data_splited []string
	if data == "" {
		data_splited = []string{"", "", "", "", "", "", ""}
	} else {
		data_splited = strings.Split(data, "|")
	}
	c.SysInfo = data_splited[0]
	c.OsInfo = data_splited[1]
	c.ComputerName = data_splited[2]
	c.UserName = data_splited[3]
	c.FileVersion = data_splited[4]
	c.BootTime = data_splited[5]
	c.KeyNum = data_splited[6]
	if c.KeyNum == "null" {
		uuid := uuid.New()
		c.KeyNum = uuid.String()
	}
}

func (c *ClientInfo) Marshal() string {
	json, err := json.Marshal(c)
	if err != nil {
		logger.Error("Error in json marshal", zap.Any("error", err.Error()))
	}
	return string(json)
}

func UnMarshal(data string) ClientInfo {
	clientinfo := ClientInfo{}
	err := json.Unmarshal([]byte(data), &clientinfo)
	if err != nil {
		logger.Error("Error in json unmarshal", zap.Any("error", err.Error()))
	}
	return clientinfo
}
