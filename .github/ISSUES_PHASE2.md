# Phase 2: Basic Resource-Aware Scheduling Issues (Week 3-4)

## Issue #13: Implement resource tracking and admission control system

**Type**: enhancement  
**Priority**: high  
**Labels**: phase2, scheduling, resources, admission-control  
**Milestone**: Phase 2 - Basic Resource-Aware Scheduling  
**Assignees**: (to be assigned)  

### Description
Implement a thread-safe resource tracking system and admission control logic to prevent resource exhaustion and enable intelligent task scheduling.

### Background
Resource-aware scheduling requires tracking allocated resources and making intelligent decisions about which tasks can be admitted based on current system capacity.

### Requirements
- [ ] Implement thread-safe resource tracking
- [ ] Add admission control logic with different strategies
- [ ] Add metrics for admission decisions (accepted/rejected/deferred)
- [ ] Add comprehensive unit tests for edge cases

### Technical Details
**Files to create:**
- `scheduler/resource_tracker.go` - Resource allocation tracking
- `scheduler/admission_controller.go` - Task admission decisions

**Key structures:**
```go
type ResourceTracker struct {
    limits      ResourceLimits
    allocated   map[string]ResourceRequirement // taskID -> resources
    mu          sync.RWMutex
}

func (rt *ResourceTracker) CanAllocate(req ResourceRequirement) bool
func (rt *ResourceTracker) Allocate(taskID string, req ResourceRequirement) error
func (rt *ResourceTracker) Release(taskID string)
func (rt *ResourceTracker) CurrentAllocation() ResourceUsage

type AdmissionController struct {
    tracker     *ResourceTracker
    monitor     ResourceMonitor
    strategy    AdmissionStrategy
}

func (ac *AdmissionController) Admit(task *Task) AdmissionDecision
type AdmissionDecision struct {
    Allowed bool
    Reason  string
    RetryAfter time.Duration
}
```

### Acceptance Criteria
- [ ] ResourceTracker maintains accurate allocation state
- [ ] Thread-safe operations with proper locking
- [ ] AdmissionController implements multiple strategies
- [ ] Metrics track all admission decisions
- [ ] Comprehensive test coverage for edge cases
- [ ] Graceful handling of resource exhaustion
- [ ] Support for task queuing when resources unavailable

### Implementation Notes
- Use RWMutex for concurrent access patterns
- Implement exponential backoff for retry scenarios
- Add metrics for monitoring admission patterns
- Consider implementing resource reservation for critical tasks
- Handle partial resource allocation gracefully

---

## Issue #14: Implement core resource-aware scheduling with multiple strategies

**Type**: enhancement  
**Priority**: high  
**Labels**: phase2, scheduling, algorithms, strategies  
**Milestone**: Phase 2 - Basic Resource-Aware Scheduling  
**Assignees**: (to be assigned)  

### Description
Implement the core resource-aware scheduling loop with multiple scheduling strategies to efficiently distribute tasks based on resource requirements and availability.

### Background
Different scheduling strategies work better for different workload patterns. This implementation provides multiple algorithms and allows runtime selection.

### Requirements
- [ ] Implement core scheduling loop with resource checks
- [ ] Add multiple scheduling strategies (FirstFit, BestFit, WorstFit)
- [ ] Implement task queueing when resources unavailable
- [ ] Add scheduling metrics and logs

### Technical Details
**Files to create:**
- `scheduler/resource_scheduler.go` - Core scheduling logic
- `scheduler/strategies.go` - Scheduling strategy implementations

**Key structures:**
```go
type ResourceAwareScheduler struct {
    workers         []*Worker
    taskQueue       PriorityQueue
    tracker         *ResourceTracker
    admissionCtrl   *AdmissionController
}

func (s *ResourceAwareScheduler) Schedule(task *Task) error
func (s *ResourceAwareScheduler) selectWorker(task *Task) *Worker
func (s *ResourceAwareScheduler) rebalance()

type SchedulingStrategy interface {
    Score(task *Task, worker *Worker, tracker *ResourceTracker) float64
}

type FirstFitStrategy struct{}
type BestFitStrategy struct{}
type WorstFitStrategy struct{}
```

### Acceptance Criteria
- [ ] Core scheduling loop handles resource constraints
- [ ] Multiple scheduling strategies implemented and tested
- [ ] Task queueing works when resources unavailable
- [ ] Scheduling metrics provide actionable insights
- [ ] Rebalancing logic improves resource utilization
- [ ] Comprehensive logging for debugging
- [ ] Performance benchmarks show acceptable overhead

### Implementation Notes
- Implement strategy pattern for easy strategy switching
- Use scoring functions for worker selection
- Consider resource fragmentation in strategy design
- Add circuit breakers for resource exhaustion scenarios
- Implement graceful degradation when resources limited

---

## Issue #15: Integrate resource scheduler with main worker pool

**Type**: enhancement  
**Priority**: high  
**Labels**: phase2, integration, pool, backward-compatibility  
**Milestone**: Phase 2 - Basic Resource-Aware Scheduling  
**Assignees**: (to be assigned)  

### Description
Integrate the resource-aware scheduler with the main worker pool while maintaining backward compatibility and providing configuration options.

### Background
The resource-aware scheduler needs to be integrated into the existing pool architecture without breaking existing functionality.

### Requirements
- [ ] Integrate resource scheduler with main pool
- [ ] Maintain backward compatibility (resource-aware is opt-in)
- [ ] Add configuration validation
- [ ] Update existing tests to cover new functionality

### Technical Details
**Files to modify:**
- `pool.go` - Integrate resource scheduler

**Key additions:**
```go
type Pool struct {
    scheduler       Scheduler
    resourceLimits  *ResourceLimits
    resourceMonitor ResourceMonitor
    options         PoolOptions
}

type PoolOptions struct {
    EnableResourceAwareScheduling bool
    ResourceLimits               ResourceLimits
    SchedulingStrategy           SchedulingStrategy
    MonitoringInterval           time.Duration
}

func NewPool(opts ...PoolOption) *Pool
func WithResourceAwareScheduling() PoolOption
func WithResourceLimits(limits ResourceLimits) PoolOption
func WithSchedulingStrategy(strategy SchedulingStrategy) PoolOption
```

### Acceptance Criteria
- [ ] Resource scheduler integrated with main pool
- [ ] Existing code continues to work without modification
- [ ] Resource-aware scheduling is opt-in feature
- [ ] Configuration validation prevents invalid setups
- [ ] All existing tests continue to pass
- [ ] New tests cover resource-aware functionality
- [ ] Configuration options are well-documented

### Implementation Notes
- Use functional options pattern for configuration
- Maintain existing pool constructors for backward compatibility
- Add feature flags for gradual rollout
- Consider migration path from existing pools
- Add configuration validation with helpful error messages
