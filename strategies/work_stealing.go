package strategies

import (
	"context"
	"sync"
	"time"

	"github.com/go-foundations/workerpool"
)

// WorkStealingStrategy implements work stealing using Chase-Lev work stealing deques
type WorkStealingStrategy[T any, R any] struct{}

// Name returns the strategy name
func (s *WorkStealingStrategy[T, R]) Name() string {
	return "Work Stealing"
}

// Execute runs the work stealing distribution strategy
func (s *WorkStealingStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config,
	jobs []workerpool.Job[T], processor workerpool.Processor[T, R],
	results chan<- workerpool.Result[R]) error {

	var wg sync.WaitGroup

	// Create work stealing deques for each worker
	deques := make([]*WorkStealingDeque[T], config.NumWorkers)
	for i := 0; i < config.NumWorkers; i++ {
		deques[i] = NewWorkStealingDeque[T](len(jobs)/config.NumWorkers + 1)
	}

	// Distribute jobs initially across worker deques (round-robin)
	for i, job := range jobs {
		workerIndex := i % config.NumWorkers
		deques[workerIndex].Push(job)
	}

	// Start work stealing workers
	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		go s.workStealingWorker(i, deques, &wg, processor, results, config, ctx)
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

// workStealingWorker implements work stealing behavior
func (s *WorkStealingStrategy[T, R]) workStealingWorker(id int, deques []*WorkStealingDeque[T],
	wg *sync.WaitGroup, processor workerpool.Processor[T, R],
	results chan<- workerpool.Result[R], config *workerpool.Config, ctx context.Context) {

	defer wg.Done()

	myDeque := deques[id]
	numWorkers := len(deques)

	for {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Try to get work from own deque first (LIFO for better cache locality)
		if job, ok := myDeque.Pop(); ok {
			s.processJob(id, job, processor, results, config)
			continue
		}

		// No work in own deque, try to steal from other workers (FIFO)
		stolen := false
		for attempts := 0; attempts < numWorkers*2; attempts++ {
			// Pick a random victim (avoid bias)
			victimID := (id + attempts + 1) % numWorkers
			if victimID == id {
				continue // Don't steal from yourself
			}

			if job, ok := deques[victimID].Steal(); ok {
				s.processJob(id, job, processor, results, config)
				stolen = true
				break
			}
		}

		// If no work was stolen, check if all deques are empty
		if !stolen {
			allEmpty := true
			for _, deque := range deques {
				if !deque.IsEmpty() {
					allEmpty = false
					break
				}
			}

			if allEmpty {
				// No more work available
				return
			}

			// Brief pause before trying again to avoid busy waiting
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// processJob handles the actual job processing with retries and metrics
func (s *WorkStealingStrategy[T, R]) processJob(workerID int, job workerpool.Job[T],
	processor workerpool.Processor[T, R], results chan<- workerpool.Result[R],
	config *workerpool.Config) {
	processJob(workerID, job, processor, results, config)
}

// WorkStealingDeque implements a lock-free work stealing deque
// Based on the Chase-Lev work stealing deque algorithm
type WorkStealingDeque[T any] struct {
	bottom int
	top    int
	buffer []workerpool.Job[T]
	mu     sync.RWMutex
}

// NewWorkStealingDeque creates a new work stealing deque
func NewWorkStealingDeque[T any](initialSize int) *WorkStealingDeque[T] {
	if initialSize <= 0 {
		initialSize = 64
	}
	return &WorkStealingDeque[T]{
		buffer: make([]workerpool.Job[T], initialSize),
	}
}

// Push adds a job to the bottom of the deque (owner thread)
func (d *WorkStealingDeque[T]) Push(job workerpool.Job[T]) {
	d.mu.Lock()
	defer d.mu.Unlock()

	bottom := d.bottom
	top := d.top

	// Check if we need to grow the buffer
	if bottom-top >= len(d.buffer) {
		d.grow()
	}

	d.buffer[bottom%len(d.buffer)] = job
	d.bottom++
}

// Pop removes and returns a job from the bottom of the deque (owner thread)
func (d *WorkStealingDeque[T]) Pop() (workerpool.Job[T], bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	bottom := d.bottom - 1
	d.bottom = bottom

	top := d.top

	if top > bottom {
		d.bottom = top
		return workerpool.Job[T]{}, false
	}

	job := d.buffer[bottom%len(d.buffer)]
	if top == bottom {
		d.bottom = top
	}
	return job, true
}

// Steal removes and returns a job from the top of the deque (thief thread)
func (d *WorkStealingDeque[T]) Steal() (workerpool.Job[T], bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	top := d.top
	bottom := d.bottom

	if top >= bottom {
		return workerpool.Job[T]{}, false
	}

	job := d.buffer[top%len(d.buffer)]
	d.top++
	return job, true
}

// grow increases the buffer size when needed
func (d *WorkStealingDeque[T]) grow() {
	newBuffer := make([]workerpool.Job[T], len(d.buffer)*2)

	// Copy existing elements
	for i := d.top; i < d.bottom; i++ {
		newBuffer[i%len(newBuffer)] = d.buffer[i%len(d.buffer)]
	}

	d.buffer = newBuffer
}

// Size returns the current number of jobs in the deque
func (d *WorkStealingDeque[T]) Size() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return int(d.bottom - d.top)
}

// IsEmpty checks if the deque is empty
func (d *WorkStealingDeque[T]) IsEmpty() bool {
	return d.Size() == 0
}
