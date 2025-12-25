package workerpool

import "sync"

func WorkerPool[T any, R any](jobCh <-chan T, handler func(T) R, workersCount int) <-chan R {
	resCh := make(chan R)

	var wg sync.WaitGroup
	wg.Add(workersCount)

	for i := range workersCount {
		go func(workerId int) {
			for job := range jobCh {
				res := handler(job)
				resCh <- res
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	return resCh
}
