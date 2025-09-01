# Custom Strategy Example

This example demonstrates how to implement custom distribution strategies for the worker pool using the strategy pattern.

## What This Example Shows

1. **Using Built-in Adaptive Strategy** - Demonstrates the new adaptive strategy that automatically selects the best distribution method
2. **Implementing Custom Strategies** - Shows how to create your own distribution strategy by implementing the `Strategy` interface
3. **Strategy Performance Comparison** - Benchmarks different strategies to help you choose the right one for your workload

## Files

- `main.go` - Complete example implementation
- `README.md` - This documentation file

## Running the Example

```bash
cd examples/custom_strategy_example
go run main.go
```

## Key Concepts Demonstrated

### 1. Strategy Interface

All strategies must implement this interface:

```go
type Strategy[T any, R any] interface {
    Execute(ctx context.Context, config *Config, jobs []Job[T],
        processor Processor[T, R], results chan<- Result[R]) error
    Name() string
}
```

### 2. Custom Strategy Implementation

The example includes a `CustomBatchStrategy` that:

- Processes jobs in configurable batches
- Uses goroutines for concurrent batch processing
- Handles context cancellation properly
- Implements proper error handling

### 3. Adaptive Strategy Usage

Shows how to use the built-in adaptive strategy:

```go
config := workerpool.Config{
    NumWorkers: 4,
    Strategy:   workerpool.Adaptive, // Automatically selects best strategy
    Timeout:    30 * time.Second,
    MaxRetries: 2,
}
```

### 4. Performance Benchmarking

Compares all available strategies:

- Round Robin
- Chunked
- Work Stealing
- Priority Based
- Adaptive

## Expected Output

```
=== Custom Strategy Implementation Example ===

1. Using Built-in Adaptive Strategy
   The adaptive strategy automatically selects the best distribution method
   based on workload characteristics.
   Processing 30 jobs with adaptive strategy...
   Completed in 123.45ms
   Processed: 30, Failed: 0

==================================================

2. Implementing a Custom Strategy
   This demonstrates how to create your own distribution strategy.
   Custom strategy completed 30 jobs

==================================================

3. Strategy Performance Comparison
   Comparing different strategies with the same workload.
   Workload: 30 jobs, 4 workers
   Strategy          | Duration    | Jobs/sec
   ------------------|-------------|----------
   Round Robin      | 125.67ms    | 238.7
   Chunked          | 118.23ms    | 253.6
   Work Stealing    | 122.89ms    | 244.1
   Priority Based   | 119.45ms    | 251.1
   Adaptive         | 115.67ms    | 259.4
```

## Customizing the Example

### 1. Modify the Custom Strategy

Edit the `CustomBatchStrategy` to implement your own distribution logic:

```go
type MyCustomStrategy[T any, R any] struct {
    // Add your custom fields
    customParam string
}

func (s *MyCustomStrategy[T, R]) Name() string {
    return "My Custom Strategy"
}

func (s *MyCustomStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config, 
    jobs []workerpool.Job[T], processor workerpool.Processor[T, R], 
    results chan<- workerpool.Result[R]) error {
    
    // Implement your custom distribution logic here
    // ...
    
    return nil
}
```

### 2. Add New Strategies

Create additional strategies and add them to the comparison:

```go
strategies := []workerpool.DistributionStrategy{
    workerpool.RoundRobin,
    workerpool.Chunked,
    workerpool.WorkStealing,
    workerpool.PriorityBased,
    workerpool.Adaptive,
    // Add your custom strategy here
}
```

### 3. Test Different Workloads

Modify the `createMixedWorkload()` function to test different job characteristics:

```go
func createMixedWorkload() []workerpool.Job[string] {
    var jobs []workerpool.Job[string]
    
    // Add different types of jobs
    // - High priority jobs
    // - Large batches
    // - Mixed priorities
    // - Different job sizes
    
    return jobs
}
```

## Learning Points

1. **Interface Implementation** - How to implement the Strategy interface
2. **Concurrency Patterns** - Using goroutines and WaitGroups for parallel processing
3. **Context Handling** - Proper context cancellation and timeout handling
4. **Error Handling** - Graceful error handling and resource cleanup
5. **Performance Measurement** - Benchmarking and comparing different approaches
6. **Generic Types** - Using Go generics for type-safe implementations

## Next Steps

After understanding this example:

1. **Read the Documentation** - Check `docs/CUSTOM_STRATEGIES.md` for comprehensive guidance
2. **Experiment** - Try different distribution algorithms and measure their performance
3. **Integrate** - Add your custom strategies to your own projects
4. **Contribute** - Consider contributing new strategies to the main project

## Troubleshooting

### Common Issues

1. **Import Errors** - Ensure you're using the correct import path for the workerpool package
2. **Type Errors** - Make sure your strategy implements the correct generic types
3. **Context Issues** - Always check for context cancellation in long-running operations
4. **Resource Leaks** - Ensure channels are properly closed and goroutines terminate

### Getting Help

- Check the main project documentation
- Review the strategy interface requirements
- Look at existing strategy implementations for reference
- Test with simple workloads before scaling up
