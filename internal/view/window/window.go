package window

import (
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/simulation"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	scale    int32
	window   *sdl.Window
	renderer *sdl.Renderer
	lines    bool
}

// Create a new window
func New(title string, scale int32, lines bool) (*Window, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, scale, scale, sdl.WINDOW_SHOWN)
	if err != nil {
		sdl.Quit()
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, 0, sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Destroy()
		sdl.Quit()
		return nil, err
	}

	return &Window{
		scale,
		window,
		renderer,
		lines,
	}, nil
}

// Destroy this window
func (w *Window) Destroy() {
	w.window.Destroy()
	sdl.Quit()
}

// Render content to the screen/window
func (w *Window) Render(entities []simulation.Drawable) {
	w.renderer.SetDrawColor(255, 255, 255, 255)
	w.renderer.Clear()
	if w.lines {
		width := w.scale / 10
		for i := int32(1); i < 10; i++ {
			gfx.LineRGBA(w.renderer, 0, i*width, w.scale, i*width, 210, 210, 210, 255)
			gfx.LineRGBA(w.renderer, i*width, 0, i*width, w.scale, 210, 210, 210, 255)

		}
	}
	for _, e := range entities {
		gfx.FilledCircleRGBA(w.renderer, int32(e.GetX()*float64(w.scale)), int32(e.GetY()*float64(w.scale)), int32(e.GetR()*float64(w.scale)), uint8((e.GetColor()>>24)&0xff), uint8((e.GetColor()>>16)&0xff), uint8((e.GetColor()>>8)&0xff), 255)
		gfx.ThickLineRGBA(w.renderer, int32(e.GetX()*float64(w.scale)), int32(e.GetY()*float64(w.scale)), int32((e.GetX()+e.GetVelocity().X)*float64(w.scale)), int32((e.GetY()+e.GetVelocity().Y)*float64(w.scale)), 2, 0, 255, 0, 255)
	}
	w.renderer.Present()
}
