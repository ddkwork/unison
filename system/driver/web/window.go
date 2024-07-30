//go:build js

package web

import (
	"syscall/js"

	"cogentcore.org/core/system/driver/base"
)

type Window struct {
	base.WindowSingle[*App]
}

func (w *Window) SetTitle(title string) {
	w.WindowSingle.SetTitle(title)
	js.Global().Get("document").Set("title", title)
}
