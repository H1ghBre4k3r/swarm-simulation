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
	consensus     bool
}

func CreateSimulationPortal(spatial *collision.SpatialHashmap, entities *entities.EntityManger, noise float64, tau float64, obstacles []*obstacles.Obstacle, consensus bool) *SimulationPortal {
	portal := &SimulationPortal{
		spatial:       spatial,
		entityManager: entities,
		noise:         noise,
		tau:           tau,
		obstacles:     obstacles,
		consensus:     consensus,
	}
	portal.Update()
	return portal
}

func (p *SimulationPortal) Update() {
	ents := p.entityManager.Get()
	p.entities = []*entities.Entity{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, e := range ents {
		newEntity := e.Copy()
		newEntity.Move()
		// Noise positions of all participants, so everyone has a "common" state
		newEntity.NoisePosition(r.NormFloat64() * p.noise)
		p.entities = append(p.entities, newEntity)
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
func (p *SimulationPortal) Consensus() bool {
	return p.consensus
}
