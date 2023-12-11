package C_AES

import (
	"bytes"
	"edetector_go/pkg/file"
	"os"
	"testing"
)

var encryptShort, decryptShort, encryptLong, decryptLong []byte

func init() {
	for i := 0; i < 2; i++ {
		file.MoveToParentDir()
	}
	var err error
	encryptShort, err = os.ReadFile("test/encrypted_short_packet")
	if err != nil {
		panic(err)
	}
	decryptShort, err = os.ReadFile("test/decrypted_short_packet")
	if err != nil {
		panic(err)
	}
	encryptLong, err = os.ReadFile("test/encrypted_long_packet")
	if err != nil {
		panic(err)
	}
	decryptLong, err = os.ReadFile("test/decrypted_long_packet")
	if err != nil {
		panic(err)
	}
}

func TestDecryptbuffer(t *testing.T) {
	tests := []struct {
		cipherText []byte
		size       int
		out        []byte
		want       []byte
	}{
		{encryptShort, len(encryptShort), make([]byte, len(decryptShort)), decryptShort},
		{encryptLong, len(encryptLong), make([]byte, len(decryptLong)), decryptLong},
	}
	for ind, tt := range tests {
		Decryptbuffer(tt.cipherText, tt.size, tt.out)
		if !bytes.Equal(tt.out, tt.want) {
			t.Errorf("Failed TestCase %v: Decryptbuffer", ind)
		}
	}
}

func TestEncryptbuffer(t *testing.T) {
	tests := []struct {
		Text []byte
		size int
		out  []byte
		want []byte
	}{
		{decryptShort, len(decryptShort), make([]byte, len(encryptShort)), encryptShort},
		{decryptLong, len(decryptLong), make([]byte, len(encryptLong)), encryptLong},
	}
	for ind, tt := range tests {
		Encryptbuffer(tt.Text, tt.size, tt.out)
		if !bytes.Equal(tt.out, tt.want) {
			t.Errorf("Failed TestCase %v: Encryptbuffer", ind)
		}
	}
}
