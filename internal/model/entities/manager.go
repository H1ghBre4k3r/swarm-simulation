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

// Get a pair of all running participants and all "auto executors".
// running participants := all participants, that did not reach their target position & do not have "IgnoreFinish" flag
// auto executors := all participants, that have the flag "IgnoreFinish" flag set
func (m *EntityManger) GetRunning() ([]*Entity, []*Entity) {
	ents := m.Get()
	active := []*Entity{}
	autoExecutors := []*Entity{}
	for _, e := range ents {
		if e.ignoreFinish {
			autoExecutors = append(autoExecutors, e)
		} else if e.IsRunning() {
			active = append(active, e)
		}
	}
	return active, autoExecutors
}
