package system

import (
	"image"
	"math"
)

var (
	LogicalDPIScale  = float32(1)
	LogicalDPIScales map[string]float32
)

type Screen struct {
	ScreenNumber       int
	Geometry           image.Rectangle
	DevicePixelRatio   float32
	PixSize            image.Point
	PhysicalSize       image.Point
	LogicalDPI         float32
	PhysicalDPI        float32
	Depth              int
	RefreshRate        float32
	Orientation        ScreenOrientation
	NativeOrientation  ScreenOrientation
	PrimaryOrientation ScreenOrientation
	Name               string
	Manufacturer       string
	Model              string
	SerialNumber       string
}

type ScreenOrientation int32

const (
	OrientationUnknown ScreenOrientation = iota

	Portrait

	Landscape
)

func LogicalFromPhysicalDPI(logScale, pdpi float32) float32 {
	idpi := int(math.Round(float64(pdpi * logScale)))
	mdpi := idpi / 6
	mdpi *= 6
	return float32(mdpi)
}

func SetLogicalDPIScale(scrnName string, dpiScale float32) {
	if LogicalDPIScales == nil {
		LogicalDPIScales = make(map[string]float32)
	}
	LogicalDPIScales[scrnName] = dpiScale
}

func (sc *Screen) UpdateLogicalDPI() {
	dpisc := LogicalDPIScale
	if LogicalDPIScales != nil {
		if dsc, has := LogicalDPIScales[sc.Name]; has {
			dpisc = dsc
		}
	}
	sc.LogicalDPI = LogicalFromPhysicalDPI(dpisc, sc.PhysicalDPI)
}

func (sc *Screen) UpdatePhysicalDPI() {
	sc.PhysicalDPI = 25.4 * (float32(sc.PixSize.X) / float32(sc.PhysicalSize.X))
}

func (sc *Screen) WinSizeToPix(sz image.Point) image.Point {
	var psz image.Point
	psz.X = int(float32(sz.X) * sc.DevicePixelRatio)
	psz.Y = int(float32(sz.Y) * sc.DevicePixelRatio)
	return psz
}

func (sc *Screen) WinSizeFromPix(sz image.Point) image.Point {
	var wsz image.Point
	wsz.X = int(float32(sz.X) / sc.DevicePixelRatio)
	wsz.Y = int(float32(sz.Y) / sc.DevicePixelRatio)
	return wsz
}

func (sc *Screen) ConstrainWinGeom(sz, pos image.Point) (csz, cpos image.Point) {
	scsz := sc.Geometry.Size()

	csz = sc.WinSizeFromPix(sz)
	cpos = pos

	if csz.X > scsz.X {
		csz.X = scsz.X - 10
	}
	if csz.Y > scsz.Y {
		csz.Y = scsz.Y - 10
	}

	if cpos.X == -32000 {
		cpos.X = 0
	}
	if cpos.Y == -32000 {
		cpos.Y = 50
	}

	if cpos.X+csz.X > scsz.X {
		cpos.X = scsz.X - csz.X
	}
	if cpos.Y+csz.Y > scsz.Y {
		cpos.Y = scsz.Y - csz.Y
	}
	if cpos.X < 0 {
		cpos.X = 0
	}
	if cpos.Y < 0 {
		cpos.Y = 0
	}

	csz = sc.WinSizeToPix(csz)
	return
}

var InitScreenLogicalDPIFunc func()
