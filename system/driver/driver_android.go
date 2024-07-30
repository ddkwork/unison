//go:build android && !offscreen

package driver

import (
	"cogentcore.org/core/system/driver/android"
)

func init() {
	android.Init()
}
