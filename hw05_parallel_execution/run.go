package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, countWorkers, countErrors int) error {
	// Place your code here.

	var wg sync.WaitGroup
	var lock sync.Mutex
	needStopByLimitErrors := false
	currentCountErrors := 0
	for i := 0; i < countWorkers; i++ {
		if needStopByLimitErrors {
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer lock.Unlock()

			for {
				if needStopByLimitErrors {
					break
				}
				var task Task
				lock.Lock()
				if len(tasks) == 0 {
					break
				}
				task = tasks[0]
				tasks = tasks[1:]
				lock.Unlock()

				go func() {
					err := task()
					if err != nil {
						lock.Lock()
						currentCountErrors++
						if countErrors <= currentCountErrors {
							needStopByLimitErrors = true
						}
						lock.Unlock()
					}
				}()

			}

		}()
	}

	wg.Wait()

	if needStopByLimitErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}
