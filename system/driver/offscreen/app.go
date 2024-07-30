package offscreen

import (
	"image"
	"os"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/events"
	"cogentcore.org/core/system/driver/base"
	"github.com/richardwilkes/unison/system"
)

func Init() {
	TheApp.Draw = &Drawer{}
	TheApp.GetScreens()
	TheApp.TempDataDir = errors.Log1(os.MkdirTemp("", "cogent-core-offscreen-data-dir-"))
	base.Init(TheApp, &TheApp.App)
}

var TheApp = &App{AppSingle: base.NewAppSingle[*Drawer, *Window]()}

type App struct {
	base.AppSingle[*Drawer, *Window]
	TempDataDir string
}

func (a *App) NewWindow(opts *system.NewWindowOptions) (system.Window, error) {
	defer func() { system.HandleRecover(recover()) }()
	a.Win = &Window{base.NewWindowSingle(a, opts)}
	a.Win.This = a.Win
	a.Scrn.PixSize = opts.Size
	a.GetScreens()

	a.Event.WindowResize()
	a.Event.Window(events.WinShow)
	a.Event.Window(events.ScreenUpdate)
	a.Event.Window(events.WinFocus)

	go a.Win.WinLoop()
	return a.Win, nil
}

func (a *App) GetScreens() {
	if a.Scrn.PixSize.X == 0 {
		a.Scrn.PixSize.X = 800
	}
	if a.Scrn.PixSize.Y == 0 {
		a.Scrn.PixSize.Y = 600
	}

	a.Scrn.DevicePixelRatio = 1
	a.Scrn.Geometry.Max = a.Scrn.PixSize
	dpi := float32(160)
	a.Scrn.PhysicalDPI = dpi
	a.Scrn.LogicalDPI = dpi

	if system.InitScreenLogicalDPIFunc != nil {
		system.InitScreenLogicalDPIFunc()
	}

	physX := 25.4 * float32(a.Scrn.PixSize.X) / dpi
	physY := 25.4 * float32(a.Scrn.PixSize.Y) / dpi
	a.Scrn.PhysicalSize = image.Pt(int(physX), int(physY))

	a.Draw.Image = image.NewRGBA(image.Rectangle{Max: a.Scrn.PixSize})
}

func (a *App) QuitClean() bool {
	if a.TempDataDir != "" {
		errors.Log(os.RemoveAll(a.TempDataDir))
	}
	return a.AppSingle.QuitClean()
}

func (a *App) DataDir() string {
	return a.TempDataDir
}

func (a *App) Platform() system.Platforms {
	return system.Offscreen
}
