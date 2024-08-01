package glfw

import (
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"golang.org/x/sys/windows"
)

//go:embed glfw.dll
var dllData []byte

func init() {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func init() {
	// runtime.LockOSThread()
	dir := mylog.Check2(os.UserCacheDir())
	dir = filepath.Join(dir, "glfw3", "dll_cache")
	stream.CreatDirectory(dir)
	mylog.Check(windows.SetDllDirectory(dir))
	sha := sha256.Sum256(dllData)
	dllName := fmt.Sprintf("glfw3-%s.dll", base64.RawURLEncoding.EncodeToString(sha[:]))
	filePath := filepath.Join(dir, dllName)
	if !stream.IsFilePath(filePath) {
		stream.WriteTruncate(filePath, dllData)
	}
	mylog.Check2(GengoLibrary.LoadFrom(filePath))
}

func TestName(t *testing.T) {
	Init()
	mylog.Info("version", BytePointerToString(GetVersionString()))
	defer Terminate()

	// WindowHint(Visible, False)
	// WindowHint(Resizable,Enable(!Resizable))
	// WindowHint(Decorated,Enable(!w.Decorated))
	// WindowHint(Floating,Enable(Floating))
	// WindowHint(AutoIconify, False)
	// WindowHint(TransparentFramebuffer, False)
	// WindowHint(FocusOnShow, False)
	// WindowHint(ScaleToMonitor, False)

	// PostEmptyEvent()

	w := CreateWindow(200, 200, StringToBytePointer("hello word"), nil, nil)

	// SetCursorEnterCallback(w, func() {})
	// SetCursorPosCallback(w, func() {})
	// SetMouseButtonCallback(w, func() {})
	// SetWindowFocusCallback(w, func() {})
	// SetWindowCloseCallback(w, func() {})
	// SetWindowSizeCallback(w, func() {})
	// SetWindowRefreshCallback(w, func() {})
	// SetScrollCallback(w, func() {})
	// SetKeyCallback(w, func() {})
	// SetCharCallback(w, func() {})
	// SetDropCallback(w, func() {})
	// SetWindowIcon(w,32, func() {})

	MakeContextCurrent(w)
	for {
		// PostEmptyEvent()
		PollEvents()
		SwapBuffers(w)
		if WindowShouldClose(w) != 0 {
			DestroyWindow(w)
			break
		}
	}
}

func StringToBytePointer(s string) *byte {
	bytes := []byte(s)
	ptr := &bytes[0]
	return ptr
}

func BytePointerToString(ptr *byte) string {
	var bytes []byte
	for *ptr != 0 {
		bytes = append(bytes, *ptr)
		ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
	}
	return string(bytes)
}
