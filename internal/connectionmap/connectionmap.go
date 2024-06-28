package connectionmap

import (
	"net"
	"sync"
)

var ConnMsgMap = sync.Map{}

type ConnMsg struct {
	Info    string
	DataLen int
}

func StoreConnMsg(conn net.Conn, data ConnMsg) {
	ConnMsgMap.Store(conn, data)
}

func GetConnMsg(conn net.Conn) (ConnMsg, bool) {
	data, ok := ConnMsgMap.Load(conn)
	if !ok {
		return ConnMsg{}, false
	}
	return data.(ConnMsg), true
}

func RemoveConnMsg(conn net.Conn) {
	ConnMsgMap.Delete(conn)
}
