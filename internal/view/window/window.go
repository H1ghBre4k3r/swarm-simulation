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
	fps      uint64
}

// Create a new window
func New(title string, scale int32, lines bool, fps uint64) (*Window, error) {
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
		fps,
	}, nil
}

func (w *Window) Scale(n float64) float64 {
	return float64(w.scale) * n
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
		gfx.FilledCircleRGBA(w.renderer, int32(w.Scale(e.GetX())), int32(w.Scale(e.GetY())), int32(w.Scale(e.GetR())), uint8(e.GetColor().Red), uint8(e.GetColor().Green), uint8(e.GetColor().Blue), 255)
		vel := e.GetVelocity()
		vel = *vel.Scale(float64(w.fps))
		gfx.ThickLineRGBA(w.renderer, int32(w.Scale(e.GetX())), int32(w.Scale(e.GetY())), int32(w.Scale(e.GetX()+vel.X)), int32(w.Scale(e.GetY()+vel.Y)), 2, 0, 255, 0, 255)

		gfx.LineRGBA(w.renderer, int32(w.Scale(e.GetX())), int32(w.Scale(e.GetY())), int32(w.Scale(e.GetTarget().X)), int32(w.Scale(e.GetTarget().Y)), 0, 0, 0, 10)
	}
	w.renderer.Present()
}
