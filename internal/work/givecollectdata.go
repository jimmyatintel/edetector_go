package work

import (
	"bytes"
	C_AES "edetector_go/internal/C_AES"
	"edetector_go/internal/checkdir"
	clientsearchsend "edetector_go/internal/clientsearch/send"
	packet "edetector_go/internal/packet"
	task "edetector_go/internal/task"
	taskservice "edetector_go/internal/taskservice"
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"errors"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"net"

	"go.uber.org/zap"
)

var collectMu sync.Mutex
var collectTotalMap = make(map[string]int)
var collectCountMap = make(map[string]int)
var collectProgressMap = make(map[string]int)
var dbWorkingPath = "dbWorking"
var dbUstagePath = "dbUnstage"

// ! tmp version
var tmpMu sync.Mutex
var lastDataTime = time.Now()

func init() {
	checkdir.CheckDir(dbWorkingPath)
	checkdir.CheckDir(dbUstagePath)
}

func ImportStartup(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("ImportStartup: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func CollectInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("CollectInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.GET_COLLECT_INFO_DATA,
		Message:    "10",
	}
	collectMu.Lock()
	collectProgressMap[p.GetRkey()] = 0
	collectMu.Unlock()
	go collectProgress(p.GetRkey())
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectProgress(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectProgress: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))

	// update progress
	parts := strings.Split(p.GetMessage(), "/")
	if len(parts) != 2 {
		return task.FAIL, errors.New("invalid progress format")
	}
	numerator, err := strconv.Atoi(parts[0])
	if err != nil {
		return task.FAIL, err
	}
	denominator, err := strconv.Atoi(parts[1])
	if err != nil {
		return task.FAIL, err
	}
	collectMu.Lock()
	collectProgressMap[p.GetRkey()] = int(math.Min((float64(numerator) / float64(denominator) * 20), 20))
	collectMu.Unlock()

	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err = clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectDataInfo(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataInfo: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	//! tmp version
	tmpMu.Lock()
	lastDataTime = time.Now()
	tmpMu.Unlock()
	go TmpEnd(p.GetRkey())
	// init collect info
	len, err := strconv.Atoi(p.GetMessage())
	if err != nil {
		return task.FAIL, err
	}
	collectCountMap[p.GetRkey()] = 0
	collectTotalMap[p.GetRkey()] = len

	// Create or truncate the db file
	path := filepath.Join(dbWorkingPath, (p.GetRkey() + ".db"))
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return task.FAIL, err
	}
	file.Close()

	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err = clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectData(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Debug("GiveCollectData: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	// write file
	dp := packet.CheckIsData(p)
	decrypt_buf := bytes.Repeat([]byte{0}, len(dp.Raw_data))
	C_AES.Decryptbuffer(dp.Raw_data, len(dp.Raw_data), decrypt_buf)
	decrypt_buf = decrypt_buf[100:]
	path := filepath.Join(dbWorkingPath, (p.GetRkey() + ".db"))
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return task.FAIL, err
	}
	_, err = file.Seek(0, 2)
	if err != nil {
		return task.FAIL, err
	}
	_, err = file.Write(decrypt_buf)
	if err != nil {
		return task.FAIL, err
	}
	file.Close()

	// update progress
	collectCountMap[p.GetRkey()] += 1
	collectMu.Lock()
	collectProgressMap[p.GetRkey()] = int(20 + float64(collectCountMap[p.GetRkey()])/(float64(collectTotalMap[p.GetRkey()]/65436))*80)
	collectMu.Unlock()
	//! tmp version
	tmpMu.Lock()
	lastDataTime = time.Now()
	tmpMu.Unlock()
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err = clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectDataEnd(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataEnd: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))

	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func GiveCollectDataError(p packet.Packet, conn net.Conn) (task.TaskResult, error) {
	logger.Info("GiveCollectDataError: ", zap.Any("message", p.GetRkey()+", Msg: "+p.GetMessage()))
	var send_packet = packet.WorkPacket{
		MacAddress: p.GetMacAddress(),
		IpAddress:  p.GetipAddress(),
		Work:       task.DATA_RIGHT,
		Message:    "",
	}
	err := clientsearchsend.SendTCPtoClient(send_packet.Fluent(), conn)
	if err != nil {
		return task.FAIL, err
	}
	return task.SUCCESS, nil
}

func collectProgress(clientid string) {
	for {
		collectMu.Lock()
		if collectProgressMap[clientid] >= 100 {
			break
		}
		rowsAffected := query.Update_progress(collectProgressMap[clientid], clientid, "StartCollect")
		collectMu.Unlock()
		if rowsAffected != 0 {
			go taskservice.RequestToUser(clientid)
		}
	}
}

func TmpEnd(key string) { //!tmp version
	for {
		tmpMu.Lock()
		if time.Since(lastDataTime) > time.Duration(120)*time.Second {
			lastDataTime = time.Now()
			tmpMu.Unlock()
			logger.Info("Collect tmp End version: ", zap.Any("message", key))
			// truncate data
			path := filepath.Join(dbWorkingPath, (key + ".db"))
			data, err := os.ReadFile(path)
			if err != nil {
				logger.Error("Read file error", zap.Any("message", err.Error()))
				continue
			}
			fileInfo, err := os.Stat(path)
			if err != nil {
				logger.Error("Stat file error", zap.Any("message", err.Error()))
				continue
			}
			realLen := fileInfo.Size()
			if int(realLen) < collectTotalMap[key] {
				logger.Error("Incomplete data")
				continue
			}
			err = os.WriteFile(path, data[:collectTotalMap[key]], 0644)
			if err != nil {
				logger.Error("Write file error", zap.Any("message", err.Error()))
				continue
			}
			// move to dbUnstage
			srcPath := filepath.Join(dbWorkingPath, (key + ".db"))
			dstPath := filepath.Join(dbUstagePath, (key + ".db"))
			err = moveFile(srcPath, dstPath)
			if err != nil {
				logger.Error("Move failed", zap.Any("message", err.Error()))
				continue
			}
			return
		}
		tmpMu.Unlock()
	}
}

func moveFile(srcPath string, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	err = os.Remove(srcPath)
	if err != nil {
		return err
	}
	return nil
}
