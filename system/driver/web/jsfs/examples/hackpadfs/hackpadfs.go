//go:build js

package main

import (
	"context"
	"syscall/js"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/system/driver/web/jsfs"
	"github.com/hack-pad/hackpadfs/indexeddb"
)

func main() {
	idb := errors.Must1(indexeddb.NewFS(context.Background(), "idb", indexeddb.Options{}))
	errors.Must(idb.MkdirAll("me", 0777))
	js.Global().Get("console").Call("log", "stat file info", jsfs.JSStat(errors.Must1(idb.Stat("me"))))
}
