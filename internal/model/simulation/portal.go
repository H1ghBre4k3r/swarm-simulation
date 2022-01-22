package simulation

import (
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
	fps           uint64
}

func CreateSimulationPortal(spatial *collision.SpatialHashmap, entities *entities.EntityManger, noise float64, fps uint64, obstacles []*obstacles.Obstacle) *SimulationPortal {
	portal := &SimulationPortal{
		spatial:       spatial,
		entityManager: entities,
		noise:         noise,
		fps:           fps,
		obstacles:     obstacles,
	}
	portal.Update()
	return portal
}

func (p *SimulationPortal) Update() {
	ents := p.entityManager.Get()
	p.entities = []*entities.Entity{}
	for _, e := range ents {
		newE := e.Copy()
		newE.Move()
		p.entities = append(p.entities, e)
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

func (p *SimulationPortal) FPS() uint64 {
	return p.fps
}

func (p *SimulationPortal) Obstacles() []*obstacles.Obstacle {
	return p.obstacles
}
