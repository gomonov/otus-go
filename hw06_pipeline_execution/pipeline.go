package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	for _, stage := range stages {
		in = stage(wrapPreStageChannel(in, done))
	}

	return in
}

func wrapPreStageChannel(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
			// вычитываем все оставшиеся сообщения из канала стейджа
			for v := range in {
				_ = v
			}
		}()
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}
