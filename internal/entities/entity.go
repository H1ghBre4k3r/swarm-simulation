package entities

// Basic entity type which can be renderred in SDL
type Entity struct {
	id     interface{}
	x      int32
	y      int32
	width  int32
	height int32
	color  uint32
}

func New(id interface{}) *Entity {
	return &Entity{
		id:     id,
		x:      100,
		y:      0,
		width:  200,
		height: 200,
		color:  0xffff0000,
	}
}

func (e *Entity) Id() interface{} {
	return e.id
}

func (e *Entity) GetX() int32 {
	return e.x
}

func (e *Entity) SetX(x int32) {
	e.x = x
}

func (e *Entity) GetY() int32 {
	return e.y
}

func (e *Entity) SetY(y int32) {
	e.y = y
}

func (e *Entity) GetWidth() int32 {
	return e.width
}

func (e *Entity) SetWidth(width int32) {
	e.width = width
}

func (e *Entity) GetHeight() int32 {
	return e.height
}

func (e *Entity) SetHeight(height int32) {
	e.height = height
}

func (e *Entity) GetColor() uint32 {
	return e.color
}

func (e *Entity) SetColor(color uint32) {
	e.color = color
}
