package workerpool

import (
	"context"
	"sync"
)

func RunWithWorkers[T any, R any](ctx context.Context, jobCh <-chan T, handler func(ctx context.Context, job T) R, workersCount int) <-chan R {
	resCh := make(chan R)

	var wg sync.WaitGroup
	wg.Add(workersCount)

	for i := range workersCount {
		go func(workerId int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobCh:
					if !ok {
						return
					}

					res := handler(ctx, job)

					select {
					case <-ctx.Done():
						return
					case resCh <- res:
					}
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	return resCh
}
