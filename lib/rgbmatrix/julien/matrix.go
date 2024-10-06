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
	"github.com/aybabtme/rgbterm"
	"image/color"
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
func (e Row) String() string {
	out := "|"

	for _, elem := range e {
		if elem.Set {
			out += rgbterm.FgString("*", elem.Color.R, elem.Color.G, elem.Color.B)
		} else {
			out += " "
		}
	}
	out += "|\n"
	return out
}

// Matrix is a logical, or (0, 1)-matrix
type Matrix struct {
	Matrix []Row
}

// Size returns the dimensions of the matrix in (Rows, Columns) order.
func (e Matrix) Size() (int, int) {
	if len(e.Matrix) == 0 {
		return 0, 0
	} else {
		return len(e.Matrix), e.Matrix[0].Size()
	}
}

// String converts the matrix to space-and-dot notation.
func (e Matrix) String() string {
	out := []rune{}
	_, b := e.Size()

	addBar := func() {
		for i := -2; i < b; i++ {
			out = append(out, '-')
		}
		out = append(out, '\n')
	}

	addBar()
	for _, row := range e.Matrix {
		out = append(out, []rune(row.String())...)
	}
	addBar()

	return string(out)
}

// GenerateEmpty generates the n-by-n matrix with all entries set to 0.
func GenerateEmpty(height, width int) Matrix {
	out := make([]Row, height)

	for i := 0; i < height; i++ {
		out[i] = NewRow(width)
	}

	return Matrix{out}
}

func (e *Matrix) Geometry() (width, height int) {
	return e.Size()
}

func (e *Matrix) position(x, y int) int {
	return x + (y * len(e.Matrix))
}

func (t *Matrix) Apply(leds []color.Color) error {
	for position, l := range leds {
		t.Set(position, l)
	}

	return t.Render()
}

func (e *Matrix) Render() error {
	fmt.Println(e.String())
	newMatrix := GenerateEmpty(e.Size())
	e.Matrix = newMatrix.Matrix
	return nil
}

func (e *Matrix) At(position int) color.Color {
	h, w := e.Size()
	posY := position / h
	posX := position % w

	return e.Matrix[posY][posX].Color
}

func (e *Matrix) Set(position int, c color.Color) {
	h, w := e.Size()
	posY := position / h
	posX := position % w

	r, g, b, a := c.RGBA()
	cc := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	e.Matrix[posY][posX] = point{Set: true, Color: cc}
}

// Those new functions have no use with the emulator
func (t *Matrix) Close() error {
	return nil
}

func (t *Matrix) GetBrightness() int {
	return 0
}

func (t *Matrix) SetBrightness(brightness int) {}
