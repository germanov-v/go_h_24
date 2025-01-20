package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
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

					_ = <-current
				}
				return
			case val, ok := <-current:
				if !ok {
					return
				}
				out <- val
				//default:
				// ?????
			}
		}

	}()

	return out
}
