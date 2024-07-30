//go:build android

package android

import "C"

import (
	"image"
	"log"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"time"
	"unsafe"

	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/system/driver/base"
	"cogentcore.org/core/system/driver/mobile/callfn"
	"cogentcore.org/core/system/driver/mobile/mobileinit"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison/system"
)

var mimeMap = map[string]string{
	".txt": "text/plain",
}

func RunOnJVM(fn func(vm, jniEnv, ctx uintptr) error) error {
	return mobileinit.RunOnJVM(fn)
}

func setCurrentContext(vm *C.JavaVM, ctx C.jobject) {
	mobileinit.SetCurrentContext(unsafe.Pointer(vm), uintptr(ctx))
}

func callMain(mainPC uintptr) {
	for _, name := range []string{"FILESDIR", "TMPDIR", "PATH", "LD_LIBRARY_PATH"} {
		n := C.CString(name)
		os.Setenv(name, C.GoString(C.getenv(n)))
		C.free(unsafe.Pointer(n))
	}

	var curtime C.time_t
	var curtm C.struct_tm
	C.time(&curtime)
	C.localtime_r(&curtime, &curtm)
	tzOffset := int(curtm.tm_gmtoff)
	tz := C.GoString(curtm.tm_zone)
	time.Local = time.FixedZone(tz, tzOffset)

	go callfn.CallFn(mainPC)
}

func onSaveInstanceState(activity *C.ANativeActivity, outSize *C.size_t) unsafe.Pointer {
	return nil
}

func onCreate(activity *C.ANativeActivity) {
	windowConfigChange <- windowConfigRead(activity)
}

func onDestroy(activity *C.ANativeActivity) {
	activityDestroyed <- struct{}{}
}

func onWindowFocusChanged(activity *C.ANativeActivity, hasFocus C.int) {
	TheApp.Mu.Lock()
	defer TheApp.Mu.Unlock()
	if hasFocus > 0 {
		TheApp.Event.Window(events.WinFocus)
	} else {
		TheApp.Event.Window(events.WinFocusLost)
	}
}

func onNativeWindowCreated(activity *C.ANativeActivity, window *C.ANativeWindow) {
	TheApp.Mu.Lock()
	defer TheApp.Mu.Unlock()
	TheApp.SetSystemWindow(uintptr(unsafe.Pointer(window)))
}

func onNativeWindowRedrawNeeded(activity *C.ANativeActivity, window *C.ANativeWindow) {
	windowRedrawNeeded <- window
}

func onNativeWindowDestroyed(activity *C.ANativeActivity, window *C.ANativeWindow) {
	windowDestroyed <- window
}

func onInputQueueCreated(activity *C.ANativeActivity, q *C.AInputQueue) {
	inputQueue <- q
	<-inputQueueDone
}

func onInputQueueDestroyed(activity *C.ANativeActivity, q *C.AInputQueue) {
	inputQueue <- nil
	<-inputQueueDone
}

func onContentRectChanged(activity *C.ANativeActivity, rect *C.ARect) {
}

func setDarkMode(dark C.bool) {
	TheApp.Dark = bool(dark)
}

type windowConfig struct {
	Orientation system.ScreenOrientation
	DPI         float32
}

func windowConfigRead(activity *C.ANativeActivity) windowConfig {
	defer func() { system.HandleRecover(recover()) }()

	aconfig := C.AConfiguration_new()
	C.AConfiguration_fromAssetManager(aconfig, activity.assetManager)
	orient := C.AConfiguration_getOrientation(aconfig)
	density := C.AConfiguration_getDensity(aconfig)
	C.AConfiguration_delete(aconfig)

	var dpi int
	switch density {
	case C.ACONFIGURATION_DENSITY_DEFAULT:
		dpi = 160
	case C.ACONFIGURATION_DENSITY_LOW,
		C.ACONFIGURATION_DENSITY_MEDIUM,
		213,
		C.ACONFIGURATION_DENSITY_HIGH,
		320,
		480,
		640:
		dpi = int(density)
	case C.ACONFIGURATION_DENSITY_NONE:
		slog.Warn("android device reports no screen density")
		dpi = 72
	default:

		slog.Warn("android device reports unknown screen density", "density", density)

		if density > 0 {
			dpi = int(density)
		} else {
			dpi = 72
		}
	}

	o := system.OrientationUnknown
	switch orient {
	case C.ACONFIGURATION_ORIENTATION_PORT:
		o = system.Portrait
	case C.ACONFIGURATION_ORIENTATION_LAND:
		o = system.Landscape
	}

	return windowConfig{
		Orientation: o,
		DPI:         float32(dpi),
	}
}

func onConfigurationChanged(activity *C.ANativeActivity) {
	windowConfigChange <- windowConfigRead(activity)
}

func onLowMemory(activity *C.ANativeActivity) {
	runtime.GC()
	debug.FreeOSMemory()
}

var (
	inputQueue         = make(chan *C.AInputQueue)
	inputQueueDone     = make(chan struct{})
	windowDestroyed    = make(chan *C.ANativeWindow)
	windowRedrawNeeded = make(chan *C.ANativeWindow)
	windowRedrawDone   = make(chan struct{})
	windowConfigChange = make(chan windowConfig)
	activityDestroyed  = make(chan struct{})
)

func (a *App) MainLoop() {
	a.MainQueue = make(chan base.FuncRun)
	a.MainDone = make(chan struct{})

	go func() {
		defer func() { system.HandleRecover(recover()) }()
		if mylog.Check(mobileinit.RunOnJVM(RunInputQueue)); err != nil {
			log.Fatalf("app: %v", err)
		}
	}()

	if mylog.Check(mobileinit.RunOnJVM(TheApp.MainUI)); err != nil {
		log.Fatalf("app: %v", err)
	}
}

func (a *App) ShowVirtualKeyboard(typ styles.VirtualKeyboards) {
	mylog.Check(mobileinit.RunOnJVM(func(vm, jniEnv, ctx uintptr) error {
		env := (*C.JNIEnv)(unsafe.Pointer(jniEnv))
		C.showKeyboard(env, C.int(int32(typ)))
		return nil
	}))
}

func (a *App) HideVirtualKeyboard() {
	if mylog.Check(mobileinit.RunOnJVM(hideSoftInput)); err != nil {
		log.Fatalf("app: %v", err)
	}
}

func hideSoftInput(vm, jniEnv, ctx uintptr) error {
	env := (*C.JNIEnv)(unsafe.Pointer(jniEnv))
	C.hideKeyboard(env)
	return nil
}

func insetsChanged(top, bottom, left, right int) {
	TheApp.Insets.Set(top, right, bottom, left)
	TheApp.Event.WindowResize()
}

func (a *App) MainUI(vm, jniEnv, ctx uintptr) error {
	defer func() { system.HandleRecover(recover()) }()

	var dpi float32
	var orientation system.ScreenOrientation

	for {
		select {
		case <-a.MainDone:
			a.FullDestroyVk()
			return nil
		case f := <-a.MainQueue:
			f.F()
			if f.Done != nil {
				f.Done <- struct{}{}
			}
		case cfg := <-windowConfigChange:
			dpi = cfg.DPI
			orientation = cfg.Orientation
		case w := <-windowRedrawNeeded:
			widthPx := int(C.ANativeWindow_getWidth(w))
			heightPx := int(C.ANativeWindow_getHeight(w))

			a.Scrn.Orientation = orientation

			a.Scrn.DevicePixelRatio = 1
			a.Scrn.PixSize = image.Pt(widthPx, heightPx)
			a.Scrn.Geometry.Max = a.Scrn.PixSize

			a.Scrn.PhysicalDPI = dpi
			a.Scrn.LogicalDPI = dpi

			if system.InitScreenLogicalDPIFunc != nil {
				system.InitScreenLogicalDPIFunc()
			}

			physX := 25.4 * float32(widthPx) / dpi
			physY := 25.4 * float32(heightPx) / dpi
			a.Scrn.PhysicalSize = image.Pt(int(physX), int(physY))

			if system.OnSystemWindowCreated != nil {
				system.OnSystemWindowCreated <- struct{}{}
			}

			a.Event.WindowResize()
		case <-windowDestroyed:

			a.Win.SetSize(image.Point{})
			a.Event.Window(events.WinMinimize)
			a.DestroyVk()
		case <-activityDestroyed:

		}
	}
}

func RunInputQueue(vm, jniEnv, ctx uintptr) error {
	env := (*C.JNIEnv)(unsafe.Pointer(jniEnv))

	l := C.ALooper_prepare(C.ALOOPER_PREPARE_ALLOW_NON_CALLBACKS)
	pending := make(chan *C.AInputQueue, 1)
	go func() {
		for q := range inputQueue {
			pending <- q
			C.ALooper_wake(l)
		}
	}()

	var q *C.AInputQueue
	for {
		if C.ALooper_pollAll(-1, nil, nil, nil) == C.ALOOPER_POLL_WAKE {
			select {
			default:
			case p := <-pending:
				if q != nil {
					ProcessEvents(env, q)
					C.AInputQueue_detachLooper(q)
				}
				q = p
				if q != nil {
					C.AInputQueue_attachLooper(q, l, 0, nil, nil)
				}
				inputQueueDone <- struct{}{}
			}
		}
		if q != nil {
			ProcessEvents(env, q)
		}
	}
}
