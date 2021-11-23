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
}

// Create a new window
func New(title string, scale int32) (*Window, error) {
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
	}, nil
}

// Destroy this window
func (w *Window) Destroy() {
	w.window.Destroy()
	sdl.Quit()
}

// Render content to the screen/window
func (w *Window) Render(entities []simulation.Drawable) {
	w.renderer.SetDrawColor(0, 0, 0, 255)
	w.renderer.Clear()
	for _, e := range entities {
		gfx.FilledCircleRGBA(w.renderer, int32(e.GetX()*float64(w.scale)), int32(e.GetY()*float64(w.scale)), int32(e.GetR()*float64(w.scale)), 255, 0, 0, 255)
	}
	w.renderer.Present()
}
