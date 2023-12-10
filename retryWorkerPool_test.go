package retryworkerpool_test

import (
	retryworkerpool "retryWorkerPool"
	"testing"
	"time"
)

func Test_Run(t *testing.T) {

	pool := retryworkerpool.CreateRetryWorkPool[int, int](
		[]int{1, 2, 3},
		func(inputArg int) int {
			return inputArg * 10
		},
		func(result int) bool {
			return false
		},
	)

	result := pool.Run(time.Second*20, 5)

	if len(result) != 3 {
		t.Fail()
	}
}
