package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}

	var countErrors int32

	taskChannel := make(chan Task)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChannel {
				if err := task(); err != nil {
					atomic.AddInt32(&countErrors, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if int(atomic.LoadInt32(&countErrors)) >= m {
			close(taskChannel)
			wg.Wait()
			return ErrErrorsLimitExceeded
		}
		taskChannel <- task
	}

	close(taskChannel)
	wg.Wait()

	return nil
}
