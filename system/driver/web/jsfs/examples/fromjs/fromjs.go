//go:build js

package main

import (
	"syscall/js"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/system/driver/web/jsfs"
)

func main() {
	fs := errors.Must1(jsfs.Config(js.Global().Get("fs")))
	errors.Must1(fs.MkdirAll([]js.Value{js.ValueOf("me"), js.ValueOf(0777)}))
	callback := js.FuncOf(func(this js.Value, args []js.Value) any {
		js.Global().Get("console").Call("log", "stat file info", args[1])
		return nil
	})
	js.Global().Get("fs").Call("stat", "me", callback)
	select {}
}
