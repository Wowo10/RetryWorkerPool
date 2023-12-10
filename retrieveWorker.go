package retryworkerpool

type retrieveWorker[T any, R any] struct {
	baseWorker
}

func (w retrieveWorker[T, R]) runRetrieving(inputChannel <-chan T,
	inputArray []T, work func(input T) R,
	errorCondition func(result R) bool,
	errorCallback func(inputArg T), resultCallback func(r result[T, R])) {

	for {
		select {
		case inputArg := <-inputChannel:
			res := work(inputArg)

			isError := errorCondition(res)

			if isError {
				errorCallback(inputArg)
			} else {
				resultCallback(result[T, R]{
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
