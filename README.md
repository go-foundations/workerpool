# WorkerPool - Generic Go Worker Pool Library

A high-performance, generic worker pool implementation for Go with multiple distribution strategies, comprehensive metrics, and enterprise-ready features.

## 🚀 Features

- **Generic Types**: Type-safe implementation using Go generics
- **Multiple Strategies**: Round-Robin, Chunked, Work Stealing, and Priority-based distribution
- **Configurable**: Customizable worker counts, timeouts, retries, and buffer sizes
- **Metrics**: Built-in performance monitoring and statistics
- **Error Handling**: Robust error handling with retry mechanisms
- **Context Support**: Full context cancellation and timeout support
- **Thread-Safe**: Concurrent access safe with proper synchronization
- **High Performance**: Optimized for both CPU and I/O-bound workloads

## 📦 Installation

```bash
go get github.com/go-foundations/workerpool
```

## 🎯 Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "strings"
    
    	"github.com/go-foundations/workerpool"
)

func main() {
    // Create a worker pool
    pool := workerpool.New[string, string]().
        WithProcessor(stringProcessor)

    // Add jobs
    jobs := []workerpool.Job[string]{
        {ID: "1", Data: "hello"},
        {ID: "2", Data: "world"},
        {ID: "3", Data: "golang"},
    }
    
    pool.AddJobs(jobs)

    // Run and get results
    results, err := pool.Run()
    if err != nil {
        panic(err)
    }

    // Process results
    for _, result := range results {
        if result.Error != nil {
            fmt.Printf("Error processing %s: %v\n", result.JobID, result.Error)
        } else {
            fmt.Printf("Processed %s: %s\n", result.JobID, result.Data)
        }
    }
}

func stringProcessor(ctx context.Context, job workerpool.Job[string]) (string, error) {
    return strings.ToUpper(job.Data), nil
}
```

### Advanced Configuration

```go
config := workerpool.Config{
    NumWorkers:    8,                    // Number of worker goroutines
    Strategy:      workerpool.Chunked,   // Distribution strategy
    BufferSize:    1000,                 // Channel buffer size
    Timeout:       5 * time.Minute,      // Overall timeout
    WorkerTimeout: 30 * time.Second,     // Per-worker timeout
    MaxRetries:    3,                    // Retry failed jobs
    EnableMetrics: true,                 // Collect performance metrics
}

pool := workerpool.NewWithConfig[string, string](config).
    WithProcessor(processor)
```

## 🔧 Distribution Strategies

### 1. Round-Robin (Default)
Distributes jobs evenly across workers in a round-robin fashion.

```go
pool := workerpool.New[string, string]().
    WithProcessor(processor)
// Automatically uses RoundRobin strategy
```

**Best for**: Even workload distribution, predictable processing order

### 2. Chunked
Divides jobs into chunks and assigns each chunk to a worker.

```go
config := workerpool.Config{
    NumWorkers: 4,
    Strategy:   workerpool.Chunked,
}
```

**Best for**: Batch processing, when jobs are related

### 3. Work Stealing 🚧 *Coming Soon*
Advanced strategy where idle workers can steal work from busy ones using the Chase-Lev work stealing deque algorithm.

```go
config := workerpool.Config{
    NumWorkers: 4,
    Strategy:   workerpool.WorkStealing, // Not yet implemented
}
```

**Best for**: Dynamic workloads, maximizing resource utilization
**Status**: Research phase - implementing Chase-Lev work stealing deque

### 4. Priority-Based 🚧 *Coming Soon*
Processes jobs based on priority levels with fair scheduling to prevent starvation.

```go
config := workerpool.Config{
    NumWorkers: 4,
    Strategy:   workerpool.PriorityBased, // Not yet implemented
}
```

**Best for**: Time-sensitive operations, SLA requirements
**Status**: Research phase - implementing priority queue with fair scheduling

## 📊 Metrics and Monitoring

The library provides comprehensive metrics for monitoring performance:

```go
results, err := pool.Run()
if err != nil {
    panic(err)
}

metrics := pool.GetMetrics()
fmt.Printf("Total jobs: %d\n", metrics.TotalJobs)
fmt.Printf("Processed: %d\n", metrics.ProcessedJobs)
fmt.Printf("Failed: %d\n", metrics.FailedJobs)
fmt.Printf("Total time: %v\n", metrics.TotalDuration)
fmt.Printf("Average time per job: %v\n", metrics.AverageDuration)
```

## 🗺️ Roadmap

### 🚧 **In Development**
- **Work Stealing Strategy**: Implementing Chase-Lev work stealing deque for dynamic load balancing
- **Priority-Based Strategy**: Priority queue with fair scheduling to prevent job starvation

### 🔬 **Research & Design**
- **Work Stealing**: Studying Chase-Lev algorithm variants and multi-level work stealing
- **Priority Scheduling**: Researching fair scheduling algorithms and dynamic priority adjustment
- **Performance Optimization**: Benchmarking against existing Go worker pool implementations

### 🎯 **Future Enhancements**
- **Adaptive Strategies**: Runtime strategy switching based on workload characteristics
- **Distributed Work Stealing**: Cross-process work distribution
- **Machine Learning Integration**: Predictive job scheduling based on historical patterns

## 🔄 Error Handling and Retries

Configure automatic retries for failed jobs:

```go
config := workerpool.Config{
    NumWorkers: 4,
    MaxRetries: 3, // Retry failed jobs up to 3 times
}

pool := workerpool.NewWithConfig[string, string](config).
    WithProcessor(processor)
```

## ⏱️ Timeout and Cancellation

Full context support for timeouts and cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Create pool with context
pool := workerpool.New[string, string]()
pool.ctx = ctx

// Or use the built-in timeout configuration
config := workerpool.Config{
    NumWorkers: 4,
    Timeout:    30 * time.Second,
}
```

## 📈 Performance Benchmarks

Run benchmarks to find optimal configuration for your workload:

```bash
cd workerpool/benchmarks
go test -bench=.
```

Benchmarks include:
- Different distribution strategies
- Various worker counts
- Different job sizes
- Processing time variations

## 🎯 Use Cases

### CPU-Intensive Tasks
```go
config := workerpool.Config{
    NumWorkers: runtime.NumCPU(), // Use all CPU cores
    Strategy:   workerpool.Chunked,
}
```

### I/O-Bound Operations
```go
config := workerpool.Config{
    NumWorkers: 16, // More workers for I/O operations
    Strategy:   workerpool.RoundRobin,
}
```

### HTTP Request Processing
```go
pool := workerpool.New[string, *http.Response]().
    WithProcessor(func(ctx context.Context, job workerpool.Job[string]) (*http.Response, error) {
        req, _ := http.NewRequestWithContext(ctx, "GET", job.Data, nil)
        return http.DefaultClient.Do(req)
    })
```

### File Processing
```go
pool := workerpool.New[string, []byte]().
    WithProcessor(func(ctx context.Context, job workerpool.Job[string]) ([]byte, error) {
        return os.ReadFile(job.Data)
    })
```

## 🏗️ Architecture

The library is built with these design principles:

- **Separation of Concerns**: Clear separation between job distribution and processing
- **Immutable Configuration**: Configuration is set once and cannot be changed
- **Builder Pattern**: Fluent API for easy configuration
- **Resource Management**: Proper cleanup and resource management
- **Error Propagation**: Clear error handling and propagation

## 🔒 Thread Safety

All operations are thread-safe:
- Multiple goroutines can safely add jobs
- Configuration is immutable after creation
- Results collection is synchronized
- Metrics are protected with read-write mutexes

## 🚨 Best Practices

1. **Choose the right strategy**:
   - Round-Robin: General purpose, even distribution
   - Chunked: Batch processing, related jobs
   - Work Stealing: Dynamic workloads
   - Priority-Based: Time-sensitive operations

2. **Size your workers appropriately**:
   - CPU-bound: `runtime.NumCPU()`
   - I/O-bound: `runtime.NumCPU() * 2` to `runtime.NumCPU() * 4`

3. **Handle errors gracefully**:
   - Always check for errors in results
   - Use retries for transient failures
   - Log failed jobs for debugging

4. **Monitor performance**:
   - Use built-in metrics
   - Set appropriate timeouts
   - Monitor worker utilization

## 🤝 Contributing

Contributions are welcome! Please see our contributing guidelines:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Go's concurrency patterns
- Built with performance and usability in mind
- Tested extensively with real-world workloads

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/go-foundations/workerpool/issues)
- **Discussions**: [GitHub Discussions](https://github.com/go-foundations/workerpool/discussions)
- **Documentation**: [GoDoc](https://pkg.go.dev/github.com/go-foundations/workerpool)

---

**Star this repository if you find it useful! ⭐**
