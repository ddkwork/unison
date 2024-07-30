//go:build js

package jsfs

import (
	"syscall"
	"syscall/js"

	"github.com/ddkwork/golibrary/mylog"
)

func Config(jfs js.Value) (*FS, error) {
	fs := mylog.Check2(NewFS())

	constants := jfs.Get("constants")
	constants.Set("O_RDONLY", syscall.O_RDONLY)
	constants.Set("O_WRONLY", syscall.O_WRONLY)
	constants.Set("O_RDWR", syscall.O_RDWR)
	constants.Set("O_CREAT", syscall.O_CREATE)
	constants.Set("O_TRUNC", syscall.O_TRUNC)
	constants.Set("O_APPEND", syscall.O_APPEND)
	constants.Set("O_EXCL", syscall.O_EXCL)

	SetFunc(jfs, "chmod", fs.Chmod)
	SetFunc(jfs, "chown", fs.Chown)
	SetFunc(jfs, "close", fs.Close)
	SetFunc(jfs, "fchmod", fs.Fchmod)
	SetFunc(jfs, "fchown", fs.Fchown)
	SetFunc(jfs, "fstat", fs.Fstat)
	SetFunc(jfs, "fsync", fs.Fsync)
	SetFunc(jfs, "ftruncate", fs.Ftruncate)
	SetFunc(jfs, "lchown", fs.Lchown)
	SetFunc(jfs, "link", fs.Link)
	SetFunc(jfs, "lstat", fs.Lstat)
	SetFunc(jfs, "mkdir", fs.Mkdir)
	SetFunc(jfs, "mkdirall", fs.MkdirAll)
	SetFunc(jfs, "open", fs.Open)
	SetFunc(jfs, "readdir", fs.Readdir)
	SetFunc(jfs, "readlink", fs.Readlink)
	SetFunc(jfs, "rename", fs.Rename)
	SetFunc(jfs, "rmdir", fs.Rmdir)
	SetFunc(jfs, "stat", fs.Stat)
	SetFunc(jfs, "symlink", fs.Symlink)
	SetFunc(jfs, "unlink", fs.Unlink)
	SetFunc(jfs, "utimes", fs.Utimes)
	SetFunc(jfs, "truncate", fs.Truncate)
	SetFunc(jfs, "read", fs.Read)
	SetFunc(jfs, "readfile", fs.ReadFile)
	SetFunc(jfs, "write", fs.Write)

	return fs, err
}

type Func interface {
	func(args []js.Value) (any, error) | func(args []js.Value) (any, any, error)
}

func SetFunc[F Func](v js.Value, name string, fn F) {
	f := FuncOf(name, fn)
	v.Set(name, f)
}

func FuncOf[F Func](name string, fn F) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		callback := args[len(args)-1]
		args = args[:len(args)-1]

		go func() {
			var res []any

			switch fn := any(fn).(type) {
			case func(args []js.Value) (any, error):
				r, e := fn(args)
				res = []any{r}
				err = e
			case func(args []js.Value) (any, any, error):
				r0, r1, e := fn(args)
				res = []any{r0, r1}
				err = e
			}

			errv := JSError(err, name, args...)
			callback.Invoke(append([]any{errv}, res...)...)
		}()
		return nil
	})
}
