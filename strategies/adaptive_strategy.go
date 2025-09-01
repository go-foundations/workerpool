package strategies

import (
	"context"
	"sync"
	"time"

	"github.com/go-foundations/workerpool"
)

// AdaptiveStrategy implements an adaptive distribution strategy that switches between
// strategies based on workload characteristics and performance metrics
type AdaptiveStrategy[T any, R any] struct {
	strategies map[string]workerpool.Strategy[T, R]
	metrics    *AdaptiveMetrics
	mu         sync.RWMutex
}

// AdaptiveMetrics tracks performance metrics for strategy selection
type AdaptiveMetrics struct {
	strategyPerformance map[string]float64
	lastSwitch          time.Time
	switchThreshold     float64
	mu                  sync.RWMutex
}

// NewAdaptiveStrategy creates a new adaptive strategy with default strategies
func NewAdaptiveStrategy[T any, R any]() *AdaptiveStrategy[T, R] {
	as := &AdaptiveStrategy[T, R]{
		strategies: make(map[string]workerpool.Strategy[T, R]),
		metrics: &AdaptiveMetrics{
			strategyPerformance: make(map[string]float64),
			switchThreshold:     0.2, // 20% performance difference threshold
		},
	}

	// Register default strategies
	as.RegisterStrategy("round_robin", &RoundRobinStrategy[T, R]{})
	as.RegisterStrategy("chunked", &ChunkedStrategy[T, R]{})
	as.RegisterStrategy("work_stealing", &WorkStealingStrategy[T, R]{})
	as.RegisterStrategy("priority_based", &PriorityBasedStrategy[T, R]{})

	return as
}

// RegisterStrategy adds a new strategy to the adaptive strategy
func (as *AdaptiveStrategy[T, R]) RegisterStrategy(name string, strategy workerpool.Strategy[T, R]) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.strategies[name] = strategy
}

// Name returns the strategy name
func (as *AdaptiveStrategy[T, R]) Name() string {
	return "Adaptive Strategy"
}

// Execute runs the adaptive distribution strategy
func (as *AdaptiveStrategy[T, R]) Execute(ctx context.Context, config *workerpool.Config,
	jobs []workerpool.Job[T], processor workerpool.Processor[T, R],
	results chan<- workerpool.Result[R]) error {

	// Analyze workload characteristics
	workloadType := as.analyzeWorkload(jobs, config)

	// Select best strategy for current workload
	selectedStrategy := as.selectStrategy(workloadType, jobs, config)

	// Execute the selected strategy
	startTime := time.Now()
	err := selectedStrategy.Execute(ctx, config, jobs, processor, results)
	duration := time.Since(startTime)

	// Update performance metrics
	as.updateMetrics(selectedStrategy.Name(), duration, len(jobs))

	return err
}

// analyzeWorkload determines the type of workload based on job characteristics
func (as *AdaptiveStrategy[T, R]) analyzeWorkload(jobs []workerpool.Job[T], config *workerpool.Config) string {
	if len(jobs) == 0 {
		return "empty"
	}

	// Analyze job priorities
	highPriorityCount := 0
	lowPriorityCount := 0
	for _, job := range jobs {
		if job.Priority > 5 {
			highPriorityCount++
		} else {
			lowPriorityCount++
		}
	}

	// Analyze job distribution
	jobCount := len(jobs)
	workerCount := config.NumWorkers

	// Determine workload type based on characteristics
	if highPriorityCount > jobCount/2 {
		return "priority_based"
	} else if jobCount < workerCount*2 {
		return "round_robin"
	} else if jobCount > workerCount*10 {
		return "chunked"
	} else {
		return "work_stealing"
	}
}

// selectStrategy chooses the best strategy based on workload analysis and historical performance
func (as *AdaptiveStrategy[T, R]) selectStrategy(workloadType string, jobs []workerpool.Job[T], config *workerpool.Config) workerpool.Strategy[T, R] {
	as.mu.RLock()
	defer as.mu.RUnlock()

	// Get the strategy for the workload type
	strategy, exists := as.strategies[workloadType]
	if !exists {
		// Fallback to round-robin if strategy not found
		strategy = as.strategies["round_robin"]
	}

	// Check if we should switch strategies based on performance metrics
	if as.shouldSwitchStrategy(workloadType) {
		// Find the best performing strategy
		bestStrategy := as.findBestStrategy()
		if bestStrategy != nil {
			strategy = bestStrategy
		}
	}

	return strategy
}

// shouldSwitchStrategy determines if we should switch strategies based on performance
func (as *AdaptiveStrategy[T, R]) shouldSwitchStrategy(currentStrategy string) bool {
	as.metrics.mu.RLock()
	defer as.metrics.mu.RUnlock()

	// Don't switch too frequently
	if time.Since(as.metrics.lastSwitch) < 5*time.Second {
		return false
	}

	currentPerf := as.metrics.strategyPerformance[currentStrategy]
	if currentPerf == 0 {
		return false
	}

	// Find the best performing strategy
	var bestStrategy string
	var bestPerf float64

	for strategy, perf := range as.metrics.strategyPerformance {
		if perf > bestPerf {
			bestPerf = perf
			bestStrategy = strategy
		}
	}

	// Switch if there's a significant performance difference
	if bestStrategy != currentStrategy &&
		(bestPerf-currentPerf)/currentPerf > as.metrics.switchThreshold {
		return true
	}

	return false
}

// findBestStrategy returns the strategy with the best performance
func (as *AdaptiveStrategy[T, R]) findBestStrategy() workerpool.Strategy[T, R] {
	as.metrics.mu.RLock()
	defer as.metrics.mu.RUnlock()

	var bestStrategy string
	var bestPerf float64

	for strategy, perf := range as.metrics.strategyPerformance {
		if perf > bestPerf {
			bestPerf = perf
			bestStrategy = strategy
		}
	}

	if bestStrategy != "" {
		return as.strategies[bestStrategy]
	}

	return nil
}

// updateMetrics updates performance metrics for the executed strategy
func (as *AdaptiveStrategy[T, R]) updateMetrics(strategyName string, duration time.Duration, jobCount int) {
	as.metrics.mu.Lock()
	defer as.metrics.mu.Unlock()

	// Calculate performance metric (jobs per second)
	performance := float64(jobCount) / duration.Seconds()

	// Use exponential moving average for smooth metrics
	alpha := 0.3 // Smoothing factor
	if currentPerf, exists := as.metrics.strategyPerformance[strategyName]; exists {
		performance = alpha*performance + (1-alpha)*currentPerf
	}

	as.metrics.strategyPerformance[strategyName] = performance
	as.metrics.lastSwitch = time.Now()
}

// GetPerformanceMetrics returns current performance metrics for all strategies
func (as *AdaptiveStrategy[T, R]) GetPerformanceMetrics() map[string]float64 {
	as.metrics.mu.RLock()
	defer as.metrics.mu.RUnlock()

	metrics := make(map[string]float64)
	for k, v := range as.metrics.strategyPerformance {
		metrics[k] = v
	}

	return metrics
}

// SetSwitchThreshold sets the performance threshold for strategy switching
func (as *AdaptiveStrategy[T, R]) SetSwitchThreshold(threshold float64) {
	as.metrics.mu.Lock()
	defer as.metrics.mu.Unlock()
	as.metrics.switchThreshold = threshold
}
