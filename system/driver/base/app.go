package base

import (
	"os"
	"path/filepath"
	"sync"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/events/key"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/system"
)

type App struct {
	This           system.App    `display:"-"`
	Mu             sync.Mutex    `display:"-"`
	MainQueue      chan FuncRun  `display:"-"`
	MainDone       chan struct{} `display:"-"`
	Nm             string        `label:"Name"`
	OpenFls        []string      `label:"Open files"`
	Quitting       bool
	QuitReqFunc    func()
	QuitCleanFuncs []func()
	Dark           bool
}

func Init(a system.App, ab *App) {
	ab.This = a
	system.TheApp = a
	key.SystemPlatform = a.SystemPlatform().String()
}

func (a *App) MainLoop() {
	a.MainQueue = make(chan FuncRun)
	a.MainDone = make(chan struct{})
	for {
		select {
		case <-a.MainDone:
			return
		case f := <-a.MainQueue:
			f.F()
			if f.Done != nil {
				f.Done <- struct{}{}
			}
		}
	}
}

func (a *App) RunOnMain(f func()) {
	if a.MainQueue == nil {
		f()
		return
	}
	a.This.SendEmptyEvent()
	done := make(chan struct{})
	a.MainQueue <- FuncRun{F: f, Done: done}
	<-done
	a.This.SendEmptyEvent()
}

func (a *App) SendEmptyEvent() {
}

func (a *App) StopMain() {
	a.MainDone <- struct{}{}
}

func (a *App) SystemPlatform() system.Platforms {
	return a.This.Platform()
}

func (a *App) SystemInfo() string {
	return ""
}

func (a *App) Name() string {
	return a.Nm
}

func (a *App) SetName(name string) {
	a.Nm = name
}

func (a *App) OpenFiles() []string {
	return a.OpenFls
}

func (a *App) AppDataDir() string {
	pdir := filepath.Join(system.TheApp.DataDir(), a.Name())
	errors.Log(os.MkdirAll(pdir, 0755))
	return pdir
}

func (a *App) CogentCoreDataDir() string {
	pdir := filepath.Join(a.This.DataDir(), "Cogent Core")
	errors.Log(os.MkdirAll(pdir, 0755))
	return pdir
}

func (a *App) SetQuitReqFunc(fun func()) {
	a.QuitReqFunc = fun
}

func (a *App) AddQuitCleanFunc(fun func()) {
	a.QuitCleanFuncs = append(a.QuitCleanFuncs, fun)
}

func (a *App) QuitReq() {
	if a.Quitting {
		return
	}
	if a.QuitReqFunc != nil {
		a.QuitReqFunc()
	} else {
		a.Quit()
	}
}

func (a *App) IsQuitting() bool {
	return a.Quitting
}

func (a *App) Quit() {
	if a.Quitting {
		return
	}
	a.Quitting = true
	if a.This.QuitClean() {
		a.StopMain()
	} else {
		a.Quitting = false
	}
}

func (a *App) IsDark() bool {
	return a.Dark
}

func (a *App) GetScreens() {
}

func (a *App) OpenURL(url string) {
}

func (a *App) Clipboard(win system.Window) system.Clipboard {
	return &system.ClipboardBase{}
}

func (a *App) Cursor(win system.Window) system.Cursor {
	return &system.CursorBase{}
}

func (a *App) ShowVirtualKeyboard(typ styles.VirtualKeyboards) {
}

func (a *App) HideVirtualKeyboard() {
}
