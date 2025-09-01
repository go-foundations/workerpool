package strategies

import (
	"context"
	"sync"

	"github.com/go-foundations/workerpool"
)

// ChunkedStrategy distributes jobs in chunks to workers
type ChunkedStrategy[T any, R any] struct{}

// Name returns the strategy name
func (s *ChunkedStrategy[T, R]) Name() string {
	return "Chunked"
}

// Execute runs the chunked distribution strategy
func (s *ChunkedStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config,
	jobs []workerpool.Job[T], processor workerpool.Processor[T, R],
	results chan<- workerpool.Result[R]) error {

	var wg sync.WaitGroup

	chunkSize := max(1, len(jobs)/config.NumWorkers)
	remainder := len(jobs) % config.NumWorkers

	start := 0
	for i := 0; i < config.NumWorkers; i++ {
		end := start + chunkSize
		if i < remainder {
			end++
		}

		if start < len(jobs) {
			wg.Add(1)
			go s.workerWithSlice(i, jobs[start:end], &wg, processor, results, config)
		}
		start = end
	}

	wg.Wait()
	close(results)

	// Check if context was cancelled during execution
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// workerWithSlice processes a slice of jobs directly
func (s *ChunkedStrategy[T, R]) workerWithSlice(id int, jobSlice []workerpool.Job[T],
	wg *sync.WaitGroup, processor workerpool.Processor[T, R],
	results chan<- workerpool.Result[R], config *workerpool.Config) {

	defer wg.Done()

	for _, job := range jobSlice {
		s.processJob(id, job, processor, results, config)
	}
}

// processJob handles the actual job processing with retries and metrics
func (s *ChunkedStrategy[T, R]) processJob(workerID int, job workerpool.Job[T],
	processor workerpool.Processor[T, R], results chan<- workerpool.Result[R],
	config *workerpool.Config) {
	processJob(workerID, job, processor, results, config)
}
