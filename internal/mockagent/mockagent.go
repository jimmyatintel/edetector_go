package mockagent

import (
	"bytes"
	"edetector_go/config"
	"edetector_go/internal/C_AES"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var serverAddr string
var numberOfAgents int

func init() {
	vp, err := config.LoadConfig()
	if vp == nil {
		logger.Panic("Error loading config file: " + err.Error())
		panic(err)
	}
	if len(os.Args) != 3 {
		err := errors.New("usage: go run mockagent/agent.go 192.168.200.163(serverIP) 10(number of agents)")
		logger.Panic(err.Error())
		panic(err)
	}
	serverAddr = os.Args[1] + ":" + config.Viper.GetString("WORKER_DEFAULT_WORKER_PORT")
	numberOfAgents, err = strconv.Atoi(os.Args[2])
	if err != nil {
		logger.Panic("Error converting number of agents: " + err.Error())
	}
	if true {
		logger.InitLogger(config.Viper.GetString("MOCK_AGENT_LOG_FILE"), "server", "SERVER")
		logger.Info("Logger is enabled please check all out info in log file: " + config.Viper.GetString("MOCK_AGENT_LOG_FILE"))
	}
	connString, err := mariadb.Connect_init()
	if err != nil {
		logger.Panic("Error connecting to mariadb: " + err.Error())
		panic(err)
	} else {
		logger.Info("Mariadb connectionString: " + connString)
	}
}

func Main() {
	var wg sync.WaitGroup
	for i := 0; i < numberOfAgents; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createAgent()
		}()
	}
	wg.Wait()
}

func createAgent() {
	// mock data
	rand.New(rand.NewSource(time.Now().UnixNano()))
	key := strings.Replace(uuid.New().String(), "-", "", -1)
	ip := "192.168.100." + fmt.Sprint(rand.Intn(101-1))
	macBytes := []byte{byte(0x02), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))}
	mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", macBytes[0], macBytes[1], macBytes[2], macBytes[3], macBytes[4], macBytes[5])
	logger.Info("MockAgentData: " + key + "|" + ip + "|" + mac)
	detectStatus := "0|0"

	for {
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			logger.Error(key + ":: Error connecting to the server:" + err.Error())
			time.Sleep(30 * time.Second)
			continue
		}
		logger.Info(key + ":: Connected to the server at " + serverAddr)
		// handshake
		timestamp := time.Now().Format("20060102150405")
		info := "x64|MockAgent|MockAgent|SYSTEM|1.0.0|" + timestamp + "|" + key
		SendTCPtoServer(task.GIVE_INFO, info, conn, ip, mac, key)
		err = handleMainConn(conn, &detectStatus, ip, mac, key)
		if err != nil {
			logger.Error(key + ":: Error handling main connection: " + err.Error())
			time.Sleep(30 * time.Second)
			continue
		}
	}
}

func receive(buf []byte, conn net.Conn, key string) packet.Packet {
	var decrypt_buf []byte
	var NewPacket packet.Packet
	reqLen, err := conn.Read(buf)
	if err != nil {
		logger.Error(key + ":: Error reading from server: " + err.Error())
		return nil
	}
	if reqLen == 1024 {
		NewPacket = new(packet.WorkPacket)
	} else {
		logger.Error(key + ":: Invalid packet (too short): " + string(decrypt_buf))
		return nil
	}
	decrypt_buf = bytes.Repeat([]byte{0}, len(buf))
	C_AES.Decryptbuffer(buf, len(buf), decrypt_buf)
	err = NewPacket.NewPacket(decrypt_buf, buf)
	if err != nil {
		logger.Error(key + ":: Error reading: " + err.Error() + ", Content: " + string(decrypt_buf))
		return nil
	}
	return NewPacket
}

func handleMainConn(conn net.Conn, detectStatus *string, ip string, mac string, key string) error {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		NewPacket := receive(buf, conn, key)
		if NewPacket == nil {
			return errors.New("Main connection closed")
		}
		logger.Info(key + ":: Received task from server: " + string(NewPacket.GetTaskType()))
		if NewPacket.GetTaskType() == task.OPEN_CHECK_THREAD {
			SendTCPtoServer(task.GIVE_DETECT_INFO_FIRST, *detectStatus, conn, ip, mac, key)
			new_conn, err := net.Dial("tcp", serverAddr)
			if err != nil {
				logger.Error(key + ":: Error connecting to the server:" + err.Error())
				return err
			}
			defer new_conn.Close()
			go agentDetect(new_conn, detectStatus, ip, mac, key)
			go func() {
				for {
					SendTCPtoServer(task.CHECK_CONNECT, "", conn, ip, mac, key)
					time.Sleep(30 * time.Second)
				}
			}()
		} else if NewPacket.GetTaskType() == task.UPDATE_DETECT_MODE {
			*detectStatus = NewPacket.GetMessage()
			SendTCPtoServer(task.GIVE_DETECT_INFO, *detectStatus, conn, ip, mac, key)
		} else if NewPacket.GetTaskType() == task.GET_DRIVE {
			SendTCPtoServer(task.GIVE_DRIVE_INFO, "C-NTFS,HDD|", conn, ip, mac, key)
		} else if NewPacket.GetTaskType() == task.GET_SCAN || NewPacket.GetTaskType() == task.GET_COLLECT_INFO || NewPacket.GetTaskType() == task.EXPLORER_INFO {
			go handleNewTask(NewPacket.GetTaskType(), ip, mac, key)
		} else if NewPacket.GetTaskType() != task.DATA_RIGHT {
			logger.Error(key + ":: Undefined task type: " + string(NewPacket.GetTaskType()))
		}
	}
}

func handleNewTask(taskType task.TaskType, ip string, mac string, key string) {
	new_conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		logger.Error(key + ":: Error connecting to the server:" + err.Error())
		return
	}
	defer new_conn.Close()
	dataRightFromServer := make(chan int)
	switch taskType {
	case task.GET_SCAN:
		go agentScan(new_conn, dataRightFromServer, ip, mac, key)
	case task.GET_COLLECT_INFO:
		go agentCollect(new_conn, dataRightFromServer, ip, mac, key)
	case task.EXPLORER_INFO:
		go agentDrive(new_conn, dataRightFromServer, ip, mac, key)
	}
	buf := make([]byte, 1024)
	for {
		NewPacket := receive(buf, new_conn, key)
		if NewPacket == nil {
			break
		} else if NewPacket.GetTaskType() == task.DATA_RIGHT {
			dataRightFromServer <- 1
		}
	}
}
