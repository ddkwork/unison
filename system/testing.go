package system

import (
	"cogentcore.org/core/base/iox/imagex"
)

func AssertCapture(t imagex.TestingT, filename string) {
	imagex.Assert(t, Capture(), filename)
}
