//go:build darwin && (arm || arm64)

package mobileinit

import (
	"log"
	"unsafe"
)

import "C"

type aslWriter struct{}

func (aslWriter) Write(p []byte) (n int, err error) {
	cstr := C.CString(string(p))
	C.log_wrap(cstr)
	C.free(unsafe.Pointer(cstr))
	return len(p), nil
}

func init() {
	log.SetOutput(aslWriter{})
}
