//go:build js

package jsfs

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"syscall/js"
	"time"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/hack-pad/hackpadfs"
	"github.com/hack-pad/hackpadfs/indexeddb"
	"github.com/hack-pad/hackpadfs/mem"
	"github.com/hack-pad/hackpadfs/mount"
)

type FS struct {
	*mount.FS

	PreviousFID uint64
	Files       map[uint64]hackpadfs.File
	Mu          sync.Mutex
}

func NewFS() (*FS, error) {
	memfs := mylog.Check2(mem.NewFS())

	monfs := mylog.Check2(mount.NewFS(memfs))

	f := &FS{
		FS:    monfs,
		Files: map[uint64]hackpadfs.File{},
	}

	_ = mylog.Check2(f.OpenImpl("/dev/stdin", syscall.O_RDONLY, 0))

	_ = mylog.Check2(f.OpenImpl("/dev/stdout", syscall.O_WRONLY, 0))

	_ = mylog.Check2(f.OpenImpl("/dev/stderr", syscall.O_WRONLY, 0))

	_ = mylog.Check2(f.MkdirAll([]js.Value{js.ValueOf("/tmp"), js.ValueOf(0755)}))
	return f, err
}

func (f *FS) ConfigUnix() error {
	perm := js.ValueOf(0755)
	_ := mylog.Check2(f.MkdirAll([]js.Value{js.ValueOf("/home/me"), perm}))

	ifs := mylog.Check2(indexeddb.NewFS(context.Background(), "/home/me", indexeddb.Options{}))
	mylog.Check(f.FS.AddMount("home/me", ifs))

	_ = mylog.Check2(f.MkdirAll([]js.Value{js.ValueOf("/home/me/.data"), perm}))

	_ = mylog.Check2(f.MkdirAll([]js.Value{js.ValueOf("/home/me/Desktop"), perm}))

	_ = mylog.Check2(f.MkdirAll([]js.Value{js.ValueOf("/home/me/Documents"), perm}))

	_ = mylog.Check2(f.MkdirAll([]js.Value{js.ValueOf("/home/me/Downloads"), perm}))

	return nil
}

func NormPath(p string) string {
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return "."
	}
	return p
}

func (f *FS) GetFile(args []js.Value) (hackpadfs.File, error) {
	fd := uint64(args[0].Int())
	fl := f.Files[fd]
	if fl == nil {
		return nil, ErrBadFileNumber(fd)
	}
	return fl, nil
}

func (f *FS) Chmod(args []js.Value) (any, error) {
	return nil, hackpadfs.Chmod(f.FS, NormPath(args[0].String()), hackpadfs.FileMode(args[1].Int()))
}

func (f *FS) Chown(args []js.Value) (any, error) {
	return nil, hackpadfs.Chown(f.FS, NormPath(args[0].String()), args[1].Int(), args[2].Int())
}

func (f *FS) Close(args []js.Value) (any, error) {
	delete(f.Files, uint64(args[0].Int()))
	return nil, nil
}

func (f *FS) Fchmod(args []js.Value) (any, error) {
	fl := mylog.Check2(f.GetFile(args))

	return nil, hackpadfs.ChmodFile(fl, hackpadfs.FileMode(args[1].Int()))
}

func (f *FS) Fchown(args []js.Value) (any, error) {
	fl := mylog.Check2(f.GetFile(args))

	return nil, hackpadfs.ChownFile(fl, args[1].Int(), args[2].Int())
}

func (f *FS) Fstat(args []js.Value) (any, error) {
	fl := mylog.Check2(f.GetFile(args))

	s := mylog.Check2(fl.Stat())

	return JSStat(s), nil
}

func (f *FS) Fsync(args []js.Value) (any, error) {
	fl := mylog.Check2(f.GetFile(args))
	mylog.Check(hackpadfs.SyncFile(fl))
	if errors.Is(err, hackpadfs.ErrNotImplemented) {
		err = nil
	}
	return nil, err
}

func (f *FS) Ftruncate(args []js.Value) (any, error) {
	fl := mylog.Check2(f.GetFile(args))

	return nil, hackpadfs.TruncateFile(fl, int64(args[1].Int()))
}

func (f *FS) Lchown(args []js.Value) (any, error) {
	return nil, hackpadfs.Chown(f.FS, NormPath(args[0].String()), args[1].Int(), args[2].Int())
}

func (f *FS) Link(args []js.Value) (any, error) {
	return nil, hackpadfs.ErrNotImplemented
}

func (f *FS) Lstat(args []js.Value) (any, error) {
	s := mylog.Check2(hackpadfs.LstatOrStat(f.FS, NormPath(args[0].String())))

	return JSStat(s), nil
}

func (f *FS) Mkdir(args []js.Value) (any, error) {
	return nil, hackpadfs.Mkdir(f.FS, NormPath(args[0].String()), hackpadfs.FileMode(args[1].Int()))
}

func (f *FS) MkdirAll(args []js.Value) (any, error) {
	return nil, hackpadfs.MkdirAll(f.FS, NormPath(args[0].String()), hackpadfs.FileMode(args[1].Int()))
}

func (f *FS) Open(args []js.Value) (any, error) {
	return f.OpenImpl(args[0].String(), args[1].Int(), hackpadfs.FileMode(args[2].Int()))
}

func (f *FS) OpenImpl(path string, flags int, mode hackpadfs.FileMode) (uint64, error) {
	path = NormPath(path)

	f.Mu.Lock()
	defer f.Mu.Unlock()

	fid := atomic.AddUint64((*uint64)(&f.PreviousFID), 1) - 1
	fl := mylog.Check2(f.NewFile(path, flags, mode))

	f.Files[fid] = fl

	return fid, nil
}

func (f *FS) NewFile(absPath string, flags int, mode os.FileMode) (hackpadfs.File, error) {
	switch absPath {
	case "dev/null":
		return NewNullFile("dev/null"), nil
	case "dev/stdin":
		return NewNullFile("dev/stdin"), nil
	case "dev/stdout":
		return Stdout, nil
	case "dev/stderr":
		return Stderr, nil
	}
	return hackpadfs.OpenFile(f.FS, absPath, flags, mode)
}

func (f *FS) Readdir(args []js.Value) (any, error) {
	des := mylog.Check2(hackpadfs.ReadDir(f.FS, NormPath(args[0].String())))

	names := make([]any, len(des))
	for i, de := range des {
		names[i] = de.Name()
	}
	return names, nil
}

func (f *FS) Readlink(args []js.Value) (any, error) {
	return nil, hackpadfs.ErrNotImplemented
}

func (f *FS) Rename(args []js.Value) (any, error) {
	return nil, hackpadfs.Rename(f.FS, NormPath(args[0].String()), NormPath(args[1].String()))
}

func (f *FS) Rmdir(args []js.Value) (any, error) {
	info := mylog.Check2(f.Stat(args))

	if !js.ValueOf(info).Call("isDirectory").Bool() {
		return nil, ErrNotDir
	}
	return nil, hackpadfs.Remove(f.FS, NormPath(args[0].String()))
}

func (f *FS) Stat(args []js.Value) (any, error) {
	s := mylog.Check2(hackpadfs.Stat(f.FS, NormPath(args[0].String())))

	return JSStat(s), nil
}

func (f *FS) Symlink(args []js.Value) (any, error) {
	return nil, hackpadfs.ErrNotImplemented
}

func (f *FS) Unlink(args []js.Value) (any, error) {
	info := mylog.Check2(f.Stat(args))

	if js.ValueOf(info).Call("isDirectory").Bool() {
		return nil, os.ErrPermission
	}
	return nil, hackpadfs.Remove(f.FS, NormPath(args[0].String()))
}

func (f *FS) Utimes(args []js.Value) (any, error) {
	path := NormPath(args[0].String())
	atime := time.Unix(int64(args[1].Int()), 0)
	mtime := time.Unix(int64(args[2].Int()), 0)

	return nil, hackpadfs.Chtimes(f.FS, path, atime, mtime)
}

func (f *FS) Truncate(args []js.Value) (any, error) {
	return nil, hackpadfs.ErrNotImplemented
}
