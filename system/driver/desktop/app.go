package desktop

import (
	"log"
	"runtime"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/system"
	"cogentcore.org/core/system/driver/base"
	"cogentcore.org/core/vgpu"
	"cogentcore.org/core/vgpu/vdraw"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison/internal/glfw"

	vk "github.com/goki/vulkan"
)

func Init() {
	runtime.LockOSThread()
	TheApp.InitVk()
	base.Init(TheApp, &TheApp.App)
}

var TheApp = &App{AppMulti: base.NewAppMulti[*Window]()}

type App struct {
	base.AppMulti[*Window]
	GPU      *vgpu.GPU
	ShareWin *glfw.Window
}

func (a *App) SendEmptyEvent() {
	glfw.PostEmptyEvent()
}

func (a *App) MainLoop() {
	a.MainQueue = make(chan base.FuncRun)
	a.MainDone = make(chan struct{})
	for {
		select {
		case <-a.MainDone:
			glfw.Terminate()
			return
		case f := <-a.MainQueue:
			f.F()
			if f.Done != nil {
				f.Done <- struct{}{}
			}
		default:
			glfw.WaitEvents()
		}
	}
}

func (a *App) InitVk() {
	if mylog.Check(glfw.Init()); err != nil {
		log.Fatalln("system/driver/desktop failed to initialize glfw:", err)
	}
	vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	vk.Init()
	glfw.SetMonitorCallback(a.MonitorChange)

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)

	a.ShareWin = mylog.Check2(glfw.CreateWindow(16, 16, "Share Window", nil, nil))

	winext := a.ShareWin.GetRequiredInstanceExtensions()
	a.GPU = vgpu.NewGPU()
	a.GPU.AddInstanceExt(winext...)
	a.GPU.Config(a.Name())

	a.GetScreens()
}

func (a *App) NewWindow(opts *system.NewWindowOptions) (system.Window, error) {
	if len(a.Windows) == 0 && system.InitScreenLogicalDPIFunc != nil {
		if MonitorDebug {
			log.Println("app first new window calling InitScreenLogicalDPIFunc")
		}
		system.InitScreenLogicalDPIFunc()
	}

	sc := a.Screens[0]

	if opts == nil {
		opts = &system.NewWindowOptions{}
	}
	opts.Fixup()

	var glw *glfw.Window

	a.RunOnMain(func() {
		glw = mylog.Check2(NewGlfwWindow(opts, sc))
	})

	w := &Window{
		WindowMulti:  base.NewWindowMulti[*App, *vdraw.Drawer](a, opts),
		Glw:          glw,
		ScreenWindow: sc.Name,
	}
	w.This = w
	w.Draw = &vdraw.Drawer{}

	a.RunOnMain(func() {
		surfPtr := errors.Log1(glw.CreateWindowSurface(a.GPU.Instance, nil))
		sf := vgpu.NewSurface(a.GPU, vk.SurfaceFromPointer(surfPtr))
		w.Draw.YIsDown = true
		w.Draw.ConfigSurface(sf, vgpu.MaxTexturesPerSet)
	})

	a.Mu.Lock()
	a.Windows = append(a.Windows, w)
	a.Mu.Unlock()

	glw.SetPosCallback(w.Moved)
	glw.SetSizeCallback(w.WinResized)
	glw.SetFramebufferSizeCallback(w.FbResized)
	glw.SetCloseCallback(w.OnCloseReq)

	glw.SetFocusCallback(w.Focused)
	glw.SetIconifyCallback(w.Iconify)

	glw.SetKeyCallback(w.KeyEvent)
	glw.SetCharModsCallback(w.CharEvent)
	glw.SetMouseButtonCallback(w.MouseButtonEvent)
	glw.SetScrollCallback(w.ScrollEvent)
	glw.SetCursorPosCallback(w.CursorPosEvent)
	glw.SetCursorEnterCallback(w.CursorEnterEvent)
	glw.SetDropCallback(w.DropEvent)

	w.Show()
	a.RunOnMain(func() {
		w.UpdateGeom()
	})

	go w.WinLoop()

	return w, nil
}

func (a *App) Clipboard(win system.Window) system.Clipboard {
	a.Mu.Lock()
	a.CtxWindow = win.(*Window)
	a.Mu.Unlock()
	return TheClipboard
}

func (a *App) Cursor(win system.Window) system.Cursor {
	a.Mu.Lock()
	a.CtxWindow = win.(*Window)
	a.Mu.Unlock()
	return TheCursor
}
