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

func (f *FS) Write(args []js.Value) (any, any, error) {
	fl := mylog.Check2(f.GetFile(args))

	flw, ok := fl.(io.Writer)
	if !ok {
		return 0, nil, hackpadfs.ErrNotImplemented
	}

	var n int

	buffer := args[1]

	iblob := mylog.Check2(idbblob.New(buffer))

	offset := args[2].Int()
	length := args[3].Int()
	position := args[4]

	if !position.IsUndefined() && !position.IsNull() {
		_ := mylog.Check2(hackpadfs.SeekFile(fl, int64(position.Int()), io.SeekStart))
	}
	dataToCopy := mylog.Check2(blob.View(iblob, int64(offset), int64(offset+length)))

	n, err = blob.Write(flw, dataToCopy)
	if err == io.EOF {
		err = nil
	}
	return n, buffer, err
}
