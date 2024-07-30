//go:build android

package android

import (
	"cogentcore.org/core/system/driver/base"
)

type Window struct {
	base.WindowSingle[*App]
}
