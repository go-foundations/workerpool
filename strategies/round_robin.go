package strategies

import (
	"context"
	"sync"

	"github.com/go-foundations/workerpool"
)

// RoundRobinStrategy distributes jobs evenly across workers in round-robin fashion
type RoundRobinStrategy[T any, R any] struct{}

// Name returns the strategy name
func (s *RoundRobinStrategy[T, R]) Name() string {
	return "Round Robin"
}

// Execute runs the round-robin distribution strategy
func (s *RoundRobinStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config,
	jobs []workerpool.Job[T], processor workerpool.Processor[T, R],
	results chan<- workerpool.Result[R]) error {

	var wg sync.WaitGroup

	// Create separate job channels for each worker
	jobChannels := make([]chan workerpool.Job[T], config.NumWorkers)
	for i := 0; i < config.NumWorkers; i++ {
		bufferSize := max(1, len(jobs)/config.NumWorkers+1)
		jobChannels[i] = make(chan workerpool.Job[T], bufferSize)
		wg.Add(1)
		go s.worker(i, jobChannels[i], &wg, processor, results, config)
	}

	// Distribute jobs round-robin
	for i, job := range jobs {
		workerIndex := i % config.NumWorkers
		select {
		case jobChannels[workerIndex] <- job:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Close all channels
	for _, ch := range jobChannels {
		close(ch)
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

// worker processes jobs from a dedicated channel
func (s *RoundRobinStrategy[T, R]) worker(id int, jobs <-chan workerpool.Job[T],
	wg *sync.WaitGroup, processor workerpool.Processor[T, R],
	results chan<- workerpool.Result[R], config *workerpool.Config) {

	defer wg.Done()

	for job := range jobs {
		s.processJob(id, job, processor, results, config)
	}
}

// processJob handles the actual job processing with retries and metrics
func (s *RoundRobinStrategy[T, R]) processJob(workerID int, job workerpool.Job[T],
	processor workerpool.Processor[T, R], results chan<- workerpool.Result[R],
	config *workerpool.Config) {
	processJob(workerID, job, processor, results, config)
}
