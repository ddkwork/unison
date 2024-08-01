package system

import (
	"image"
	"image/color"
	"image/draw"

	"cogentcore.org/core/math32"
	"github.com/richardwilkes/unison/internal/skia"
)

const (
	MaxTexturesPerSet = 16
	MaxImageLayers    = 128
	FlipY             = true
	NoFlipY           = false
)

type Drawer interface {
	SetMaxTextures(maxTextures int)
	MaxTextures() int
	DestBounds() image.Rectangle
	SetGoImage(idx, layer int, img image.Image, flipY bool)
	ConfigImageDefaultFormat(idx int, width int, height int, layers int)
	SyncImages()
	Copy(idx, layer int, dp image.Point, sr image.Rectangle, op draw.Op, flipY bool) error
	Scale(idx, layer int, dr image.Rectangle, sr image.Rectangle, op draw.Op, flipY bool, rotDeg float32) error
	UseTextureSet(descIndex int)
	StartDraw(descIndex int) bool
	EndDraw()
	Fill(clr color.Color, src2dst math32.Matrix3, reg image.Rectangle, op draw.Op) error
	StartFill() bool
	EndFill()
	Surface() skia.Surface
	SetFrameImage(idx int, fb any)
}

type DrawerBase struct {
	MaxTxts int
	Image   *image.RGBA
	Images  [][]*image.RGBA
}

func (dw *DrawerBase) SetMaxTextures(maxTextures int) {
	dw.MaxTxts = maxTextures
}

func (dw *DrawerBase) MaxTextures() int {
	return dw.MaxTxts
}

func (dw *DrawerBase) SetGoImage(idx, layer int, img image.Image, flipY bool) {
	for len(dw.Images) <= idx {
		dw.Images = append(dw.Images, nil)
	}
	imgs := &dw.Images[idx]
	for len(*imgs) <= layer {
		*imgs = append(*imgs, nil)
	}
	(*imgs)[layer] = img.(*image.RGBA)
}

func (dw *DrawerBase) ConfigImageDefaultFormat(idx int, width int, height int, layers int) {
}

func (dw *DrawerBase) SyncImages() {
}

func (dw *DrawerBase) Copy(idx, layer int, dp image.Point, sr image.Rectangle, op draw.Op, flipY bool) error {
	img := dw.Images[idx][layer]
	draw.Draw(dw.Image, image.Rectangle{dp, dp.Add(img.Rect.Size())}, img, sr.Min, op)
	return nil
}

func (dw *DrawerBase) Scale(idx, layer int, dr image.Rectangle, sr image.Rectangle, op draw.Op, flipY bool, rotDeg float32) error {
	img := dw.Images[idx][layer]

	draw.Draw(dw.Image, dr, img, sr.Min, op)
	return nil
}

func (dw *DrawerBase) UseTextureSet(descIndex int) {
}

func (dw *DrawerBase) StartDraw(descIndex int) bool {
	return true
}

func (dw *DrawerBase) Fill(clr color.Color, src2dst math32.Matrix3, reg image.Rectangle, op draw.Op) error {
	draw.Draw(dw.Image, reg, image.NewUniform(clr), image.Point{}, op)
	return nil
}

func (dw *DrawerBase) StartFill() bool {
	return true
}

func (dw *DrawerBase) EndFill() {
}

func (dw *DrawerBase) Surface() any {
	return nil
}

func (dw *DrawerBase) SetFrameImage(idx int, fb any) {
}
