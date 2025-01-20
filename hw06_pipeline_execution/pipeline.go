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
		select {
		case <-done:
			return
		default:
			for _ = range in {
				for _, stage := range stages {
					out <- stage(in)
				}
			}
		}
	}()

	return out
}
