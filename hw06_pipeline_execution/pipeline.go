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
		for range in {
			for _, stage := range stages {
				current = stage(current)
			}
		}

		select {
		case <-done:
			return
		case out <- current:
		default:
			// ?????
		}
	}()

	return out
}
