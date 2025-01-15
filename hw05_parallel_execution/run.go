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
		lock.Lock()
		if needStopByLimitErrors {
			lock.Unlock()
			break
		}
		lock.Unlock()
		wg.Add(1)
		go func() {
			defer wg.Done()
			//defer lock.Unlock()

			for {

				lock.Lock()
				if needStopByLimitErrors || len(tasks) == 0 {
					lock.Unlock()
					break
				}
				task := tasks[0]
				tasks = tasks[1:]
				lock.Unlock()

				err := task()
				if err != nil {
					lock.Lock()
					currentCountErrors++
					if countErrors <= currentCountErrors {
						needStopByLimitErrors = true
					}
					lock.Unlock()
				}

			}

		}()
	}

	wg.Wait()

	if needStopByLimitErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func RunCh(tasks []Task, countWorkers, countErrors int) error {
	// Place your code here.

	var wg sync.WaitGroup
	var lock sync.Mutex
	needStopByLimitErrors := false
	currentCountErrors := 0

	ch := make(chan Task, countWorkers)

	go func() {
		//	for i := 0; i < len(tasks); i++ {
		for _, item := range tasks {
			ch <- item
		}
		close(ch)
	}()

	for i := 0; i < countWorkers; i++ {
		lock.Lock()
		if needStopByLimitErrors {
			lock.Unlock()
			break
		}
		lock.Unlock()
		wg.Add(1)
		go func() {
			defer wg.Done()
			//defer lock.Unlock()

			for t := range ch {

				lock.Lock()
				if needStopByLimitErrors || len(tasks) == 0 {
					lock.Unlock()
					break
				}
				lock.Unlock()

				if err := t(); err != nil {
					lock.Lock()
					currentCountErrors++
					if countErrors <= currentCountErrors {
						needStopByLimitErrors = true
					}
					lock.Unlock()
				}

			}

		}()
	}

	wg.Wait()

	if needStopByLimitErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func RunChNonBuffer(tasks []Task, countWorkers, countErrors int) error {
	var wg sync.WaitGroup
	var lock sync.Mutex
	currentCountErrors := 0

	ch := make(chan Task)

	for i := 0; i < countWorkers; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range ch {
				if err := t(); err != nil {
					lock.Lock()
					currentCountErrors++
					lock.Unlock()
				}

			}

		}()
	}
	for _, item := range tasks {
		lock.Lock()
		if countErrors <= currentCountErrors {
			lock.Unlock()
			break
		}
		lock.Unlock()
		ch <- item
	}
	close(ch)

	wg.Wait()

	if countErrors <= currentCountErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}
