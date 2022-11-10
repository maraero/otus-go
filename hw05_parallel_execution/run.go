package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run executes functions in n goroutines.
// If the number of errors reaches `m`, no more than `m+n` functions will be executed.
// If `m == 0` all functions will be executed in spite of errors.
func Run(tasks []Task, n, m int) error {
	ch := make(chan Task)
	wg := sync.WaitGroup{}
	errCounter := newErrCounter(m)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range ch {
				if err := task(); err != nil {
					errCounter.inc()
				}
			}
		}()
	}

	for _, task := range tasks {
		if errCounter.exceedsLimit() {
			break
		}
		ch <- task
	}

	close(ch)
	wg.Wait()

	if errCounter.exceedsLimit() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
