//go:build windows

package desktop

import (
	"fmt"
	"log/slog"
	"syscall"
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sys/windows/registry"
)

const (
	themeRegKey  = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`
	themeRegName = `AppsUseLightTheme`
)

func (app *App) IsDark() bool {
	k := mylog.Check2(registry.OpenKey(registry.CURRENT_USER, themeRegKey, registry.QUERY_VALUE))
	defer k.Close()
	val, _ := mylog.Check3(k.GetIntegerValue(themeRegName))
	return val == 0
}

func (app *App) IsDarkMonitor(fn func(isDark bool), done chan struct{}) (chan error, error) {
	var regNotifyChangeKeyValue *syscall.Proc
	if advapi32 := mylog.Check2(syscall.LoadDLL("Advapi32.dll")); err == nil {
		if p := mylog.Check2(advapi32.FindProc("RegNotifyChangeKeyValue")); err == nil {
			regNotifyChangeKeyValue = p
		} else {
			return nil, fmt.Errorf("error finding function RegNotifyChangeKeyValue in Advapi32.dll: %w", err)
		}
	}

	ec := make(chan error)
	if regNotifyChangeKeyValue != nil {
		go func() {
			k := mylog.Check2(registry.OpenKey(registry.CURRENT_USER, themeRegKey, syscall.KEY_NOTIFY|registry.QUERY_VALUE))

			var wasDark, haveSetWasDark bool
			for {
				select {
				case <-done:

					return
				default:
					regNotifyChangeKeyValue.Call(uintptr(k), 0, 0x00000001|0x00000004, 0, 0)
					val, _ := mylog.Check3(k.GetIntegerValue(themeRegName))

					isDark := val == 0

					if isDark != wasDark || !haveSetWasDark {
						fn(isDark)
						wasDark = isDark
						haveSetWasDark = true
					}
				}
			}
		}()
	}
	return ec, nil
}

func (w *Window) SetTitleBarIsDark(isDark bool) {
	if !w.IsVisible() {
		return
	}
	hwnd := w.Glw.GetWin32Window()
	dwm := syscall.NewLazyDLL("dwmapi.dll")
	setAtt := dwm.NewProc("DwmSetWindowAttribute")
	ret, _ := mylog.Check3(setAtt.Call(uintptr(unsafe.Pointer(hwnd)),
		20,
		uintptr(unsafe.Pointer(&isDark)),
		8))
	if ret != 0 {
		slog.Error("failed to set window title bar color", "err", err)
	}
}
