package main

import (
	rgbmatrix "github.com/jmaitrehenry/go-rpi-rgb-led-matrix"
)

var apps = []App{
	{
		path:    "packages/tidbyt-community/apps/burgeroftheday/burgerotd.star",
		config:  map[string]string{},
		timeout: 30000,
	},
	{
		path:    "packages/tidbyt-community/apps/snake/snake.star",
		config:  map[string]string{},
		timeout: 30000,
	},
	{
		path:    "packages/tidbyt-community/apps/amazing/amazing.star",
		config:  map[string]string{},
		timeout: 30000,
	},
}

func main() {
	config := &rgbmatrix.HardwareConfig{
		Rows:                   32,
		Cols:                   64,
		Parallel:               1,
		ChainLength:            1,
		Brightness:             100,
		HardwareMapping:        "regular",
		ShowRefreshRate:        false,
		InverseColors:          false,
		DisableHardwarePulsing: true,
		// PanelType:              "FM6127",
	}

	m, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		panic(err)
	}

	c := rgbmatrix.NewCanvas(m)
	defer c.Close()

	for {
		for _, app := range apps {
			app.Render(c)
		}
	}
}
