package work

import (
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"math"
	"net"
	"strconv"
)

func ReadyImage(p packet.Packet, conn net.Conn, dataRight chan net.Conn) (task.TaskResult, error) {
	logger.Info("ReadyImage: " + p.GetRkey() + "::" + p.GetMessage())
	imageList, err := query.Load_key_image(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	imageStr := ""
	for _, image := range imageList {
		imageStr = imageStr + image[1] + "|" + image[0] + "|" + image[2] + "\n"
	}
	logger.Info("Image Path: " + p.GetRkey() + "::" + imageStr)
	fileLen := len(imageStr)
	logger.Info("ServerSend GiveImagePathInfo: " + p.GetRkey() + "::" + strconv.Itoa(fileLen))
	err = clientsearchsend.SendTCPtoClient(p, task.GIVE_IMAGE_PATH_INFO, strconv.Itoa(fileLen), conn)
	if err != nil {
		return task.FAIL, err
	}
	go GiveImagePath(p, len(imageStr), imageStr, dataRight)
	return task.SUCCESS, nil
}

func GiveImagePath(p packet.Packet, fileLen int, content string, dataRight chan net.Conn) {
	start := 0
	for {
		conn := <-dataRight
		if start >= fileLen {
			logger.Info("ServerSend GiveImagePathEnd: " + p.GetRkey())
			err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_IMAGE_PATH_END, []byte{}, conn)
			if err != nil {
				logger.Error("Send GiveImagePathEnd error: " + err.Error())
				query.Failed_task(p.GetRkey(), "StartGetImage", 6)
				return
			}
			<-dataRight
			break
		}
		end := int(math.Min(float64(fileLen), float64(start+65436)))
		data := []byte(content[start:end])
		logger.Info("ServerSend GiveImagePath: " + p.GetRkey())
		err := clientsearchsend.SendDataTCPtoClient(p, task.GIVE_IMAGE_PATH, data, conn)
		if err != nil {
			logger.Error("Send GiveImagePath error: " + err.Error())
			query.Failed_task(p.GetRkey(), "StartGetImage", 6)
			return
		}
		start += 65436
	}
}
