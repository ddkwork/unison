//go:build js

package web

import (
	"image"
	"syscall/js"
	"unsafe"

	"github.com/richardwilkes/unison/system"
)

type Drawer struct {
	system.DrawerBase
}

func (dw *Drawer) DestBounds() image.Rectangle {
	return TheApp.Scrn.Geometry
}

func (dw *Drawer) EndDraw() {
	sz := dw.Image.Bounds().Size()
	ptr := uintptr(unsafe.Pointer(&dw.Image.Pix[0]))
	js.Global().Call("displayImage", ptr, len(dw.Image.Pix), sz.X, sz.Y)
}
