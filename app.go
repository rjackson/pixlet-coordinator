package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"time"

	rgbmatrix "github.com/jmaitrehenry/go-rpi-rgb-led-matrix"
	"tidbyt.dev/pixlet/render"
	"tidbyt.dev/pixlet/runtime"
	"tidbyt.dev/pixlet/tools"
)

type App struct {
	path    string
	config  map[string]string
	timeout int
}

func (a *App) Render(c *rgbmatrix.Canvas) {
	ctx, _ := context.WithTimeoutCause(
		context.Background(),
		time.Duration(a.timeout)*time.Millisecond,
		fmt.Errorf("timeout after %dms", a.timeout),
	)

	fs := tools.NewSingleFileFS(a.path)
	applet, _ := runtime.NewAppletFromFS("app-id", fs)
	roots, _ := applet.RunWithConfig(ctx, a.config)

	images := render.PaintRoots(true, roots...)

	var delay int32 = 50
	if len(roots) > 0 && roots[0].Delay > 0 {
		delay = roots[0].Delay
	}

	fmt.Println("delay", delay)

	for _, im := range images {
		frameDuration := time.Duration(delay) * time.Millisecond
		draw.Draw(c, c.Bounds(), im, image.Point{}, draw.Src)
		c.Render()
		time.Sleep(frameDuration)
	}
}
