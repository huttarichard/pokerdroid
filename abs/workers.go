package abs

import (
	"runtime"
	"sync/atomic"

	"golang.org/x/exp/constraints"
	"golang.org/x/sync/errgroup"
)

type IndexFn[T constraints.Integer] func(id T, done uint64) error

func IndexWorkers[T constraints.Integer](size int, callback IndexFn[T]) error {
	ww := runtime.NumCPU()
	epw := size / ww

	var counter uint64
	var wg errgroup.Group

	// Worker goroutines
	for w := 0; w < ww; w++ {
		id := w
		wg.Go(func() error {
			start := id * epw
			end := start + epw

			if id == ww-1 {
				end = size
			}

			for i := start; i < end; i++ {
				err := callback(T(i), counter)
				if err != nil {
					return err
				}
				atomic.AddUint64(&counter, 1)
			}
			return nil
		})
	}

	// Wait for completion
	err := wg.Wait()
	if err != nil {
		return err
	}

	return nil
}
