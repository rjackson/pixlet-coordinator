package main

import (
	"pixlet-coordinator/lib/rgbmatrix"
)

var apps = []App{
	{
		title:   "Burger of the Day",
		path:    "packages/tidbyt-community/apps/burgeroftheday/burgerotd.star",
		config:  map[string]string{},
		timeout: 30000,
	},
	{
		title:   "Snake",
		path:    "packages/tidbyt-community/apps/snake/snake.star",
		config:  map[string]string{},
		timeout: 30000,
	},
	{
		title:   "Amazing",
		path:    "packages/tidbyt-community/apps/amazing/amazing.star",
		config:  map[string]string{},
		timeout: 30000,
	},
	{
		title:   "Arcade Classics",
		path:    "packages/tidbyt-community/apps/arcadeclassics/arcade_classics.star",
		config:  map[string]string{},
		timeout: 30000,
	},
	{
		title:   "Runescape",
		path:    "packages/tidbyt-community/apps/runescape/runescape.star",
		config:  map[string]string{},
		timeout: 30000,
	},
	{
		title:   "Starfield",
		path:    "packages/tidbyt-community/apps/starfield/starfield.star",
		config:  map[string]string{},
		timeout: 30000,
	},
	{
		title:   "DaY nIgHt MaP",
		path:    "packages/tidbyt-community/apps/daynightmap/day_night_map.star",
		config:  map[string]string{},
		timeout: 30000,
	},
}

func createTextApp(text string) App {
	return App{
		title: "Blah",
		path:  "apps/text/text.star",
		config: map[string]string{
			"text": text,
		},
		timeout: 30000,
	}
}

type BuiltApp struct {
	Title App
	App   App
}

func main() {
	// go server.Serve()

	config := &rgbmatrix.HardwareConfig{
		Rows:                   32,
		Cols:                   64,
		Parallel:               1,
		ChainLength:            1,
		Brightness:             50,
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

	var builtApps = []BuiltApp{}

	// Prebuild apps
	for _, app := range apps {
		go func(app App) {
			titleApp := createTextApp(app.title)

			app.Build()
			titleApp.Build()

			builtApps = append(builtApps, BuiltApp{
				Title: titleApp,
				App:   app,
			})
		}(app)
	}

	// Interrupt rendering frequently to switch apps
	// go func() {
	// for {
	// time.Sleep(2 * time.Second)
	// interrupt <- true
	// }
	// }()

	for {
		for _, builtApp := range builtApps {
			builtApp.Title.RenderToDisplay(c, interrupt)
			builtApp.App.RenderToDisplay(c, interrupt)
		}
	}
}
