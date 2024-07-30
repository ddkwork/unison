package base

import (
	"image"

	"cogentcore.org/core/events"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/system"
)

type WindowSingle[A AppSingler] struct {
	Window[A]
}

func NewWindowSingle[A AppSingler](a A, opts *system.NewWindowOptions) WindowSingle[A] {
	return WindowSingle[A]{
		NewWindow(a, opts),
	}
}

func (w *WindowSingle[A]) Events() *events.Source {
	return w.App.Events()
}

func (w *WindowSingle[A]) Drawer() system.Drawer {
	return w.App.Drawer()
}

func (w *WindowSingle[A]) Screen() *system.Screen {
	return w.App.Screen(0)
}

func (w *WindowSingle[A]) Size() image.Point {
	return w.Screen().PixSize
}

func (w *WindowSingle[A]) WinSize() image.Point {
	return w.Screen().PixSize
}

func (w *WindowSingle[A]) Position() image.Point {
	return image.Point{}
}

func (w *WindowSingle[A]) PhysicalDPI() float32 {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	return w.Screen().PhysicalDPI
}

func (w *WindowSingle[A]) LogicalDPI() float32 {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	return w.Screen().LogicalDPI
}

func (w *WindowSingle[A]) SetLogicalDPI(dpi float32) {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	w.Screen().LogicalDPI = dpi
}

func (w *WindowSingle[A]) SetWinSize(sz image.Point) {
	if w.This.IsClosed() {
		return
	}
	w.Screen().PixSize = sz
}

func (w *WindowSingle[A]) SetSize(sz image.Point) {
	if w.This.IsClosed() {
		return
	}
	w.Screen().PixSize = sz
}

func (w *WindowSingle[A]) SetPos(pos image.Point) {
}

func (w *WindowSingle[A]) SetGeom(pos image.Point, sz image.Point) {
	if w.This.IsClosed() {
		return
	}
	w.Screen().PixSize = sz
}

func (w *WindowSingle[A]) Raise() {
}

func (w *WindowSingle[A]) Minimize() {
}

func (w *WindowSingle[A]) RenderGeom() math32.Geom2DInt {
	return w.App.RenderGeom()
}
