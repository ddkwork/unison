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

package gl

import (
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/sys/windows"
)

var (
	opengl32              = windows.NewLazySystemDLL("opengl32")
	procWglGetProcAddress = opengl32.NewProc("wglGetProcAddress")
)

func (c *defaultContext) init() error {
	return nil
}

func (c *defaultContext) getProcAddress(namea string) (uintptr, error) {
	cname := mylog.Check2(windows.BytePtrFromString(namea))

	r, _ := mylog.Check3(procWglGetProcAddress.Call(uintptr(unsafe.Pointer(cname))))
	if r != 0 {
		return r, nil
	}
	p := opengl32.NewProc(namea)
	mylog.Check(p.Find())
	return p.Addr(), nil
}
