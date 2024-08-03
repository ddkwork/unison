// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom

import (
	"fmt"
	"math"
	"reflect"

	"github.com/richardwilkes/toolbox/xmath"
)

// Pt32 is an alias for the float32 version of Point.
type Pt32 = Point[float32]

// Pt64 is an alias for the float64 version of Point.
type Pt64 = Point[float64]

// Point defines a location.
type Point[T xmath.Numeric] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// NewPoint creates a new Point.
func NewPoint[T xmath.Numeric](x, y T) Point[T] {
	return Point[T]{
		X: x,
		Y: y,
	}
}

// NewPointPtr creates a new *Point.
func NewPointPtr[T xmath.Numeric](x, y T) *Point[T] {
	p := NewPoint[T](x, y)
	return &p
}

// Align modifies this Point to align with integer coordinates. Returns itself for easy chaining.
func (p *Point[T]) Align() *Point[T] {
	switch reflect.TypeOf(p.X).Kind() {
	case reflect.Float32, reflect.Float64:
		p.X = T(math.Floor(reflect.ValueOf(p.X).Float()))
		p.Y = T(math.Floor(reflect.ValueOf(p.Y).Float()))
	}
	return p
}

// Add modifies this Point by adding the supplied coordinates. Returns itself for easy chaining.
//func (p *Point[T]) Add(pt Point[T]) *Point[T] {
//	p.X += pt.X
//	p.Y += pt.Y
//	return p
//}

// Subtract modifies this Point by subtracting the supplied coordinates. Returns itself for easy chaining.
func (p *Point[T]) Subtract(pt Point[T]) *Point[T] {
	p.X -= pt.X
	p.Y -= pt.Y
	return p
}

// Negate modifies this Point by negating both the X and Y coordinates.
func (p *Point[T]) Negate() *Point[T] {
	p.X = -p.X
	p.Y = -p.Y
	return p
}

// String implements the fmt.Stringer interface.
func (p *Point[T]) String() string {
	return fmt.Sprintf("%v,%v", p.X, p.Y)
}

func (p Point[T]) toPoint64() Point[float64] {
	return Point[float64]{
		X: reflect.ValueOf(p.X).Float(),
		Y: reflect.ValueOf(p.Y).Float(),
	}
}

// ConvertPoint converts a Point of type F into one of type T.
func ConvertPoint[T, F xmath.Numeric](pt Point[F]) Point[T] {
	return NewPoint(T(pt.X), T(pt.Y))
}

// Add returns a new Point which is the result of adding this Point with the provided Point.
func (p Point[T]) Add(pt Point[T]) Point[T] {
	return Point[T]{X: p.X + pt.X, Y: p.Y + pt.Y}
}

// Sub returns a new Point which is the result of subtracting the provided Point from this Point.
func (p Point[T]) Sub(pt Point[T]) Point[T] {
	return Point[T]{X: p.X - pt.X, Y: p.Y - pt.Y}
}

// Mul returns a new Point which is the result of multiplying the coordinates of this point by the value.
func (p Point[T]) Mul(value T) Point[T] {
	return Point[T]{X: p.X * value, Y: p.Y * value}
}

// Div returns a new Point which is the result of dividing the coordinates of this point by the value.
func (p Point[T]) Div(value T) Point[T] {
	return Point[T]{X: p.X / value, Y: p.Y / value}
}

// Neg returns a new Point that holds the negated coordinates of this Point.
func (p Point[T]) Neg() Point[T] {
	return Point[T]{X: -p.X, Y: -p.Y}
}

// Floor returns a new Point which is aligned to integer coordinates by using Floor on them.
func (p Point[T]) Floor() Point[T] {
	return Point[T]{X: xmath.Floor(p.X), Y: xmath.Floor(p.Y)}
}

// Ceil returns a new Point which is aligned to integer coordinates by using Ceil() on them.
func (p Point[T]) Ceil() Point[T] {
	return Point[T]{X: xmath.Ceil(p.X), Y: xmath.Ceil(p.Y)}
}

// Dot returns the dot product of the two Points.
func (p Point[T]) Dot(pt Point[T]) T {
	return p.X*pt.X + p.Y*pt.Y
}

// Cross returns the cross product of the two Points.
func (p Point[T]) Cross(pt Point[T]) T {
	return p.X*pt.Y - p.Y*pt.X
}

// In returns true if this Point is within the Rect.
func (p Point[T]) In(r Rect[T]) bool {
	if r.Empty() {
		return false
	}
	return r.X <= p.X && r.Y <= p.Y && p.X < r.Right() && p.Y < r.Bottom()
}
