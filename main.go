package main

import (
	"math/rand"
	"pixlet-coordinator/lib/rgbmatrix"
	"time"
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
	{
		path:    "packages/tidbyt-community/apps/arcadeclassics/arcade_classics.star",
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

	interrupt := make(chan bool)

	c := rgbmatrix.NewCanvas(m)
	defer c.Close()

	for i := range apps {
		go func(i int) {
			apps[i].Build()
		}(i)
	}

	// Interrupt rendering frequently to switch apps
	go func() {
		for {
			time.Sleep(2 * time.Second)
			interrupt <- true
		}
	}()

	lastI := -1

	for {
		i := rand.Intn(len(apps))
		app := apps[i]

		if !app.Ready() || lastI == i {
				continue
			}

		lastI = i
		app.RenderToDisplay(c, interrupt)
	}
}
