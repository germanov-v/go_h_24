package hw06pipelineexecution

func ExecutePipelineV1(in In, done In, stages ...Stage) Out {
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

					<-current
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

func ExecutePipelineV2(in In, done In, stages ...Stage) Out {
	out := make(Bi) //make(chan<- interface{})
	go func() {
		defer close(out)
		current := in
		for _, stage := range stages {
			nextCh := make(Bi)

			go func(in In, out Out, stage Stage) {
				defer close(nextCh)
				for val := range in {
					select {
					case <-done:
						return
					case nextCh <- stage(val.(In)):

					}
				}
			}(current, nextCh, stage)
			current = nextCh // out

		}
		for {
			select {
			case <-done:
				for range current {

					<-current
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
