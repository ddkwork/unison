// Copyright ©2021-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package unison

import (
	_ "embed"

	"github.com/ddkwork/golibrary/mylog"
)

var (
	moveCursorImage *Image
	//go:embed resources/images/move.png
	moveCursorImageData []byte
)

// MoveCursorImage returns the standard move cursor image.
func MoveCursorImage() *Image {
	if moveCursorImage == nil {
		moveCursorImage = mylog.Check2(NewImageFromBytes(moveCursorImageData, 0.5))
	}
	return moveCursorImage
}

var (
	resizeHorizontalCursorImage *Image
	//go:embed resources/images/resize_horizontal.png
	resizeHorizontalCursorImageData []byte
)

// ResizeHorizontalCursorImage returns the standard horizontal resize cursor image.
func ResizeHorizontalCursorImage() *Image {
	if resizeHorizontalCursorImage == nil {
		resizeHorizontalCursorImage = mylog.Check2(NewImageFromBytes(resizeHorizontalCursorImageData, 0.5))
	}
	return resizeHorizontalCursorImage
}

var (
	resizeLeftDiagonalCursorImage *Image
	//go:embed resources/images/resize_left_diagonal.png
	resizeLeftDiagonalCursorImageData []byte
)

// ResizeLeftDiagonalCursorImage returns the standard left diagonal resize cursor image.
func ResizeLeftDiagonalCursorImage() *Image {
	if resizeLeftDiagonalCursorImage == nil {
		resizeLeftDiagonalCursorImage = mylog.Check2(NewImageFromBytes(resizeLeftDiagonalCursorImageData, 0.5))
	}
	return resizeLeftDiagonalCursorImage
}

var (
	resizeRightDiagonalCursorImage *Image
	//go:embed resources/images/resize_right_diagonal.png
	resizeRightDiagonalCursorImageData []byte
)

// ResizeRightDiagonalCursorImage returns the standard right diagonal resize cursor image.
func ResizeRightDiagonalCursorImage() *Image {
	if resizeRightDiagonalCursorImage == nil {
		resizeRightDiagonalCursorImage = mylog.Check2(NewImageFromBytes(resizeRightDiagonalCursorImageData, 0.5))
	}
	return resizeRightDiagonalCursorImage
}

var (
	resizeVerticalCursorImage *Image
	//go:embed resources/images/resize_vertical.png
	resizeVerticalCursorImageData []byte
)

// ResizeVerticalCursorImage returns the standard vertical resize cursor image.
func ResizeVerticalCursorImage() *Image {
	if resizeVerticalCursorImage == nil {
		resizeVerticalCursorImage = mylog.Check2(NewImageFromBytes(resizeVerticalCursorImageData, 0.5))
	}
	return resizeVerticalCursorImage
}
