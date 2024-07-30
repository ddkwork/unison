package base

import (
	"image"

	"cogentcore.org/core/events"
	"github.com/richardwilkes/unison/system"
)

type WindowMulti[A system.App, D system.Drawer] struct {
	Window[A]
	Event            events.Source `label:"Event manger"`
	Draw             D             `label:"Drawer"`
	Pos              image.Point   `label:"Position"`
	WnSize           image.Point   `label:"Window manager size"`
	PixSize          image.Point   `label:"Pixel size"`
	DevicePixelRatio float32
	PhysDPI          float32 `label:"Physical DPI"`
	LogDPI           float32 `label:"Logical DPI"`
}

func NewWindowMulti[A system.App, D system.Drawer](a A, opts *system.NewWindowOptions) WindowMulti[A, D] {
	return WindowMulti[A, D]{
		Window: NewWindow(a, opts),
	}
}

func (w *WindowMulti[A, D]) Events() *events.Source {
	return &w.Event
}

func (w *WindowMulti[A, D]) Drawer() system.Drawer {
	return w.Draw
}

func (w *WindowMulti[A, D]) Size() image.Point {
	return w.PixSize
}

func (w *WindowMulti[A, D]) WinSize() image.Point {
	return w.WnSize
}

func (w *WindowMulti[A, D]) Position() image.Point {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	return w.Pos
}

func (w *WindowMulti[A, D]) PhysicalDPI() float32 {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	return w.PhysDPI
}

func (w *WindowMulti[A, D]) LogicalDPI() float32 {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	return w.LogDPI
}

func (w *WindowMulti[A, D]) SetLogicalDPI(dpi float32) {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	w.LogDPI = dpi
}

func (w *WindowMulti[A, D]) SetWinSize(sz image.Point) {
	if w.This.IsClosed() {
		return
	}
	w.WnSize = sz
}

func (w *WindowMulti[A, D]) SetSize(sz image.Point) {
	if w.This.IsClosed() {
		return
	}
	sc := w.This.Screen()
	sz = sc.WinSizeFromPix(sz)
	w.SetWinSize(sz)
}

func (w *WindowMulti[A, D]) SetPos(pos image.Point) {
	if w.This.IsClosed() {
		return
	}
	w.Pos = pos
}

func (w *WindowMulti[A, D]) SetGeom(pos image.Point, sz image.Point) {
	if w.This.IsClosed() {
		return
	}
	sc := w.This.Screen()
	sz = sc.WinSizeFromPix(sz)
	w.SetWinSize(sz)
	w.Pos = pos
}

func (w *WindowMulti[A, D]) IsVisible() bool {
	return w.Window.IsVisible() && w.App.NScreens() != 0
}
