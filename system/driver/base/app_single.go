package base

import (
	"image"

	"cogentcore.org/core/events"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/system"
)

type AppSingle[D system.Drawer, W system.Window] struct {
	App
	Event  events.Source  `label:"Events"`
	Draw   D              `label:"Drawer"`
	Win    W              `label:"Window"`
	Scrn   *system.Screen `label:"Screen"`
	Insets styles.Sides[int]
}

type AppSingler interface {
	system.App

	Events() *events.Source

	Drawer() system.Drawer

	RenderGeom() math32.Geom2DInt
}

func NewAppSingle[D system.Drawer, W system.Window]() AppSingle[D, W] {
	return AppSingle[D, W]{
		Scrn: &system.Screen{},
	}
}

func (a *AppSingle[D, W]) Events() *events.Source {
	return &a.Event
}

func (a *AppSingle[D, W]) Drawer() system.Drawer {
	return a.Draw
}

func (a *AppSingle[D, W]) RenderGeom() math32.Geom2DInt {
	pos := image.Pt(a.Insets.Left, a.Insets.Top)
	return math32.Geom2DInt{
		Pos:  pos,
		Size: a.Scrn.PixSize.Sub(pos).Sub(image.Pt(a.Insets.Right, a.Insets.Bottom)),
	}
}

func (a *AppSingle[D, W]) NScreens() int {
	if a.Scrn != nil {
		return 1
	}
	return 0
}

func (a *AppSingle[D, W]) Screen(n int) *system.Screen {
	if n == 0 {
		return a.Scrn
	}
	return nil
}

func (a *AppSingle[D, W]) ScreenByName(name string) *system.Screen {
	if a.Scrn.Name == name {
		return a.Scrn
	}
	return nil
}

func (a *AppSingle[D, W]) NWindows() int {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	if system.Window(a.Win) != nil {
		return 1
	}
	return 0
}

func (a *AppSingle[D, W]) Window(win int) system.Window {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	if win == 0 {
		return a.Win
	}
	return nil
}

func (a *AppSingle[D, W]) WindowByName(name string) system.Window {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	if a.Win.Name() == name {
		return a.Win
	}
	return nil
}

func (a *AppSingle[D, W]) WindowInFocus() system.Window {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	if a.Win.Is(system.Focused) {
		return a.Win
	}
	return nil
}

func (a *AppSingle[D, W]) ContextWindow() system.Window {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	return a.Win
}

func (a *AppSingle[D, W]) RemoveWindow(w system.Window) {
}

func (a *AppSingle[D, W]) QuitClean() bool {
	a.Quitting = true
	for _, qf := range a.QuitCleanFuncs {
		qf()
	}
	a.Mu.Lock()
	defer a.Mu.Unlock()
	a.Win.Close()
	return a.Win.IsClosed()
}
