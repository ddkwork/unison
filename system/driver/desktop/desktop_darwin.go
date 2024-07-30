//go:build darwin

package desktop

import "C"

import (
	"fmt"
	"os/exec"
	"os/user"
	"path/filepath"
	"unsafe"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/base/fileinfo/mimedata"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison/system"
)

func SetThreadPri(p float64) error {
	rv := C.setThreadPri(C.double(p))
	if rv != 0 {
		mylog.Check(fmt.Errorf("SetThreadPri failed: %v\n", rv))
		fmt.Println(err)
		return err
	}
	return nil
}

func (a *App) Platform() system.Platforms {
	return system.MacOS
}

func (a *App) OpenURL(url string) {
	cmd := exec.Command("open", url)
	errors.Log(cmd.Run())
}

func (a *App) DataDir() string {
	usr := mylog.Check2(user.Current())
	if errors.Log(err) != nil {
		return "/tmp"
	}
	return filepath.Join(usr.HomeDir, "Library")
}

var TheClipboard = &Clipboard{}

type Clipboard struct {
	Data mimedata.Mimes
}

var CurMimeData *mimedata.Mimes

func (cl *Clipboard) IsEmpty() bool {
	ise := C.clipIsEmpty()
	return bool(ise)
}

func (cl *Clipboard) Read(types []string) mimedata.Mimes {
	if len(types) == 0 {
		return nil
	}
	cl.Data = nil
	CurMimeData = &cl.Data

	wantText := mimedata.IsText(types[0])

	if wantText {
		C.clipReadText()
		if len(cl.Data) == 0 {
			return nil
		}
		dat := cl.Data[0].Data
		isMulti, mediaType, boundary, body := mimedata.IsMultipart(dat)
		if isMulti {
			return mimedata.FromMultipart(body, boundary)
		} else {
			if mediaType != "" {
				return mimedata.NewMime(mediaType, dat)
			} else {
				return mimedata.NewMime(types[0], dat)
			}
		}
	} else {
	}
	return cl.Data
}

func (cl *Clipboard) WriteText(b []byte) {
	sz := len(b)
	cdata := C.malloc(C.size_t(sz))
	copy((*[1 << 30]byte)(cdata)[0:sz], b)
	C.pasteWriteAddText((*C.char)(cdata), C.int(sz))
	C.free(unsafe.Pointer(cdata))
}

func (cl *Clipboard) Write(data mimedata.Mimes) error {
	cl.Clear()
	if len(data) > 1 {
		mpd := data.ToMultipart()
		cl.WriteText(mpd)
	} else {
		d := data[0]
		if mimedata.IsText(d.Type) {
			cl.WriteText(d.Data)
		}
	}
	C.clipWrite()
	return nil
}

func (cl *Clipboard) Clear() {
	C.clipClear()
}

func addMimeText(cdata *C.char, datalen C.int) {
	if *CurMimeData == nil {
		*CurMimeData = make(mimedata.Mimes, 1)
		(*CurMimeData)[0] = &mimedata.Data{Type: mimedata.TextPlain}
	}
	md := (*CurMimeData)[0]
	if len(md.Type) == 0 {
		md.Type = mimedata.TextPlain
	}
	data := C.GoBytes(unsafe.Pointer(cdata), datalen)
	md.Data = append(md.Data, data...)
}

func addMimeData(ctyp *C.char, typlen C.int, cdata *C.char, datalen C.int) {
	if *CurMimeData == nil {
		*CurMimeData = make(mimedata.Mimes, 0)
	}
	typ := C.GoStringN(ctyp, typlen)
	data := C.GoBytes(unsafe.Pointer(cdata), datalen)
	*CurMimeData = append(*CurMimeData, &mimedata.Data{typ, data})
}

func macOpenFile(fname *C.char, flen C.int) {
}
