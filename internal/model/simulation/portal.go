package simulation

import (
	"math/rand"
	"time"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/collision"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/entities"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/obstacles"
)

type SimulationPortal struct {
	spatial       *collision.SpatialHashmap
	entityManager *entities.EntityManger
	entities      []*entities.Entity
	obstacles     []*obstacles.Obstacle
	noise         float64
	tau           float64
}

func CreateSimulationPortal(spatial *collision.SpatialHashmap, entities *entities.EntityManger, noise float64, tau float64, obstacles []*obstacles.Obstacle) *SimulationPortal {
	portal := &SimulationPortal{
		spatial:       spatial,
		entityManager: entities,
		noise:         noise,
		tau:           tau,
		obstacles:     obstacles,
	}
	portal.Update()
	return portal
}

func (p *SimulationPortal) Update() {
	ents := p.entityManager.Get()
	p.entities = []*entities.Entity{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, e := range ents {
		newE := e.Copy()
		newE.Move()
		// Noise positions of all participants, so everyone has a "common" state
		newE.NoisePosition(r.NormFloat64() * p.noise)
		p.entities = append(p.entities, newE)
	}
}

func (p *SimulationPortal) Insert(entity *entities.Entity) {
	p.spatial.Insert(entity)
}

func (p *SimulationPortal) Remove(entity *entities.Entity) {
	p.spatial.Remove(entity)
}

func (p *SimulationPortal) Participants() []*entities.Entity {
	return p.entities
}

func (p *SimulationPortal) Noise() float64 {
	return p.noise
}

func (p *SimulationPortal) TAU() float64 {
	return p.tau
}

func (p *SimulationPortal) Obstacles() []*obstacles.Obstacle {
	return p.obstacles
}
