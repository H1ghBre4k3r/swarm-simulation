package window

import "github.com/veandco/go-sdl2/sdl"

type Window struct {
	window  *sdl.Window
	surface *sdl.Surface
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
func (w *Window) Render(rect *sdl.Rect) {
	w.surface.FillRect(nil, 0)
	w.surface.FillRect(rect, 0xffff0000)
	w.window.UpdateSurface()
}
