package collision

import (
	"sync"
)

type Insertable interface {
	GetX() int32
	GetY() int32
}

// Map for storing information about entities in 2d space
type SpatialHashmap struct {
	cellSize int32
	contents sync.Map
}

func New(cellSize int32) *SpatialHashmap {
	return &SpatialHashmap{
		cellSize: cellSize,
		// contents: make(map[int]map[int][]interface{}),
	}
}

func (m *SpatialHashmap) hash(point Insertable) (int, int) {
	return int(point.GetX() / m.cellSize), int(point.GetY() / m.cellSize)
}

func (m *SpatialHashmap) Insert(entity Insertable) {
	x, y := m.hash(entity)
	ys, _ := m.contents.LoadOrStore(x, &sync.Map{})
	list, _ := ys.(*sync.Map).LoadOrStore(y, []interface{}{})
	ys.(*sync.Map).Store(y, append(list.([]interface{}), entity))
}

func (m *SpatialHashmap) Remove(entity Insertable) interface{} {
	x, y := m.hash(entity)
	ys, _ := m.contents.LoadOrStore(x, &sync.Map{})
	list, _ := ys.(*sync.Map).LoadOrStore(y, []interface{}{})

	for i, e := range list.([]interface{}) {
		if e == entity {
			list.([]interface{})[i] = list.([]interface{})[len(list.([]interface{}))-1]
			ys.(*sync.Map).Store(y, list.([]interface{})[:len(list.([]interface{}))-1])
			return e
		}
	}
	return nil
}
