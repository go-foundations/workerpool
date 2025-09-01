# Implementing Custom Distribution Strategies

This document provides a comprehensive guide on how to implement custom distribution strategies for the worker pool. The strategy pattern allows you to create new ways of distributing jobs across workers while maintaining the same interface.

## Table of Contents

1. [Strategy Interface](#strategy-interface)
2. [Basic Strategy Implementation](#basic-strategy-implementation)
3. [Advanced Strategy Features](#advanced-strategy-features)
4. [Testing Your Strategy](#testing-your-strategy)
5. [Performance Considerations](#performance-considerations)
6. [Real-World Examples](#real-world-examples)
7. [Best Practices](#best-practices)

## Strategy Interface

All strategies must implement the `Strategy[T, R]` interface:

```go
type Strategy[T any, R any] interface {
    // Execute runs the strategy with the given configuration
    Execute(ctx context.Context, config *Config, jobs []Job[T],
        processor Processor[T, R], results chan<- Result[R]) error
    
    // Name returns the human-readable name of the strategy
    Name() string
}
```

### Interface Methods

- **`Execute`**: The main method that implements your distribution logic
- **`Name`**: Returns a descriptive name for your strategy

## Basic Strategy Implementation

Here's a minimal example of a custom strategy:

```go
package mystrategies

import (
    "context"
    "sync"
    "time"
    
    "github.com/go-foundations/workerpool"
)

// SimpleBatchStrategy processes jobs in simple batches
type SimpleBatchStrategy[T any, R any] struct {
    batchSize int
}

// NewSimpleBatchStrategy creates a new batch strategy
func NewSimpleBatchStrategy[T any, R any](batchSize int) *SimpleBatchStrategy[T, R] {
    if batchSize <= 0 {
        batchSize = 10
    }
    return &SimpleBatchStrategy[T, R]{batchSize: batchSize}
}

// Name returns the strategy name
func (s *SimpleBatchStrategy[T, R]) Name() string {
    return "Simple Batch Strategy"
}

// Execute implements the batch processing logic
func (s *SimpleBatchStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config, 
    jobs []workerpool.Job[T], processor workerpool.Processor[T, R], 
    results chan<- workerpool.Result[R]) error {
    
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
    
    return nil
}

// processBatch handles a batch of jobs
func (s *SimpleBatchStrategy[T, R]) processBatch(ctx context.Context, batch []workerpool.Job[T], 
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
func (s *SimpleBatchStrategy[T, R]) processJob(job workerpool.Job[T], 
    processor workerpool.Processor[T, R], results chan<- workerpool.Result[R]) {
    
    startTime := time.Now()
    
    // Process the job
    result, err := processor(context.Background(), job)
    
    completed := time.Now()
    duration := completed.Sub(startTime)
    
    // Send result
    results <- workerpool.Result[R]{
        JobID:     job.ID,
        Data:      result,
        Error:     err,
        Worker:    0, // Custom strategies can use their own worker IDs
        Started:   startTime,
        Completed: completed,
        Duration:  duration,
    }
}
```

## Advanced Strategy Features

### 1. Workload Analysis

Advanced strategies can analyze job characteristics to make intelligent decisions:

```go
// analyzeWorkload determines optimal strategy parameters
func (s *AdvancedStrategy[T, R]) analyzeWorkload(jobs []workerpool.Job[T], config *workerpool.Config) WorkloadProfile {
    profile := WorkloadProfile{}
    
    // Analyze job priorities
    for _, job := range jobs {
        if job.Priority > 7 {
            profile.HighPriorityCount++
        } else if job.Priority < 3 {
            profile.LowPriorityCount++
        } else {
            profile.MediumPriorityCount++
        }
    }
    
    // Analyze job sizes (if applicable)
    profile.TotalJobs = len(jobs)
    profile.WorkerCount = config.NumWorkers
    
    return profile
}

type WorkloadProfile struct {
    HighPriorityCount   int
    MediumPriorityCount int
    LowPriorityCount    int
    TotalJobs           int
    WorkerCount         int
}
```

### 2. Dynamic Strategy Selection

Some strategies can switch between different approaches based on conditions:

```go
// selectOptimalMethod chooses the best processing method
func (s *AdaptiveStrategy[T, R]) selectOptimalMethod(profile WorkloadProfile) string {
    if profile.HighPriorityCount > profile.TotalJobs/2 {
        return "priority_first"
    } else if profile.TotalJobs > profile.WorkerCount*20 {
        return "large_batch"
    } else if profile.TotalJobs < profile.WorkerCount*2 {
        return "small_batch"
    }
    return "balanced"
}
```

### 3. Performance Metrics

Track and use performance data:

```go
type PerformanceMetrics struct {
    JobsProcessed int
    TotalDuration time.Duration
    AverageTime   time.Duration
    ErrorCount    int
    LastUpdated   time.Time
}

func (s *MetricsAwareStrategy[T, R]) updateMetrics(duration time.Duration, jobCount int, errorCount int) {
    s.metrics.mu.Lock()
    defer s.metrics.mu.Unlock()
    
    s.metrics.JobsProcessed += jobCount
    s.metrics.TotalDuration += duration
    s.metrics.ErrorCount += errorCount
    s.metrics.AverageTime = s.metrics.TotalDuration / time.Duration(s.metrics.JobsProcessed)
    s.metrics.LastUpdated = time.Now()
}
```

## Testing Your Strategy

### 1. Unit Tests

```go
func TestSimpleBatchStrategy(t *testing.T) {
    strategy := NewSimpleBatchStrategy[string, string](5)
    
    // Test Name method
    if strategy.Name() != "Simple Batch Strategy" {
        t.Errorf("Expected 'Simple Batch Strategy', got '%s'", strategy.Name())
    }
    
    // Test Execute method
    jobs := []workerpool.Job[string]{
        {ID: "1", Data: "test1", Priority: 1},
        {ID: "2", Data: "test2", Priority: 1},
    }
    
    config := &workerpool.Config{NumWorkers: 2, BufferSize: 10}
    results := make(chan workerpool.Result[string], len(jobs))
    
    processor := func(ctx context.Context, job workerpool.Job[string]) (string, error) {
        return "processed: " + job.Data, nil
    }
    
    err := strategy.Execute(context.Background(), config, jobs, processor, results)
    if err != nil {
        t.Errorf("Execute failed: %v", err)
    }
    
    // Verify results
    close(results)
    count := 0
    for range results {
        count++
    }
    if count != len(jobs) {
        t.Errorf("Expected %d results, got %d", len(jobs), count)
    }
}
```

### 2. Integration Tests

```go
func TestStrategyWithWorkerPool(t *testing.T) {
    // Create custom strategy
    strategy := NewSimpleBatchStrategy[string, string](3)
    
    // Create worker pool with custom strategy
    config := workerpool.Config{
        NumWorkers: 4,
        Strategy:   workerpool.Custom, // You'd need to add this
        Timeout:    10 * time.Second,
    }
    
    pool := workerpool.NewWithConfig[string, string](config)
    pool.WithProcessor(stringProcessor)
    
    // Add jobs and test
    jobs := createTestJobs(10)
    pool.AddJobs(jobs)
    
    results, err := pool.Run()
    if err != nil {
        t.Errorf("Pool execution failed: %v", err)
    }
    
    if len(results) != len(jobs) {
        t.Errorf("Expected %d results, got %d", len(jobs), len(results))
    }
}
```

## Performance Considerations

### 1. Memory Management

```go
// Good: Reuse buffers when possible
func (s *EfficientStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config, 
    jobs []workerpool.Job[T], processor workerpool.Processor[T, R], 
    results chan<- workerpool.Result[R]) error {
    
    // Reuse worker pool
    if s.workerPool == nil {
        s.workerPool = make([]chan workerpool.Job[T], config.NumWorkers)
    }
    
    // Process jobs...
    return nil
}

// Bad: Creating new buffers every time
func (s *InefficientStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config, 
    jobs []workerpool.Job[T], processor workerpool.Processor[T, R], 
    results chan<- workerpool.Result[R]) error {
    
    // This creates new channels every time
    workerPool := make([]chan workerpool.Job[T], config.NumWorkers)
    
    // Process jobs...
    return nil
}
```

### 2. Context Handling

```go
// Good: Proper context handling
func (s *GoodStrategy[T, R]) processJobs(ctx context.Context, jobs []workerpool.Job[T], 
    processor workerpool.Processor[T, R], results chan<- workerpool.Result[R]) {
    
    for _, job := range jobs {
        select {
        case <-ctx.Done():
            return // Exit early on cancellation
        default:
            // Process job
            s.processJob(job, processor, results)
        }
    }
}

// Bad: Ignoring context
func (s *BadStrategy[T, R]) processJobs(ctx context.Context, jobs []workerpool.Job[T], 
    processor workerpool.Processor[T, R], results chan<- workerpool.Result[R]) {
    
    // This ignores context cancellation
    for _, job := range jobs {
        s.processJob(job, processor, results)
    }
}
```

## Real-World Examples

### 1. Load-Aware Strategy

```go
// LoadAwareStrategy adjusts worker allocation based on system load
type LoadAwareStrategy[T any, R any] struct {
    loadThreshold float64
    metrics       *SystemMetrics
}

func (s *LoadAwareStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config, 
    jobs []workerpool.Job[T], processor workerpool.Processor[T, R], 
    results chan<- workerpool.Result[R]) error {
    
    // Get current system load
    currentLoad := s.metrics.GetCurrentLoad()
    
    // Adjust worker count based on load
    effectiveWorkers := s.calculateEffectiveWorkers(config.NumWorkers, currentLoad)
    
    // Use effective worker count for distribution
    return s.distributeWithWorkers(ctx, effectiveWorkers, jobs, processor, results)
}

func (s *LoadAwareStrategy[T, R]) calculateEffectiveWorkers(baseWorkers int, load float64) int {
    if load > s.loadThreshold {
        return max(1, baseWorkers/2) // Reduce workers under high load
    }
    return baseWorkers
}
```

### 2. Priority-Weighted Strategy

```go
// PriorityWeightedStrategy gives more workers to high-priority jobs
type PriorityWeightedStrategy[T any, R any] struct {
    priorityWeights map[int]float64
}

func (s *PriorityWeightedStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config, 
    jobs []workerpool.Job[T], processor workerpool.Processor[T, R], 
    results chan<- workerpool.Result[R]) error {
    
    // Group jobs by priority
    priorityGroups := s.groupJobsByPriority(jobs)
    
    // Allocate workers based on priority weights
    workerAllocation := s.allocateWorkers(config.NumWorkers, priorityGroups)
    
    // Process each priority group with allocated workers
    return s.processPriorityGroups(ctx, workerAllocation, priorityGroups, processor, results)
}
```

## Best Practices

### 1. Error Handling

```go
// Always handle errors gracefully
func (s *RobustStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config, 
    jobs []workerpool.Job[T], processor workerpool.Processor[T, R], 
    results chan<- workerpool.Result[R]) error {
    
    defer func() {
        if r := recover(); r != nil {
            // Log panic and return error
            log.Printf("Strategy panic: %v", r)
            // Don't close results channel on panic
        }
    }()
    
    // Your strategy logic here
    return nil
}
```

### 2. Resource Cleanup

```go
// Ensure proper cleanup
func (s *CleanStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config, 
    jobs []workerpool.Job[T], processor workerpool.Processor[T, R], 
    results chan<- workerpool.Result[R]) error {
    
    // Create resources
    workerPool := make([]chan workerpool.Job[T], config.NumWorkers)
    
    // Ensure cleanup happens
    defer func() {
        for _, ch := range workerPool {
            close(ch)
        }
    }()
    
    // Your strategy logic here
    return nil
}
```

### 3. Logging and Observability

```go
// Add logging for debugging and monitoring
func (s *ObservableStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config, 
    jobs []workerpool.Job[T], processor workerpool.Processor[T, R], 
    results chan<- workerpool.Result[R]) error {
    
    startTime := time.Now()
    log.Printf("Strategy %s starting with %d jobs, %d workers", 
        s.Name(), len(jobs), config.NumWorkers)
    
    defer func() {
        duration := time.Since(startTime)
        log.Printf("Strategy %s completed in %v", s.Name(), duration)
    }()
    
    // Your strategy logic here
    return nil
}
```

## Integration with Worker Pool

To use your custom strategy with the main worker pool, you have several options:

### Option 1: Extend the Worker Pool

```go
// Add your strategy to the main worker pool
func (wp *WorkerPool[T, R]) runCustomStrategy(ctx context.Context) error {
    customStrategy := &MyCustomStrategy[T, R]{}
    return customStrategy.Execute(ctx, &wp.config, wp.jobs, wp.processor, wp.results)
}
```

### Option 2: Use Strategy Factory

```go
// Create a strategy factory that includes your custom strategy
type ExtendedStrategyFactory[T any, R any] struct {
    strategies map[string]workerpool.Strategy[T, R]
}

func NewExtendedStrategyFactory[T any, R any]() *ExtendedStrategyFactory[T, R] {
    factory := &ExtendedStrategyFactory[T, R]{
        strategies: make(map[string]workerpool.Strategy[T, R]),
    }
    
    // Register built-in strategies
    factory.RegisterStrategy("round_robin", &workerpool.RoundRobinStrategy[T, R]{})
    factory.RegisterStrategy("custom", &MyCustomStrategy[T, R]{})
    
    return factory
}
```

## Conclusion

Implementing custom strategies gives you the flexibility to optimize job distribution for your specific use case. Remember to:

1. **Follow the interface contract** - Implement all required methods
2. **Handle errors gracefully** - Don't let panics crash the system
3. **Test thoroughly** - Unit tests and integration tests are essential
4. **Monitor performance** - Track metrics to ensure your strategy is effective
5. **Document clearly** - Explain how and when to use your strategy

The strategy pattern makes it easy to experiment with different approaches and find the optimal solution for your workload characteristics.
