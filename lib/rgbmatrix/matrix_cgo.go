//go:build linux
// +build linux

/*
Copyright (c) 2016 Máximo Cuadros
Copyright (c) 2020 TFK1410
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

/*
#cgo CFLAGS: -std=c99 -I${SRCDIR}/lib/rpi-rgb-led-matrix/include -DSHOW_REFRESH_RATE
#cgo LDFLAGS: -lrgbmatrix -L${SRCDIR}/lib/rpi-rgb-led-matrix/lib -lstdc++ -lm
#include "lib/rpi-rgb-led-matrix/include/led-matrix-c.h"

void led_matrix_swap(struct RGBLedMatrix *matrix, struct LedCanvas *offscreen_canvas,
                     int width, int height, const uint32_t pixels[]) {


  int i, x, y;
  uint32_t color;
  for (x = 0; x < width; ++x) {
    for (y = 0; y < height; ++y) {
      i = x + (y * width);
      color = pixels[i];

      led_canvas_set_pixel(offscreen_canvas, x, y,
        (color >> 16) & 255, (color >> 8) & 255, color & 255);
    }
  }

  offscreen_canvas = led_matrix_swap_on_vsync(matrix, offscreen_canvas);
}

void set_show_refresh_rate(struct RGBLedMatrixOptions *o, int show_refresh_rate) {
  o->show_refresh_rate = show_refresh_rate != 0 ? 1 : 0;
}

void set_disable_hardware_pulsing(struct RGBLedMatrixOptions *o, int disable_hardware_pulsing) {
  o->disable_hardware_pulsing = disable_hardware_pulsing != 0 ? 1 : 0;
}

void set_inverse_colors(struct RGBLedMatrixOptions *o, int inverse_colors) {
  o->inverse_colors = inverse_colors != 0 ? 1 : 0;
}
*/
import "C"
import (
	"fmt"
	"unsafe"

	"image/color"
)

type uint32_t C.uint32_t

func (c *HardwareConfig) toC() *C.struct_RGBLedMatrixOptions {
	o := &C.struct_RGBLedMatrixOptions{}
	o.rows = C.int(c.Rows)
	o.cols = C.int(c.Cols)
	o.chain_length = C.int(c.ChainLength)
	o.parallel = C.int(c.Parallel)
	o.pwm_bits = C.int(c.PWMBits)
	o.pwm_lsb_nanoseconds = C.int(c.PWMLSBNanoseconds)
	o.brightness = C.int(c.Brightness)
	o.scan_mode = C.int(c.ScanMode)
	o.hardware_mapping = C.CString(c.HardwareMapping)
	o.pixel_mapper_config = C.CString(c.PixelMapperConfig)

	if c.ShowRefreshRate == true {
		C.set_show_refresh_rate(o, C.int(1))
	} else {
		C.set_show_refresh_rate(o, C.int(0))
	}

	if c.DisableHardwarePulsing == true {
		C.set_disable_hardware_pulsing(o, C.int(1))
	} else {
		C.set_disable_hardware_pulsing(o, C.int(0))
	}

	if c.InverseColors == true {
		C.set_inverse_colors(o, C.int(1))
	} else {
		C.set_inverse_colors(o, C.int(0))
	}

	return o
}

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

	m := C.led_matrix_create_from_options(config.toC(), nil, nil)
	b := C.led_matrix_create_offscreen_canvas(m)

	var w, h C.int
	C.led_canvas_get_size(b, &w, &h)

	c = &RGBLedMatrix{
		Config: config,
		width:  int(w), height: int(h),
		matrix: m,
		buffer: b,
		leds:   make([]uint32_t, int(w)*int(h)),
	}
	if m == nil {
		return nil, fmt.Errorf("unable to allocate memory")
	}

	return c, nil
}

// Render update the display with the data from the LED buffer
func (c *RGBLedMatrix) Render() error {
	w, h := c.Geometry()

	C.led_matrix_swap(
		c.matrix,
		c.buffer,
		C.int(w), C.int(h),
		(*C.uint32_t)(unsafe.Pointer(&c.leds[0])),
	)

	c.leds = make([]uint32_t, w*h)
	return nil
}

// At return a Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *RGBLedMatrix) At(x, y int) color.Color {
	return uint32ToColor(c.leds[c.position(x, y)])
}

// Set LED at position x,y to the provided 24-bit color value.
func (c *RGBLedMatrix) Set(x, y int, color color.Color) {
	c.leds[c.position(x, y)] = uint32_t(colorToUint32(color))
}

// Close finalizes the ws281x interface
func (c *RGBLedMatrix) Close() error {
	C.led_matrix_delete(c.matrix)
	return nil
}

// GetBrightness returns the current brightness setting of the matrix
func (c *RGBLedMatrix) GetBrightness() int {
	return int(C.led_matrix_get_brightness(c.matrix))
}

// SetBrightness sets a new brightness setting to the matrix
func (c *RGBLedMatrix) SetBrightness(brightness int) {
	C.led_matrix_set_brightness(c.matrix, C.uchar(brightness))
}

// Apply set all the pixels to the values contained in leds
func (c *RGBLedMatrix) Apply(leds []color.Color) error {
	for position, l := range leds {
		x, y = position%c.w, position/c.w
		c.Set(x, y, l)
	}

	return c.Render()
}

func (c *RGBLedMatrix) position(x, y int) int {
	return x + (y * c.w)
}
