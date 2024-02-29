package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
)

func StartGetImage(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartGetImage: " + p.GetRkey() + "::" + p.GetMessage())
	// imageStr := "root:\\windows\\system32\\sru\\||srudb.dat,\\ConnectedDevicesPlatform\\*\\|LOCALAPPDATA|ActivitiesCache.db,root:\\Users\\*\\||NTUSER.DAT"
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_IMAGE, p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}
