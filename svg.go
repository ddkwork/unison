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
	"encoding/xml"
	"io"
	"strconv"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

var _ Drawable = &DrawableSVG{}

// Pre-defined SVG images used by Unison.
var (
	//go:embed resources/images/broken_image.svg
	brokenImageSVG string
	BrokenImageSVG = MustSVGFromContentString(brokenImageSVG)

	//go:embed resources/images/circled_chevron_right.svg
	circledChevronRightSVG string
	CircledChevronRightSVG = MustSVGFromContentString(circledChevronRightSVG)

	//go:embed resources/images/circled_exclamation.svg
	circledExclamationSVG string
	CircledExclamationSVG = MustSVGFromContentString(circledExclamationSVG)

	//go:embed resources/images/circled_question.svg
	circledQuestionSVG string
	CircledQuestionSVG = MustSVGFromContentString(circledQuestionSVG)

	//go:embed resources/images/checkmark.svg
	checkmarkSVG string
	CheckmarkSVG = MustSVGFromContentString(checkmarkSVG)

	//go:embed resources/images/chevron_right.svg
	chevronRightSVG string
	ChevronRightSVG = MustSVGFromContentString(chevronRightSVG)

	//go:embed resources/images/circled_x.svg
	circledXSVG string
	CircledXSVG = MustSVGFromContentString(circledXSVG)

	//go:embed resources/images/dash.svg
	dashSVG string
	DashSVG = MustSVGFromContentString(dashSVG)

	//go:embed resources/images/document.svg
	documentSVG string
	DocumentSVG = MustSVGFromContentString(documentSVG)

	//go:embed resources/images/sort_ascending.svg
	sortAscendingSVG string
	SortAscendingSVG = MustSVGFromContentString(sortAscendingSVG)

	//go:embed resources/images/sort_descending.svg
	sortDescendingSVG string
	SortDescendingSVG = MustSVGFromContentString(sortDescendingSVG)

	//go:embed resources/images/triangle_exclamation.svg
	triangleExclamationSVG string
	TriangleExclamationSVG = MustSVGFromContentString(triangleExclamationSVG)

	//go:embed resources/images/window_maximize.svg
	windowMaximizeSVG string
	WindowMaximizeSVG = MustSVGFromContentString(windowMaximizeSVG)

	//go:embed resources/images/window_restore.svg
	windowRestoreSVG string
	WindowRestoreSVG = MustSVGFromContentString(windowRestoreSVG)
)

// DrawableSVG makes an SVG conform to the Drawable interface.
type DrawableSVG struct {
	SVG  *SVG
	Size Size
}

// SVG holds an SVG.
type SVG struct {
	size          Size
	unscaledPath  *Path
	scaledPathMap map[Size]*Path
}

// MustSVG creates a new SVG the given svg path string (the contents of a single "d" attribute from an SVG "path"
// element) and panics if an error would be generated. The 'size' should be gotten from the original SVG's 'viewBox'
// parameter.
func MustSVG(size Size, svg string) *SVG {
	return mylog.Check2(NewSVG(size, svg))
}

// NewSVG creates a new SVG the given svg path string (the contents of a single "d" attribute from an SVG "path"
// element). The 'size' should be gotten from the original SVG's 'viewBox' parameter.
func NewSVG(size Size, svg string) (*SVG, error) {
	unscaledPath := mylog.Check2(NewPathFromSVGString(svg))

	return &SVG{
		size:          size,
		unscaledPath:  unscaledPath,
		scaledPathMap: make(map[Size]*Path),
	}, nil
}

// MustSVGFromContentString creates a new SVG and panics if an error would be generated. The content should contain
// valid SVG file data. Note that this only reads a very small subset of an SVG currently. Specifically, the "viewBox"
// attribute and any "d" attributes from enclosed SVG "path" elements.
func MustSVGFromContentString(content string) *SVG {
	return mylog.Check2(NewSVGFromContentString(content))
}

// NewSVGFromContentString creates a new SVG. The content should contain valid SVG file data. Note that this only reads
// a very small subset of an SVG currently. Specifically, the "viewBox" attribute and any "d" attributes from enclosed
// SVG "path" elements.
func NewSVGFromContentString(content string) (*SVG, error) {
	return NewSVGFromReader(strings.NewReader(content))
}

// MustSVGFromReader creates a new SVG and panics if an error would be generated. The reader should contain valid SVG
// file data. Note that this only reads a very small subset of an SVG currently. Specifically, the "viewBox" attribute
// and any "d" attributes from enclosed SVG "path" elements.
func MustSVGFromReader(r io.Reader) *SVG {
	return mylog.Check2(NewSVGFromReader(r))
}

// NewSVGFromReader creates a new SVG. The reader should contain valid SVG file data. Note that this only reads a very
// small subset of an SVG currently. Specifically, the "viewBox" attribute and any "d" attributes from enclosed SVG
// "path" elements.
func NewSVGFromReader(r io.Reader) (*SVG, error) {
	var svgXML struct {
		ViewBox string `xml:"viewBox,attr"`
		Paths   []struct {
			Path string `xml:"d,attr"`
		} `xml:"path"`
	}
	mylog.Check(xml.NewDecoder(r).Decode(&svgXML))
	svg := &SVG{scaledPathMap: make(map[Size]*Path)}
	var width, height string
	if parts := strings.Split(svgXML.ViewBox, " "); len(parts) == 4 {
		width = parts[2]
		height = parts[3]
	}
	v := mylog.Check2(strconv.ParseFloat(width, 64))
	if v < 1 || v > 4096 {
		mylog.Check("unable to determine SVG width")
	}
	svg.size.Width = float32(v)
	v = mylog.Check2(strconv.ParseFloat(height, 64))
	if v < 1 || v > 4096 {
		mylog.Check("unable to determine SVG height")
	}
	svg.size.Height = float32(v)
	for _, svgPath := range svgXML.Paths {
		p := mylog.Check2(NewPathFromSVGString(svgPath.Path))
		if svg.unscaledPath == nil {
			svg.unscaledPath = p
		} else {
			svg.unscaledPath.Path(p, false)
		}
	}
	return svg, nil
}

// Size returns the original size.
func (s *SVG) Size() Size {
	return s.size
}

// OffsetToCenterWithinScaledSize returns the scaled offset values to use to keep the image centered within the given
// size.
func (s *SVG) OffsetToCenterWithinScaledSize(size Size) Point {
	scale := min(size.Width/s.size.Width, size.Height/s.size.Height)
	return NewPoint((size.Width-s.size.Width*scale)/2, (size.Height-s.size.Height*scale)/2)
}

// PathScaledTo returns the path with the specified scaling. You should not modify this path, as it is cached.
func (s *SVG) PathScaledTo(scale float32) *Path {
	if scale == 1 {
		return s.unscaledPath
	}
	scaledSize := NewSize(scale, scale)
	p, ok := s.scaledPathMap[scaledSize]
	if !ok {
		p = s.unscaledPath.NewScaled(scale, scale)
		s.scaledPathMap[scaledSize] = p
	}
	return p
}

// PathForSize returns the path scaled to fit in the specified size. You should not modify this path, as it is cached.
func (s *SVG) PathForSize(size Size) *Path {
	return s.PathScaledTo(min(size.Width/s.size.Width, size.Height/s.size.Height))
}

// LogicalSize implements the Drawable interface.
func (s *DrawableSVG) LogicalSize() Size {
	return s.Size
}

// DrawInRect implements the Drawable interface.
func (s *DrawableSVG) DrawInRect(canvas *Canvas, rect Rect, _ *SamplingOptions, paint *Paint) {
	canvas.Save()
	defer canvas.Restore()
	offset := s.SVG.OffsetToCenterWithinScaledSize(rect.Size)
	canvas.Translate(rect.X+offset.X, rect.Y+offset.Y)
	canvas.DrawPath(s.SVG.PathForSize(rect.Size), paint)
}
