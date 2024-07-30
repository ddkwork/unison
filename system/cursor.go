package system

import (
	"cogentcore.org/core/enums"
)

type Cursor interface {
	Current() enums.Enum
	Set(cursor enums.Enum) error
	IsVisible() bool
	Hide()
	Show()
	SetSize(size int)
}

type CursorBase struct {
	Cur  enums.Enum
	Vis  bool
	Size int
}

var _ Cursor = (*CursorBase)(nil)

func (c *CursorBase) Current() enums.Enum {
	return c.Cur
}

func (c *CursorBase) Set(cursor enums.Enum) error {
	c.Cur = cursor
	return nil
}

func (c *CursorBase) IsVisible() bool {
	return c.Vis
}

func (c *CursorBase) Hide() {
	c.Vis = false
}

func (c *CursorBase) Show() {
	c.Vis = true
}

func (c *CursorBase) SetSize(size int) {
	c.Size = size
}
