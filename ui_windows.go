// Copyright 2016 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unison

import (
	"errors"
	"fmt"
	"runtime"
	"syscall"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison/internal/graphicsdriver"
	"github.com/richardwilkes/unison/internal/graphicsdriver/directx"
	"github.com/richardwilkes/unison/internal/graphicsdriver/opengl"
	"github.com/richardwilkes/unison/internal/microsoftgdk"
	"github.com/richardwilkes/unison/internal/winver"

	"golang.org/x/sys/windows"

	"github.com/richardwilkes/unison/internal/glfw"
)

func (u *UserInterface) initializePlatform() error {
	return nil
}

type graphicsDriverCreatorImpl struct {
	transparent bool
}

func (g *graphicsDriverCreatorImpl) newAuto() (graphicsdriver.Graphics, GraphicsLibrary) {
	if winver.IsWindows10OrGreater() {
		d := g.newDirectX()
		if d != nil {
			return d, GraphicsLibraryDirectX
		}

		o := g.newOpenGL()
		if o != nil {
			return o, GraphicsLibraryOpenGL
		}
	} else {
		// Creating a swap chain on an older machine than Windows 10 might fail (#2613).
		// Prefer OpenGL to DirectX.
		o := g.newOpenGL()
		if o != nil {
			return o, GraphicsLibraryOpenGL
		}

		// Initializing OpenGL can fail, though this is pretty rare.
		d := g.newDirectX()
		if d != nil {
			return d, GraphicsLibraryDirectX
		}
	}

	return nil, GraphicsLibraryUnknown //, fmt.Errorf("ui: failed to choose graphics drivers: DirectX: %v, OpenGL: %v", dxErr, glErr)
}

func (*graphicsDriverCreatorImpl) newOpenGL() graphicsdriver.Graphics {
	return opengl.NewGraphics()
}

func (g *graphicsDriverCreatorImpl) newDirectX() graphicsdriver.Graphics {
	if g.transparent {
		return nil //, errors.New("ui: DirectX is not available with a transparent window")
	}
	return directx.NewGraphics()
}

func (*graphicsDriverCreatorImpl) newMetal() (graphicsdriver.Graphics, error) {
	return nil, errors.New("ui: Metal is not supported in this environment")
}

func (*graphicsDriverCreatorImpl) newPlayStation5() (graphicsdriver.Graphics, error) {
	return nil, errors.New("ui: PlayStation 5 is not supported in this environment")
}

// glfwMonitorSizeInGLFWPixels must be called from the main thread.
func glfwMonitorSizeInGLFWPixels(m *glfw.Monitor) (int, int) {
	vm := m.GetVideoMode()
	return vm.Width, vm.Height
}

func dipFromGLFWPixel(x float64, deviceScaleFactor float64) float64 {
	return x / deviceScaleFactor
}

func dipToGLFWPixel(x float64, deviceScaleFactor float64) float64 {
	return x * deviceScaleFactor
}

func (u *UserInterface) adjustWindowPosition(x, y int, monitor *Monitor) (int, int) {
	if microsoftgdk.IsXbox() {
		return x, y
	}

	mx := monitor.boundsInGLFWPixels.Min.X
	my := monitor.boundsInGLFWPixels.Min.Y
	// As the video width/height might be wrong,
	// adjust x/y at least to enable to handle the window (#328)
	if x < mx {
		x = mx
	}
	t := mylog.Check2(_GetSystemMetrics(_SM_CYCAPTION))

	if y < my+int(t) {
		y = my + int(t)
	}
	return x, y
}

func initialMonitorByOS() *Monitor {
	if microsoftgdk.IsXbox() {
		return theMonitors.primaryMonitor()
	}

	px, py := mylog.Check3(_GetCursorPos())

	x, y := int(px), int(py)

	// Find the monitor including the cursor.
	return theMonitors.monitorFromPosition(x, y)
}

func monitorFromWindowByOS(w *glfw.Window) *Monitor {
	if microsoftgdk.IsXbox() {
		return theMonitors.primaryMonitor()
	}
	window := w.GetWin32Window()
	return monitorFromWin32Window(window)
}

func monitorFromWin32Window(w windows.HWND) *Monitor {
	// Get the current monitor by the window handle instead of the window position. It is because the window
	// position is not reliable in some cases e.g. when the window is put across multiple monitors.

	m := _MonitorFromWindow(w, _MONITOR_DEFAULTTONEAREST)
	if m == 0 {
		// monitorFromWindow can return error on Wine. Ignore this.
		return nil
	}

	mi := mylog.Check2(_GetMonitorInfoW(m))

	x, y := int(mi.rcMonitor.left), int(mi.rcMonitor.top)
	for _, m := range theMonitors.append(nil) {
		mx := m.boundsInGLFWPixels.Min.X
		my := m.boundsInGLFWPixels.Min.Y
		if mx == x && my == y {
			return m
		}
	}
	return nil
}

func (u *UserInterface) nativeWindow() uintptr {
	w := u.window.GetWin32Window()
	return uintptr(w)
}

func (u *UserInterface) isNativeFullscreen() bool {
	return false
}

func (u *UserInterface) isNativeFullscreenAvailable() bool {
	return false
}

func (u *UserInterface) setNativeFullscreen(fullscreen bool) error {
	panic(fmt.Sprintf("ui: setNativeFullscreen is not implemented in this environment: %s", runtime.GOOS))
}

func (u *UserInterface) adjustViewSizeAfterFullscreen() error {
	return nil
}

func (u *UserInterface) setWindowResizingModeForOS(mode WindowResizingMode) error {
	return nil
}

func initializeWindowAfterCreation(w *glfw.Window) error {
	return nil
}

func (u *UserInterface) skipTaskbar() error {
	// S_FALSE is returned when CoInitializeEx is nested. This is a successful case.
	mylog.Check(windows.CoInitializeEx(0, windows.COINIT_MULTITHREADED))
	err != nil && !errors.Is(err, syscall.Errno(windows.S_FALSE))
	{
		return err
	}
	// CoUninitialize should be called even when CoInitializeEx returns S_FALSE.
	defer windows.CoUninitialize()

	ptr := _CoCreateInstance(&_CLSID_TaskbarList, nil, _CLSCTX_SERVER, &_IID_ITaskbarList)
	t := (*_ITaskbarList)(ptr)
	defer t.Release()

	w := u.window.GetWin32Window()
	mylog.Check(t.DeleteTab(w))


	return nil
}

func (u *UserInterface) setDocumentEdited(edited bool) error {
	return nil
}

func init() {
	if microsoftgdk.IsXbox() {
		// TimeBeginPeriod might not be defined in Xbox.
		return
	}
	// Use a better timer resolution (golang/go#44343).
	// An error is ignored. The application is still valid even if a higher resolution timer is not available.
	// TODO: This might not be necessary from Go 1.23.
	_ = windows.TimeBeginPeriod(1)
}
