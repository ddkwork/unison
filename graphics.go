// Copyright 2022 The Ebiten Authors
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

package unison

import (
	"fmt"
	"os"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison/internal/graphicsdriver"
)

type graphicsDriverCreator interface {
	newAuto() (graphicsdriver.Graphics, GraphicsLibrary, error)
	newOpenGL() (graphicsdriver.Graphics, error)
	newDirectX() (graphicsdriver.Graphics, error)
	newMetal() (graphicsdriver.Graphics, error)
	newPlayStation5() (graphicsdriver.Graphics, error)
}

func newGraphicsDriver(creator graphicsDriverCreator, graphicsLibrary GraphicsLibrary) (graphicsdriver.Graphics, GraphicsLibrary, error) {
	if graphicsLibrary == GraphicsLibraryAuto {
		envName := "EBITENGINE_GRAPHICS_LIBRARY"
		env := os.Getenv(envName)
		if env == "" {
			// For backward compatibility, read the EBITEN_ version.
			envName = "EBITEN_GRAPHICS_LIBRARY"
			env = os.Getenv(envName)
		}

		switch env {
		case "", "auto":
			// Keep the automatic choosing.
		case "opengl":
			graphicsLibrary = GraphicsLibraryOpenGL
		case "directx":
			graphicsLibrary = GraphicsLibraryDirectX
		case "metal":
			graphicsLibrary = GraphicsLibraryMetal
		case "playstation5":
			graphicsLibrary = GraphicsLibraryPlayStation5
		default:
			return nil, 0, fmt.Errorf("ui: an unsupported graphics library is specified by the environment variable: %s", env)
		}
	}

	switch graphicsLibrary {
	case GraphicsLibraryAuto:
		g, lib := mylog.Check3(creator.newAuto())

		if g == nil {
			return nil, 0, fmt.Errorf("ui: no graphics library is available")
		}
		return g, lib, nil
	case GraphicsLibraryOpenGL:
		g := mylog.Check2(creator.newOpenGL())

		return g, GraphicsLibraryOpenGL, nil
	case GraphicsLibraryDirectX:
		g := mylog.Check2(creator.newDirectX())

		return g, GraphicsLibraryDirectX, nil
	case GraphicsLibraryMetal:
		g := mylog.Check2(creator.newMetal())

		return g, GraphicsLibraryMetal, nil
	case GraphicsLibraryPlayStation5:
		g := mylog.Check2(creator.newPlayStation5())

		return g, GraphicsLibraryPlayStation5, nil
	default:
		return nil, 0, fmt.Errorf("ui: an unsupported graphics library is specified: %d", graphicsLibrary)
	}
}

func (u *UserInterface) GraphicsDriverForTesting() graphicsdriver.Graphics {
	return u.graphicsDriver
}

type GraphicsLibrary int

const (
	GraphicsLibraryAuto GraphicsLibrary = iota
	GraphicsLibraryUnknown
	GraphicsLibraryOpenGL
	GraphicsLibraryDirectX
	GraphicsLibraryMetal
	GraphicsLibraryPlayStation5
)

func (g GraphicsLibrary) String() string {
	switch g {
	case GraphicsLibraryAuto:
		return "Auto"
	case GraphicsLibraryUnknown:
		return "Unknown"
	case GraphicsLibraryOpenGL:
		return "OpenGL"
	case GraphicsLibraryDirectX:
		return "DirectX"
	case GraphicsLibraryMetal:
		return "Metal"
	case GraphicsLibraryPlayStation5:
		return "PlayStation 5"
	default:
		return fmt.Sprintf("GraphicsLibrary(%d)", g)
	}
}
