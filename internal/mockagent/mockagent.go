package mockagent

import (
	"bytes"
	"edetector_go/config"
	"edetector_go/internal/C_AES"
	"edetector_go/internal/packet"
	"edetector_go/internal/task"
	"edetector_go/pkg/fflag"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb"
	"net"
	"os"
	"time"
)

var serverAddr string
var detectStatus string
var mockagentKey string

func init() {
	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Panic("Error loading feature flag")
		panic("Error loading feature flag")
	}
	vp, err := config.LoadConfig()
	if vp == nil {
		logger.Panic("Error loading config file: " + err.Error())
		panic(err)
	}
	if enable, err := fflag.FFLAG.FeatureEnabled("logger_enable"); enable && err == nil {
		logger.InitLogger(config.Viper.GetString("MOCK_AGENT_LOG_FILE"), "server", "SERVER")
		logger.Info("Logger is enabled please check all out info in log file: " + config.Viper.GetString("WORKER_LOG_FILE"))
	}
	connString, err := mariadb.Connect_init()
	if err != nil {
		logger.Panic("Error connecting to mariadb: " + err.Error())
		panic(err)
	} else {
		logger.Info("Mariadb connectionString: " + connString)
	}
    if len(os.Args) != 3 {
		logger.Panic("Usage: go run mockagent/agent.go 1(agentID) 163(serverIP)...")
		panic(err)
    }
    mockagentKey = "mockagent" + os.Args[1]
	serverAddr = config.Viper.GetString("WORKING_SERVER_IP") + os.Args[2] + ":" + config.Viper.GetString("WORKER_DEFAULT_WORKER_PORT")
	detectStatus = "0|0"
}

func Main() {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		logger.Panic("Error connecting to the server:" + err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	logger.Info("Connected to the server at " + serverAddr)

	// handshake
	timestamp := time.Now().Format("20060102150405")
	info := "x64|MockAgent|MockAgent|SYSTEM|1.0.0|" + timestamp + "|" + mockagentKey
	SendTCPtoServer(task.GIVE_INFO, info, conn)
	go handleMainConn(conn)
	for {
		SendTCPtoServer(task.CHECK_CONNECT, "", conn)
		time.Sleep(30 * time.Second)
	}
}

func receive(buf []byte, conn net.Conn) packet.Packet {
	var decrypt_buf []byte
	var NewPacket packet.Packet
	reqLen, err := conn.Read(buf)
	if err != nil {
		logger.Error("Error reading from server: " + err.Error())
		return nil
	}
	if reqLen == 1024 {
		NewPacket = new(packet.WorkPacket)
	} else {
		logger.Error("Invalid packet (too short): " + string(decrypt_buf))
		return nil
	}
	decrypt_buf = bytes.Repeat([]byte{0}, len(buf))
	C_AES.Decryptbuffer(buf, len(buf), decrypt_buf)
	err = NewPacket.NewPacket(decrypt_buf, buf)
	if err != nil {
		logger.Error("Error reading: " + err.Error() + ", Content: " + string(decrypt_buf))
		return nil
	}
	return NewPacket
}

func handleMainConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		NewPacket := receive(buf, conn)
		logger.Info("Received task from server: " + string(NewPacket.GetTaskType()))
		if NewPacket.GetTaskType() == task.OPEN_CHECK_THREAD {
			SendTCPtoServer(task.GIVE_DETECT_INFO_FIRST, detectStatus, conn)
			new_conn, err := net.Dial("tcp", serverAddr)
			if err != nil {
				logger.Error("Error connecting to the server:" + err.Error())
				return
			}
			defer new_conn.Close()
			go agentDetect(new_conn)
		} else if NewPacket.GetTaskType() == task.UPDATE_DETECT_MODE {
			detectStatus = NewPacket.GetMessage()
			SendTCPtoServer(task.GIVE_DETECT_INFO, detectStatus, conn)
		} else if NewPacket.GetTaskType() == task.GET_DRIVE {
			SendTCPtoServer(task.GIVE_DRIVE_INFO, "C-NTFS,HDD|", conn)
		} else if NewPacket.GetTaskType() == task.GET_SCAN || NewPacket.GetTaskType() == task.GET_COLLECT_INFO || NewPacket.GetTaskType() == task.EXPLORER_INFO {
			go handleNewTask(NewPacket.GetTaskType())
		} else if NewPacket.GetTaskType() != task.DATA_RIGHT {
			logger.Error("Undefined task type: " + string(NewPacket.GetTaskType()))
		}
	}
}

func handleNewTask(taskType task.TaskType) {
	new_conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		logger.Error("Error connecting to the server:" + err.Error())
		return
	}
	defer new_conn.Close()
	dataRightFromServer := make(chan int)
	switch taskType {
	case task.GET_SCAN:
		go agentScan(new_conn, dataRightFromServer)
	case task.GET_COLLECT_INFO:
		go agentCollect(new_conn, dataRightFromServer)
	case task.EXPLORER_INFO:
		go agentDrive(new_conn, dataRightFromServer)
	}
	buf := make([]byte, 1024)
	for {
		NewPacket := receive(buf, new_conn)
		if NewPacket.GetTaskType() == task.DATA_RIGHT {
			dataRightFromServer <- 1
		}
	}
}
