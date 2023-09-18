package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"

	"go.uber.org/zap"
)

func StartGetImage(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartGetImage: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	settingType := "root:\\windows\\system32\\sru\\||srudb.dat,\\ConnectedDevicesPlatform\\*\\|LOCALAPPDATA|ActivitiesCache.db,\\Microsoft\\windows\\|APPDATA|recent"
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_IMAGE, settingType)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}
