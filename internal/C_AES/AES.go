package C_AES

// #cgo LDFLAGS: -L. -l:mylib.so
// #include <stdlib.h>
// #include "warp_AES.hxx"
import "C"
import "unsafe"

func Decryptbuffer(cipherText []byte, size int, out []byte) {
	cChars := (*C.char)(C.malloc(C.size_t(size) * C.sizeof_char))
	defer C.free(unsafe.Pointer(cChars))
	inData := (*C.char)(unsafe.Pointer(C.CBytes(cipherText)))
	csize := C.int(size)
	defer C.free(unsafe.Pointer(inData))
	// pt_char := unsafe.Pointer(cChars)
	C.DecryptBuffer_cplus(inData, csize, cChars)
	// fmt.Println("pt_char: ", pt_char)
	outbyte := C.GoBytes(unsafe.Pointer(cChars), C.int(size))
	for i := 0; i < size; i++ {
		out[i] = outbyte[i]
	}
	// return str
}

func Encryptbuffer(Text []byte, size int, out []byte) {
	cChars := (*C.char)(C.malloc(C.size_t(size) * C.sizeof_char))
	defer C.free(unsafe.Pointer(cChars))
	inData := (*C.char)(unsafe.Pointer(C.CBytes(Text)))
	csize := C.int(size)
	defer C.free(unsafe.Pointer(inData))
	// pt_char := unsafe.Pointer(cChars)
	C.EncryptBuffer_cplus(inData, csize, cChars)
	// fmt.Println("pt_char: ", pt_char)
	outbyte := C.GoBytes(unsafe.Pointer(cChars), C.int(size))
	for i := 0; i < size; i++ {
		out[i] = outbyte[i]
	}
	// return str
}
