//go:build !linux
// +build !linux

/*
Copyright (c) 2016 Máximo Cuadros
Copyright (c) 2023 Julien Maitrehenry

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
*/

package rgbmatrix

import "C"
import "fmt"

type uint32_t uint32

// NewRGBLedMatrix returns a new matrix using the given size and config
func NewRGBLedMatrix(config *HardwareConfig) (c Matrix, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("error creating matrix: %v", r)
			}
		}
	}()

	if isJulienEmulator() {
		return buildJulienMatrix(config), nil
	}

	panic("No emulator found")
}
