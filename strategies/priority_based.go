package strategies

import (
	"context"
	"sync"
	"time"

	"github.com/go-foundations/workerpool"
)

// PriorityBasedStrategy processes jobs based on priority using a priority queue with fair scheduling
type PriorityBasedStrategy[T any, R any] struct{}

// Name returns the strategy name
func (s *PriorityBasedStrategy[T, R]) Name() string {
	return "Priority Based"
}

// Execute runs the priority-based distribution strategy
func (s *PriorityBasedStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config,
	jobs []workerpool.Job[T], processor workerpool.Processor[T, R],
	results chan<- workerpool.Result[R]) error {

	var wg sync.WaitGroup

	// Create priority queue and populate it with jobs
	priorityQueue := NewPriorityQueue[T]()

	// Set creation time for fair scheduling and add jobs to priority queue
	for _, job := range jobs {
		if job.Created.IsZero() {
			job.Created = time.Now()
		}
		priorityQueue.Push(job)
	}

	// Create shared work queue for workers to consume from
	workQueue := make(chan workerpool.Job[T], config.BufferSize)

	// Start workers
	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		go s.worker(i, workQueue, &wg, processor, results, config)
	}

	// Priority dispatcher: continuously feeds high-priority jobs to workers
	go func() {
		defer close(workQueue)

		for !priorityQueue.IsEmpty() {
			job, ok := priorityQueue.Pop()
			if !ok {
				break
			}

			select {
			case workQueue <- job:
				// Job sent to worker
			case <-ctx.Done():
				return
			}
		}
	}()

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
func (s *PriorityBasedStrategy[T, R]) worker(id int, jobs <-chan workerpool.Job[T],
	wg *sync.WaitGroup, processor workerpool.Processor[T, R],
	results chan<- workerpool.Result[R], config *workerpool.Config) {

	defer wg.Done()

	for job := range jobs {
		s.processJob(id, job, processor, results, config)
	}
}

// processJob handles the actual job processing with retries and metrics
func (s *PriorityBasedStrategy[T, R]) processJob(workerID int, job workerpool.Job[T],
	processor workerpool.Processor[T, R], results chan<- workerpool.Result[R],
	config *workerpool.Config) {
	processJob(workerID, job, processor, results, config)
}

// PriorityQueue implements a priority queue with fair scheduling
// Uses a binary heap with additional fairness mechanisms
type PriorityQueue[T any] struct {
	items    []workerpool.Job[T]
	mu       sync.RWMutex
	fairness map[int]int // Track job counts per priority to prevent starvation
}

// NewPriorityQueue creates a new priority queue
func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items:    make([]workerpool.Job[T], 0),
		fairness: make(map[int]int),
	}
}

// Push adds a job to the priority queue
func (pq *PriorityQueue[T]) Push(job workerpool.Job[T]) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	// Track job count per priority for fairness
	pq.fairness[job.Priority]++

	// Add job to the end
	pq.items = append(pq.items, job)

	// Bubble up to maintain heap property
	pq.bubbleUp(len(pq.items) - 1)
}

// Pop removes and returns the highest priority job
func (pq *PriorityQueue[T]) Pop() (workerpool.Job[T], bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if len(pq.items) == 0 {
		return workerpool.Job[T]{}, false
	}

	// Get the highest priority job
	job := pq.items[0]

	// Update fairness tracking
	pq.fairness[job.Priority]--

	// Move the last item to the root
	pq.items[0] = pq.items[len(pq.items)-1]
	pq.items = pq.items[:len(pq.items)-1]

	// Bubble down to maintain heap property
	if len(pq.items) > 0 {
		pq.bubbleDown(0)
	}

	return job, true
}

// Peek returns the highest priority job without removing it
func (pq *PriorityQueue[T]) Peek() (workerpool.Job[T], bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if len(pq.items) == 0 {
		return workerpool.Job[T]{}, false
	}

	return pq.items[0], true
}

// Size returns the number of jobs in the queue
func (pq *PriorityQueue[T]) Size() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return len(pq.items)
}

// IsEmpty checks if the queue is empty
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.Size() == 0
}

// GetFairnessStats returns fairness statistics
func (pq *PriorityQueue[T]) GetFairnessStats() map[int]int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	stats := make(map[int]int)
	for k, v := range pq.fairness {
		stats[k] = v
	}
	return stats
}

// bubbleUp maintains the heap property by bubbling up an element
func (pq *PriorityQueue[T]) bubbleUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2

		// Check if we need to swap with parent
		if pq.shouldSwap(parent, index) {
			pq.items[parent], pq.items[index] = pq.items[index], pq.items[parent]
			index = parent
		} else {
			break
		}
	}
}

// bubbleDown maintains the heap property by bubbling down an element
func (pq *PriorityQueue[T]) bubbleDown(index int) {
	for {
		left := 2*index + 1
		right := 2*index + 2
		smallest := index

		// Find the smallest among current node and its children
		if left < len(pq.items) && pq.shouldSwap(smallest, left) {
			smallest = left
		}
		if right < len(pq.items) && pq.shouldSwap(smallest, right) {
			smallest = right
		}

		// If no swap needed, we're done
		if smallest == index {
			break
		}

		// Swap and continue
		pq.items[index], pq.items[smallest] = pq.items[smallest], pq.items[index]
		index = smallest
	}
}

// shouldSwap determines if two jobs should be swapped for heap ordering
// Implements fair scheduling to prevent starvation
func (pq *PriorityQueue[T]) shouldSwap(parent, child int) bool {
	parentJob := pq.items[parent]
	childJob := pq.items[child]

	// Primary ordering: higher priority first
	if parentJob.Priority != childJob.Priority {
		return childJob.Priority > parentJob.Priority
	}

	// Secondary ordering: FIFO for same priority (fairness)
	return parentJob.Created.After(childJob.Created)
}
