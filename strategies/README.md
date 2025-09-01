# Worker Pool Strategies

This package contains the distribution strategies for the worker pool. Each strategy implements a different approach to distributing jobs across workers.

## Available Strategies

### 1. Round Robin (`round_robin.go`)
- Distributes jobs evenly across workers in round-robin fashion
- Each worker gets jobs in sequence
- Good for balanced load distribution
- Simple and predictable

### 2. Chunked (`chunked.go`)
- Divides jobs into chunks and assigns chunks to workers
- Each worker processes a contiguous slice of jobs
- Good for cache locality when jobs are related
- Efficient for large numbers of jobs

### 3. Work Stealing (`work_stealing.go`)
- Implements the Chase-Lev work stealing deque algorithm
- Workers steal work from other workers when their queue is empty
- Excellent for dynamic workloads with varying job processing times
- Prevents worker starvation

### 4. Priority Based (`priority_based.go`)
- Processes jobs based on priority using a binary heap
- Implements fair scheduling to prevent starvation
- Higher priority jobs are processed first
- Good for time-sensitive workloads

## Common Utilities

### `common.go`
Contains shared functionality used by all strategies:
- `processJob()`: Common job processing logic with retries and metrics
- `max()`: Utility function for finding the larger of two integers

### `strategy.go`
Defines the strategy interface and factory pattern:
- `Strategy[T, R]` interface that all strategies must implement
- `StrategyFactory[T, R]` for creating strategy instances
- Strategy constants and helper functions

## Usage

Strategies are automatically selected based on the `Config.Strategy` field when creating a worker pool:

```go
config := workerpool.Config{
    NumWorkers: 8,
    Strategy:   workerpool.WorkStealing, // or other strategies
    // ... other config options
}

pool := workerpool.NewWithConfig[string, string](config)
```

## Architecture

The strategy pattern allows for:
- Easy addition of new distribution strategies
- Runtime strategy selection
- Clean separation of concerns
- Testability of individual strategies
- Performance comparison between strategies

## Future Enhancements

Potential new strategies could include:
- Adaptive strategies that change based on workload
- Machine learning-based job distribution
- Resource-aware scheduling
- Batch processing strategies
