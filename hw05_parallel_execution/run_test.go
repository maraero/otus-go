package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("if errors < max errors amount, then finished all the tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)

				if i == 10 || i == 20 || i == 30 { // 3 errors, max — 4
					return err
				}

				return nil
			})
		}

		workersCount := 10
		maxErrorsCount := 4
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(tasksCount), "extra tasks were started")
	})

	t.Run("if errors are not limited (m = 0), then finished all the tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all the tasks were started")
	})

	t.Run("if errors are not limited (m < 0), then finished all the tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := -1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all the tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		var sumTime time.Duration
		workersCount := 5
		maxErrorsCount := 1

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks without errors (without sleep)", func(t *testing.T) {
		const tasksCount = 50
		const workersCount = 5
		tasks := make([]Task, 0, tasksCount)
		waitCh := make(chan struct{})
		var runTasksCount int32
		runErr := make(chan error)

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				<-waitCh // to hold the function until the channel is closed
				return nil
			})
		}

		go func() {
			runErr <- Run(tasks, workersCount, tasksCount)
		}()

		require.Eventually(t, func() bool {
			return atomic.LoadInt32(&runTasksCount) == workersCount
		}, time.Second, time.Millisecond)

		close(waitCh)
		require.NoError(t, <-runErr)
	})
}
