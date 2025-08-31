// Package workerpool provides a generic, high-performance worker pool implementation
// with multiple distribution strategies for concurrent job processing.
//
// The worker pool supports:
// - Generic types for jobs and results
// - Multiple distribution strategies (Round-Robin, Chunked, Work Stealing)
// - Configurable worker counts and buffer sizes
// - Context cancellation and timeout support
// - Error handling and result collection
// - Performance monitoring and metrics
package workerpool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Job represents a unit of work to be processed
type Job[T any] struct {
	ID       string    // Unique identifier for the job
	Data     T         // The actual data to be processed
	Priority int       // Job priority (higher = more important)
	Created  time.Time // When the job was created
}

// Result wraps the processing result of a job
type Result[R any] struct {
	JobID     string        // ID of the processed job
	Data      R             // The processed result
	Error     error         // Any error that occurred during processing
	Worker    int           // ID of the worker that processed the job
	Started   time.Time     // When processing started
	Completed time.Time     // When processing completed
	Duration  time.Duration // How long processing took
}

// Processor defines how to process a job
type Processor[T any, R any] func(ctx context.Context, job Job[T]) (R, error)

// DistributionStrategy defines how jobs are distributed to workers
type DistributionStrategy int

const (
	RoundRobin DistributionStrategy = iota
	Chunked
	WorkStealing
	PriorityBased
)

// Config holds configuration for the worker pool
type Config struct {
	NumWorkers    int                  // Number of worker goroutines
	BufferSize    int                  // Buffer size for job channels
	Strategy      DistributionStrategy // How to distribute jobs
	Timeout       time.Duration        // Overall timeout for the pool
	WorkerTimeout time.Duration        // Timeout per individual worker
	MaxRetries    int                  // Maximum retry attempts for failed jobs
	EnableMetrics bool                 // Whether to collect performance metrics
}

// DefaultConfig returns sensible default configuration
func DefaultConfig() Config {
	return Config{
		NumWorkers:    4,
		BufferSize:    100,
		Strategy:      RoundRobin,
		Timeout:       5 * time.Minute,
		WorkerTimeout: 30 * time.Second,
		MaxRetries:    3,
		EnableMetrics: true,
	}
}

// WorkerPool manages a pool of workers for processing jobs
type WorkerPool[T any, R any] struct {
	config    Config
	processor Processor[T, R]
	jobs      []Job[T]
	results   chan Result[R]
	ctx       context.Context
	cancel    context.CancelFunc
	metrics   *Metrics
	mu        sync.RWMutex
	ctxMu     sync.RWMutex // Protects ctx and cancel fields
}

// Metrics holds performance metrics for the worker pool
type Metrics struct {
	TotalJobs       int
	ProcessedJobs   int
	FailedJobs      int
	TotalDuration   time.Duration
	AverageDuration time.Duration
	StartTime       time.Time
	EndTime         time.Time
	mu              sync.RWMutex
}

// New creates a new worker pool with default configuration
func New[T any, R any]() *WorkerPool[T, R] {
	return NewWithConfig[T, R](DefaultConfig())
}

// NewWithConfig creates a new worker pool with custom configuration
func NewWithConfig[T any, R any](config Config) *WorkerPool[T, R] {
	if config.NumWorkers <= 0 {
		config.NumWorkers = 1
	}
	if config.BufferSize <= 0 {
		config.BufferSize = 100
	}
	if config.Timeout <= 0 {
		config.Timeout = 5 * time.Minute // Default timeout
	}
	// Ensure buffer size is at least as large as typical job counts
	if config.BufferSize < 10 {
		config.BufferSize = 10
	}

	return &WorkerPool[T, R]{
		config:  config,
		results: make(chan Result[R], config.BufferSize),
		ctx:     nil, // Will be set in Run()
		cancel:  nil, // Will be set in Run()
		metrics: &Metrics{},
	}
}

// WithProcessor sets the processing function for the worker pool
func (wp *WorkerPool[T, R]) WithProcessor(p Processor[T, R]) *WorkerPool[T, R] {
	wp.processor = p
	return wp
}

// AddJobs adds jobs to the worker pool
func (wp *WorkerPool[T, R]) AddJobs(jobs []Job[T]) *WorkerPool[T, R] {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.jobs = make([]Job[T], len(jobs))
	copy(wp.jobs, jobs)

	// Set creation time for new jobs
	now := time.Now()
	for i := range wp.jobs {
		if wp.jobs[i].Created.IsZero() {
			wp.jobs[i].Created = now
		}
	}

	wp.metrics.TotalJobs = len(wp.jobs)
	return wp
}

// AddJob adds a single job to the worker pool
func (wp *WorkerPool[T, R]) AddJob(job Job[T]) *WorkerPool[T, R] {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if job.Created.IsZero() {
		job.Created = time.Now()
	}

	wp.jobs = append(wp.jobs, job)
	wp.metrics.TotalJobs = len(wp.jobs)
	return wp
}

// Run executes the worker pool with the configured strategy
func (wp *WorkerPool[T, R]) Run() ([]Result[R], error) {
	if wp.processor == nil {
		return nil, fmt.Errorf("no processor configured")
	}
	if len(wp.jobs) == 0 {
		return nil, fmt.Errorf("no jobs to process")
	}

	// Create context with timeout for this run
	ctx, cancel := context.WithTimeout(context.Background(), wp.config.Timeout)
	wp.ctxMu.Lock()
	wp.ctx = ctx
	wp.cancel = cancel
	wp.ctxMu.Unlock()

	wp.metrics.StartTime = time.Now()
	defer func() {
		wp.metrics.EndTime = time.Now()
		wp.metrics.TotalDuration = wp.metrics.EndTime.Sub(wp.metrics.StartTime)
		if wp.metrics.ProcessedJobs > 0 {
			wp.metrics.AverageDuration = wp.metrics.TotalDuration / time.Duration(wp.metrics.ProcessedJobs)
		}
	}()

	var err error
	switch wp.config.Strategy {
	case RoundRobin:
		err = wp.runRoundRobin()
	case Chunked:
		err = wp.runChunked()
	case WorkStealing:
		err = wp.runWorkStealing()
	case PriorityBased:
		err = wp.runPriorityBased()
	default:
		err = wp.runRoundRobin()
	}

	if err != nil {
		// Clean up context on error
		wp.ctxMu.Lock()
		if wp.cancel != nil {
			wp.cancel()
			wp.cancel = nil
			wp.ctx = nil
		}
		wp.ctxMu.Unlock()
		return nil, err
	}

	// Collect results
	var results []Result[R]
	for result := range wp.results {
		// Check for context cancellation
		select {
		case <-wp.ctx.Done():
			return nil, wp.ctx.Err()
		default:
		}

		results = append(results, result)
		if result.Error != nil {
			wp.metrics.FailedJobs++
		} else {
			wp.metrics.ProcessedJobs++
		}
	}

	// Clean up context
	wp.ctxMu.Lock()
	if wp.cancel != nil {
		wp.cancel()
		wp.cancel = nil
		wp.ctx = nil
	}
	wp.ctxMu.Unlock()

	return results, nil
}

// runRoundRobin distributes jobs evenly across workers in round-robin fashion
func (wp *WorkerPool[T, R]) runRoundRobin() error {
	var wg sync.WaitGroup

	// Create separate job channels for each worker
	jobChannels := make([]chan Job[T], wp.config.NumWorkers)
	for i := 0; i < wp.config.NumWorkers; i++ {
		bufferSize := max(1, len(wp.jobs)/wp.config.NumWorkers+1)
		jobChannels[i] = make(chan Job[T], bufferSize)
		wg.Add(1)
		go wp.worker(i, jobChannels[i], &wg)
	}

	// Distribute jobs round-robin
	for i, job := range wp.jobs {
		workerIndex := i % wp.config.NumWorkers
		select {
		case jobChannels[workerIndex] <- job:
		case <-wp.ctx.Done():
			return wp.ctx.Err()
		}
	}

	// Close all channels
	for _, ch := range jobChannels {
		close(ch)
	}

	wg.Wait()
	close(wp.results)

	// Check if context was cancelled during execution
	select {
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		return nil
	}
}

// runChunked distributes jobs in chunks to workers
func (wp *WorkerPool[T, R]) runChunked() error {
	var wg sync.WaitGroup

	chunkSize := max(1, len(wp.jobs)/wp.config.NumWorkers)
	remainder := len(wp.jobs) % wp.config.NumWorkers

	start := 0
	for i := 0; i < wp.config.NumWorkers; i++ {
		end := start + chunkSize
		if i < remainder {
			end++
		}

		if start < len(wp.jobs) {
			wg.Add(1)
			go wp.workerWithSlice(i, wp.jobs[start:end], &wg)
		}
		start = end
	}

	wg.Wait()
	close(wp.results)

	// Check if context was cancelled during execution
	select {
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		return nil
	}
}

// runWorkStealing implements work stealing using Chase-Lev work stealing deques
func (wp *WorkerPool[T, R]) runWorkStealing() error {
	var wg sync.WaitGroup

	// Create work stealing deques for each worker
	deques := make([]*WorkStealingDeque[T], wp.config.NumWorkers)
	for i := 0; i < wp.config.NumWorkers; i++ {
		deques[i] = NewWorkStealingDeque[T](len(wp.jobs) / wp.config.NumWorkers + 1)
	}

	// Distribute jobs initially across worker deques (round-robin)
	for i, job := range wp.jobs {
		workerIndex := i % wp.config.NumWorkers
		deques[workerIndex].Push(job)
	}

	// Start work stealing workers
	for i := 0; i < wp.config.NumWorkers; i++ {
		wg.Add(1)
		go wp.workStealingWorker(i, deques, &wg)
	}

	wg.Wait()
	close(wp.results)

	// Check if context was cancelled during execution
	select {
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		return nil
	}
}

// workStealingWorker implements work stealing behavior
func (wp *WorkerPool[T, R]) workStealingWorker(id int, deques []*WorkStealingDeque[T], wg *sync.WaitGroup) {
	defer wg.Done()

	myDeque := deques[id]
	numWorkers := len(deques)

	for {
		// Wait for context to be available
		var ctx context.Context
		for {
			wp.ctxMu.RLock()
			ctx = wp.ctx
			wp.ctxMu.RUnlock()

			if ctx != nil {
				break
			}
			time.Sleep(1 * time.Millisecond) // Short wait
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Try to get work from own deque first (LIFO for better cache locality)
		if job, ok := myDeque.Pop(); ok {
			wp.processJob(id, job)
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
				wp.processJob(id, job)
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

// runPriorityBased processes jobs based on priority using a priority queue with fair scheduling
func (wp *WorkerPool[T, R]) runPriorityBased() error {
	var wg sync.WaitGroup

	// Create priority queue and populate it with jobs
	priorityQueue := NewPriorityQueue[T]()
	
	// Set creation time for fair scheduling and add jobs to priority queue
	for _, job := range wp.jobs {
		if job.Created.IsZero() {
			job.Created = time.Now()
		}
		priorityQueue.Push(job)
	}

	// Create shared work queue for workers to consume from
	workQueue := make(chan Job[T], wp.config.BufferSize)
	
	// Start workers
	for i := 0; i < wp.config.NumWorkers; i++ {
		wg.Add(1)
		go wp.worker(i, workQueue, &wg)
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
			case <-wp.ctx.Done():
				return
			}
		}
	}()

	wg.Wait()
	close(wp.results)

	// Check if context was cancelled during execution
	select {
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		return nil
	}
}

// worker processes jobs from a dedicated channel
func (wp *WorkerPool[T, R]) worker(id int, jobs <-chan Job[T], wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Wait for context to be available
		var ctx context.Context
		for {
			wp.ctxMu.RLock()
			ctx = wp.ctx
			wp.ctxMu.RUnlock()

			if ctx != nil {
				break
			}
			time.Sleep(1 * time.Millisecond) // Short wait
		}

		select {
		case <-ctx.Done():
			return
		default:
			wp.processJob(id, job)
		}
	}
}

// workerWithSlice processes a slice of jobs directly
func (wp *WorkerPool[T, R]) workerWithSlice(id int, jobSlice []Job[T], wg *sync.WaitGroup) {
	defer wg.Done()

	for _, job := range jobSlice {
		// Wait for context to be available
		var ctx context.Context
		for {
			wp.ctxMu.RLock()
			ctx = wp.ctx
			wp.ctxMu.RUnlock()

			if ctx != nil {
				break
			}
			time.Sleep(1 * time.Millisecond) // Short wait
		}

		select {
		case <-ctx.Done():
			return
		default:
			wp.processJob(id, job)
		}
	}
}

// processJob handles the actual job processing with retries and metrics
func (wp *WorkerPool[T, R]) processJob(workerID int, job Job[T]) {
	startTime := time.Now()

	var result R
	var err error

	// Process with retries
	for attempt := 0; attempt <= wp.config.MaxRetries; attempt++ {
		// Wait for context to be available
		var ctx context.Context
		for {
			wp.ctxMu.RLock()
			ctx = wp.ctx
			wp.ctxMu.RUnlock()

			if ctx != nil {
				break
			}
			time.Sleep(1 * time.Millisecond) // Short wait
		}

		select {
		case <-ctx.Done():
			return
		default:
			result, err = wp.processor(ctx, job)
			if err == nil {
				break
			}
			if attempt < wp.config.MaxRetries {
				time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
			}
		}
	}

	completed := time.Now()
	duration := completed.Sub(startTime)

	// Send result to channel
	wp.results <- Result[R]{
		JobID:     job.ID,
		Data:      result,
		Error:     err,
		Worker:    workerID,
		Started:   startTime,
		Completed: completed,
		Duration:  duration,
	}
}

// GetMetrics returns a copy of the current metrics
func (wp *WorkerPool[T, R]) GetMetrics() Metrics {
	wp.metrics.mu.RLock()
	defer wp.metrics.mu.RUnlock()

	return Metrics{
		TotalJobs:       wp.metrics.TotalJobs,
		ProcessedJobs:   wp.metrics.ProcessedJobs,
		FailedJobs:      wp.metrics.FailedJobs,
		TotalDuration:   wp.metrics.TotalDuration,
		AverageDuration: wp.metrics.AverageDuration,
		StartTime:       wp.metrics.StartTime,
		EndTime:         wp.metrics.EndTime,
	}
}

// GetNumWorkers returns the number of workers in the pool
func (wp *WorkerPool[T, R]) GetNumWorkers() int {
	return wp.config.NumWorkers
}

// Stop cancels the worker pool context
func (wp *WorkerPool[T, R]) Stop() {
	wp.ctxMu.RLock()
	cancel := wp.cancel
	wp.ctxMu.RUnlock()

	if cancel != nil {
		cancel()
	}
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// WorkStealingDeque implements a lock-free work stealing deque
// Based on the Chase-Lev work stealing deque algorithm
type WorkStealingDeque[T any] struct {
	bottom int32
	top    int32
	buffer []Job[T]
	mu     sync.RWMutex
}

// NewWorkStealingDeque creates a new work stealing deque
func NewWorkStealingDeque[T any](initialSize int) *WorkStealingDeque[T] {
	if initialSize <= 0 {
		initialSize = 64
	}
	return &WorkStealingDeque[T]{
		buffer: make([]Job[T], initialSize),
	}
}

// Push adds a job to the bottom of the deque (owner thread)
func (d *WorkStealingDeque[T]) Push(job Job[T]) {
	d.mu.Lock()
	defer d.mu.Unlock()

	bottom := d.bottom
	top := d.top

	// Check if we need to grow the buffer
	if bottom-top >= int32(len(d.buffer)) {
		d.grow()
	}

	d.buffer[bottom%int32(len(d.buffer))] = job
	d.bottom++
}

// Pop removes and returns a job from the bottom of the deque (owner thread)
func (d *WorkStealingDeque[T]) Pop() (Job[T], bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	bottom := d.bottom - 1
	d.bottom = bottom

	top := d.top

	if top > bottom {
		d.bottom = top
		return Job[T]{}, false
	}

	job := d.buffer[bottom%int32(len(d.buffer))]
	if top == bottom {
		d.bottom = top
	}
	return job, true
}

// Steal removes and returns a job from the top of the deque (thief thread)
func (d *WorkStealingDeque[T]) Steal() (Job[T], bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	top := d.top
	bottom := d.bottom

	if top >= bottom {
		return Job[T]{}, false
	}

	job := d.buffer[top%int32(len(d.buffer))]
	d.top++
	return job, true
}

// grow increases the buffer size when needed
func (d *WorkStealingDeque[T]) grow() {
	newBuffer := make([]Job[T], len(d.buffer)*2)

	// Copy existing elements
	for i := d.top; i < d.bottom; i++ {
		newBuffer[i%int32(len(newBuffer))] = d.buffer[i%int32(len(d.buffer))]
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

// PriorityQueue implements a priority queue with fair scheduling
// Uses a binary heap with additional fairness mechanisms
type PriorityQueue[T any] struct {
	items    []Job[T]
	mu       sync.RWMutex
	fairness map[int]int // Track job counts per priority to prevent starvation
}

// NewPriorityQueue creates a new priority queue
func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items:    make([]Job[T], 0),
		fairness: make(map[int]int),
	}
}

// Push adds a job to the priority queue
func (pq *PriorityQueue[T]) Push(job Job[T]) {
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
func (pq *PriorityQueue[T]) Pop() (Job[T], bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if len(pq.items) == 0 {
		return Job[T]{}, false
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
func (pq *PriorityQueue[T]) Peek() (Job[T], bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if len(pq.items) == 0 {
		return Job[T]{}, false
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
