package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-foundations/workerpool"
)

// This example demonstrates how to implement a custom distribution strategy
// by creating a new strategy that implements the Strategy interface

func main() {
	fmt.Println("=== Custom Strategy Implementation Example ===")
	fmt.Println()

	// Example 1: Using the built-in adaptive strategy
	demoAdaptiveStrategy()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Example 2: Implementing a custom strategy
	demoCustomStrategy()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Example 3: Strategy performance comparison
	demoStrategyComparison()
}

// demoAdaptiveStrategy shows how to use the built-in adaptive strategy
func demoAdaptiveStrategy() {
	fmt.Println("1. Using Built-in Adaptive Strategy")
	fmt.Println("   The adaptive strategy automatically selects the best distribution method")
	fmt.Println("   based on workload characteristics.")

	config := workerpool.Config{
		NumWorkers: 4,
		Strategy:   workerpool.Adaptive, // Use the new adaptive strategy
		Timeout:    30 * time.Second,
		MaxRetries: 2,
	}

	pool := workerpool.NewWithConfig[string, string](config).
		WithProcessor(stringProcessor)

	// Create jobs with different characteristics
	jobs := createMixedWorkload()

	pool.AddJobs(jobs)

	fmt.Printf("   Processing %d jobs with adaptive strategy...\n", len(jobs))

	start := time.Now()
	_, err := pool.Run()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}

	fmt.Printf("   Completed in %v\n", duration)
	fmt.Printf("   Processed: %d, Failed: %d\n",
		pool.GetMetrics().ProcessedJobs,
		pool.GetMetrics().FailedJobs)
}

// demoCustomStrategy shows how to implement a custom strategy
func demoCustomStrategy() {
	fmt.Println("2. Implementing a Custom Strategy")
	fmt.Println("   This demonstrates how to create your own distribution strategy.")

	// Create a custom strategy
	customStrategy := &CustomBatchStrategy[string, string]{}

	// Use the custom strategy in a custom worker pool
	results := runWithCustomStrategy(customStrategy, createMixedWorkload(), stringProcessor)

	fmt.Printf("   Custom strategy completed %d jobs\n", len(results))
}

// demoStrategyComparison compares different strategies
func demoStrategyComparison() {
	fmt.Println("3. Strategy Performance Comparison")
	fmt.Println("   Comparing different strategies with the same workload.")

	jobs := createMixedWorkload()
	strategies := []workerpool.DistributionStrategy{
		workerpool.RoundRobin,
		workerpool.Chunked,
		workerpool.WorkStealing,
		workerpool.PriorityBased,
		workerpool.Adaptive,
	}

	fmt.Printf("   Workload: %d jobs, 4 workers\n", len(jobs))
	fmt.Println("   Strategy          | Duration    | Jobs/sec")
	fmt.Println("   ------------------|-------------|----------")

	for _, strategy := range strategies {
		duration := benchmarkStrategy(strategy, jobs)
		jobsPerSec := float64(len(jobs)) / duration.Seconds()
		fmt.Printf("   %-17s | %-11v | %.1f\n",
			getStrategyName(strategy), duration, jobsPerSec)
	}
}

// CustomBatchStrategy implements a custom distribution strategy
// This demonstrates how to implement the Strategy interface
type CustomBatchStrategy[T any, R any] struct {
	batchSize int
}

// Name returns the strategy name
func (s *CustomBatchStrategy[T, R]) Name() string {
	return "Custom Batch Strategy"
}

// Execute runs the custom batch strategy
func (s *CustomBatchStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config,
	jobs []workerpool.Job[T], processor workerpool.Processor[T, R],
	results chan<- workerpool.Result[R]) error {

	if s.batchSize == 0 {
		s.batchSize = 5 // Default batch size
	}

	var wg sync.WaitGroup

	// Process jobs in batches
	for i := 0; i < len(jobs); i += s.batchSize {
		end := i + s.batchSize
		if end > len(jobs) {
			end = len(jobs)
		}

		batch := jobs[i:end]
		wg.Add(1)
		go s.processBatch(ctx, batch, processor, results, &wg)
	}

	wg.Wait()
	close(results)

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// processBatch processes a batch of jobs
func (s *CustomBatchStrategy[T, R]) processBatch(ctx context.Context, batch []workerpool.Job[T],
	processor workerpool.Processor[T, R], results chan<- workerpool.Result[R],
	wg *sync.WaitGroup) {

	defer wg.Done()

	for _, job := range batch {
		select {
		case <-ctx.Done():
			return
		default:
			s.processJob(job, processor, results)
		}
	}
}

// processJob handles individual job processing
func (s *CustomBatchStrategy[T, R]) processJob(job workerpool.Job[T],
	processor workerpool.Processor[T, R], results chan<- workerpool.Result[R]) {

	startTime := time.Now()

	// Create a context for this job
	jobCtx := context.Background()
	result, err := processor(jobCtx, job)

	completed := time.Now()
	duration := completed.Sub(startTime)

	// Send result to channel
	results <- workerpool.Result[R]{
		JobID:     job.ID,
		Data:      result,
		Error:     err,
		Worker:    0, // Custom strategy doesn't use worker IDs
		Started:   startTime,
		Completed: completed,
		Duration:  duration,
	}
}

// runWithCustomStrategy demonstrates using a custom strategy
func runWithCustomStrategy[T any, R any](strategy workerpool.Strategy[T, R], jobs []workerpool.Job[T], processor workerpool.Processor[T, R]) []workerpool.Result[R] {
	config := &workerpool.Config{
		NumWorkers: 4,
		BufferSize: 100,
		Timeout:    30 * time.Second,
		MaxRetries: 2,
	}

	// Create a results channel
	results := make(chan workerpool.Result[R], len(jobs))

	// Execute the custom strategy
	err := strategy.Execute(context.Background(), config, jobs, processor, results)
	if err != nil {
		fmt.Printf("   Custom strategy error: %v\n", err)
		return nil
	}

	// Collect results
	var allResults []workerpool.Result[R]
	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults
}

// Helper functions

// stringProcessor is a simple string processor for examples
func stringProcessor(ctx context.Context, job workerpool.Job[string]) (string, error) {
	// Simulate some processing time
	time.Sleep(10 * time.Millisecond)

	// Process the string (convert to uppercase)
	result := fmt.Sprintf("PROCESSED: %s", job.Data)

	// Simulate occasional errors for priority 1 jobs
	if job.Priority == 1 {
		return result, fmt.Errorf("simulated error for priority 1 job")
	}

	return result, nil
}

// createMixedWorkload creates a variety of jobs for testing
func createMixedWorkload() []workerpool.Job[string] {
	var jobs []workerpool.Job[string]

	// High priority jobs
	for i := 0; i < 5; i++ {
		jobs = append(jobs, workerpool.Job[string]{
			ID:       fmt.Sprintf("high_%d", i),
			Data:     fmt.Sprintf("high priority task %d", i),
			Priority: 8,
			Created:  time.Now(),
		})
	}

	// Medium priority jobs
	for i := 0; i < 10; i++ {
		jobs = append(jobs, workerpool.Job[string]{
			ID:       fmt.Sprintf("med_%d", i),
			Data:     fmt.Sprintf("medium priority task %d", i),
			Priority: 5,
			Created:  time.Now(),
		})
	}

	// Low priority jobs
	for i := 0; i < 15; i++ {
		jobs = append(jobs, workerpool.Job[string]{
			ID:       fmt.Sprintf("low_%d", i),
			Data:     fmt.Sprintf("low priority task %d", i),
			Priority: 2,
			Created:  time.Now(),
		})
	}

	return jobs
}

// benchmarkStrategy measures the performance of a strategy
func benchmarkStrategy(strategy workerpool.DistributionStrategy, jobs []workerpool.Job[string]) time.Duration {
	config := workerpool.Config{
		NumWorkers: 4,
		Strategy:   strategy,
		Timeout:    30 * time.Second,
		MaxRetries: 2,
	}

	pool := workerpool.NewWithConfig[string, string](config).
		WithProcessor(stringProcessor)

	pool.AddJobs(jobs)

	start := time.Now()
	_, err := pool.Run()
	if err != nil {
		return 0
	}

	return time.Since(start)
}

// getStrategyName returns a human-readable name for a strategy
func getStrategyName(strategy workerpool.DistributionStrategy) string {
	switch strategy {
	case workerpool.RoundRobin:
		return "Round Robin"
	case workerpool.Chunked:
		return "Chunked"
	case workerpool.WorkStealing:
		return "Work Stealing"
	case workerpool.PriorityBased:
		return "Priority Based"
	case workerpool.Adaptive:
		return "Adaptive"
	default:
		return "Unknown"
	}
}
