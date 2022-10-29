package hw05parallelexecution

import "sync"

type errCounter struct {
	mu      sync.Mutex
	counter int
	limit int
}

func (c *errCounter) inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counter++
}

func (c *errCounter) exceedsLimit() bool {
	if c.limit == 0 {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counter >= c.limit
}

func newErrCouner(limit int) *errCounter {
	return &errCounter{limit: limit}
}