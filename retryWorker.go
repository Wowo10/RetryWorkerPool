package retryworkerpool

type retryWorker[T any] struct {
	baseWorker
}

func (w retryWorker[T]) runRetry(inputChannel chan<- T, inputArg T) {
	select {
	case <-w.baseWorker.quitChannel:
		return
	case inputChannel <- inputArg:
	}
}

func (w *retryWorker[T]) prepare() *retryWorker[T] {
	w.quitChannel = make(chan bool, 1)

	return w
}
