package entities_test

import (
	"testing"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/entities"
)

func TestNew(t *testing.T) {
	entity := entities.Create("someRandomId", entities.Rect{}, 0x0)
	if entity.Id() != "someRandomId" {
		t.Fatal("Initial ID is not set correctly")
	}
}

func TestSetters(t *testing.T) {
	entity := entities.Create("", entities.Rect{}, 0xffff0000)

	entity.SetX(200)
	if entity.GetX() != 200 {
		t.Fatal("SetX does not set the value returned by GetX")
	}

	entity.SetY(42)
	if entity.GetY() != 42 {
		t.Fatal("SetY does not set the value returned by SetY")
	}

	entity.SetWidth(100)
	if entity.GetWidth() != 100 {
		t.Fatal("SetWidth does not set the value returned by GetWidth")
	}

	entity.SetHeight(360)
	if entity.GetHeight() != 360 {
		t.Fatal("SetHeight does not set the value returned by GetHeight")
	}

	entity.SetColor(0x20202020)
	if entity.GetColor() != 0x20202020 {
		t.Fatal("SetColor does not set the value returned by GetColor")
	}
}
