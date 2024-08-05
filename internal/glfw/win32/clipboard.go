//go:build windows

/*
 * Based on code originally from https://github.com/atotto/clipboard. Copyright (c) 2013 Ato Araki. All rights reserved.
 */

package win32

import (
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
)

const (
	cfUnicodetext = 13
	gmemMoveable  = 0x0002
)

func waitOpenClipboard() {
	started := time.Now()
	limit := started.Add(time.Second)
	var r uintptr
	for time.Now().Before(limit) {
		r, _ = mylog.Check3(procOpenClipboard.Call(0))
		if r != 0 {
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func GetClipboardText() string {
	// LockOSThread ensure that the whole method will keep executing on the same thread from begin to end (it actually locks the goroutine thread attribution).
	// Otherwise if the goroutine switch thread during execution (which is a common practice), the OpenClipboard and CloseClipboard will happen on two different threads, and it will result in a clipboard deadlock.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if formatAvailable, _ := mylog.Check3(procIsClipboardFormatAvailable.Call(cfUnicodetext)); formatAvailable == 0 {
		return ""
	}
	waitOpenClipboard()

	h, _ := mylog.Check3(procGetClipboardData.Call(cfUnicodetext))
	if h == 0 {
		_, _, _ = procCloseClipboard.Call()
		return ""
	}

	l, _ := mylog.Check3(kernelGlobalLock.Call(h))
	if l == 0 {
		_, _, _ = procCloseClipboard.Call()
		return ""
	}

	text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(l))[:])

	r, _ := mylog.Check3(kernelGlobalUnlock.Call(h))
	if r == 0 {
		_, _, _ = procCloseClipboard.Call()
		return ""
	}

	closed, _ := mylog.Check3(procCloseClipboard.Call())
	if closed == 0 {
		return ""
	}
	return text
}

func SetClipboardText(text string) {
	// LockOSThread ensure that the whole method will keep executing on the same thread from begin to end (it actually locks the goroutine thread attribution).
	// Otherwise if the goroutine switch thread during execution (which is a common practice), the OpenClipboard and CloseClipboard will happen on two different threads, and it will result in a clipboard deadlock.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	waitOpenClipboard()

	r, _ := mylog.Check3(procEmptyClipboard.Call(0))
	if r == 0 {
		_, _, _ = procCloseClipboard.Call()
		return
	}

	data := mylog.Check2(syscall.UTF16FromString(text))

	// "If the hMem parameter identifies a memory object, the object must have
	// been allocated using the function with the GMEM_MOVEABLE flag."
	h, _ := mylog.Check3(kernelGlobalAlloc.Call(gmemMoveable, uintptr(len(data)*int(unsafe.Sizeof(data[0])))))
	if h == 0 {
		_, _, _ = procCloseClipboard.Call()
		return
	}
	defer func() {
		if h != 0 {
			kernelGlobalFree.Call(h)
		}
	}()

	l, _ := mylog.Check3(kernelGlobalLock.Call(h))
	if l == 0 {
		_, _, _ = procCloseClipboard.Call()
		return
	}

	r, _ = mylog.Check3(kernelLstrcpy.Call(l, uintptr(unsafe.Pointer(&data[0]))))
	if r == 0 {
		_, _, _ = procCloseClipboard.Call()
		return
	}

	r, _, e := kernelGlobalUnlock.Call(h)
	if r == 0 {
		if e.(syscall.Errno) != 0 {
			_, _, _ = procCloseClipboard.Call()
			return
		}
	}

	r, _ = mylog.Check3(procSetClipboardData.Call(cfUnicodetext, h))
	if r == 0 {
		_, _, _ = procCloseClipboard.Call()
		return
	}
	h = 0 // suppress deferred cleanup
	closed, _ := mylog.Check3(procCloseClipboard.Call())
	if closed == 0 {
		return
	}
}
