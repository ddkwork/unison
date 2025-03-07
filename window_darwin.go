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
	"github.com/ddkwork/unison/internal/ns"
)

func (w *Window) frameRect() Rect {
	if w.IsValid() {
		left, top, right, bottom := w.wnd.GetFrameSize()
		return NewRect(float32(left), float32(top), float32(right-left), float32(bottom-top))
	}
	return NewRect(0, 0, 1, 1)
}

// ContentRect returns the boundaries in display coordinates of the window's content area.
func (w *Window) ContentRect() Rect {
	if w.IsValid() {
		x, y := w.wnd.GetPos()
		width, height := w.wnd.GetSize()
		return NewRect(float32(x), float32(y), float32(width), float32(height))
	}
	return NewRect(0, 0, 1, 1)
}

// SetContentRect sets the boundaries of the frame of this window by converting the content rect into a suitable frame
// rect and then applying it to the window.
func (w *Window) SetContentRect(rect Rect) {
	if w.IsValid() {
		rect = w.adjustContentRectForMinMax(rect)
		w.wnd.SetPos(int(rect.X), int(rect.Y))
		tx := int(rect.Width)
		ty := int(rect.Height)
		w.wnd.SetSize(tx, ty)
	}
}

func (w *Window) convertMouseLocation(x, y float64) Point {
	return Point{X: float32(x), Y: float32(y)}
}

func (w *Window) keyCallbackForGLFW(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, mods glfw.ModifierKey) {
	if w.okToProcess() {
		w.commonKeyCallbackForGLFW(key, action, mods)
	}
}

// CurrentKeyModifiers returns the current key modifiers, which is usually the same as calling .LastKeyModifiers(),
// however, on platforms that are using native menus, this will also capture modifier changes that occurred while the
// menu is being displayed.
func (w *Window) CurrentKeyModifiers() Modifiers {
	return modifiersFromEventModifierFlags(ns.CurrentModifierFlags())
}
