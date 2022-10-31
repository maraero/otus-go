package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func interrupter(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
			for range in {}
		}()

		for {
			select {
			case v, ok := <- in:
				if !ok {
					return
				}
				out <- v
			case <- done:
				return
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out = in

	for _, stage := range stages {
		out = stage(interrupter(out, done))
	}

	return out
}
