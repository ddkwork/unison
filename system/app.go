package system

import "cogentcore.org/core/styles"

//go:generate core generate

var (
	TheApp               App
	AppVersion           = "dev"
	CoreVersion          = "dev"
	ReservedWebShortcuts = []string{"Command+r", "Command+Shift+r", "Command+w", "Command+t", "Command+Shift+t", "Command+q", "Command+n", "Command+m", "Command+l", "Command+h", "Command+Shift+j", "Command+Alt+j", "Command+Alt+∆", "Command+1", "Command+2", "Command+3", "Command+4", "Command+5", "Command+6", "Command+7", "Command+8", "Command+9"}
)

type App interface {
	Platform() Platforms
	SystemPlatform() Platforms
	SystemInfo() string
	Name() string
	SetName(name string)
	GetScreens()
	NScreens() int
	Screen(n int) *Screen
	ScreenByName(name string) *Screen
	NWindows() int
	Window(win int) Window
	WindowByName(name string) Window
	WindowInFocus() Window
	ContextWindow() Window
	NewWindow(opts *NewWindowOptions) (Window, error)
	RemoveWindow(win Window)
	Clipboard(win Window) Clipboard
	Cursor(win Window) Cursor
	DataDir() string
	AppDataDir() string
	CogentCoreDataDir() string
	OpenURL(url string)
	OpenFiles() []string
	SetQuitReqFunc(fun func())
	AddQuitCleanFunc(fun func())
	QuitReq()
	IsQuitting() bool
	QuitClean() bool
	Quit()
	MainLoop()
	RunOnMain(f func())
	SendEmptyEvent()
	ShowVirtualKeyboard(typ styles.VirtualKeyboards)
	HideVirtualKeyboard()
	IsDark() bool
}

var OnSystemWindowCreated chan struct{}

type Platforms int32

const (
	MacOS Platforms = iota
	Linux
	Windows
	IOS
	Android
	Web
	Offscreen
)

func (p Platforms) IsMobile() bool {
	return p == IOS || p == Android || p == Web || p == Offscreen
}
