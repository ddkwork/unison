package base

import (
	"fmt"
	"image"
	"sync"
	"time"

	"cogentcore.org/core/events"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/system"
)

type Window[A system.App] struct {
	This           system.Window `display:"-"`
	App            A
	Mu             sync.Mutex    `display:"-"`
	WinClose       chan struct{} `display:"-"`
	CloseReqFunc   func(win system.Window)
	CloseCleanFunc func(win system.Window)
	Nm             string             `label:"Name"`
	Titl           string             `label:"Title"`
	Flgs           system.WindowFlags `label:"Flags" table:"-"`
	FPS            int
	DestroyGPUFunc func()
	CursorEnabled  bool
}

func NewWindow[A system.App](a A, opts *system.NewWindowOptions) Window[A] {
	return Window[A]{
		WinClose:      make(chan struct{}),
		App:           a,
		Titl:          opts.GetTitle(),
		Flgs:          opts.Flags,
		FPS:           60,
		CursorEnabled: true,
	}
}

func (w *Window[A]) WinLoop() {
	defer func() { system.HandleRecover(recover()) }()

	var winPaint *time.Ticker
	if w.FPS > 0 {
		winPaint = time.NewTicker(time.Second / time.Duration(w.FPS))
	} else {
		winPaint = &time.Ticker{C: make(chan time.Time)}
	}
outer:
	for {
		select {
		case <-w.WinClose:
			winPaint.Stop()
			break outer
		case <-winPaint.C:
			if w.This.IsClosed() {
				fmt.Println("win IsClosed in paint:", w.Name())
				break outer
			}
			w.This.Events().WindowPaint()
		}
	}
}

func (w *Window[A]) Lock() bool {
	if w.This.IsClosed() {
		return false
	}
	w.Mu.Lock()
	return true
}

func (w *Window[A]) Unlock() {
	w.Mu.Unlock()
}

func (w *Window[A]) Name() string {
	return w.Nm
}

func (w *Window[A]) SetName(name string) {
	w.Nm = name
}

func (w *Window[A]) Title() string {
	return w.Titl
}

func (w *Window[A]) SetTitle(title string) {
	if w.This.IsClosed() {
		return
	}
	w.Titl = title
}

func (w *Window[A]) SetIcon(images []image.Image) {
}

func (w *Window[A]) Flags() system.WindowFlags {
	return w.Flgs
}

func (w *Window[A]) Is(flag system.WindowFlags) bool {
	return w.Flgs.HasFlag(flag)
}

func (w *Window[A]) IsClosed() bool {
	return w == nil || w.This == nil || w.This.Drawer() == nil
}

func (w *Window[A]) IsVisible() bool {
	return !w.This.IsClosed() && !w.Is(system.Minimized)
}

func (w *Window[A]) SetFPS(fps int) {
	w.FPS = fps
}

func (w *Window[A]) SetDestroyGPUResourcesFunc(f func()) {
	w.DestroyGPUFunc = f
}

func (w *Window[A]) RenderGeom() math32.Geom2DInt {
	return math32.Geom2DInt{Size: w.This.Size()}
}

func (w *Window[A]) SetCloseReqFunc(fun func(win system.Window)) {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	w.CloseReqFunc = fun
}

func (w *Window[A]) SetCloseCleanFunc(fun func(win system.Window)) {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	w.CloseCleanFunc = fun
}

func (w *Window[A]) CloseReq() {
	if w.CloseReqFunc != nil {
		w.CloseReqFunc(w.This)
	} else {
		w.This.Close()
	}
}

func (w *Window[A]) CloseClean() {
	if w.CloseCleanFunc != nil {
		w.CloseCleanFunc(w.This)
	}
}

func (w *Window[A]) Close() {
	w.This.Events().Window(events.WinClose)

	w.WinClose <- struct{}{}

	w.Mu.Lock()
	defer w.Mu.Unlock()

	w.CloseClean()
	w.App.RemoveWindow(w.This)
}

func (w *Window[A]) SetCursorEnabled(enabled, raw bool) {
	w.CursorEnabled = enabled
}

func (w *Window[A]) IsCursorEnabled() bool {
	return w.CursorEnabled
}

func (w *Window[A]) SetMousePos(x, y float64) {
}

func (w *Window[A]) SetTitleBarIsDark(isDark bool) {
}
