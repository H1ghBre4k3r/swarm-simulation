package window

import (
	"github.com/H1ghBre4k3r/swarm-simulation/internal/simulation"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

// Create a new window
func New(title string, width int32, height int32) (*Window, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
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
		gfx.FilledCircleRGBA(w.renderer, e.GetX(), e.GetY(), e.GetR(), 255, 0, 0, 255)
	}
	w.renderer.Present()
}
