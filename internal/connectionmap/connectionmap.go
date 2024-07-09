package connectionmap

import (
	"net"
	"sync"
)

var ConnInfoMap = sync.Map{}

type ConnInfo struct {
	TaskId  string
	Msg     string
	DataLen int
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
