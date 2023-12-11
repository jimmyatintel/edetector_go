package clientsearchsend

import (
	"bytes"
	"testing"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func TestAppendByteMsg(t *testing.T) {
	longMsg := make([]byte, 70000)
	longMsg[0] = 't'

	tests := []struct {
		data []byte
		msg  []byte
	}{
		{make([]byte, 200), []byte("test content")},
		{[]byte("data header"), []byte("t")},
		{make([]byte, 10), longMsg},
	}

	for ind, tt := range tests {
		out := AppendByteMsg(tt.data, tt.msg)
		header := min(len(tt.data), 100)
		if len(out) != 65536 || !bytes.Equal(out[:header], tt.data[:header]) || !bytes.Equal(out[100:min(100+len(tt.msg), 65536)], tt.msg[:min(len(tt.msg), 65536-100)]) {
			t.Errorf("Failed TestCase %v: AppendByteMsg", ind)
		}
	}
}
