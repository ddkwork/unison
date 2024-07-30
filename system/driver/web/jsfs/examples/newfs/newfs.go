//go:build js

package main

import (
	"syscall/js"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/system/driver/web/jsfs"
)

func main() {
	fs := errors.Must1(jsfs.NewFS())
	errors.Must1(fs.MkdirAll([]js.Value{js.ValueOf("me"), js.ValueOf(0777)}))
	js.Global().Get("console").Call("log", "stat file info", errors.Must1(fs.Stat([]js.Value{js.ValueOf("me")})))
}
