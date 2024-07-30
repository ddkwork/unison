//go:build js && !offscreen

package driver

import (
	"cogentcore.org/core/system/driver/web"
)

func init() {
	web.Init()
}
