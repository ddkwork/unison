//go:build offscreen

package driver

import (
	"cogentcore.org/core/system/driver/offscreen"
)

func init() {
	offscreen.Init()
}
