// Copyright 2022 The Ebitengine Authors
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

package opengl

import (
	"runtime"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison/internal/graphicsdriver"
	"github.com/richardwilkes/unison/internal/graphicsdriver/opengl/gl"
	"github.com/richardwilkes/unison/internal/microsoftgdk"

	"github.com/richardwilkes/unison/internal/glfw"
)

type graphicsPlatform struct {
	window *glfw.Window
}

// NewGraphics creates an implementation of graphicsdriver.Graphics for OpenGL.
// The returned graphics value is nil iff the error is not nil.
func NewGraphics() graphicsdriver.Graphics {
	if microsoftgdk.IsXbox() {
		return nil //, fmt.Errorf("opengl: OpenGL is not supported on Xbox")
	}

	ctx := mylog.Check2(gl.NewDefaultContext())

	mylog.Check(setGLFWClientAPI(ctx.IsES()))
	err != nil{
		return nil
	}

	return newGraphics(ctx)
}

func setGLFWClientAPI(isES bool) error {
	if isES {
		mylog.Check(glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI))

		mylog.Check(glfw.WindowHint(glfw.ContextVersionMajor, 3))

		mylog.Check(glfw.WindowHint(glfw.ContextVersionMinor, 0))

		mylog.Check(glfw.WindowHint(glfw.ContextCreationAPI, glfw.EGLContextAPI))

		return nil
	}

	mylog.Check(glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLAPI))

	mylog.Check(glfw.WindowHint(glfw.ContextVersionMajor, 3))

	mylog.Check(glfw.WindowHint(glfw.ContextVersionMinor, 2))

	// macOS requires forward-compatible and a core profile.
	if runtime.GOOS == "darwin" {
		mylog.Check(glfw.WindowHint(glfw.OpenGLForwardCompat, glfw.True))

		mylog.Check(glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile))

	}
	return nil
}

func (g *Graphics) SetGLFWWindow(window *glfw.Window) {
	g.window = window
}

func (g *Graphics) makeContextCurrent() error {
	return g.window.MakeContextCurrent()
}

func (g *Graphics) swapBuffers() error {
	// Call SwapIntervals even though vsync is not changed.
	// When toggling to fullscreen, vsync state might be reset unexpectedly (#1787).

	// SwapInterval is affected by the current monitor of the window.
	// This needs to be called at least after SetMonitor.
	// Without SwapInterval after SetMonitor, vsynch doesn't work (#375).
	if g.vsync {
		mylog.Check(glfw.SwapInterval(1))

	} else {
		mylog.Check(glfw.SwapInterval(0))

	}

	mylog.Check(g.window.SwapBuffers())

	return nil
}
