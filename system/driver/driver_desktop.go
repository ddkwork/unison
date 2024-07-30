//go:build !(android || ios || js || offscreen)

package driver

import (
	"os"
	"slices"
	"testing"

	"cogentcore.org/core/system/driver/desktop"
	"cogentcore.org/core/system/driver/offscreen"
)

func init() {
	if testing.Testing() || slices.Contains(os.Args, "-nogui") {
		offscreen.Init()
		return
	}
	desktop.Init()
}
