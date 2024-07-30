//go:build ios

package ios

import "C"

import (
	"unsafe"

	"cogentcore.org/core/base/fileinfo/mimedata"
	"github.com/richardwilkes/unison/system"
)

var TheClipboard = &Clipboard{}

type Clipboard struct {
	system.ClipboardBase
}

func (cl *Clipboard) Read(types []string) mimedata.Mimes {
	cstr := C.getClipboardContent()
	str := C.GoString(cstr)
	return mimedata.NewText(str)
}

func (cl *Clipboard) Write(data mimedata.Mimes) error {
	str := ""
	if len(data) > 1 {
		mpd := data.ToMultipart()
		str = string(mpd)
	} else {
		d := data[0]
		if mimedata.IsText(d.Type) {
			str = string(d.Data)
		}
	}

	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	C.setClipboardContent(cstr)
	return nil
}
