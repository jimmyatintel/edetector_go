package workfromapi

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
)

func StartGetImage(p packet.UserPacket) (task.TaskResult, error) {
	logger.Info("StartGetImage: " + p.GetRkey() + "::" + p.GetMessage())
	// imageList, err := query.Load_key_image(p.GetMessage())
	// if err != nil {
	// 	return task.FAIL, err
	// }
	imageStr := "root:\\windows\\system32\\sru\\||srudb.dat,\\ConnectedDevicesPlatform\\*\\|LOCALAPPDATA|ActivitiesCache.db,root:\\Users\\*\\||NTUSER.DAT"  
	// for _, image := range imageList {
	// 	imageStr = imageStr + image[1] + "|" + image[0] + "|" + image[2] + ","
	// }
	err := clientsearchsend.SendUserTCPtoClient(p, task.GET_IMAGE, imageStr)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}
