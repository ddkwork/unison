//go:build ios && !offscreen

package driver

import (
	"cogentcore.org/core/system/driver/ios"
)

func init() {
	ios.Init()
}
