package window

import "github.com/veandco/go-sdl2/sdl"

type Window struct {
	window  *sdl.Window
	surface *sdl.Surface
}

type Drawable interface {
	GetX() int32
	GetY() int32
	GetWidth() int32
	GetHeight() int32
	GetColor() uint32
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

	surface, err := window.GetSurface()
	if err != nil {
		window.Destroy()
		sdl.Quit()
		return nil, err
	}

	return &Window{
		window,
		surface,
	}, nil
}

// Destroy this window
func (w *Window) Destroy() {
	w.window.Destroy()
	sdl.Quit()
}

// Render content to the screen/window
func (w *Window) Render(entities []Drawable) {
	w.surface.FillRect(nil, 0)
	for _, e := range entities {
		w.surface.FillRect(&sdl.Rect{
			X: e.GetX(),
			Y: e.GetY(),
			W: e.GetWidth(),
			H: e.GetHeight(),
		}, e.GetColor())
	}
	w.window.UpdateSurface()
}
