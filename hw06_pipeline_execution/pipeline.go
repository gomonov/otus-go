package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(wrapChannel(in, done))
	}

	return in
}

func wrapChannel(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				for v := range in {
					_ = v
				}
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					for v := range in {
						_ = v
					}
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}
