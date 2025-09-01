package strategies

import (
	"context"
	"time"

	"github.com/go-foundations/workerpool"
)

// processJob handles the actual job processing with retries and metrics
// This is a common implementation used by all strategies
func processJob[T any, R any](workerID int, job workerpool.Job[T],
	processor workerpool.Processor[T, R], results chan<- workerpool.Result[R],
	config *workerpool.Config) {

	startTime := time.Now()

	var result R
	var err error

	// Process with retries
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Create a context for this job processing
		jobCtx := context.Background()
		if config.WorkerTimeout > 0 {
			var cancel context.CancelFunc
			jobCtx, cancel = context.WithTimeout(context.Background(), config.WorkerTimeout)
			defer cancel()
		}

		result, err = processor(jobCtx, job)
		if err == nil {
			break
		}
		if attempt < config.MaxRetries {
			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
		}
	}

	completed := time.Now()
	duration := completed.Sub(startTime)

	// Send result to channel
	results <- workerpool.Result[R]{
		JobID:     job.ID,
		Data:      result,
		Error:     err,
		Worker:    workerID,
		Started:   startTime,
		Completed: completed,
		Duration:  duration,
	}
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
