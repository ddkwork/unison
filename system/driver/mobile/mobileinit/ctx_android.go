package mobileinit

import "C"

import (
	"errors"
	"runtime"
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
)

var currentVM *C.JavaVM

var currentCtx C.jobject

func SetCurrentContext(vm unsafe.Pointer, ctx uintptr) {
	currentVM = (*C.JavaVM)(vm)
	currentCtxPrev := currentCtx
	currentCtx = (C.jobject)(ctx)
	RunOnJVM(func(vm, jniEnv, ctx uintptr) error {
		env := (*C.JNIEnv)(unsafe.Pointer(jniEnv))
		C.deletePrevCtx(env, C.jobject(currentCtxPrev))
		return nil
	})
}

func RunOnJVM(fn func(vm, env, ctx uintptr) error) error {
	errch := make(chan error)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		env := C.uintptr_t(0)
		attached := C.int(0)
		if errStr := C.lockJNI(currentVM, &env, &attached); errStr != nil {
			errch <- errors.New(C.GoString(errStr))
			return
		}
		if attached != 0 {
			defer C.unlockJNI(currentVM)
		}

		vm := uintptr(unsafe.Pointer(currentVM))
		if mylog.Check(fn(vm, uintptr(env), uintptr(currentCtx))); err != nil {
			errch <- err
			return
		}

		if exc := C.checkException(env); exc != nil {
			errch <- errors.New(C.GoString(exc))
			C.free(unsafe.Pointer(exc))
			return
		}
		errch <- nil
	}()
	return <-errch
}
