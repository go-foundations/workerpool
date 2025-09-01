package strategies

import (
	"context"

	"github.com/go-foundations/workerpool"
)

// DistributionStrategy defines how jobs are distributed to workers
type DistributionStrategy int

const (
	RoundRobin DistributionStrategy = iota
	Chunked
	WorkStealing
	PriorityBased
)

// Strategy defines the interface for job distribution strategies
type Strategy[T any, R any] interface {
	// Execute runs the strategy with the given configuration
	Execute(ctx context.Context, config *workerpool.Config, jobs []workerpool.Job[T],
		processor workerpool.Processor[T, R], results chan<- workerpool.Result[R]) error

	// Name returns the human-readable name of the strategy
	Name() string
}

// StrategyFactory creates strategy instances
type StrategyFactory[T any, R any] struct{}

// NewStrategyFactory creates a new strategy factory
func NewStrategyFactory[T any, R any]() *StrategyFactory[T, R] {
	return &StrategyFactory[T, R]{}
}

// CreateStrategy creates a strategy instance based on the distribution type
func (f *StrategyFactory[T, R]) CreateStrategy(strategy DistributionStrategy) Strategy[T, R] {
	switch strategy {
	case RoundRobin:
		return &RoundRobinStrategy[T, R]{}
	case Chunked:
		return &ChunkedStrategy[T, R]{}
	case WorkStealing:
		return &WorkStealingStrategy[T, R]{}
	case PriorityBased:
		return &PriorityBasedStrategy[T, R]{}
	default:
		return &RoundRobinStrategy[T, R]{} // Default fallback
	}
}

// GetStrategyName returns the human-readable name for a strategy
func GetStrategyName(strategy DistributionStrategy) string {
	switch strategy {
	case RoundRobin:
		return "Round Robin"
	case Chunked:
		return "Chunked"
	case WorkStealing:
		return "Work Stealing"
	case PriorityBased:
		return "Priority Based"
	default:
		return "Unknown"
	}
}
