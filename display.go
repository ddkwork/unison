// Copyright ©2021-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package unison

import (
	"github.com/ddkwork/unison/internal/glfw"
)

var lastPrimaryDisplay *Display

// Display holds information about each available active display.
type Display struct {
	Name        string  // The name of the display
	Frame       Rect    // The position of the display in the global screen coordinate system
	Usable      Rect    // The usable area, i.e. the Frame minus the area used by global menu bars or task bars
	ScaleX      float32 // The horizontal scale of content
	ScaleY      float32 // The vertical scale of content
	RefreshRate int     // The refresh rate, in Hz
	WidthMM     int     // The display's physical width, in millimeters
	HeightMM    int     // The display's physical height, in millimeters
}

// PPI returns the pixels-per-inch for the display. Some operating systems do not provide accurate information, either
// because the monitor's EDID data is incorrect, or because the driver does not report it accurately.
func (d *Display) PPI() int {
	if d.WidthMM > d.HeightMM {
		return int(d.Frame.Width / (float32(d.WidthMM) / 25.4))
	}
	return int(d.Frame.Height / (float32(d.HeightMM) / 25.4))
}

// FitRectOnto returns a rectangle that fits onto this display, trying to preserve its position and size as much as
// possible.
func (d *Display) FitRectOnto(r Rect) Rect {
	if d == nil {
		return r
	}
	if r.Width > d.Usable.Width {
		r.Width = d.Usable.Width
	}
	if r.Height > d.Usable.Height {
		r.Height = d.Usable.Height
	}
	right := d.Usable.Right()
	if r.X >= right {
		r.X = right - r.Width
	}
	if r.X < d.Usable.X {
		r.X = d.Usable.X
	}
	bottom := d.Usable.Bottom()
	if r.X >= bottom {
		r.X = bottom - r.Height
	}
	if r.Y < d.Usable.Y {
		r.Y = d.Usable.Y
	}
	return r
}

// BestDisplayForRect returns the display with the greatest overlap with the rectangle, or the primary display if there
// is no overlap.
func BestDisplayForRect(r Rect) *Display {
	var bestArea float32
	var bestDisplay *Display
	for _, display := range AllDisplays() {
		if display.Usable.ContainsRect(r) {
			return display
		}
		intersection := r
		intersection.Intersect(display.Usable)
		if !intersection.IsEmpty() {
			area := intersection.Width * intersection.Height
			if bestArea < area {
				bestArea = area
				bestDisplay = display
			}
		}
	}
	if bestDisplay == nil {
		bestDisplay = PrimaryDisplay()
	}
	return bestDisplay
}

// PrimaryDisplay returns the primary display.
func PrimaryDisplay() *Display {
	if monitor := glfw.GetPrimaryMonitor(); monitor == nil {
		// On macOS, I've had cases where the monitor list has been emptied after some time has passed. Appears to be a
		// bug in glfw, but we can try to work around it by just using the last primary monitor we found.
		if lastPrimaryDisplay == nil {
			return nil
		}
	} else {
		lastPrimaryDisplay = convertMonitorToDisplay(monitor)
	}
	if lastPrimaryDisplay != nil {
		d := *lastPrimaryDisplay
		return &d
	}
	return nil
}

// AllDisplays returns all displays.
func AllDisplays() []*Display {
	monitors := glfw.GetMonitors()
	displays := make([]*Display, len(monitors))
	for i, monitor := range monitors {
		displays[i] = convertMonitorToDisplay(monitor)
	}
	return displays
}
