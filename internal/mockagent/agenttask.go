package mockagent

import (
	"edetector_go/internal/task"
	"edetector_go/pkg/file"
	"edetector_go/pkg/logger"
	"math"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func agentDetect(conn net.Conn, detectStatus *string, info []string) {
	for {
		if *detectStatus == "1|1" {
			SendTCPtoServer(task.GIVE_DETECT_PROCESS, "conhost.exe|1698219996|\\??\\C:\\Windows\\system32\\conhost.exe 0x4|9430b20076a19e6ed9084530ddcc8caa|C:\\Windows\\System32\\conhost.exe|43760|ClientSearch.exe|C:\\Program Files (x86)\\eDetectorClient\\ClientSearch.exe|null|43416|0,0|0|0,0|0,0|null|NlsAnsiCodePage:0x0000FFFD0000FDE9 -> 0x0000003F000003B6;", conn, info)
		}
		time.Sleep(180 * time.Second)
	}
}

func agentScan(conn net.Conn, dataRight chan int, info []string) {
	SendTCPtoServer(task.GIVE_SCAN_INFO, "187", conn, info)
	<-dataRight
	for i := 1; i <= 187; i++ {
		SendTCPtoServer(task.GIVE_SCAN_PROGRESS, strconv.Itoa(i)+"/187", conn, info)
		<-dataRight
		randomSleep := (rand.Intn(10) + 1) * 100 // 0.1 ~ 1
		time.Sleep(time.Duration(randomSleep) * time.Millisecond)
	}
	sendZipFile("scan.zip", task.GIVE_SCAN_DATA_INFO, task.GIVE_SCAN, task.GIVE_SCAN_END, conn, dataRight, info)
}

func agentCollect(conn net.Conn, dataRight chan int, info []string) {
	for i := 1; i <= 48; i++ {
		SendTCPtoServer(task.GIVE_COLLECT_PROGRESS, strconv.Itoa(i)+"/48", conn, info)
		<-dataRight
		randomSleep := (rand.Intn(5) + 1) // 1 ~ 5
		time.Sleep(time.Duration(randomSleep) * time.Second)
	}
	sendZipFile("collect.zip", task.GIVE_COLLECT_DATA_INFO, task.GIVE_COLLECT_DATA, task.GIVE_COLLECT_DATA_END, conn, dataRight, info)
}

func agentDrive(conn net.Conn, dataRight chan int, info []string) {
	SendTCPtoServer(task.EXPLORER, "C|NTFS", conn, info)
	<-dataRight
	for i := 1; i <= 7531; i = i + 100 {
		SendTCPtoServer(task.GIVE_EXPLORER_PROGRESS, strconv.Itoa(i)+"/7531", conn, info)
		<-dataRight
		randomSleep := (rand.Intn(15) + 1) * 100 // 0.1 ~ 1.5
		time.Sleep(time.Duration(randomSleep) * time.Millisecond)
	}
	sendZipFile("explorer.zip", task.GIVE_EXPLORER_INFO, task.GIVE_EXPLORER_DATA, task.GIVE_EXPLORER_END, conn, dataRight, info)
}

func agentImage(conn net.Conn, dataRight chan int, info []string) {
	sendZipFile("explorer.zip", task.GIVE_IMAGE_INFO, task.GIVE_IMAGE, task.GIVE_IMAGE_END, conn, dataRight, info)
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
