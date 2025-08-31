# Phase 1: Foundation Issues (Week 1-2)

## Issue #10: Add resource type definitions and monitoring interfaces

**Type**: enhancement  
**Priority**: high  
**Labels**: phase1, foundation, resources, types  
**Milestone**: Phase 1 - Foundation  
**Assignees**: (to be assigned)  

### Description
Transform the workerpool library by adding comprehensive resource type definitions and monitoring interfaces to support resource-aware scheduling.

### Background
Current worker pools don't understand resource requirements, leading to resource exhaustion in production. This issue establishes the foundation for intelligent resource management.

### Requirements
- [ ] Define all resource type structures (CPU, Memory, IOPS, Network, GPU)
- [ ] Create monitoring interfaces with ResourceMonitor
- [ ] Add unit tests for resource calculations
- [ ] Document resource units and measurements

### Technical Details
**Files to create/modify:**
- `resources/types.go` - Core resource type definitions
- `resources/monitor.go` - Monitoring interfaces

**Key structures:**
```go
type ResourceType int
const (
    CPU ResourceType = iota
    Memory
    IOPS
    NetworkBandwidth
    GPU
)

type ResourceRequirement struct {
    Type     ResourceType
    Amount   float64
    Unit     string
}

type ResourceLimits struct {
    MaxCPU           float64
    MaxMemory        uint64
    MaxIOPS          int
    MaxNetworkBPS    uint64
}

type ResourceMonitor interface {
    CurrentUsage() ResourceUsage
    Available() ResourceAvailability
    StartMonitoring(ctx context.Context)
}
```

### Acceptance Criteria
- [ ] All resource types defined with proper constants
- [ ] ResourceRequirement struct supports all resource types
- [ ] ResourceLimits struct provides comprehensive limits
- [ ] ResourceMonitor interface is well-defined and extensible
- [ ] Unit tests cover all resource calculations
- [ ] Documentation explains resource units and measurements
- [ ] No breaking changes to existing API

### Implementation Notes
- Start with types and interfaces first
- Use Go idioms for resource representation
- Consider future extensibility for new resource types
- Add comprehensive validation for resource values
- Use context.Context for monitoring lifecycle

---

## Issue #11: Implement cross-platform system resource monitoring

**Type**: enhancement  
**Priority**: high  
**Labels**: phase1, foundation, monitoring, cross-platform  
**Milestone**: Phase 1 - Foundation  
**Assignees**: (to be assigned)  

### Description
Implement cross-platform system resource monitoring to gather real-time CPU, memory, and I/O statistics for resource-aware scheduling decisions.

### Background
Resource-aware scheduling requires accurate, real-time information about system resource availability. This must work across Linux, macOS, and Windows environments.

### Requirements
- [ ] Implement cross-platform resource monitoring (Linux, macOS, Windows)
- [ ] Add caching layer for metrics (avoid excessive syscalls)
- [ ] Benchmark monitoring overhead (<0.1% CPU usage)
- [ ] Add integration tests for each platform

### Technical Details
**Files to create:**
- `resources/system_monitor.go` - Common interface and logic
- `resources/system_monitor_linux.go` - Linux implementation using /proc
- `resources/system_monitor_darwin.go` - macOS implementation using syscall
- `resources/system_monitor_windows.go` - Windows implementation

**Key implementation details:**
```go
type SystemMonitor struct {
    interval time.Duration
    metrics  chan ResourceUsage
    cache    *ResourceCache
}

func (s *SystemMonitor) gatherCPUStats() float64
func (s *SystemMonitor) gatherMemoryStats() uint64
func (s *SystemMonitor) gatherIOStats() int
```

### Acceptance Criteria
- [ ] Linux monitoring works using /proc/stat, /proc/meminfo, /proc/diskstats
- [ ] macOS monitoring works using syscall and IOKit
- [ ] Windows monitoring works using WMI or Performance Counters
- [ ] Caching layer prevents excessive syscalls
- [ ] Monitoring overhead is <0.1% CPU usage
- [ ] Integration tests pass on all supported platforms
- [ ] Graceful fallback when platform-specific features unavailable

### Implementation Notes
- Use build tags for platform-specific code
- Implement exponential backoff for failed monitoring calls
- Cache results with appropriate TTL based on resource type
- Handle monitoring failures gracefully
- Consider using existing libraries like gopsutil for complex platforms

---

## Issue #12: Extend Task struct with resource requirements and builder API

**Type**: enhancement  
**Priority**: high  
**Labels**: phase1, foundation, task, api, builder  
**Milestone**: Phase 1 - Foundation  
**Assignees**: (to be assigned)  

### Description
Extend the existing Task struct to support resource requirements and provide a fluent builder API for easy task configuration.

### Background
Tasks need to specify their resource requirements for intelligent scheduling. A builder API makes this configuration intuitive and prevents errors.

### Requirements
- [ ] Extend existing Task struct with resource requirements
- [ ] Implement fluent builder API for task configuration
- [ ] Add validation for resource requirements
- [ ] Maintain backward compatibility with existing API

### Technical Details
**Files to modify/create:**
- `task.go` - Extend existing Task struct
- `task_profile.go` - Add task profiling support

**Key additions:**
```go
type Task struct {
    ID              string
    Fn              func() error
    Resources       *ResourceRequirement  // NEW
    ProfileData     *TaskProfile         // NEW
    // ... existing fields
}

type TaskBuilder struct {
    task *Task
}

func (t *TaskBuilder) WithCPU(cores float64) *TaskBuilder
func (t *TaskBuilder) WithMemory(size string) *TaskBuilder
func (t *TaskBuilder) WithIOIntensive() *TaskBuilder
func (t *TaskBuilder) WithResourceProfile(profile ResourceProfile) *TaskBuilder
func (t *TaskBuilder) Build() *Task
```

### Acceptance Criteria
- [ ] Task struct extended with Resources and ProfileData fields
- [ ] TaskBuilder provides fluent API for resource configuration
- [ ] Resource requirements are validated before task creation
- [ ] Existing code continues to work without modification
- [ ] Builder API supports all resource types
- [ ] Unit tests cover all builder methods
- [ ] Examples demonstrate builder usage

### Implementation Notes
- Make Resources field optional (nil means no specific requirements)
- Use pointer for Resources to distinguish "no requirements" from "zero requirements"
- Validate resource values are positive and within reasonable bounds
- Consider adding convenience methods for common resource patterns
- Maintain existing Task constructors for backward compatibility
