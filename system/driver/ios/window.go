//go:build ios

package ios

import (
	"cogentcore.org/core/system/driver/base"
)

type Window struct {
	base.WindowSingle[*App]
}
