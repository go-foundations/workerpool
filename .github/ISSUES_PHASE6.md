# Phase 6: Cloud Integration Issues (Week 11-12)

## Issue #25: Add automatic container resource limit detection

**Type**: enhancement  
**Priority**: medium  
**Labels**: phase6, cloud, containers, kubernetes, docker  
**Milestone**: Phase 6 - Cloud Integration  
**Assignees**: (to be assigned)  

### Description
Implement automatic detection of container resource limits from various container platforms to enable optimal resource-aware scheduling in containerized environments.

### Background
Containerized applications need to respect resource limits set by the orchestration platform. Automatic detection ensures the worker pool operates within these constraints.

### Requirements
- [ ] Implement cgroups v1/v2 detection
- [ ] Add Kubernetes resource detection
- [ ] Add Docker resource detection
- [ ] Add tests with container environments

### Technical Details
**Files to create:**
- `cloud/container.go` - Common container interface
- `cloud/cgroups_linux.go` - Linux cgroups implementation
- `cloud/kubernetes.go` - Kubernetes resource detection

**Key structures:**
```go
type ContainerResourceDetector struct {
    platform ContainerPlatform
}

type ContainerPlatform interface {
    DetectLimits() (ResourceLimits, error)
    IsRunningInContainer() bool
}

// File: cloud/cgroups_linux.go
type CGroupsDetector struct{}
func (c *CGroupsDetector) ReadCPULimit() (float64, error)
func (c *CGroupsDetector) ReadMemoryLimit() (uint64, error)

// File: cloud/kubernetes.go
type KubernetesDetector struct{}
func (k *KubernetesDetector) ReadPodResources() (ResourceLimits, error)
```

### Acceptance Criteria
- [ ] cgroups v1/v2 detection works on Linux systems
- [ ] Kubernetes resource limits are detected automatically
- [ ] Docker resource limits are detected when available
- [ ] Graceful fallback when container detection fails
- [ ] Tests pass in containerized environments
- [ ] Performance impact of detection is minimal
- [ ] Configuration allows manual override of detected limits

### Implementation Notes
- Use build tags for platform-specific code
- Implement caching for resource limit detection
- Handle different cgroup versions gracefully
- Consider using Downward API for Kubernetes
- Add health checks for detection reliability

---

## Issue #26: Integrate cost calculation and optimization reporting

**Type**: enhancement  
**Priority**: low  
**Labels**: phase6, cloud, cost, optimization, reporting  
**Milestone**: Phase 6 - Cloud Integration  
**Assignees**: (to be assigned)  

### Description
Integrate cost calculation and optimization reporting to help users understand the financial impact of their resource usage and identify optimization opportunities.

### Background
Cloud costs are a significant concern for many organizations. Understanding the cost implications of resource allocation helps optimize spending.

### Requirements
- [ ] Implement cost calculation logic
- [ ] Add provider-specific pricing (AWS, GCP, Azure)
- [ ] Generate cost optimization reports
- [ ] Add cost prediction for tasks

### Technical Details
**Files to create:**
- `cloud/cost/calculator.go` - Core cost calculation
- `cloud/cost/providers.go` - Provider-specific implementations

**Key structures:**
```go
type CostCalculator struct {
    provider CloudProvider
    rates    PricingRates
}

func (cc *CostCalculator) EstimateCost(usage ResourceUsage, duration time.Duration) Cost
func (cc *CostCalculator) GenerateReport(start, end time.Time) CostReport

// File: cloud/cost/providers.go
type AWSCostProvider struct{}
type GCPCostProvider struct{}
type AzureCostProvider struct{}

type CostReport struct {
    TotalCost       float64
    CPUCost         float64
    MemoryCost      float64
    OptimizationTips []string
}
```

### Acceptance Criteria
- [ ] Cost calculation provides accurate estimates
- [ ] Multiple cloud providers are supported
- [ ] Cost reports identify optimization opportunities
- [ ] Cost prediction works for individual tasks
- [ ] Pricing data can be updated without code changes
- [ ] Cost calculations handle different resource types
- [ ] Comprehensive test coverage for cost scenarios

### Implementation Notes
- Use external pricing APIs when available
- Implement pricing data caching to reduce API calls
- Consider spot/preemptible instance pricing
- Add cost alerts and thresholds
- Implement cost optimization recommendations

---

## Issue #27: Create Kubernetes operator for WorkerPool management

**Type**: enhancement  
**Priority**: low  
**Labels**: phase6, cloud, kubernetes, operator, crd  
**Milestone**: Phase 6 - Cloud Integration  
**Assignees**: (to be assigned)  

### Description
Create a Kubernetes operator that manages WorkerPool resources, enabling declarative configuration and integration with Kubernetes ecosystem.

### Background
Kubernetes operators provide a way to manage complex applications declaratively. A WorkerPool operator would enable seamless integration with Kubernetes workflows.

### Requirements
- [ ] Create Kubernetes CRD
- [ ] Implement operator controller
- [ ] Add HPA integration
- [ ] Add example deployments

### Technical Details
**Files to create:**
- `operator/controller.go` - Operator controller logic
- `operator/crd.go` - Custom Resource Definition

**Key structures:**
```go
type WorkerPoolController struct {
    client kubernetes.Interface
}

func (w *WorkerPoolController) Reconcile(req ctrl.Request) (ctrl.Result, error)

// File: operator/crd.go
// Define Custom Resource Definition for WorkerPool
type WorkerPoolSpec struct {
    ResourceLimits   ResourceLimits
    SchedulingStrategy string
    AutoScaling      AutoScalingConfig
}
```

### Acceptance Criteria
- [ ] CRD is properly defined and validated
- [ ] Operator controller reconciles desired state
- [ ] HPA integration works with resource metrics
- [ ] Example deployments demonstrate usage
- [ ] Operator handles scaling and updates gracefully
- [ ] Comprehensive test coverage for operator logic
- [ ] Documentation explains operator usage

### Implementation Notes
- Use kubebuilder or operator-sdk for scaffolding
- Implement proper error handling and retry logic
- Consider using existing HPA patterns
- Add operator metrics and health checks
- Implement proper RBAC and security
