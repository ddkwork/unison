// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2002-2006 Marcus Geelnard
// SPDX-FileCopyrightText: 2006-2018 Camilla Löwy <elmindreda@glfw.org>
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

package glfw

import (
	"errors"

	"github.com/ddkwork/golibrary/mylog"
)

func terminate() error {
	for _, w := range _glfw.windows {
		mylog.Check(w.Destroy())
	}
	for _, c := range _glfw.cursors {
		mylog.Check(c.Destroy())
	}
	_glfw.monitors = nil
	mylog.Check(platformTerminate())
	_glfw.initialized = false
	mylog.Check(_glfw.contextSlot.destroy())
	return nil
}

func Init() (ferr error) {
	defer func() {
		if ferr != nil {
			// InvalidValue can happen when specific joysticks are used. This issue
			// will be fixed in GLFW 3.3.5. As a temporary fix, ignore this error.
			// See go-gl/glfw#292, go-gl/glfw#324, and glfw/glfw#1763
			// (#1229).
			if errors.Is(ferr, InvalidValue) {
				ferr = nil
				return
			}
			mylog.Check(terminate())
		}
	}()
	if _glfw.initialized {
		return nil
	}
	_glfw.hints.init.hatButtons = true
	mylog.Check(platformInit())
	mylog.Check(_glfw.contextSlot.create())
	_glfw.initialized = true
	mylog.Check(defaultWindowHints())
	return nil
}

func Terminate() error {
	if !_glfw.initialized {
		return nil
	}
	mylog.Check(terminate())
	return nil
}
