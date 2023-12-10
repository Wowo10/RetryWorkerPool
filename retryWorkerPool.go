package retryworkerpool

import (
	"sync"
	"time"
)

type result[T any, R any] struct {
	input  T
	output R
}

type retryWorkerPool[T any, R any] struct {
	workers  []workerInterface
	results  []result[T, R]
	inputArr []T
	mutex    sync.Mutex

	inputChannel  chan T
	finishChannel chan bool

	work           func(inputArg T) R
	errorCondition func(result R) bool
	quit           bool
}

func CreateRetryWorkPool[T any, R any](inputArr []T, workFn func(inputArg T) R, errorConditionFn func(result R) bool) retryWorkerPool[T, R] {
	return retryWorkerPool[T, R]{
		inputArr:       inputArr,
		work:           workFn,
		errorCondition: errorConditionFn,
	}
}

func (d *retryWorkerPool[T, R]) Run(timeout time.Duration, threads int) []result[T, R] {
	d.inputChannel = make(chan T, threads)
	d.finishChannel = make(chan bool)

	inputWorker := new(inputWorker[T]).prepare()
	go inputWorker.runInput(d.inputChannel, d.inputArr)
	d.workers = append(d.workers, inputWorker)

	for i := 0; i < threads; i++ {
		retriever := new(retrieveWorker[T, R]).prepare()
		go retriever.runRetrieving(d.inputChannel, d.inputArr,
			d.work, d.errorCondition, d.retry, d.saveResult)

		d.workers = append(d.workers, retriever)
	}

	select {
	case <-d.finishChannel:
	case <-time.After(timeout):
	}

	d.quit = true

	for i := range d.workers {
		d.workers[i].stop()
	}

	defer close(d.inputChannel)
	defer close(d.finishChannel)

	return d.results
}

func (d *retryWorkerPool[T, R]) saveResult(r result[T, R]) {
	d.mutex.Lock()

	d.results = append(d.results, r)

	if len(d.results) == len(d.inputArr) && !d.quit {
		d.finishChannel <- true
	}
	d.mutex.Unlock()
}

func (d *retryWorkerPool[T, R]) retry(inputArg T) {
	if d.quit {
		return
	}

	retryWorker := new(retryWorker[T]).prepare()
	d.workers = append(d.workers, retryWorker)
	go retryWorker.runRetry(d.inputChannel, inputArg)
}
