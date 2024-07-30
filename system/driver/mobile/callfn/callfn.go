//go:build android && (arm || 386 || amd64 || arm64)

package callfn

func CallFn(fn uintptr)
