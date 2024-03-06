package clientinfo

import (
	logger "edetector_go/pkg/logger"
	"encoding/json"
	"errors"
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

func (c *ClientInfo) Load_data(data string) error {
	data_splited := strings.Split(data, "|@|")
	if len(data_splited) < 7 {
		return errors.New("error in GiveInfo format, version conflicted")
	}
	c.SysInfo = data_splited[0]
	c.OsInfo = data_splited[1]
	c.ComputerName = data_splited[2]
	c.UserName = data_splited[3]
	c.FileVersion = data_splited[4]
	c.BootTime = data_splited[5]
	c.KeyNum = data_splited[6]
	return nil
}

func (c *ClientInfo) Marshal() string {
	json, err := json.Marshal(c)
	if err != nil {
		logger.Error("Error in json marshal: " + err.Error())
	}
	return string(json)
}

func UnMarshal(data string) ClientInfo {
	clientinfo := ClientInfo{}
	err := json.Unmarshal([]byte(data), &clientinfo)
	if err != nil {
		logger.Error("Error in json unmarshal: " + err.Error())
	}
	return clientinfo
}
