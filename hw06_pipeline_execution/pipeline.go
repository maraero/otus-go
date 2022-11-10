package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func interrupter(done In, in In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
			for range in { // drain the channel to let the previous stage finish
			}
		}()

		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				out <- v
			case <-done:
				return
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		out = stage(interrupter(done, out))
	}

	return interrupter(done, out)
}
