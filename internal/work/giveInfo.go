package work

import (
	"edetector_go/config"
	client "edetector_go/internal/client"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	mq "edetector_go/pkg/mariadb/query"
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func GiveInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) { // the first packet: insert user info
	logger.Info("GiveInfo: " + p.GetRkey() + "::" + p.GetMessage())
	logger.Info("ip: " + p.GetipAddress() + ", mac: " + p.GetMacAddress())
	np := packet.CheckIsWork(p)
	ClientInfo, err := client.PacketClientInfo(np)
	if err != nil {
		clientsearchsend.SendTCPtoClient(p, task.REJECT_AGENT, "", conn)
		return task.FAIL, err
	}
	info := strings.Split(ClientInfo.FileVersion, ",")
	if len(info) != 3 {
		logger.Error("Error FileVersion format: " + ClientInfo.FileVersion)
		clientsearchsend.SendTCPtoClient(p, task.REJECT_AGENT, "", conn)
		return task.FAIL, err
	}
	minVersion := config.Viper.GetString("MIN_AGENT_VERSION")
	if checkInvalidVersion(info[0], minVersion) {
		logger.Error("Version Conflict: " + ClientInfo.FileVersion)
		clientsearchsend.SendTCPtoClient(p, task.REJECT_AGENT, "", conn)
		return task.FAIL, err
	}
	if (ClientInfo.KeyNum == "") || (ClientInfo.KeyNum == "null") || (ClientInfo.KeyNum == "NoKey") { // assign a new key(uuid)
		ClientInfo.KeyNum = strings.Replace(uuid.New().String(), "-", "", -1)
		logger.Info("New key: " + ClientInfo.KeyNum)
	}
	mq.Checkindex(ClientInfo.KeyNum, p.GetipAddress(), p.GetMacAddress())
	mq.Addmachine(ClientInfo)
	err = clientsearchsend.SendTCPtoClient(p, task.OPEN_CHECK_THREAD, ClientInfo.KeyNum, conn) // send ack(OPEN_CHECK_THREAD) to client
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func checkInvalidVersion(version string, minVersion string) bool {
	minVersions := strings.Split(minVersion, ".")
	versions := strings.Split(version, ".")
	if len(minVersions) != 3 || len(versions) != 3 {
		logger.Error("Invalid version format: " + version + " " + minVersion)
		return true
	}
	for i := 0; i < 3; i++ {
		v, err := strconv.Atoi(versions[i])
		if err != nil {
			logger.Error("Invalid version format: " + version)
			return true
		}
		mv, err := strconv.Atoi(minVersions[i])
		if err != nil {
			logger.Error("Invalid version format: " + minVersion)
			return true
		}
		if v < mv {
			return true
		} else if v > mv {
			return false
		}
	}
	return false
}
