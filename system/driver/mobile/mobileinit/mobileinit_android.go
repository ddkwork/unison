package mobileinit

import "C"

import (
	"bufio"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
)

var (
	ctag           = C.CString("GoLog")
	stderr, stdout *os.File
)

type infoWriter struct{}

func (infoWriter) Write(p []byte) (n int, err error) {
	cstr := C.CString(string(p))
	C.__android_log_write(C.ANDROID_LOG_INFO, ctag, cstr)
	C.free(unsafe.Pointer(cstr))
	return len(p), nil
}

func lineLog(f *os.File, priority C.int) {
	const logSize = 1024
	r := bufio.NewReaderSize(f, logSize)
	for {
		line, _ := mylog.Check3(r.ReadLine())
		str := string(line)

		cstr := C.CString(str)
		C.__android_log_write(priority, ctag, cstr)
		C.free(unsafe.Pointer(cstr))

	}
}

func init() {
	log.SetOutput(infoWriter{})

	log.SetFlags(log.Flags() &^ log.LstdFlags)

	r, w := mylog.Check3(os.Pipe())

	stderr = w
	if mylog.Check(syscall.Dup3(int(w.Fd()), int(os.Stderr.Fd()), 0)); err != nil {
		panic(err)
	}
	go lineLog(r, C.ANDROID_LOG_ERROR)

	r, w = mylog.Check3(os.Pipe())

	stdout = w
	if mylog.Check(syscall.Dup3(int(w.Fd()), int(os.Stdout.Fd()), 0)); err != nil {
		panic(err)
	}
	go lineLog(r, C.ANDROID_LOG_INFO)
}
