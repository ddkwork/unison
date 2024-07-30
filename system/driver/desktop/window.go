package desktop

import (
	"image"
	"log"

	"cogentcore.org/core/events"
	"cogentcore.org/core/system"
	"cogentcore.org/core/system/driver/base"
	"cogentcore.org/core/vgpu/vdraw"
	"github.com/ddkwork/golibrary/mylog"
	vk "github.com/goki/vulkan"
	"github.com/richardwilkes/unison/internal/glfw"
)

type Window struct {
	base.WindowMulti[*App, *vdraw.Drawer]
	Glw          *glfw.Window
	ScreenWindow string
}

func (w *Window) IsVisible() bool {
	return w.WindowMulti.IsVisible() && w.Glw != nil
}

func (w *Window) Activate() bool {
	if w == nil || w.Glw == nil {
		return false
	}
	w.Glw.MakeContextCurrent()
	return true
}

func (w *Window) DeActivate() {
	glfw.DetachCurrentContext()
}

func NewGlfwWindow(opts *system.NewWindowOptions, sc *system.Screen) (*glfw.Window, error) {
	_, _, tool, fullscreen := system.WindowFlagsToBool(opts.Flags)

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.Focused, glfw.True)

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)

	if fullscreen {
		glfw.WindowHint(glfw.Maximized, glfw.True)
	}
	if tool {
		glfw.WindowHint(glfw.Decorated, glfw.False)
	} else {
		glfw.WindowHint(glfw.Decorated, glfw.True)
	}

	sz := opts.Size
	if TheApp.Platform() == system.MacOS {
		sz = sc.WinSizeFromPix(opts.Size)
	}
	win := mylog.Check2(glfw.CreateWindow(sz.X, sz.Y, opts.GetTitle(), nil, nil))

	if !fullscreen {
		win.SetPos(opts.Pos.X, opts.Pos.Y)
	}
	if opts.Icon != nil {
		win.SetIcon(opts.Icon)
	}
	return win, err
}

func (w *Window) Screen() *system.Screen {
	if w == nil || w.Glw == nil {
		return TheApp.Screens[0]
	}
	w.Mu.Lock()
	defer w.Mu.Unlock()

	var sc *system.Screen
	mon := w.Glw.GetMonitor()

	if mon != nil {
		if MonitorDebug {
			log.Printf("MonitorDebug: desktop.Window.Screen: %v: got screen: %v\n", w.Nm, mon.GetName())
		}
		sc = TheApp.ScreenByName(mon.GetName())
		if sc == nil {
			log.Printf("MonitorDebug: desktop.Window.Screen: could not find screen of name: %v\n", mon.GetName())
			sc = TheApp.Screens[0]
		}
		goto setScreen
	}
	sc = w.GetScreenOverlap()

setScreen:
	w.ScreenWindow = sc.Name
	w.PhysDPI = sc.PhysicalDPI
	w.DevicePixelRatio = sc.DevicePixelRatio
	if w.LogDPI == 0 {
		w.LogDPI = sc.LogicalDPI
	}
	return sc
}

func (w *Window) GetScreenOverlap() *system.Screen {
	var wgeom image.Rectangle
	wgeom.Min.X, wgeom.Min.Y = w.Glw.GetPos()
	var sz image.Point
	sz.X, sz.Y = w.Glw.GetSize()
	wgeom.Max = wgeom.Min.Add(sz)

	var csc *system.Screen
	var ovlp int
	for _, sc := range TheApp.Screens {
		isect := sc.Geometry.Intersect(wgeom).Size()
		ov := isect.X * isect.Y
		if ov > ovlp || ovlp == 0 {
			csc = sc
			ovlp = ov
		}
	}
	return csc
}

func (w *Window) Position() image.Point {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	if w.Glw == nil {
		return w.Pos
	}
	var ps image.Point
	ps.X, ps.Y = w.Glw.GetPos()
	w.Pos = ps
	return ps
}

func (w *Window) SetTitle(title string) {
	if w.IsClosed() {
		return
	}
	w.Titl = title
	w.App.RunOnMain(func() {
		if w.Glw == nil {
			return
		}
		w.Glw.SetTitle(title)
	})
}

func (w *Window) SetIcon(images []image.Image) {
	if w.IsClosed() {
		return
	}
	w.App.RunOnMain(func() {
		if w.Glw == nil {
			return
		}
		w.Glw.SetIcon(images)
	})
}

func (w *Window) SetWinSize(sz image.Point) {
	if w.IsClosed() || w.Is(system.Fullscreen) {
		return
	}

	w.App.RunOnMain(func() {
		if w.Glw == nil {
			return
		}
		w.Glw.SetSize(sz.X, sz.Y)
	})
}

func (w *Window) SetPos(pos image.Point) {
	if w.IsClosed() || w.Is(system.Fullscreen) {
		return
	}

	w.App.RunOnMain(func() {
		if w.Glw == nil {
			return
		}
		w.Glw.SetPos(pos.X, pos.Y)
	})
}

func (w *Window) SetGeom(pos image.Point, sz image.Point) {
	if w.IsClosed() || w.Is(system.Fullscreen) {
		return
	}
	sc := w.Screen()
	sz = sc.WinSizeFromPix(sz)

	w.App.RunOnMain(func() {
		if w.Glw == nil {
			return
		}
		w.Glw.SetSize(sz.X, sz.Y)
		w.Glw.SetPos(pos.X, pos.Y)
	})
}

func (w *Window) Show() {
	if w.IsClosed() {
		return
	}

	w.App.RunOnMain(func() {
		if w.Glw == nil {
			return
		}
		w.Glw.Show()
	})
}

func (w *Window) Raise() {
	if w.IsClosed() {
		return
	}

	w.App.RunOnMain(func() {
		if w.Glw == nil {
			return
		}
		if w.Is(system.Minimized) {
			w.Glw.Restore()
		} else {
			w.Glw.Focus()
		}
	})
}

func (w *Window) Minimize() {
	if w.IsClosed() {
		return
	}

	w.App.RunOnMain(func() {
		if w.Glw == nil {
			return
		}
		w.Glw.Iconify()
	})
}

func (w *Window) Close() {
	if w == nil {
		return
	}
	w.Window.Close()

	w.Mu.Lock()
	defer w.Mu.Unlock()

	w.App.RunOnMain(func() {
		vk.DeviceWaitIdle(w.Draw.Surf.Device.Device)
		if w.DestroyGPUFunc != nil {
			w.DestroyGPUFunc()
		}
		w.Draw.Destroy()
		w.Draw.Surf.Destroy()
		w.Glw.Destroy()
		w.Glw = nil
		w.Draw = nil
	})
}

func (w *Window) SetMousePos(x, y float64) {
	if !w.IsVisible() {
		return
	}
	w.Mu.Lock()
	defer w.Mu.Unlock()
	if TheApp.Platform() == system.MacOS {
		w.Glw.SetCursorPos(x/float64(w.DevicePixelRatio), y/float64(w.DevicePixelRatio))
	} else {
		w.Glw.SetCursorPos(x, y)
	}
}

func (w *Window) SetCursorEnabled(enabled, raw bool) {
	w.CursorEnabled = enabled
	if enabled {
		w.Glw.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	} else {
		w.Glw.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		if raw && glfw.RawMouseMotionSupported() {
			w.Glw.SetInputMode(glfw.RawMouseMotion, glfw.True)
		}
	}
}

func (w *Window) Moved(gw *glfw.Window, x, y int) {
	w.Mu.Lock()
	w.Pos = image.Pt(x, y)
	w.Mu.Unlock()

	w.Screen()
	w.UpdateFullscreen()
	w.Event.Window(events.WinMove)
}

func (w *Window) WinResized(gw *glfw.Window, width, height int) {
	w.UpdateFullscreen()
	w.UpdateGeom()
}

func (w *Window) UpdateFullscreen() {
	w.Flgs.SetFlag(w.Glw.GetAttrib(glfw.Maximized) == glfw.True, system.Fullscreen)
}

func (w *Window) UpdateGeom() {
	w.Mu.Lock()
	cursc := w.ScreenWindow
	w.Mu.Unlock()
	sc := w.Screen()
	w.Mu.Lock()
	var wsz image.Point
	wsz.X, wsz.Y = w.Glw.GetSize()

	w.WnSize = wsz
	var fbsz image.Point
	fbsz.X, fbsz.Y = w.Glw.GetFramebufferSize()
	w.PixSize = fbsz
	w.PhysDPI = sc.PhysicalDPI
	w.LogDPI = sc.LogicalDPI
	w.Mu.Unlock()

	if cursc != w.ScreenWindow {
		if MonitorDebug {
			log.Printf("desktop.Window.UpdateGeom: %v: got new screen: %v (was: %v)\n", w.Nm, w.ScreenWindow, cursc)
		}
	}
	w.Event.WindowResize()
}

func (w *Window) FbResized(gw *glfw.Window, width, height int) {
	fbsz := image.Point{width, height}
	if w.PixSize != fbsz {
		w.UpdateGeom()
	}
}

func (w *Window) OnCloseReq(gw *glfw.Window) {
	go w.CloseReq()
}

func (w *Window) Focused(gw *glfw.Window, focused bool) {
	if focused {
		w.Event.Window(events.WinFocus)
	} else {
		w.Event.Last.MousePos = image.Point{-1, -1}
		w.Event.Window(events.WinFocusLost)
	}
}

func (w *Window) Iconify(gw *glfw.Window, iconified bool) {
	w.Flgs.SetFlag(iconified, system.Minimized)
	if iconified {
		w.Event.Window(events.WinMinimize)
	} else {
		w.Screen()
		w.Event.Window(events.WinMinimize)
	}
}
