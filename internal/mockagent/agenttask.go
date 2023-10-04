package mockagent

import (
	"edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"math"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func agentDetect(conn net.Conn, detectStatus *string, info []string) {
	for {
		if *detectStatus == "1|1" {
			SendTCPtoServer(task.GIVE_DETECT_NETWORK, "10924|192.168.200.109:56421|1695608694|0|1|7680\n11032|192.168.200.167:8080|1695608694|1695018229|0|59304\n", conn, info)
			time.Sleep(5 * time.Second)
			SendTCPtoServer(task.GIVE_DETECT_PROCESS, "13.exe|1694863834|\"fontdrvhost.exe\"|rrr|C:\\Windows\\application|1160|wininit.exe|null|MS|1496|0,0|0|0,0|0,0|null|null", conn, info)
		}
		time.Sleep(10 * time.Second)
	}
}

func agentScan(conn net.Conn, dataRight chan int, info []string) {
	SendTCPtoServer(task.GIVE_SCAN_INFO, "300", conn, info)
	<-dataRight
	for i := 1; i <= 300; i++ {
		SendTCPtoServer(task.GIVE_SCAN_PROGRESS, strconv.Itoa(i)+"/300", conn, info)
		<-dataRight
	}
	sendZipFile("scan.zip", task.GIVE_SCAN_DATA_INFO, task.GIVE_SCAN, task.GIVE_SCAN_END, conn, dataRight, info)
}

func agentCollect(conn net.Conn, dataRight chan int, info []string) {
	for i := 1; i <= 48; i++ {
		SendTCPtoServer(task.GIVE_COLLECT_PROGRESS, strconv.Itoa(i)+"/48", conn, info)
		<-dataRight
	}
	sendZipFile("collect.zip", task.GIVE_COLLECT_DATA_INFO, task.GIVE_COLLECT_DATA, task.GIVE_COLLECT_DATA_END, conn, dataRight, info)
}

func agentDrive(conn net.Conn, dataRight chan int, info []string) {
	SendTCPtoServer(task.EXPLORER, "C|NTFS", conn, info)
	<-dataRight
	for i := 1; i <= 7000; i = i + 100 {
		SendTCPtoServer(task.GIVE_EXPLORER_PROGRESS, strconv.Itoa(i)+"/7000", conn, info)
		<-dataRight
	}
	sendZipFile("explorer.zip", task.GIVE_EXPLORER_INFO, task.GIVE_EXPLORER_DATA, task.GIVE_EXPLORER_END, conn, dataRight, info)
}

func sendZipFile(zipPath string, taskInfo task.TaskType, taskData task.TaskType, taskEnd task.TaskType, conn net.Conn, dataRight chan int, info []string) {
	path := filepath.Join("mockFiles", zipPath)
	fileLen, err := file.GetFileSize(path)
	if err != nil {
		logger.Error(info[0] + ":: Error getting file size: " + err.Error())
	}
	SendTCPtoServer(taskInfo, strconv.Itoa(fileLen), conn, info)
	<-dataRight
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error(info[0] + ":: Read file error: " + err.Error())
	}
	start := 0
	for {
		end := int(math.Min(float64(fileLen), float64(start+65436)))
		data := content[start:end]
		SendDataTCPtoServer(taskData, data, conn, info)
		start += 65436
		if start >= fileLen {
			SendTCPtoServer(taskEnd, "", conn, info)
			<-dataRight
			break
		}
		<-dataRight
	}
	conn.Close()
}
