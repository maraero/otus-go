package hw05parallelexecution

import (
	"sync/atomic"
)

type errCounter struct {
	counter int32
	limit   int32
}

func (c *errCounter) inc() {
	atomic.AddInt32(&c.counter, 1)
}

func (c *errCounter) exceedsLimit() bool {
	if c.limit <= 0 {
		return false
	}

	return atomic.LoadInt32(&c.counter) >= c.limit
}

func newErrCounter(limit int) *errCounter {
	return &errCounter{limit: int32(limit)}
}
