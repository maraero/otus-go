package hw05parallelexecution

import "sync"

type errCounter struct {
	sync.Mutex
	counter int
	limit   int
}

func (c *errCounter) inc() {
	c.Lock()
	defer c.Unlock()
	c.counter++
}

func (c *errCounter) exceedsLimit() bool {
	if c.limit <= 0 {
		return false
	}

	c.Lock()
	defer c.Unlock()
	return c.counter >= c.limit
}

func newErrCounter(limit int) *errCounter {
	return &errCounter{limit: limit}
}
