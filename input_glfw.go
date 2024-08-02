// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !android && !ios && !js && !nintendosdk && !playstation5

package unison

import (
	"math"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison/internal/glfw"
)

var glfwMouseButtonToMouseButton = map[glfw.MouseButton]MouseButton{
	glfw.MouseButtonLeft:   MouseButton0,
	glfw.MouseButtonMiddle: MouseButton1,
	glfw.MouseButtonRight:  MouseButton2,
	glfw.MouseButton4:      MouseButton3,
	glfw.MouseButton5:      MouseButton4,
}

func (u *UserInterface) registerInputCallbacks() error {
	if _ := mylog.Check2(u.window.SetCharModsCallback(func(w *glfw.Window, char rune, mods glfw.ModifierKey) {
		// As this function is called from GLFW callbacks, the current thread is main.
		u.m.Lock()
		defer u.m.Unlock()
		u.inputState.appendRune(char)
	})); 

	if _ := mylog.Check2(u.window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
		// As this function is called from GLFW callbacks, the current thread is main.
		u.m.Lock()
		defer u.m.Unlock()
		u.inputState.WheelX += xoff
		u.inputState.WheelY += yoff
	})); 

	return nil
}

func (u *UserInterface) updateInputState() error {
	u.mainThread.Call(func() {
		mylog.Check(u.updateInputStateImpl())
	})
	return err
}

// updateInputStateImpl must be called from the main thread.
func (u *UserInterface) updateInputStateImpl() error {
	u.m.Lock()
	defer u.m.Unlock()

	for uk, gk := range uiKeyToGLFWKey {
		s := mylog.Check2(u.window.GetKey(gk))

		u.inputState.KeyPressed[uk] = s == glfw.Press
	}
	for gb, ub := range glfwMouseButtonToMouseButton {
		s := mylog.Check2(u.window.GetMouseButton(gb))

		u.inputState.MouseButtonPressed[ub] = s == glfw.Press
	}

	m := u.currentMonitor()
	s := m.DeviceScaleFactor()

	cx, cy := u.savedCursorX, u.savedCursorY
	defer func() {
		u.savedCursorX = math.NaN()
		u.savedCursorY = math.NaN()
	}()

	if !math.IsNaN(cx) && !math.IsNaN(cy) {
		cx2, cy2 := u.context.logicalPositionToClientPosition(cx, cy, s)
		cx2 = dipToGLFWPixel(cx2, s)
		cy2 = dipToGLFWPixel(cy2, s)
		mylog.Check(u.window.SetCursorPos(cx2, cy2))

	} else {
		cx2, cy2 := u.window.GetCursorPos()
		cx2 = dipFromGLFWPixel(cx2, s)
		cy2 = dipFromGLFWPixel(cy2, s)
		cx, cy = u.context.clientPositionToLogicalPosition(cx2, cy2, s)
	}

	// AdjustPosition can return NaN at the initialization.
	if !math.IsNaN(cx) && !math.IsNaN(cy) {
		u.inputState.CursorX, u.inputState.CursorY = cx, cy
	}

	mylog.Check(gamepad.Update())

	return nil
}

func (u *UserInterface) KeyName(key Key) string {
	if !u.isRunning() {
		return ""
	}

	gk, ok := uiKeyToGLFWKey[key]
	if !ok {
		return ""
	}

	var name string
	u.mainThread.Call(func() {
		if u.isTerminated() {
			return
		}
		n := mylog.Check2(glfw.GetKeyName(gk, 0))

		name = n
	})
	return name
}

func (u *UserInterface) saveCursorPosition() {
	u.m.Lock()
	defer u.m.Unlock()

	u.savedCursorX = u.inputState.CursorX
	u.savedCursorY = u.inputState.CursorY
}
