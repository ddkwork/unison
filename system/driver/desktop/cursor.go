package desktop

import (
	"sync"

	"cogentcore.org/core/cursors/cursorimg"
	"cogentcore.org/core/enums"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison/internal/glfw"
	"github.com/richardwilkes/unison/system"
)

var TheCursor = &Cursor{CursorBase: system.CursorBase{Vis: true, Size: 32}, Cursors: map[enums.Enum]map[int]*glfw.Cursor{}}

type Cursor struct {
	system.CursorBase
	Cursors  map[enums.Enum]map[int]*glfw.Cursor
	Mu       sync.Mutex
	PrevSize int
}

func (cu *Cursor) Set(cursor enums.Enum) error {
	cu.Mu.Lock()
	defer cu.Mu.Unlock()
	if cursor == cu.Cur && cu.Size == cu.PrevSize {
		return nil
	}
	sm := cu.Cursors[cursor]
	if sm == nil {
		sm = map[int]*glfw.Cursor{}
		cu.Cursors[cursor] = sm
	}
	if cur, ok := sm[cu.Size]; ok {
		TheApp.CtxWindow.Glw.SetCursor(cur)
		cu.PrevSize = cu.Size
		cu.Cur = cursor
		return nil
	}

	ci := mylog.Check2(cursorimg.Get(cursor, cu.Size))

	h := ci.Hotspot
	gc := glfw.CreateCursor(ci.Image, h.X, h.Y)
	sm[cu.Size] = gc
	TheApp.CtxWindow.Glw.SetCursor(gc)
	cu.PrevSize = cu.Size
	cu.Cur = cursor
	return nil
}
