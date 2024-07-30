package base

import (
	"slices"

	"github.com/richardwilkes/unison/system"
)

type AppMulti[W system.Window] struct {
	App
	Windows    []W
	Screens    []*system.Screen
	AllScreens []*system.Screen
	CtxWindow  W `label:"Context window"`
}

func NewAppMulti[W system.Window]() AppMulti[W] {
	return AppMulti[W]{}
}

func (a *AppMulti[W]) NScreens() int {
	return len(a.Screens)
}

func (a *AppMulti[W]) Screen(n int) *system.Screen {
	if n < len(a.Screens) {
		return a.Screens[n]
	}
	return nil
}

func (a *AppMulti[W]) ScreenByName(name string) *system.Screen {
	for _, sc := range a.Screens {
		if sc.Name == name {
			return sc
		}
	}
	return nil
}

func (a *AppMulti[W]) NWindows() int {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	return len(a.Windows)
}

func (a *AppMulti[W]) Window(win int) system.Window {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	if win < len(a.Windows) {
		return a.Windows[win]
	}
	return nil
}

func (a *AppMulti[W]) WindowByName(name string) system.Window {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	for _, win := range a.Windows {
		if win.Name() == name {
			return win
		}
	}
	return nil
}

func (a *AppMulti[W]) WindowInFocus() system.Window {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	for _, win := range a.Windows {
		if win.Is(system.Focused) {
			return win
		}
	}
	return nil
}

func (a *AppMulti[W]) ContextWindow() system.Window {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	return a.CtxWindow
}

func (a *AppMulti[W]) RemoveWindow(w system.Window) {
	a.Windows = slices.DeleteFunc(a.Windows, func(ew W) bool {
		return system.Window(ew) == w
	})
}

func (a *AppMulti[W]) QuitClean() bool {
	for _, qf := range a.QuitCleanFuncs {
		qf()
	}
	nwin := len(a.Windows)
	for i := nwin - 1; i >= 0; i-- {
		win := a.Windows[i]
		win.CloseReq()
	}
	return len(a.Windows) == 0
}
