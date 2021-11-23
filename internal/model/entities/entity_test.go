package entities_test

import (
	"testing"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/entities"
)

func TestNew(t *testing.T) {
	entity := entities.Create("someRandomId", entities.Position{}, 0x0, nil, nil, "")
	if entity.Id() != "someRandomId" {
		t.Fatal("Initial ID is not set correctly")
	}
}

func TestSetters(t *testing.T) {
	entity := entities.Create("", entities.Position{}, 0xffff0000, nil, nil, "")

	entity.SetX(200)
	if entity.GetX() != 200 {
		t.Fatal("SetX does not set the value returned by GetX")
	}

	entity.SetY(42)
	if entity.GetY() != 42 {
		t.Fatal("SetY does not set the value returned by SetY")
	}

	entity.SetR(100)
	if entity.GetR() != 100 {
		t.Fatal("SetR does not set the value returned by GetR")
	}

	entity.SetColor(0x20202020)
	if entity.GetColor() != 0x20202020 {
		t.Fatal("SetColor does not set the value returned by GetColor")
	}
}
