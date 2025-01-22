package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {

	var out = in // вместо создания make(Bi) - берем просто на чтение

	for _, stageItem := range stages {
		// перезаписываем переменную результат передаем дальше ниже в аргумент по функции
		out = func(in In, stage Stage) Out {
			outIntern := make(Bi) // !!!! обычный канал, его будем возвращать
			go func() {
				defer close(outIntern)
				stageIntern := stage(in)
				for {
					select {
					case <-done:
						return
					case val, ok := <-stageIntern:
						if !ok { // страхуемся на закрытие канала
							return
						}
						select { // и еще раз страхуемся, а вдруг  done в это время пришло
						case <-done:
							return
						case outIntern <- val:
						}

					}
				}
			}()
			return outIntern
		}(out, stageItem)
	}

	return out
}

func ExecutePipelineV5(in In, done In, stages ...Stage) Out {

	var wg sync.WaitGroup

	shadowStream := make(Bi)
	out := in
	shadowProc := func(ch <-chan interface{}, stage Stage) {
		defer wg.Done()
		for i := range ch {
			select {
			case <-done:
				return
			case shadowStream <- i:
				out = stage(shadowStream)
			}
		}
	}
	go shadowProc(done, nil)
	for _, stage := range stages {
		wg.Add(1)
		go shadowProc(out, stage)
	}

	go func() {
		wg.Wait()
		close(shadowStream)
	}()

	return out
}

func ExecutePipelineV4(in In, done In, stages ...Stage) Out {

	var wg sync.WaitGroup

	shadowStream := make(Bi)
	out := make(In)
	shadowProc := func(ch <-chan interface{}, stage Stage) {
		defer wg.Done()
		for i := range ch {
			select {
			case <-done:
				return
			case shadowStream <- i: // продолжаем работу
				out = stage(shadowStream)
			}
		}
	}

	for _, stage := range stages {
		wg.Add(1)
		shadowProc(in, stage)
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
