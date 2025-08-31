package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-foundations/workerpool"
)

// This example demonstrates how to use a worker pool for HTTP requests
// Note: Import paths would be adjusted based on your module structure

func main() {
	fmt.Println("=== HTTP Requests with Worker Pool ===")

	// Example configuration for HTTP processing
	config := workerpool.Config{
		NumWorkers: 8, // More workers for I/O-bound tasks
		Strategy:   workerpool.RoundRobin,
		Timeout:    30 * time.Second,
		MaxRetries: 2,
	}

	pool := workerpool.NewWithConfig[string, *http.Response](config).
		WithProcessor(httpProcessor)

	// Add URLs to process
	urls := []workerpool.Job[string]{
		{ID: "1", Data: "https://httpbin.org/get", Priority: 1},
		{ID: "2", Data: "https://httpbin.org/status/200", Priority: 1},
		{ID: "3", Data: "https://httpbin.org/delay/1", Priority: 2},
		{ID: "4", Data: "https://httpbin.org/delay/2", Priority: 2},
		{ID: "5", Data: "https://httpbin.org/status/404", Priority: 3},
	}

	pool.AddJobs(urls)

	fmt.Printf("Processing %d URLs with %d workers...\n\n", len(urls), config.NumWorkers)

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
			fmt.Printf("%d. [Worker %d] %s → Status: %d (took %v)\n",
				i+1, result.Worker, result.JobID, result.Data.StatusCode, result.Duration)
		}
	}

	fmt.Printf("\nTotal time: %v\n", duration)
}

// httpProcessor makes HTTP requests
func httpProcessor(ctx context.Context, job workerpool.Job[string]) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", job.Data, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}
