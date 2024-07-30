//go:build js

package web

import (
	"image"
	"slices"
	"strings"
	"syscall/js"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/events"
	"cogentcore.org/core/events/key"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/system"
	"github.com/ddkwork/golibrary/mylog"
)

func (a *App) AddEventListeners() {
	g := js.Global()
	AddEventListener(g, "mousedown", a.OnMouseDown)

	AddEventListener(g, "touchstart", a.OnTouchStart, map[string]any{"passive": false})
	AddEventListener(g, "mouseup", a.OnMouseUp)
	AddEventListener(g, "touchend", a.OnTouchEnd)
	AddEventListener(g, "mousemove", a.OnMouseMove)

	AddEventListener(g, "touchmove", a.OnTouchMove, map[string]any{"passive": false})

	AddEventListener(g, "wheel", a.OnWheel, map[string]any{"passive": false})
	AddEventListener(g, "contextmenu", a.OnContextMenu)
	AddEventListener(g, "keydown", a.OnKeyDown)
	AddEventListener(g, "keyup", a.OnKeyUp)
	AddEventListener(g, "beforeinput", a.OnBeforeInput)
	AddEventListener(g.Get("visualViewport"), "resize", a.OnResize)
	AddEventListener(g, "blur", a.OnBlur)
}

func AddEventListener(v js.Value, nm string, fn func(this js.Value, args []js.Value) any, opts ...map[string]any) {
	if len(opts) > 0 {
		v.Call("addEventListener", nm, js.FuncOf(fn), opts[0])
	} else {
		v.Call("addEventListener", nm, js.FuncOf(fn))
	}
}

func (a *App) EventPos(e js.Value) image.Point {
	return a.EventPosFor(e.Get("clientX"), e.Get("clientY"))
}

func (a *App) EventPosFor(x, y js.Value) image.Point {
	xi, yi := x.Int(), y.Int()
	xi = int(float32(xi) * a.Scrn.DevicePixelRatio)
	yi = int(float32(yi) * a.Scrn.DevicePixelRatio)
	return image.Pt(xi, yi)
}

func (a *App) OnMouseDown(this js.Value, args []js.Value) any {
	e := args[0]
	but := e.Get("button").Int()
	var ebut events.Buttons
	switch but {
	case 0:
		ebut = events.Left
	case 1:
		ebut = events.Middle
	case 2:
		ebut = events.Right
	}
	where := a.EventPos(e)
	a.Event.MouseButton(events.MouseDown, ebut, where, a.KeyMods)
	e.Call("preventDefault")
	return nil
}

func (a *App) OnTouchStart(this js.Value, args []js.Value) any {
	e := args[0]
	touches := e.Get("changedTouches")
	for i := 0; i < touches.Length(); i++ {
		touch := touches.Index(i)
		where := a.EventPos(touch)
		a.Event.MouseButton(events.MouseDown, events.Left, where, 0)
	}
	e.Call("preventDefault")
	return nil
}

func (a *App) OnMouseUp(this js.Value, args []js.Value) any {
	e := args[0]
	but := e.Get("button").Int()
	var ebut events.Buttons
	switch but {
	case 0:
		ebut = events.Left
	case 1:
		ebut = events.Middle
	case 2:
		ebut = events.Right
	}
	where := a.EventPos(e)
	a.Event.MouseButton(events.MouseUp, ebut, where, a.KeyMods)
	e.Call("preventDefault")
	return nil
}

func (a *App) OnTouchEnd(this js.Value, args []js.Value) any {
	e := args[0]
	touches := e.Get("changedTouches")
	for i := 0; i < touches.Length(); i++ {
		touch := touches.Index(i)
		where := a.EventPos(touch)
		a.Event.MouseButton(events.MouseUp, events.Left, where, 0)
	}
	e.Call("preventDefault")
	return nil
}

func (a *App) OnMouseMove(this js.Value, args []js.Value) any {
	e := args[0]
	where := a.EventPos(e)
	a.Event.MouseMove(where)
	e.Call("preventDefault")
	return nil
}

func (a *App) OnTouchMove(this js.Value, args []js.Value) any {
	e := args[0]
	touches := e.Get("changedTouches")
	for i := 0; i < touches.Length(); i++ {
		touch := touches.Index(i)
		where := a.EventPos(touch)
		a.Event.MouseMove(where)
	}
	e.Call("preventDefault")
	return nil
}

func (a *App) OnWheel(this js.Value, args []js.Value) any {
	e := args[0]
	delta := a.EventPosFor(e.Get("deltaX"), e.Get("deltaY"))
	a.Event.Scroll(a.EventPos(e), math32.Vector2FromPoint(delta).DivScalar(8))
	e.Call("preventDefault")
	return nil
}

func (a *App) OnContextMenu(this js.Value, args []js.Value) any {
	e := args[0]
	e.Call("preventDefault")
	return nil
}

func (a *App) ShouldProcessKey(k string) bool {
	if k == "Unidentified" {
		return false
	}
	k = a.KeyMods.ModifiersString() + k
	if a.SystemPlatform() == system.MacOS {
		k = strings.ReplaceAll(k, "Meta", "Command")
	}
	if slices.Contains(system.ReservedWebShortcuts, k) {
		return false
	}
	if a.SystemPlatform() != system.MacOS {

		k = strings.ReplaceAll(k, "Control", "Command")
		if slices.Contains(system.ReservedWebShortcuts, k) {
			return false
		}
	}
	return true
}

func (a *App) RuneAndCodeFromKey(k string, down bool) (rune, key.Codes) {
	switch k {
	case "Shift":
		a.KeyMods.SetFlag(down, key.Shift)
		return 0, key.CodeLeftShift
	case "Control":
		a.KeyMods.SetFlag(down, key.Control)
		return 0, key.CodeLeftControl
	case "Alt":
		a.KeyMods.SetFlag(down, key.Alt)
		return 0, key.CodeLeftAlt
	case "Meta":
		a.KeyMods.SetFlag(down, key.Meta)
		return 0, key.CodeLeftMeta
	case "Enter":
		return 0, key.CodeReturnEnter
	case "ArrowDown":
		return 0, key.CodeDownArrow
	case "ArrowLeft":
		return 0, key.CodeLeftArrow
	case "ArrowRight":
		return 0, key.CodeRightArrow
	case "ArrowUp":
		return 0, key.CodeUpArrow
	case "Spacebar":
		return ' ', 0
	default:
		r := []rune(k)

		if len(r) > 1 {
			kc := key.Codes(0)
			mylog.Check(kc.SetString(k))
			if errors.Log(err) == nil {
				return 0, kc
			}
		}
		return r[0], 0
	}
}

func (a *App) OnKeyDown(this js.Value, args []js.Value) any {
	e := args[0]
	k := e.Get("key").String()
	if !a.ShouldProcessKey(k) {
		return nil
	}
	r, c := a.RuneAndCodeFromKey(k, true)
	a.Event.Key(events.KeyDown, r, c, a.KeyMods)
	e.Call("preventDefault")
	return nil
}

func (a *App) OnKeyUp(this js.Value, args []js.Value) any {
	e := args[0]
	k := e.Get("key").String()
	if !a.ShouldProcessKey(k) {
		return nil
	}
	r, c := a.RuneAndCodeFromKey(k, false)
	a.Event.Key(events.KeyUp, r, c, a.KeyMods)
	e.Call("preventDefault")
	return nil
}

func (a *App) OnBeforeInput(this js.Value, args []js.Value) any {
	e := args[0]
	data := e.Get("data").String()
	if data == "" {
		return nil
	}
	for _, r := range data {
		a.Event.KeyChord(r, 0, a.KeyMods)
	}
	e.Call("preventDefault")
	return nil
}

func (a *App) OnResize(this js.Value, args []js.Value) any {
	a.Resize()
	return nil
}

func (a *App) OnBlur(this js.Value, args []js.Value) any {
	a.KeyMods = 0
	return nil
}
