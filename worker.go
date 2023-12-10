package retryworkerpool

type workerInterface interface {
	stop()
}

type baseWorker struct {
	quitChannel chan bool
}

func (w baseWorker) stop() {
	w.quitChannel <- true
	close(w.quitChannel)
}
