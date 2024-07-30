//go:build js

package web

import (
	"fmt"
	"image"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"syscall/js"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/base/fileinfo"
	"cogentcore.org/core/events"
	"cogentcore.org/core/events/key"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/system/driver/base"
	"cogentcore.org/core/system/driver/web/jsfs"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison/system"
)

func Init() {
	TheApp.Draw = &Drawer{}
	mylog.Check(os.Setenv("HOME", "/home/me"))

	fs := mylog.Check2(jsfs.Config(js.Global().Get("fs")))

	TheApp.SetSystemWindow()

	base.Init(TheApp, &TheApp.App)
}

var TheApp = &App{AppSingle: base.NewAppSingle[*Drawer, *Window]()}

type App struct {
	base.AppSingle[*Drawer, *Window]

	UnderlyingPlatform system.Platforms

	KeyMods key.Modifiers
}

func (a *App) NewWindow(opts *system.NewWindowOptions) (system.Window, error) {
	defer func() { system.HandleRecover(recover()) }()

	a.Win = &Window{base.NewWindowSingle(a, opts)}
	a.Win.This = a.Win

	go a.Win.WinLoop()

	return a.Win, nil
}

func (a *App) SetSystemWindow() {
	defer func() { system.HandleRecover(recover()) }()

	a.AddEventListeners()

	ua := js.Global().Get("navigator").Get("userAgent").String()
	a.UnderlyingPlatform = UserAgentToOS(ua)

	a.Resize()
	a.Event.Window(events.WinShow)
	a.Event.Window(events.ScreenUpdate)
	a.Event.Window(events.WinFocus)
}

func UserAgentToOS(ua string) system.Platforms {
	lua := strings.ToLower(ua)
	switch {
	case strings.Contains(lua, "android"):
		return system.Android
	case strings.Contains(lua, "ipad"),
		strings.Contains(lua, "iphone"),
		strings.Contains(lua, "ipod"):
		return system.IOS
	case strings.Contains(lua, "mac"):
		return system.MacOS
	case strings.Contains(lua, "win"):
		return system.Windows
	default:
		return system.Linux
	}
}

func (a *App) Resize() {
	a.Scrn.DevicePixelRatio = float32(js.Global().Get("devicePixelRatio").Float())
	dpi := 160 * a.Scrn.DevicePixelRatio
	a.Scrn.PhysicalDPI = dpi
	a.Scrn.LogicalDPI = dpi

	if system.InitScreenLogicalDPIFunc != nil {
		system.InitScreenLogicalDPIFunc()
	}

	vv := js.Global().Get("visualViewport")
	w, h := vv.Get("width").Int(), vv.Get("height").Int()
	sz := image.Pt(w, h)
	a.Scrn.Geometry.Max = sz
	a.Scrn.PixSize = image.Pt(int(math32.Ceil(float32(sz.X)*a.Scrn.DevicePixelRatio)), int(math32.Ceil(float32(sz.Y)*a.Scrn.DevicePixelRatio)))
	physX := 25.4 * float32(w) / dpi
	physY := 25.4 * float32(h) / dpi
	a.Scrn.PhysicalSize = image.Pt(int(physX), int(physY))

	canvas := js.Global().Get("document").Call("getElementById", "app")
	canvas.Set("width", a.Scrn.PixSize.X)
	canvas.Set("height", a.Scrn.PixSize.Y)

	cstyle := canvas.Get("style")
	cstyle.Set("width", fmt.Sprintf("%gpx", float32(a.Scrn.PixSize.X)/a.Scrn.DevicePixelRatio))
	cstyle.Set("height", fmt.Sprintf("%gpx", float32(a.Scrn.PixSize.Y)/a.Scrn.DevicePixelRatio))

	a.Draw.Image = image.NewRGBA(image.Rectangle{Max: a.Scrn.PixSize})

	a.Event.WindowResize()
}

func (a *App) DataDir() string {
	return "/home/me/.data"
}

func (a *App) Platform() system.Platforms {
	return system.Web
}

func (a *App) SystemPlatform() system.Platforms {
	return a.UnderlyingPlatform
}

func (a *App) SystemInfo() string {
	return "User agent: " + js.Global().Get("navigator").Get("userAgent").String()
}

func (a *App) OpenURL(url string) {
	if !strings.HasPrefix(url, "file://") {
		js.Global().Call("open", url)
		return
	}
	filename := strings.TrimPrefix(url, "file://")
	b := mylog.Check2(os.ReadFile(filename))

	jb := js.Global().Get("Uint8ClampedArray").New(len(b))
	js.CopyBytesToJS(jb, b)
	mtype, _ := mylog.Check3(fileinfo.MimeFromFile(filename))
	if errors.Log(err) != nil {
		mtype = "text/plain"
	}
	blob := js.Global().Get("Blob").New([]any{jb}, map[string]any{"type": mtype})
	objectURL := js.Global().Get("URL").Call("createObjectURL", blob)
	anchor := js.Global().Get("document").Call("createElement", "a")
	anchor.Set("style", "display: none;")
	anchor.Set("href", objectURL)
	anchor.Set("download", filepath.Base(filename))
	js.Global().Get("document").Get("body").Call("appendChild", anchor)
	anchor.Call("click")
	js.Global().Get("document").Get("body").Call("removeChild", anchor)
	js.Global().Get("URL").Call("revokeObjectURL", objectURL)
}

func (a *App) Clipboard(win system.Window) system.Clipboard {
	return TheClipboard
}

func (a *App) Cursor(win system.Window) system.Cursor {
	return TheCursor
}

func (a *App) IsDark() bool {
	return js.Global().Get("matchMedia").Truthy() &&
		js.Global().Call("matchMedia", "(prefers-color-scheme: dark)").Get("matches").Truthy()
}

func (a *App) ShowVirtualKeyboard(typ styles.VirtualKeyboards) {
	tf := js.Global().Get("document").Call("getElementById", "text-field")
	switch typ {
	case styles.KeyboardNumber, styles.KeyboardPassword, styles.KeyboardEmail, styles.KeyboardURL:
		tf.Set("type", typ.String())
	case styles.KeyboardPhone:
		tf.Set("type", "tel")
	default:
		tf.Set("type", "text")
	}
	tf.Call("focus")
}

func (a *App) HideVirtualKeyboard() {
	js.Global().Get("document").Call("getElementById", "text-field").Call("blur")
}
