package util

import "sync"

// Barrier for handling "complex" process synchronization
type Barrier struct {
	execWg    *sync.WaitGroup
	waitWg    *sync.WaitGroup
	resolveWg *sync.WaitGroup
	raceWg    *sync.WaitGroup
}

// Create a new barrier.
func NewBarrier() *Barrier {
	b := &Barrier{
		execWg:    &sync.WaitGroup{},
		waitWg:    &sync.WaitGroup{},
		resolveWg: &sync.WaitGroup{},
		raceWg:    &sync.WaitGroup{},
	}
	b.waitWg.Add(1)
	return b
}

// Let the barrier tick. This means, that all waiting processes can pass the barrier.
// Note: They have to call Barrier.Resolve() after finishing their work, otherwise
// this barrier gets stuck.
func (b *Barrier) Tick(size int) {
	b.execWg.Add(size)
	b.resolveWg.Add(1)
	b.waitWg.Done()
	b.execWg.Wait()
	b.waitWg.Add(1)
	b.raceWg.Add(size)
	b.resolveWg.Done()
	b.raceWg.Wait()
}

// Wait for this barrier to "drop". After this you can perform computing action.
func (b *Barrier) Wait() {
	b.waitWg.Wait()
}

// Resolve the current execution of this process.
func (b *Barrier) Resolve() {
	b.execWg.Done()
	b.resolveWg.Wait()
	b.raceWg.Done()
}
