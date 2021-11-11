package entities

type Rect struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

// Basic entity type which can be renderred in SDL
type Entity struct {
	id    string
	rect  Rect
	color uint32
}

func Create(id string, rect Rect, color uint32) *Entity {
	return &Entity{
		id:    id,
		color: color,
		rect:  rect,
	}
}

func (e *Entity) Id() string {
	return e.id
}

func (e *Entity) GetX() int32 {
	return e.rect.X
}

func (e *Entity) SetX(x int32) {
	e.rect.X = x
}

func (e *Entity) GetY() int32 {
	return e.rect.Y
}

func (e *Entity) SetY(y int32) {
	e.rect.Y = y
}

func (e *Entity) GetWidth() int32 {
	return e.rect.Width
}

func (e *Entity) SetWidth(width int32) {
	e.rect.Width = width
}

func (e *Entity) GetHeight() int32 {
	return e.rect.Height
}

func (e *Entity) SetHeight(height int32) {
	e.rect.Height = height
}

func (e *Entity) GetColor() uint32 {
	return e.color
}

func (e *Entity) SetColor(color uint32) {
	e.color = color
}
