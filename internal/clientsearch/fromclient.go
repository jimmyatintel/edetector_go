package clientsearch

import (
	"bytes"
	config "edetector_go/config"
	C_AES "edetector_go/internal/C_AES"
	"edetector_go/internal/task"
	"edetector_go/internal/taskservice"
	"fmt"
	"strings"

	channelmap "edetector_go/internal/channelmap"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	work "edetector_go/internal/work"
	logger "edetector_go/pkg/logger"
	mq "edetector_go/pkg/mariadb/query"
	"edetector_go/pkg/redis"
	rq "edetector_go/pkg/redis/query"
	"edetector_go/pkg/request"
	"net"
)

func handleTCPRequest(conn net.Conn, task_chan chan packet.Packet, port string) {
	firstCheckConnect := true
	logger.Info("Worker port accepted, IP: " + conn.RemoteAddr().String())
	defer conn.Close()
	buf := make([]byte, 1024)
	key := "unknown"
	agentTaskType := "unknown"
	lastTask := "unknown"
	var updateDataRightChan chan net.Conn
	var yaraDataRightChan chan net.Conn
	var imageDataRightChan chan net.Conn
	closeConn := make(chan bool)
	for {
		var NewPacket packet.Packet
		var decrypt_buf []byte
		reqLen, err := conn.Read(buf)
		if err != nil {
			connectionClosedByAgent(key, agentTaskType, lastTask, err)
			close(closeConn)
			return
		}
		if reqLen < 1024 {
			logger.Error("Invalid packet (too short): " + string(buf[:reqLen]))
			continue
		}
		agentCount, err := redis.RedisGetInt("OnlineClientCount")
		if err != nil {
			logger.Error("Error getting online client count: " + err.Error())
			continue
		}
		Data_acache := make([]byte, 0)
		Data_acache = append(Data_acache, buf[:reqLen]...)
		decrypt_buf = bytes.Repeat([]byte{0}, len(Data_acache))
		C_AES.Decryptbuffer(Data_acache, len(Data_acache), decrypt_buf)
		NewPacket = new(packet.WorkPacket)
		err = NewPacket.NewPacket(decrypt_buf, Data_acache)
		if err != nil {
			logger.Error("Error reading: " + err.Error())
			logger.Info("Content: " + fmt.Sprintf("%x", decrypt_buf))
			continue
		}
		t := NewPacket.GetTaskType()
		if t != "GiveInfo" && t != "GiveDetectInfoFirst" && t != "GiveDetectInfo" && t != "CheckConnect" && t != "ReadyUpdateAgent" && t != "ReadyYaraRule" && t != "ReadyImage" && t != "DataRight" && t != "GiveDriveInfo" {
			for len(Data_acache) < 65535 {
				reqLen, err := conn.Read(buf)
				if err != nil {
					connectionClosedByAgent(key, agentTaskType, lastTask, err)
					close(closeConn)
					return
				}
				Data_acache = append(Data_acache, buf[:reqLen]...)
			}
			NewPacket = new(packet.DataPacket)
			decrypt_buf = bytes.Repeat([]byte{0}, len(Data_acache))
			C_AES.Decryptbuffer(Data_acache, len(Data_acache), decrypt_buf)
			err = NewPacket.NewPacket(decrypt_buf, Data_acache)
			if err != nil {
				logger.Error("Error reading: " + err.Error())
				logger.Info("Content: " + fmt.Sprintf("%x", decrypt_buf))
				continue
			}
		}
		if key == "unknown" || key == "NoKey" || key == "" || key == "null" || len(key) != 32 {
			logger.Debug("Get key: " + NewPacket.GetRkey())
			key = NewPacket.GetRkey()
		}
		if agentTaskType == "unknown" {
			taskType, ok := task.TaskTypeMap[NewPacket.GetTaskType()]
			if ok {
				agentTaskType = taskType
			}
		}
		if NewPacket.GetTaskType() == "Undefine" {
			nullIndex := bytes.IndexByte(decrypt_buf[76:100], 0)
			logger.Error("Undefine TaskType: " + string(decrypt_buf[76:76+nullIndex]))
			continue
		}
		if NewPacket.GetTaskType() == task.GIVE_INFO &&
			(agentCount >= config.Viper.GetInt("AGENT_LIMIT") || len(mq.Load_all_client()) >= config.Viper.GetInt("TOTAL_AGENT_LIMIT")) {
			logger.Error("Too many clients, reject: " + string(NewPacket.GetRkey()))
			clientsearchsend.SendTCPtoClient(NewPacket, task.REJECT_AGENT, "", conn)
			close(closeConn)
			return
		} else if NewPacket.GetTaskType() == task.GIVE_DETECT_INFO_FIRST {
			rq.Online(key)
			redis.RedisSet_AddInteger("OnlineClientCount", 1)
			request.RequestToUser(key)
			logger.Info("add online client: " + key + "-" + fmt.Sprint(agentCount))
			channelmap.AssignTaskChannel(key, &task_chan)
			logger.Info("Set key-channel mapping: " + key)

			// send retry task after agent online (for task status 2, 7)
			status2TaskLists, err := mq.Load_stored_task("nil", key, 2, "nil")
			if err != nil {
				logger.Error("Get stored task failed: " + err.Error())
			}

			for _, taskInfo := range status2TaskLists {
				taskResendMap := map[string]task.TaskType{
					"StartCollect":  task.RESEND_COLLECT,
					"StartGetDrive": task.RESEND_DRIVE,
					"StartGetImage": task.RESEND_IMAGE,
					"StartYaraRule": task.RESEND_YAYA,
				}
				resendTaskType, ok := taskResendMap[taskInfo[3]]
				if !ok {
					continue
				}
				err := clientsearchsend.SendTCPtoClient(NewPacket, resendTaskType, "", conn)
				if err != nil {
					logger.Error("Error sending retry task: " + err.Error())
				} else {
					logger.Info("Retry task sent: " + string(resendTaskType))
				}
			}

			go func() {
				for {
					select {
					case message := <-task_chan:
						data := message.Fluent()
						logger.Info("Get task msg: " + string(data))
						if err := clientsearchsend.SendTaskTCPtoClient(data, conn); err != nil {
							logger.Error("Error Sending: " + err.Error())
						}
					case <-closeConn:
						return
					}
				}
			}()
		}
		if NewPacket.GetTaskType() == task.CHECK_CONNECT { // update online status only when CHECK_CONNECT
			if rq.GetStatus(key) == 0 { // offline
				if !firstCheckConnect {
					logger.Error("Agent already offline: " + string(key))
					close(closeConn)
					return
				}
				firstCheckConnect = false
			} else {
				rq.Online(key)
			}
		}
		if NewPacket.GetTaskType() == task.READY_UPDATE_AGENT {
			updateDataRightChan = make(chan net.Conn)
			_, err = work.ReadyUpdateAgent(NewPacket, conn, updateDataRightChan)
			if err != nil {
				logger.Error("Task " + string(NewPacket.GetTaskType()) + " failed: " + err.Error())
				mq.Failed_task(NewPacket.GetRkey(), agentTaskType, 6)
			}
		} else if NewPacket.GetTaskType() == task.READY_YARA_RULE {
			yaraDataRightChan = make(chan net.Conn)
			_, err = work.ReadyYaraRule(NewPacket, conn, yaraDataRightChan)
			if err != nil {
				logger.Error("Task " + string(NewPacket.GetTaskType()) + " failed: " + err.Error())
				mq.Failed_task(NewPacket.GetRkey(), agentTaskType, 6)
			}
		} else if NewPacket.GetTaskType() == task.READY_IMAGE {
			imageDataRightChan = make(chan net.Conn)
			_, err = work.ReadyImage(NewPacket, conn, imageDataRightChan)
			if err != nil {
				logger.Error("Task " + string(NewPacket.GetTaskType()) + " failed: " + err.Error())
				mq.Failed_task(NewPacket.GetRkey(), agentTaskType, 6)
			}
		} else if agentTaskType == "StartUpdate" && NewPacket.GetTaskType() == task.DATA_RIGHT {
			logger.Info("UpdateDataRight: " + NewPacket.GetRkey())
			updateDataRightChan <- conn
		} else if agentTaskType == "StartYaraRule" && NewPacket.GetTaskType() == task.DATA_RIGHT {
			logger.Info("YaraDataRight: " + NewPacket.GetRkey())
			yaraDataRightChan <- conn
		} else if agentTaskType == "StartGetImage" && NewPacket.GetTaskType() == task.DATA_RIGHT {
			logger.Info("ImageDataRight: " + NewPacket.GetRkey())
			imageDataRightChan <- conn
		} else {
			taskFunc, ok := work.WorkMap[NewPacket.GetTaskType()]
			if !ok {
				logger.Error("Function notfound: " + string(NewPacket.GetTaskType()))
				continue
			}
			_, err = taskFunc(NewPacket, conn)
			if err != nil {
				logger.Error("Task " + string(NewPacket.GetTaskType()) + " failed: " + err.Error())
				if agentTaskType != "unknown" {
					mq.Failed_task(NewPacket.GetRkey(), agentTaskType, 6)
				}
			}
		}
		lastTask = string(NewPacket.GetTaskType())
	}
}

func handleUDPRequest(addr net.Addr, buf []byte) {
	logger.Info("UDP")
}

// To-Do (TBD)
func connectionClosedByAgent(key string, agentTaskType string, lastTask string, err error) {
	logger.Warn("Connection close: " + key + "|" + agentTaskType + ", Error: " + err.Error())
	if agentTaskType == "StartScan" && lastTask == "ReadyScan" {
		logger.Error("Scan failed: " + key)
		mq.Update_task_status(key, agentTaskType, 2, 0)
	} else if agentTaskType == "StartCollect" && (lastTask == "GiveCollectProgress" || lastTask == "GiveCollectDataInfo" || lastTask == "GiveCollectData") {
		if err := taskservice.RetryTask(key, agentTaskType, task.RESEND_COLLECT); err != nil {
			logger.Error("ResendCollect failed: " + err.Error())
			mq.Failed_task(key, agentTaskType, 7)
		}
	} else if agentTaskType == "StartGetDrive" && (lastTask == "GiveExplorerProgress" || lastTask == "GiveExplorerInfo" || lastTask == "GiveExplorerData") {
		if err := taskservice.RetryTask(key, agentTaskType, task.RESEND_DRIVE); err != nil {
			logger.Error("ResendDrive failed: " + err.Error())
			mq.Failed_task(key, agentTaskType, 7)
		}
	} else if agentTaskType == "StartGetImage" && (lastTask == "GiveImageProgress" || lastTask == "GiveImageInfo" || lastTask == "GiveImage") {
		if err := taskservice.RetryTask(key, agentTaskType, task.RESEND_IMAGE); err != nil {
			logger.Error("ResendImage failed: " + err.Error())
			mq.Failed_task(key, agentTaskType, 7)
		}
	} else if agentTaskType == "StartYaraRule" && (lastTask == "GiveYaraProgress" || lastTask == "GiveRuleMatchInfo" || lastTask == "GiveRuleMatch") {
		if err := taskservice.RetryTask(key, agentTaskType, task.RESEND_YAYA); err != nil {
			logger.Error("ResendYara failed: " + err.Error())
			mq.Failed_task(key, agentTaskType, 7)
		}
	} else if agentTaskType == "Main" {
		removeTasks, err := mq.Load_stored_task("nil", key, 2, "StartRemove")
		if err != nil {
			logger.Error("Get StartRemove tasks failed: " + err.Error())
		}
		rq.Offline(key)
		if err := channelmap.RemoveTaskChannel(key); err != nil {
			logger.Error("Remove key-channel mapping failed: " + err.Error())
		}
		if len(removeTasks) != 0 {
			taskservice.DeleteAgentData(key)
			logger.Info("Finish remove agent: " + key)
		}
	} else if agentTaskType != "unknown" {
		if !strings.Contains(lastTask, "End") {
			mq.Failed_task(key, agentTaskType, 7)
		}
	}
}
