package entities

type EntityManger struct {
	entities []*Entity
}

func Manager() *EntityManger {
	return &EntityManger{
		entities: []*Entity{},
	}
}

func (m *EntityManger) Add(entities ...*Entity) {
	m.entities = append(m.entities, entities...)
}

func (m *EntityManger) Get() []*Entity {
	return m.entities
}
