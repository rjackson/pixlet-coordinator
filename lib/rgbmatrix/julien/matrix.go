/*
Copyright (c) 2016 Máximo Cuadros

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
*/

package julien

import (
	"fmt"
	"image/color"

	"github.com/aybabtme/rgbterm"
)

type point struct {
	Set   bool
	Color color.RGBA
}

// A binary row / vector in GF(2)^n.
type Row []point

// NewRow returns an empty n-component row.
func NewRow(n int) Row {
	return Row(make([]point, n))
}

// Size returns the dimension of the vector.
func (e Row) Size() int {
	return len(e)
}

// String converts the row into space-and-dot notation.
func (e Row) String() (out string) {
	for _, elem := range e {
		if elem.Set {
			out += rgbterm.FgString("██", elem.Color.R, elem.Color.G, elem.Color.B)
		} else {
			out += "██"
		}
	}
	out += "\n"
	return out
}

// Matrix is a logical, or (0, 1)-matrix
type Matrix struct {
	Matrix []Row
}

// Size returns the dimensions of the matrix in (Width, Height) order.
func (e Matrix) Size() (int, int) {
	if len(e.Matrix) == 0 {
		return 0, 0
	} else {
		return e.Matrix[0].Size(), len(e.Matrix)
	}
}

// String converts the matrix to space-and-dot notation.
func (e Matrix) String() string {
	out := []rune{}
	out = append(out, []rune("\033[H")...)

	for _, row := range e.Matrix {
		out = append(out, []rune(row.String())...)
	}

	return string(out)
}

// GenerateEmpty generates the n-by-n matrix with all entries set to 0.
func GenerateEmpty(width, height int) Matrix {
	out := make([]Row, height)

	for i := 0; i < height; i++ {
		out[i] = NewRow(width)
	}

	return Matrix{out}
}

func (e *Matrix) Geometry() (width, height int) {
	return e.Size()
}

func (e *Matrix) Apply(leds []color.Color) error {
	w, h := e.Size()
	for position, l := range leds {
		y, x := position/h, position%w
		e.Set(x, y, l)
	}

	return e.Render()
}

func (e *Matrix) Render() error {
	fmt.Println(e.String())
	newMatrix := GenerateEmpty(e.Size())
	e.Matrix = newMatrix.Matrix
	return nil
}

func (e *Matrix) At(x, y int) color.Color {
	return e.Matrix[y][x].Color
}

func (e *Matrix) Set(x, y int, c color.Color) {
	// Out of bounds. Just happens apparently and we have to handle it
	if y >= len(e.Matrix) || x >= len(e.Matrix[0]) {
		return
	}

	r, g, b, a := c.RGBA()
	cc := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	e.Matrix[y][x] = point{Set: true, Color: cc}
}

// Those new functions have no use with the emulator
func (t *Matrix) Close() error {
	return nil
}

func (t *Matrix) GetBrightness() int {
	return 0
}

func (t *Matrix) SetBrightness(brightness int) {}
