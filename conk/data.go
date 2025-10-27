package conk

import (
	"sync"

	"github.com/roidaradal/rdb/ze"
)

type (
	DataWorkerFn[I any, O any]        = func(I) (O, error)
	RequestDataWorkerFn[I any, O any] = func(*ze.Request, I) (O, error)
)

// Perform data work (func(I) (O, error)) sequentially
func DataWorkersLinear[I any, O any](items []I, fn DataWorkerFn[I, O]) *DataResult[O] {
	result := newDataResult[O]()
	for i, item := range items {
		out, err := fn(item)
		if err == nil {
			result.Success += 1
			result.Output[i] = out
		} else {
			result.Errors[i] = err
		}
	}
	return result
}

// Perform data work (func(I) (O, error)) concurrently
func DataWorkers[I any, O any](items []I, fn DataWorkerFn[I, O], numWorkers uint) *DataResult[O] {
	// Worker function
	worker := func(inputCh <-chan inputData[I], outputCh chan<- outputData[O]) {
		for input := range inputCh {
			out, err := fn(input.item)
			outputCh <- outputData[O]{input.index, out, err}
		}
	}
	return dataWorkerPool(items, worker, numWorkers)
}

// Perform request data work (func(*Request, I) (O, error)) sequentially
func RequestDataWorkersLinear[I any, O any](rq *ze.Request, items []I, fn RequestDataWorkerFn[I, O]) *DataResult[O] {
	result := newDataResult[O]()
	for i, item := range items {
		out, err := fn(rq, item)
		if err == nil {
			result.Success += 1
			result.Output[i] = out
		} else {
			result.Errors[i] = err
		}
	}
	return result
}

// Perform request data work (func(I) (O, error)) concurrently
func RequestDataWorkers[I any, O any](rq *ze.Request, items []I, fn RequestDataWorkerFn[I, O], numWorkers uint) *DataResult[O] {
	// Worker function
	worker := func(inputCh <-chan inputData[I], outputCh chan<- outputData[O]) {
		srq := rq.SubRequest()
		for input := range inputCh {
			out, err := fn(srq, input.item)
			outputCh <- outputData[O]{input.index, out, err}
		}
		rq.MergeLogs(srq)
	}
	return dataWorkerPool(items, worker, numWorkers)
}

// Common: data worker pool
type dataWorkerFn[I any, O any] = func(<-chan inputData[I], chan<- outputData[O])

func dataWorkerPool[I any, O any](items []I, worker dataWorkerFn[I, O], numWorkers uint) *DataResult[O] {
	numWorkers = max(numWorkers, 1) // lower-bound: 1 worker
	inputCh := make(chan inputData[I])
	outputCh := make(chan outputData[O], numWorkers) // need buffered, otherwise deadlocks

	// Spawn workers
	var wg sync.WaitGroup
	for range numWorkers {
		wg.Go(func() {
			worker(inputCh, outputCh)
		})
	}

	// Feed input items to input channel
	go func() {
		for i, item := range items {
			inputCh <- inputData[I]{i, item}
		}
		close(inputCh)
	}()

	// Wait for all workers to finish and close output channel
	go func() {
		wg.Wait()
		close(outputCh)
	}()

	// Get results
	result := newDataResult[O]()
	for out := range outputCh {
		if out.err == nil {
			result.Success += 1
			result.Output[out.index] = out.item
		} else {
			result.Errors[out.index] = out.err
		}
	}
	return result
}
