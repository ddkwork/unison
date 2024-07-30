package system

import (
	"image"
	"unicode/utf8"

	"cogentcore.org/core/events"
	"cogentcore.org/core/math32"
)

type Window interface {
	Name() string
	SetName(name string)
	Title() string
	SetTitle(title string)
	SetIcon(images []image.Image)
	Size() image.Point
	WinSize() image.Point
	Position() image.Point
	RenderGeom() math32.Geom2DInt
	SetWinSize(sz image.Point)
	SetSize(sz image.Point)
	SetPos(pos image.Point)
	SetGeom(pos image.Point, sz image.Point)
	Raise()
	Minimize()
	PhysicalDPI() float32
	LogicalDPI() float32
	SetLogicalDPI(dpi float32)
	Screen() *Screen
	Flags() WindowFlags
	Is(flag WindowFlags) bool
	IsClosed() bool
	IsVisible() bool
	SetCloseReqFunc(fun func(win Window))
	SetCloseCleanFunc(fun func(win Window))
	CloseReq()
	CloseClean()
	Close()
	SetMousePos(x, y float64)
	SetCursorEnabled(enabled, raw bool)
	IsCursorEnabled() bool
	Drawer() Drawer
	Lock() bool
	Unlock()
	SetDestroyGPUResourcesFunc(f func())
	SetFPS(fps int)
	SetTitleBarIsDark(isDark bool)
	Events() *events.Source
}

type WindowFlags int64

const (
	Dialog WindowFlags = iota
	Modal
	Tool
	Fullscreen
	Minimized
	Focused
)

type NewWindowOptions struct {
	Size      image.Point
	StdPixels bool
	Pos       image.Point
	Title     string
	Icon      []image.Image
	Flags     WindowFlags
}

func (o *NewWindowOptions) SetDialog() {
	o.Flags.SetFlag(true, Dialog)
}

func (o *NewWindowOptions) SetModal() {
	o.Flags.SetFlag(true, Modal)
}

func (o *NewWindowOptions) SetTool() {
	o.Flags.SetFlag(true, Tool)
}

func (o *NewWindowOptions) SetFullscreen() {
	o.Flags.SetFlag(true, Fullscreen)
}

func WindowFlagsToBool(flags WindowFlags) (dialog, modal, tool, fullscreen bool) {
	dialog = flags.HasFlag(Dialog)
	modal = flags.HasFlag(Modal)
	tool = flags.HasFlag(Tool)
	fullscreen = flags.HasFlag(Fullscreen)
	return
}

func (o *NewWindowOptions) GetTitle() string {
	if o == nil {
		return ""
	}
	return sanitizeUTF8(o.Title, 4096)
}

func sanitizeUTF8(s string, n int) string {
	if n < len(s) {
		s = s[:n]
	}
	i := 0
	for i < len(s) {
		r, n := utf8.DecodeRuneInString(s[i:])
		if r == 0 || (r == utf8.RuneError && n == 1) {
			break
		}
		i += n
	}
	return s[:i]
}

func (o *NewWindowOptions) Fixup() {
	sc := TheApp.Screen(0)
	scsz := sc.Geometry.Size()

	if o.Size.X <= 0 {
		o.StdPixels = false
		o.Size.X = int(0.8 * float32(scsz.X) * sc.DevicePixelRatio)
	}
	if o.Size.Y <= 0 {
		o.StdPixels = false
		o.Size.Y = int(0.8 * float32(scsz.Y) * sc.DevicePixelRatio)
	}

	o.Size, o.Pos = sc.ConstrainWinGeom(o.Size, o.Pos)
	if o.Pos.X == 0 && o.Pos.Y == 0 {
		wsz := sc.WinSizeFromPix(o.Size)
		dialog, modal, _, _ := WindowFlagsToBool(o.Flags)
		nw := TheApp.NWindows()
		if nw > 0 {
			lastw := TheApp.Window(nw - 1)
			lsz := lastw.WinSize()
			lp := lastw.Position()

			nwbig := wsz.X > lsz.X || wsz.Y > lsz.Y

			if modal || dialog || !nwbig {
				ctrx := lp.X + (lsz.X / 2)
				ctry := lp.Y + (lsz.Y / 2)
				o.Pos.X = ctrx - wsz.X/2
				o.Pos.Y = ctry - wsz.Y/2
			} else {
				o.Pos.X = lp.X + lsz.X
				o.Pos.Y = lp.Y + 72
			}
		} else {
			o.Pos.X = scsz.X/2 - wsz.X/2
			o.Pos.Y = scsz.Y/2 - wsz.Y/2
		}
		o.Size, o.Pos = sc.ConstrainWinGeom(o.Size, o.Pos)
	}
}
