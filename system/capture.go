package system

import (
	"image"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/base/iox/imagex"
)

func Capture() *image.RGBA {
	NeedsCapture = true
	TheApp.Window(0).Drawer().EndDraw()
	return CaptureImage
}

func CaptureAs(filename string) error {
	return errors.Log(imagex.Save(Capture(), filename))
}

var (
	NeedsCapture bool
	CaptureImage *image.RGBA
)
