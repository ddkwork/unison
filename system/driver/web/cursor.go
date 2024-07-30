//go:build js

package web

import (
	"strings"
	"syscall/js"

	"cogentcore.org/core/cursors"
	"cogentcore.org/core/enums"
	"github.com/richardwilkes/unison/system"
)

var TheCursor = &Cursor{CursorBase: system.CursorBase{Vis: true, Size: 32}}

type Cursor struct {
	system.CursorBase
}

func (cu *Cursor) Set(cursor enums.Enum) error {
	s := cursor.String()

	if cursor == cursors.Arrow {
		s = "default"
	}

	if strings.HasPrefix(s, "resize-") {
		s = strings.TrimPrefix(s, "resize-")
		s += "-resize"
	}
	js.Global().Get("document").Get("body").Get("style").Set("cursor", s)
	return nil
}
