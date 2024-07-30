package system

import (
	"cogentcore.org/core/base/fileinfo/mimedata"
)

type Clipboard interface {
	IsEmpty() bool
	Read(types []string) mimedata.Mimes
	Write(data mimedata.Mimes) error
	Clear()
}

type ClipboardBase struct{}

var _ Clipboard = &ClipboardBase{}

func (bb *ClipboardBase) IsEmpty() bool                      { return false }
func (bb *ClipboardBase) Read(types []string) mimedata.Mimes { return nil }
func (bb *ClipboardBase) Write(data mimedata.Mimes) error    { return nil }
func (bb *ClipboardBase) Clear()                             {}
