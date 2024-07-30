//go:build js

package jsfs

import (
	"fmt"
	"io"
	"syscall/js"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/hack-pad/hackpadfs"
	"github.com/hack-pad/hackpadfs/indexeddb/idbblob"
	"github.com/hack-pad/hackpadfs/keyvalue/blob"
)

func (f *FS) Read(args []js.Value) (any, any, error) {
	fl := mylog.Check2(f.GetFile(args))

	var readBuf blob.Blob
	var n int

	buffer := args[1]

	iblob := mylog.Check2(idbblob.New(buffer))

	offset := args[2].Int()
	length := args[3].Int()
	position := args[4]

	if position.IsUndefined() || position.IsNull() {
		readBuf, n = mylog.Check3(blob.Read(fl, length))
	} else {
		readerAt, ok := fl.(io.ReaderAt)
		if ok {
			readBuf, n = mylog.Check3(blob.ReadAt(readerAt, length, int64(position.Int())))
		} else {
			err = &hackpadfs.PathError{Op: "read", Path: fmt.Sprint(fl), Err: hackpadfs.ErrNotImplemented}
		}
	}
	if err == io.EOF {
		err = nil
	}
	if readBuf != nil {
		_, setErr := blob.Set(iblob, readBuf, int64(offset))
		if err == nil && setErr != nil {
			err = &hackpadfs.PathError{Op: "read", Path: fmt.Sprint(fl), Err: setErr}
		}
	}
	return n, buffer, err
}

func (f *FS) ReadFile(args []js.Value) (any, error) {
	fdu := mylog.Check2(f.OpenImpl(NormPath(args[0].String()), 0, 0))

	fd := js.ValueOf(fdu)
	defer f.Close([]js.Value{fd})

	infoa := mylog.Check2(f.Fstat([]js.Value{fd}))

	info := js.ValueOf(infoa)

	buf := js.Global().Get("Uint8Array").New(info.Get("size"))
	_, _ = mylog.Check3(f.Read([]js.Value{fd, buf, js.ValueOf(0), js.ValueOf(buf.Length())}))
	return buf, err
}
