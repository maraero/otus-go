package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	ch := make(chan Task)
	wg := sync.WaitGroup{}
	errCounter := newErrCouner(m)

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
			break;
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
