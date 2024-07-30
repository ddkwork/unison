//go:build ios

package ios

import "C"

import (
	"image"
	"log"
	"log/slog"
	"strings"
	"unsafe"

	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/system"
)

func (a *App) MainLoop() {
	go a.AppSingle.MainLoop()
	C.runApp()
	log.Fatalln("unexpected return from runApp")
}

var DisplayMetrics struct {
	WidthPx int

	HeightPx int

	DPI float32

	ScreenScale int
}

func setWindowPtr(window *C.void) {
	TheApp.Mu.Lock()
	defer TheApp.Mu.Unlock()
	TheApp.SetSystemWindow(uintptr(unsafe.Pointer(window)))
}

func setDisplayMetrics(width, height int, scale int) {
	DisplayMetrics.WidthPx = width
	DisplayMetrics.HeightPx = height
	DisplayMetrics.ScreenScale = scale
}

func setScreen(scale int) {
	C.uname(&C.sysInfo)
	name := C.GoString(&C.sysInfo.machine[0])

	var v float32

	switch {
	case strings.HasPrefix(name, "iPhone"):
		v = 163
	case strings.HasPrefix(name, "iPad"):

		switch name {
		case "iPad2,5", "iPad2,6", "iPad2,7", "iPad4,4", "iPad4,5", "iPad4,6", "iPad4,7":
			v = 163
		default:
			v = 132
		}
	default:
		v = 163
	}

	if v == 0 {
		slog.Warn("unknown machine: %s", name)
		v = 163
	}

	DisplayMetrics.DPI = v * float32(scale)
	DisplayMetrics.ScreenScale = scale
}

func updateConfig(width, height, orientation int32) {
	TheApp.Mu.Lock()
	defer TheApp.Mu.Unlock()
	TheApp.Scrn.Orientation = system.OrientationUnknown
	switch orientation {
	case C.UIDeviceOrientationPortrait, C.UIDeviceOrientationPortraitUpsideDown:
		TheApp.Scrn.Orientation = system.Portrait
	case C.UIDeviceOrientationLandscapeLeft, C.UIDeviceOrientationLandscapeRight:
		TheApp.Scrn.Orientation = system.Landscape
		width, height = height, width
	}
	insets := C.getDevicePadding()
	s := DisplayMetrics.ScreenScale
	TheApp.Insets.Set(
		int(insets.top)*s,
		int(insets.right)*s,
		int(insets.bottom)*s,
		int(insets.left)*s,
	)

	TheApp.Scrn.DevicePixelRatio = float32(s)
	TheApp.Scrn.PixSize = image.Pt(int(width), int(height))
	TheApp.Scrn.Geometry.Max = TheApp.Scrn.PixSize

	TheApp.Scrn.PhysicalDPI = DisplayMetrics.DPI
	TheApp.Scrn.LogicalDPI = DisplayMetrics.DPI

	if system.InitScreenLogicalDPIFunc != nil {
		system.InitScreenLogicalDPIFunc()
	}

	physX := 25.4 * float32(width) / DisplayMetrics.DPI
	physY := 25.4 * float32(height) / DisplayMetrics.DPI
	TheApp.Scrn.PhysicalSize = image.Pt(int(physX), int(physY))

	TheApp.Dark = bool(C.isDark())

	if system.OnSystemWindowCreated != nil {
		system.OnSystemWindowCreated <- struct{}{}
	}
	TheApp.Event.WindowResize()
}

func lifecycleDead() {
	TheApp.FullDestroyVk()
}

func lifecycleAlive() {
}

func lifecycleVisible() {
	if TheApp.Win != nil {
		TheApp.Event.Window(events.WinShow)
	}
}

func lifecycleFocused() {
	if TheApp.Win != nil {
		TheApp.Event.Window(events.WinFocus)
	}
}

func (a *App) ShowVirtualKeyboard(typ styles.VirtualKeyboards) {
	C.showKeyboard(C.int(int32(typ)))
}

func (a *App) HideVirtualKeyboard() {
	C.hideKeyboard()
}
