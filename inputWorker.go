package retryworkerpool

type inputWorker[T any] struct {
	baseWorker
}

func (w inputWorker[T]) runInput(inputChannel chan<- T, inputArray []T) {
	for i := range inputArray {
		select {
		case <-w.baseWorker.quitChannel:
			return
		case inputChannel <- inputArray[i]:
		}
	}
}

func (w *inputWorker[T]) prepare() *inputWorker[T] {
	w.quitChannel = make(chan bool, 1)

	return w
}
