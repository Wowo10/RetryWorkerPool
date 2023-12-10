package retryworkerpool

type retrieveWorker[T any, R any] struct {
	baseWorker
}

func (w retrieveWorker[T, R]) runRetrieving(inputChannel <-chan T,
	inputArray []T, workFn func(input T) R,
	errorCheckFn func(result R) bool,
	errorCallbackFn func(inputArg T),
	resultCallbackFn func(r result[T, R])) {

	for {
		select {
		case inputArg := <-inputChannel:
			res := workFn(inputArg)

			isError := errorCheckFn(res)

			if isError {
				errorCallbackFn(inputArg)
			} else {
				resultCallbackFn(result[T, R]{
					input:  inputArg,
					output: res,
				})
			}

		case <-w.baseWorker.quitChannel:
			return
		}
	}
}

func (w *retrieveWorker[T, R]) prepare() *retrieveWorker[T, R] {
	w.quitChannel = make(chan bool, 1)

	return w
}
