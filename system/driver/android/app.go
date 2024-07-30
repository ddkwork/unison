//go:build android

package android

import (
	"log"

	"cogentcore.org/core/events"
	"cogentcore.org/core/system/driver/base"
	"cogentcore.org/core/vgpu"
	"cogentcore.org/core/vgpu/vdraw"
	"github.com/ddkwork/golibrary/mylog"
	vk "github.com/goki/vulkan"
	"github.com/richardwilkes/unison/system"
)

func Init() {
	system.OnSystemWindowCreated = make(chan struct{})
	TheApp.InitVk()
	base.Init(TheApp, &TheApp.App)
}

var TheApp = &App{AppSingle: base.NewAppSingle[*vdraw.Drawer, *Window]()}

type App struct {
	base.AppSingle[*vdraw.Drawer, *Window]

	GPU *vgpu.GPU

	Winptr uintptr
}

func (a *App) InitVk() {
	mylog.Check(vk.SetDefaultGetInstanceProcAddr())
	mylog.Check(vk.Init())

	winext := vk.GetRequiredInstanceExtensions()
	a.GPU = vgpu.NewGPU()
	a.GPU.AddInstanceExt(winext...)
	a.GPU.Config(a.Name())
}

func (a *App) DestroyVk() {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	vk.DeviceWaitIdle(a.Draw.Surf.Device.Device)
	a.Draw.Destroy()
	a.Draw.Surf.Destroy()
	a.Draw = nil
}

func (a *App) FullDestroyVk() {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	a.GPU.Destroy()
}

func (a *App) NewWindow(opts *system.NewWindowOptions) (system.Window, error) {
	defer func() { system.HandleRecover(recover()) }()

	a.Mu.Lock()
	defer a.Mu.Unlock()
	a.Win = &Window{base.NewWindowSingle(a, opts)}
	a.Win.This = a.Win
	a.Event.Window(events.WinShow)
	a.Event.Window(events.WinFocus)

	go a.Win.WinLoop()

	return a.Win, nil
}

func (a *App) SetSystemWindow(winptr uintptr) error {
	defer func() { system.HandleRecover(recover()) }()
	var vsf vk.Surface

	ret := vk.CreateWindowSurface(a.GPU.Instance, winptr, nil, &vsf)
	if mylog.Check(vk.Error(ret)); err != nil {
		return err
	}
	sf := vgpu.NewSurface(a.GPU, vsf)

	sys := a.GPU.NewGraphicsSystem(a.Name(), &sf.Device)
	sys.ConfigRender(&sf.Format, vgpu.UndefType)
	sf.SetRender(&sys.Render)

	sys.Config()
	a.Draw = &vdraw.Drawer{
		Sys:     *sys,
		YIsDown: true,
	}

	a.Draw.ConfigSurface(sf, vgpu.MaxTexturesPerSet)

	a.Winptr = winptr

	if a.Win != nil {
		a.Event.Window(events.WinShow)
		a.Event.Window(events.ScreenUpdate)
	}
	return nil
}

func (a *App) DataDir() string {
	return "/data/data"
}

func (a *App) Platform() system.Platforms {
	return system.Android
}

func (a *App) OpenURL(url string) {
}

func (a *App) Clipboard(win system.Window) system.Clipboard {
	return &Clipboard{}
}
