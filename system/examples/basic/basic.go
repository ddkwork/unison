package main

import (
	"fmt"
	"image"

	"cogentcore.org/core/system"
	_ "cogentcore.org/core/system/driver"
	"github.com/ddkwork/golibrary/mylog"
)

func main() {
	fmt.Println("Hello, world!")
	opts := &system.NewWindowOptions{
		Size:      image.Pt(1024, 768),
		StdPixels: true,
		Title:     "System Test Window",
	}
	w := mylog.Check2(system.TheApp.NewWindow(opts))

	fmt.Println("got new window", w)
	system.TheApp.MainLoop()
}
