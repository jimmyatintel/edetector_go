package connectionmap

import (
	"net"
	"sync"
)

var ConnInfoMap = sync.Map{}

type ConnInfo struct {
	TaskId     string
	Msg        string
	DataLen    int
	CurDataLen int
}

func StoreConnInfo(conn net.Conn, data ConnInfo) {
	ConnInfoMap.Store(conn, data)
}

func GetConnInfo(conn net.Conn) (ConnInfo, bool) {
	data, ok := ConnInfoMap.Load(conn)
	if !ok {
		return ConnInfo{}, false
	}
	return data.(ConnInfo), true
}

func RemoveConnInfo(conn net.Conn) {
	ConnInfoMap.Delete(conn)
}

func UpdateConnInfoCurDataLen(conn net.Conn, dataLen int) int {
	data, ok := ConnInfoMap.Load(conn)
	if !ok {
		return 0
	}
	connInfo := data.(ConnInfo)
	connInfo.CurDataLen += dataLen
	ConnInfoMap.Store(conn, connInfo)

	if connInfo.CurDataLen >= connInfo.DataLen {
		return connInfo.DataLen
	} else {
		return connInfo.CurDataLen
	}
}

func UpdateConnInfoDataLen(conn net.Conn, dataLen int) {
	data, ok := ConnInfoMap.Load(conn)
	if !ok {
		return
	}
	connInfo := data.(ConnInfo)
	connInfo.DataLen = dataLen
	connInfo.CurDataLen = 0
	ConnInfoMap.Store(conn, connInfo)
}
