package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-foundations/workerpool"
)

func main() {
	fmt.Println("=== String Processing with Worker Pool ===")

	// Create a worker pool for string processing
	pool := workerpool.New[string, string]().
		WithProcessor(stringProcessor)

	// Add jobs
	jobs := []workerpool.Job[string]{
		{ID: "1", Data: "hello world", Priority: 1},
		{ID: "2", Data: "golang programming", Priority: 2},
		{ID: "3", Data: "concurrent processing", Priority: 3},
		{ID: "4", Data: "worker pool pattern", Priority: 1},
		{ID: "5", Data: "generic types", Priority: 2},
		{ID: "6", Data: "high performance", Priority: 3},
	}

	pool.AddJobs(jobs)

	fmt.Printf("Processing %d strings with %d workers...\n\n", len(jobs), pool.GetNumWorkers())

	// Run the worker pool
	start := time.Now()
	results, err := pool.Run()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Display results
	fmt.Println("Results:")
	fmt.Println("--------")
	for i, result := range results {
		if result.Error != nil {
			fmt.Printf("%d. [ERROR] %s: %v\n", i+1, result.JobID, result.Error)
		} else {
			fmt.Printf("%d. [Worker %d] %s → %s (took %v)\n",
				i+1, result.Worker, result.JobID, result.Data, result.Duration)
		}
	}

	// Display metrics
	metrics := pool.GetMetrics()
	fmt.Printf("\nMetrics:\n")
	fmt.Printf("--------\n")
	fmt.Printf("Total jobs: %d\n", metrics.TotalJobs)
	fmt.Printf("Processed: %d\n", metrics.ProcessedJobs)
	fmt.Printf("Failed: %d\n", metrics.FailedJobs)
	fmt.Printf("Total time: %v\n", metrics.TotalDuration)
	fmt.Printf("Average time per job: %v\n", metrics.AverageDuration)
	fmt.Printf("Overall duration: %v\n", duration)
}

// stringProcessor processes string jobs
func stringProcessor(ctx context.Context, job workerpool.Job[string]) (string, error) {
	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	// Process the string
	result := strings.ToUpper(job.Data)

	// Add some transformation based on priority
	switch job.Priority {
	case 1:
		result = "LOW: " + result
	case 2:
		result = "MED: " + result
	case 3:
		result = "HIGH: " + result
	}

	return result, nil
}
