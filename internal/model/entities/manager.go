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

func (m *EntityManger) GetRunning() []*Entity {
	ents := m.Get()
	active := []*Entity{}
	for _, e := range ents {
		if e.IsRunning() {
			active = append(active, e)
		}
	}
	return active
}
