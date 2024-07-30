package offscreen

import (
	"image"

	"github.com/richardwilkes/unison/system"
)

type Drawer struct {
	system.DrawerBase
}

func (dw *Drawer) DestBounds() image.Rectangle {
	return TheApp.Scrn.Geometry
}

func (dw *Drawer) EndDraw() {
	if !system.NeedsCapture {
		return
	}
	system.NeedsCapture = false
	system.CaptureImage = dw.Image
}
