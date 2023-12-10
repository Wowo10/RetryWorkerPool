# Retry Worker Thread Pool
Abstraction over running multiple workers with retry option and a timeout.
Created as example implementation of Worker Thread Pool pattern with Retry and Timeout.
Using Go and go routines and channels.

Example usage:
```
retryWorkerPool := retryworkerpool.CreateRetryWorkPool[int, int](
	[]int{1,2,3,4,5},
	func(inputArg int) int {	
		time.Sleep(time.Second)
		rand := rand.Intn(100)

		fmt.Println("finished for input:", inputArg, "result:", rand)

		return rand
	},
	func(result int) bool {
		if result%2 == 0 {
			fmt.Println("error for result:", result)

			return  true
		}
	
		return  false
	},
)

result := retryWorkerPool.Run(time.Second*5, 5)

fmt.Println("ResultArr:", result)
fmt.Println("ResultArrLength:", len(result))
```

Mostly created for education purposes. But might be somebody would benefit from the package.