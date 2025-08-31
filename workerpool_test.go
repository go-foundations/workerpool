package workerpool

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// WorkerPoolTestSuite holds test utilities and state
type WorkerPoolTestSuite struct {
	suite.Suite
}

// TestWorkerPoolTestSuite runs all tests in the suite
func TestWorkerPoolTestSuite(t *testing.T) {
	suite.Run(t, new(WorkerPoolTestSuite))
}

func (ts *WorkerPoolTestSuite) TestNewWorkerPool() {
	pool := New[string, string]()

	ts.NotNil(pool)
	ts.Equal(4, pool.config.NumWorkers)
	ts.Equal(RoundRobin, pool.config.Strategy)
	ts.Equal(5*time.Minute, pool.config.Timeout)
}

func (ts *WorkerPoolTestSuite) TestNewWithConfig() {
	config := Config{
		NumWorkers: 8,
		Strategy:   Chunked,
		Timeout:    1 * time.Minute,
	}

	pool := NewWithConfig[string, string](config)

	ts.NotNil(pool)
	ts.Equal(8, pool.config.NumWorkers)
	ts.Equal(Chunked, pool.config.Strategy)
	ts.Equal(1*time.Minute, pool.config.Timeout)
}

func (ts *WorkerPoolTestSuite) TestWithProcessor() {
	pool := New[string, string]()

	processor := func(ctx context.Context, job Job[string]) (string, error) {
		return strings.ToUpper(job.Data), nil
	}

	pool.WithProcessor(processor)
	ts.NotNil(pool.processor)
}

func (ts *WorkerPoolTestSuite) TestAddJobs() {
	pool := New[string, string]()

	jobs := []Job[string]{
		{ID: "1", Data: "hello"},
		{ID: "2", Data: "world"},
	}

	pool.AddJobs(jobs)
	ts.Len(pool.jobs, 2)
	ts.Equal("hello", pool.jobs[0].Data)
	ts.Equal("world", pool.jobs[1].Data)
	ts.Equal(2, pool.metrics.TotalJobs)
}

func (ts *WorkerPoolTestSuite) TestAddJob() {
	pool := New[string, string]()

	job := Job[string]{ID: "1", Data: "hello"}
	pool.AddJob(job)

	ts.Len(pool.jobs, 1)
	ts.Equal("hello", pool.jobs[0].Data)
	ts.Equal(1, pool.metrics.TotalJobs)
}

func (ts *WorkerPoolTestSuite) TestRunWithoutProcessor() {
	pool := New[string, string]()
	pool.AddJob(Job[string]{ID: "1", Data: "hello"})

	_, err := pool.Run()
	ts.Error(err)
	ts.Contains(err.Error(), "no processor configured")
}

func (ts *WorkerPoolTestSuite) TestRunWithoutJobs() {
	pool := New[string, string]()
	pool.WithProcessor(func(ctx context.Context, job Job[string]) (string, error) {
		return strings.ToUpper(job.Data), nil
	})

	_, err := pool.Run()
	ts.Error(err)
	ts.Contains(err.Error(), "no jobs to process")
}

func (ts *WorkerPoolTestSuite) TestRoundRobinStrategy() {
	pool := New[string, string]()
	pool.config.NumWorkers = 2

	processor := func(ctx context.Context, job Job[string]) (string, error) {
		return strings.ToUpper(job.Data), nil
	}

	jobs := []Job[string]{
		{ID: "1", Data: "hello"},
		{ID: "2", Data: "world"},
		{ID: "3", Data: "test"},
		{ID: "4", Data: "data"},
	}

	pool.WithProcessor(processor).AddJobs(jobs)

	results, err := pool.Run()
	ts.NoError(err)
	ts.Len(results, 4)

	// Verify all jobs were processed
	processedData := make(map[string]bool)
	for _, result := range results {
		ts.NoError(result.Error)
		processedData[result.Data] = true
	}

	ts.True(processedData["HELLO"])
	ts.True(processedData["WORLD"])
	ts.True(processedData["TEST"])
	ts.True(processedData["DATA"])
}

func (ts *WorkerPoolTestSuite) TestChunkedStrategy() {
	config := Config{
		NumWorkers: 2,
		Strategy:   Chunked,
		BufferSize: 100,
	}

	pool := NewWithConfig[string, string](config)

	processor := func(ctx context.Context, job Job[string]) (string, error) {
		return strings.ToUpper(job.Data), nil
	}

	jobs := []Job[string]{
		{ID: "1", Data: "hello"},
		{ID: "2", Data: "world"},
		{ID: "3", Data: "test"},
		{ID: "4", Data: "data"},
	}

	pool.WithProcessor(processor).AddJobs(jobs)

	results, err := pool.Run()
	ts.NoError(err)
	ts.Len(results, 4)

	// Verify all jobs were processed
	processedData := make(map[string]bool)
	for _, result := range results {
		ts.NoError(result.Error)
		processedData[result.Data] = true
	}

	ts.True(processedData["HELLO"])
	ts.True(processedData["WORLD"])
	ts.True(processedData["TEST"])
	ts.True(processedData["DATA"])
}

func (ts *WorkerPoolTestSuite) TestErrorHandling() {
	pool := New[string, string]()

	processor := func(ctx context.Context, job Job[string]) (string, error) {
		if job.Data == "error" {
			return "", fmt.Errorf("processing error")
		}
		return strings.ToUpper(job.Data), nil
	}

	jobs := []Job[string]{
		{ID: "1", Data: "hello"},
		{ID: "2", Data: "error"},
		{ID: "3", Data: "world"},
	}

	pool.WithProcessor(processor).AddJobs(jobs)

	results, err := pool.Run()
	ts.NoError(err)
	ts.Len(results, 3)

	// Check that error job failed
	errorJob := false
	successJobs := 0

	for _, result := range results {
		if result.Error != nil {
			errorJob = true
		} else {
			successJobs++
		}
	}

	ts.True(errorJob)
	ts.Equal(2, successJobs)
}

func (ts *WorkerPoolTestSuite) TestRetryMechanism() {
	config := Config{
		NumWorkers: 1,
		MaxRetries: 2,
	}

	pool := NewWithConfig[string, string](config)

	attempts := 0
	processor := func(ctx context.Context, job Job[string]) (string, error) {
		attempts++
		if attempts < 3 {
			return "", fmt.Errorf("attempt %d failed", attempts)
		}
		return strings.ToUpper(job.Data), nil
	}

	jobs := []Job[string]{
		{ID: "1", Data: "hello"},
	}

	pool.WithProcessor(processor).AddJobs(jobs)

	results, err := pool.Run()
	ts.NoError(err)
	ts.Len(results, 1)
	ts.NoError(results[0].Error)
	ts.Equal("HELLO", results[0].Data)
	ts.Equal(3, attempts) // Should have retried twice
}

func (ts *WorkerPoolTestSuite) TestContextCancellation() {
	pool := New[string, string]()

	processor := func(ctx context.Context, job Job[string]) (string, error) {
		time.Sleep(200 * time.Millisecond)
		return strings.ToUpper(job.Data), nil
	}

	jobs := []Job[string]{
		{ID: "1", Data: "hello"},
		{ID: "2", Data: "world"},
	}

	pool.WithProcessor(processor).AddJobs(jobs)

	// Cancel after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		pool.Stop()
	}()

	_, err := pool.Run()
	ts.Error(err)
	ts.Contains(err.Error(), "context canceled")
}

func (ts *WorkerPoolTestSuite) TestMetrics() {
	pool := New[string, string]()

	processor := func(ctx context.Context, job Job[string]) (string, error) {
		time.Sleep(10 * time.Millisecond)
		return strings.ToUpper(job.Data), nil
	}

	jobs := []Job[string]{
		{ID: "1", Data: "hello"},
		{ID: "2", Data: "world"},
	}

	pool.WithProcessor(processor).AddJobs(jobs)

	_, err := pool.Run()
	ts.NoError(err)

	metrics := pool.GetMetrics()
	ts.Equal(2, metrics.TotalJobs)
	ts.Equal(2, metrics.ProcessedJobs)
	ts.Equal(0, metrics.FailedJobs)
	ts.True(metrics.TotalDuration > 0)
	ts.True(metrics.AverageDuration > 0)
}

func (ts *WorkerPoolTestSuite) TestConcurrentAccess() {
	pool := New[string, string]()

	// Add jobs from multiple goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			pool.AddJob(Job[string]{ID: fmt.Sprintf("%d", id), Data: fmt.Sprintf("data%d", id)})
		}(i)
	}
	wg.Wait()

	ts.Equal(10, pool.metrics.TotalJobs)
}

func (ts *WorkerPoolTestSuite) TestJobPriority() {
	pool := New[string, string]()

	processor := func(ctx context.Context, job Job[string]) (string, error) {
		return fmt.Sprintf("%s_priority_%d", strings.ToUpper(job.Data), job.Priority), nil
	}

	jobs := []Job[string]{
		{ID: "1", Data: "low", Priority: 1},
		{ID: "2", Data: "high", Priority: 10},
		{ID: "3", Data: "medium", Priority: 5},
	}

	pool.WithProcessor(processor).AddJobs(jobs)

	results, err := pool.Run()
	ts.NoError(err)
	ts.Len(results, 3)

	// Verify priorities were preserved
	for _, result := range results {
		ts.NoError(result.Error)
		ts.Contains(result.Data, "priority")
	}
}

func (ts *WorkerPoolTestSuite) TestBufferOverflow() {
	config := Config{
		NumWorkers: 1,
		BufferSize: 1, // Very small buffer
	}

	pool := NewWithConfig[string, string](config)

	processor := func(ctx context.Context, job Job[string]) (string, error) {
		time.Sleep(50 * time.Millisecond) // Slow processing
		return strings.ToUpper(job.Data), nil
	}

	jobs := []Job[string]{
		{ID: "1", Data: "hello"},
		{ID: "2", Data: "world"},
		{ID: "3", Data: "test"},
	}

	pool.WithProcessor(processor).AddJobs(jobs)

	results, err := pool.Run()
	ts.NoError(err)
	ts.Len(results, 3)
}

func (ts *WorkerPoolTestSuite) TestZeroWorkers() {
	config := Config{
		NumWorkers: 0, // Should default to 1
	}

	pool := NewWithConfig[string, string](config)
	ts.Equal(1, pool.config.NumWorkers)
}

func (ts *WorkerPoolTestSuite) TestNegativeWorkers() {
	config := Config{
		NumWorkers: -5, // Should default to 1
	}

	pool := NewWithConfig[string, string](config)
	ts.Equal(1, pool.config.NumWorkers)
}

func (ts *WorkerPoolTestSuite) TestJobCreationTime() {
	pool := New[string, string]()

	before := time.Now()
	pool.AddJob(Job[string]{ID: "1", Data: "hello"})
	after := time.Now()

	ts.True(pool.jobs[0].Created.After(before) || pool.jobs[0].Created.Equal(before))
	ts.True(pool.jobs[0].Created.Before(after) || pool.jobs[0].Created.Equal(after))
}

func (ts *WorkerPoolTestSuite) TestResultMetadata() {
	pool := New[string, string]()

	processor := func(ctx context.Context, job Job[string]) (string, error) {
		return strings.ToUpper(job.Data), nil
	}

	jobs := []Job[string]{
		{ID: "1", Data: "hello"},
	}

	pool.WithProcessor(processor).AddJobs(jobs)

	results, err := pool.Run()
	ts.NoError(err)
	ts.Len(results, 1)

	result := results[0]
	ts.Equal("1", result.JobID)
	ts.Equal("HELLO", result.Data)
	ts.NoError(result.Error)
	ts.True(result.Started.After(result.Completed.Add(-time.Second)))
	ts.True(result.Duration >= 0)
}
