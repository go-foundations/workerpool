package benchmarks

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-foundations/workerpool"
)

// Benchmark different worker pool configurations
func BenchmarkRoundRobin(b *testing.B) {
	benchmarkStrategy(b, workerpool.RoundRobin)
}

func BenchmarkChunked(b *testing.B) {
	benchmarkStrategy(b, workerpool.Chunked)
}

func BenchmarkWorkStealing(b *testing.B) {
	benchmarkStrategy(b, workerpool.WorkStealing)
}

func benchmarkStrategy(b *testing.B, strategy workerpool.DistributionStrategy) {
	config := workerpool.Config{
		NumWorkers: 4,
		Strategy:   strategy,
		BufferSize: 1000,
		Timeout:    1 * time.Minute,
	}

	pool := workerpool.NewWithConfig[string, string](config).
		WithProcessor(benchmarkProcessor)

	// Create test jobs
	jobs := make([]workerpool.Job[string], 100)
	for i := 0; i < 100; i++ {
		jobs[i] = workerpool.Job[string]{
			ID:       fmt.Sprintf("job_%d", i),
			Data:     fmt.Sprintf("data_%d", i),
			Priority: i % 3,
		}
	}

	pool.AddJobs(jobs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pool.Run()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark different worker counts
func BenchmarkWorkerCounts(b *testing.B) {
	workerCounts := []int{1, 2, 4, 8, 16}

	for _, numWorkers := range workerCounts {
		b.Run(fmt.Sprintf("Workers_%d", numWorkers), func(b *testing.B) {
			config := workerpool.Config{
				NumWorkers: numWorkers,
				Strategy:   workerpool.RoundRobin,
				BufferSize: 1000,
				Timeout:    1 * time.Minute,
			}

			pool := workerpool.NewWithConfig[string, string](config).
				WithProcessor(benchmarkProcessor)

			jobs := make([]workerpool.Job[string], 100)
			for i := 0; i < 100; i++ {
				jobs[i] = workerpool.Job[string]{
					ID:       fmt.Sprintf("job_%d", i),
					Data:     fmt.Sprintf("data_%d", i),
					Priority: i % 3,
				}
			}

			pool.AddJobs(jobs)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := pool.Run()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Benchmark different job sizes
func BenchmarkJobSizes(b *testing.B) {
	jobSizes := []int{10, 100, 1000, 10000}

	for _, jobSize := range jobSizes {
		b.Run(fmt.Sprintf("Jobs_%d", jobSize), func(b *testing.B) {
			config := workerpool.Config{
				NumWorkers: 4,
				Strategy:   workerpool.RoundRobin,
				BufferSize: jobSize,
				Timeout:    1 * time.Minute,
			}

			pool := workerpool.NewWithConfig[string, string](config).
				WithProcessor(benchmarkProcessor)

			jobs := make([]workerpool.Job[string], jobSize)
			for i := 0; i < jobSize; i++ {
				jobs[i] = workerpool.Job[string]{
					ID:       fmt.Sprintf("job_%d", i),
					Data:     fmt.Sprintf("data_%d", i),
					Priority: i % 3,
				}
			}

			pool.AddJobs(jobs)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := pool.Run()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Benchmark with different processing times
func BenchmarkProcessingTimes(b *testing.B) {
	processingTimes := []time.Duration{
		0, // No delay
		1 * time.Microsecond,
		10 * time.Microsecond,
		100 * time.Microsecond,
		1 * time.Millisecond,
	}

	for _, procTime := range processingTimes {
		b.Run(fmt.Sprintf("ProcTime_%v", procTime), func(b *testing.B) {
			config := workerpool.Config{
				NumWorkers: 4,
				Strategy:   workerpool.RoundRobin,
				BufferSize: 100,
				Timeout:    1 * time.Minute,
			}

			pool := workerpool.NewWithConfig[string, string](config).
				WithProcessor(func(ctx context.Context, job workerpool.Job[string]) (string, error) {
					if procTime > 0 {
						time.Sleep(procTime)
					}
					return strings.ToUpper(job.Data), nil
				})

			jobs := make([]workerpool.Job[string], 100)
			for i := 0; i < 100; i++ {
				jobs[i] = workerpool.Job[string]{
					ID:       fmt.Sprintf("job_%d", i),
					Data:     fmt.Sprintf("data_%d", i),
					Priority: i % 3,
				}
			}

			pool.AddJobs(jobs)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := pool.Run()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// benchmarkProcessor is a simple processor for benchmarking
func benchmarkProcessor(ctx context.Context, job workerpool.Job[string]) (string, error) {
	// Simulate some minimal processing
	return strings.ToUpper(job.Data), nil
}
