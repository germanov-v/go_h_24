package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// after race detected
func ExecutePipeline(in In, done In, stages ...Stage) Out {

	var wg sync.WaitGroup

	shadowStream := make(Bi)
	out := make(Bi)
	shadowProc := func(ch <-chan interface{}) {
		defer wg.Done()
		for i := range ch {
			select {
			case <-done:
				return
			case shadowStream <- i: // продолжаем работу
				out <- shadowStream
			}
		}
	}

	for _, stage := range stages {
		wg.Add(1)
		shadowProc(stage(in))
	}

	go func() {
		wg.Wait()
		close(shadowStream)
	}()

	return out
}

// alarm: race detecteted
func ExecutePipelineV3(in In, done In, stages ...Stage) Out {
	out := make(Bi) //make(chan<- interface{})

	go func() {
		defer close(out)
		current := in
		for _, stage := range stages {
			current = stage(current)
		}

		for {
			select {
			case <-done:
				for range current {
				}
				return
			case val, ok := <-current:
				if !ok {
					return
				}
				out <- val
			}
		}

	}()
	return out
}
