#!/bin/bash

# Auto-create GitHub Issues from documentation
# Prerequisites: 
# 1. Install GitHub CLI: https://cli.github.com/
# 2. Authenticate: gh auth login
# 3. Run from repository root

set -e

REPO="go-foundations/workerpool"
MILESTONE_PREFIX="Phase"

echo "🚀 Creating GitHub Issues for Resource-Aware Scheduling Implementation"
echo "Repository: $REPO"
echo ""

# Function to create milestone if it doesn't exist
create_milestone() {
    local phase=$1
    local title="$MILESTONE_PREFIX $phase"
    local description=""
    
    case $phase in
        "1")
            description="Foundation - Resource types, monitoring, task API"
            ;;
        "2") 
            description="Basic Resource-Aware Scheduling - Resource tracking, admission control"
            ;;
        "3")
            description="Learning and Profiling - Profiling, prediction, auto-learning"
            ;;
        "4")
            description="Advanced Scheduling Algorithms - Bin packing, optimization"
            ;;
        "5")
            description="Observability and Monitoring - Metrics, dashboard, tracing"
            ;;
        "6")
            description="Cloud Integration - Containers, cost, Kubernetes"
            ;;
        "7")
            description="Testing and Documentation - Comprehensive testing, documentation"
            ;;
    esac
    
    echo "Checking milestone: $title"
    # Try to create milestone, but don't fail if it already exists
    gh api repos/$REPO/milestones --method POST --field title="$title" --field description="$description" --field state="open" 2>/dev/null || echo "Milestone already exists: $title"
}

# Function to create issue
create_issue() {
    local issue_num=$1
    local title=$2
    local body=$3
    local labels=$4
    local milestone=$5
    
    echo "Creating Issue #$issue_num: $title"
    
    # Create the issue
    gh issue create \
        --repo "$REPO" \
        --title "$title" \
        --body "$body" \
        --label "$labels" \
        --milestone "$MILESTONE_PREFIX $milestone"
    
    echo "✅ Issue #$issue_num created successfully"
    echo ""
}

# Create all milestones first
echo "📋 Creating milestones..."
for phase in {1..7}; do
    create_milestone $phase
done
echo "✅ All milestones created"
echo ""

# Phase 1 Issues
echo "🌱 Creating Phase 1 Issues..."
create_issue "10" "Add resource type definitions and monitoring interfaces" \
"## Description
Transform the workerpool library by adding comprehensive resource type definitions and monitoring interfaces to support resource-aware scheduling.

## Background
Current worker pools don't understand resource requirements, leading to resource exhaustion in production. This issue establishes the foundation for intelligent resource management.

## Requirements
- [ ] Define all resource type structures (CPU, Memory, IOPS, Network, GPU)
- [ ] Create monitoring interfaces with ResourceMonitor
- [ ] Add unit tests for resource calculations
- [ ] Document resource units and measurements

## Technical Details
**Files to create/modify:**
- \`resources/types.go\` - Core resource type definitions
- \`resources/monitor.go\` - Monitoring interfaces

## Acceptance Criteria
- [ ] All resource types defined with proper constants
- [ ] ResourceRequirement struct supports all resource types
- [ ] ResourceLimits struct provides comprehensive limits
- [ ] ResourceMonitor interface is well-defined and extensible
- [ ] Unit tests cover all resource calculations
- [ ] Documentation explains resource units and measurements
- [ ] No breaking changes to existing API

## Implementation Notes
- Start with types and interfaces first
- Use Go idioms for resource representation
- Consider future extensibility for new resource types
- Add comprehensive validation for resource values
- Use context.Context for monitoring lifecycle" \
"phase1,foundation,resources,types,enhancement" "1"

create_issue "11" "Implement cross-platform system resource monitoring" \
"## Description
Implement cross-platform system resource monitoring to gather real-time CPU, memory, and I/O statistics for resource-aware scheduling decisions.

## Background
Resource-aware scheduling requires accurate, real-time information about system resource availability. This must work across Linux, macOS, and Windows environments.

## Requirements
- [ ] Implement cross-platform resource monitoring (Linux, macOS, Windows)
- [ ] Add caching layer for metrics (avoid excessive syscalls)
- [ ] Benchmark monitoring overhead (<0.1% CPU usage)
- [ ] Add integration tests for each platform

## Technical Details
**Files to create:**
- \`resources/system_monitor.go\` - Common interface and logic
- \`resources/system_monitor_linux.go\` - Linux implementation using /proc
- \`resources/system_monitor_darwin.go\` - macOS implementation using syscall
- \`resources/system_monitor_windows.go\` - Windows implementation

## Acceptance Criteria
- [ ] Linux monitoring works using /proc/stat, /proc/meminfo, /proc/diskstats
- [ ] macOS monitoring works using syscall and IOKit
- [ ] Windows monitoring works using WMI or Performance Counters
- [ ] Caching layer prevents excessive syscalls
- [ ] Monitoring overhead is <0.1% CPU usage
- [ ] Integration tests pass on all supported platforms
- [ ] Graceful fallback when platform-specific features unavailable

## Implementation Notes
- Use build tags for platform-specific code
- Implement exponential backoff for failed monitoring calls
- Cache results with appropriate TTL based on resource type
- Handle monitoring failures gracefully
- Consider using existing libraries like gopsutil for complex platforms" \
"phase1,foundation,monitoring,cross-platform,enhancement" "1"

create_issue "12" "Extend Task struct with resource requirements and builder API" \
"## Description
Extend the existing Task struct to support resource requirements and provide a fluent builder API for easy task configuration.

## Background
Tasks need to specify their resource requirements for intelligent scheduling. A builder API makes this configuration intuitive and prevents errors.

## Requirements
- [ ] Extend existing Task struct with resource requirements
- [ ] Implement fluent builder API for task configuration
- [ ] Add validation for resource requirements
- [ ] Maintain backward compatibility with existing API

## Technical Details
**Files to modify/create:**
- \`task.go\` - Extend existing Task struct
- \`task_profile.go\` - Add task profiling support

## Acceptance Criteria
- [ ] Task struct extended with Resources and ProfileData fields
- [ ] TaskBuilder provides fluent API for resource configuration
- [ ] Resource requirements are validated before task creation
- [ ] Existing code continues to work without modification
- [ ] Builder API supports all resource types
- [ ] Unit tests cover all builder methods
- [ ] Examples demonstrate builder usage

## Implementation Notes
- Make Resources field optional (nil means no specific requirements)
- Use pointer for Resources to distinguish \"no requirements\" from \"zero requirements\"
- Validate resource values are positive and within reasonable bounds
- Consider adding convenience methods for common resource patterns
- Maintain existing Task constructors for backward compatibility" \
"phase1,foundation,task,api,builder,enhancement" "1"

# Phase 2 Issues
echo "⚡ Creating Phase 2 Issues..."
create_issue "13" "Implement resource tracking and admission control system" \
"## Description
Implement a thread-safe resource tracking system and admission control logic to prevent resource exhaustion and enable intelligent task scheduling.

## Background
Resource-aware scheduling requires tracking allocated resources and making intelligent decisions about which tasks can be admitted based on current system capacity.

## Requirements
- [ ] Implement thread-safe resource tracking
- [ ] Add admission control logic with different strategies
- [ ] Add metrics for admission decisions (accepted/rejected/deferred)
- [ ] Add comprehensive unit tests for edge cases

## Technical Details
**Files to create:**
- \`scheduler/resource_tracker.go\` - Resource allocation tracking
- \`scheduler/admission_controller.go\` - Task admission decisions

## Acceptance Criteria
- [ ] ResourceTracker maintains accurate allocation state
- [ ] Thread-safe operations with proper locking
- [ ] AdmissionController implements multiple strategies
- [ ] Metrics track all admission decisions
- [ ] Comprehensive test coverage for edge cases
- [ ] Graceful handling of resource exhaustion
- [ ] Support for task queuing when resources unavailable

## Implementation Notes
- Use RWMutex for concurrent access patterns
- Implement exponential backoff for retry scenarios
- Add metrics for monitoring admission patterns
- Consider implementing resource reservation for critical tasks
- Handle partial resource allocation gracefully" \
"phase2,scheduling,resources,admission-control,enhancement" "2"

create_issue "14" "Implement core resource-aware scheduling with multiple strategies" \
"## Description
Implement the core resource-aware scheduling loop with multiple scheduling strategies to efficiently distribute tasks based on resource requirements and availability.

## Background
Different scheduling strategies work better for different workload patterns. This implementation provides multiple algorithms and allows runtime selection.

## Requirements
- [ ] Implement core scheduling loop with resource checks
- [ ] Add multiple scheduling strategies (FirstFit, BestFit, WorstFit)
- [ ] Implement task queueing when resources unavailable
- [ ] Add scheduling metrics and logs

## Technical Details
**Files to create:**
- \`scheduler/resource_scheduler.go\` - Core scheduling logic
- \`scheduler/strategies.go\` - Scheduling strategy implementations

## Acceptance Criteria
- [ ] Core scheduling loop handles resource constraints
- [ ] Multiple scheduling strategies implemented and tested
- [ ] Task queueing works when resources unavailable
- [ ] Scheduling metrics provide actionable insights
- [ ] Rebalancing logic improves resource utilization
- [ ] Comprehensive logging for debugging
- [ ] Performance benchmarks show acceptable overhead

## Implementation Notes
- Implement strategy pattern for easy strategy switching
- Use scoring functions for worker selection
- Consider resource fragmentation in strategy design
- Add circuit breakers for resource exhaustion scenarios
- Implement graceful degradation when resources limited" \
"phase2,scheduling,algorithms,strategies,enhancement" "2"

create_issue "15" "Integrate resource scheduler with main worker pool" \
"## Description
Integrate the resource-aware scheduler with the main worker pool while maintaining backward compatibility and providing configuration options.

## Background
The resource-aware scheduler needs to be integrated into the existing pool architecture without breaking existing functionality.

## Requirements
- [ ] Integrate resource scheduler with main pool
- [ ] Maintain backward compatibility (resource-aware is opt-in)
- [ ] Add configuration validation
- [ ] Update existing tests to cover new functionality

## Technical Details
**Files to modify:**
- \`pool.go\` - Integrate resource scheduler

## Acceptance Criteria
- [ ] Resource scheduler integrated with main pool
- [ ] Existing code continues to work without modification
- [ ] Resource-aware scheduling is opt-in feature
- [ ] Configuration validation prevents invalid setups
- [ ] All existing tests continue to pass
- [ ] New tests cover resource-aware functionality
- [ ] Configuration options are well-documented

## Implementation Notes
- Use functional options pattern for configuration
- Maintain existing pool constructors for backward compatibility
- Add feature flags for gradual rollout
- Consider migration path from existing pools
- Add configuration validation with helpful error messages" \
"phase2,integration,pool,backward-compatibility,enhancement" "2"

# Phase 3 Issues
echo "🧠 Creating Phase 3 Issues..."
create_issue "16" "Implement comprehensive task profiling and storage system" \
"## Description
Implement a comprehensive task profiling system that tracks resource usage patterns and stores historical data for predictive scheduling decisions.

## Background
Learning-based resource estimation requires historical data about task execution patterns. This system provides the foundation for intelligent resource prediction.

## Requirements
- [ ] Implement profiling hooks in worker execution
- [ ] Add rolling statistics calculation
- [ ] Implement multiple storage backends (InMemory, File, Redis)
- [ ] Add profile export/import functionality

## Technical Details
**Files to create:**
- \`profiling/profiler.go\` - Core profiling logic
- \`profiling/storage.go\` - Storage interface and implementations

## Acceptance Criteria
- [ ] Profiling hooks capture resource usage during execution
- [ ] Rolling statistics provide accurate historical data
- [ ] Multiple storage backends work correctly
- [ ] Profile export/import supports multiple formats
- [ ] Thread-safe operations with proper locking
- [ ] Comprehensive test coverage for all storage backends
- [ ] Performance impact of profiling is minimal

## Implementation Notes
- Use sampling for high-frequency tasks to reduce overhead
- Implement circular buffers for memory-efficient storage
- Consider compression for long-term storage
- Add TTL for old profiles to prevent unbounded growth
- Implement profile versioning for schema evolution" \
"phase3,learning,profiling,storage,enhancement" "3"

create_issue "17" "Add predictive resource estimation with multiple strategies" \
"## Description
Implement multiple resource estimation strategies that can predict task resource requirements based on historical execution patterns.

## Background
Different estimation strategies work better for different types of workloads. This implementation provides multiple algorithms and confidence scoring.

## Requirements
- [ ] Implement multiple estimation strategies (Average, P95, ML-based)
- [ ] Add confidence scoring for predictions
- [ ] Add fallback for tasks without profiles
- [ ] Benchmark prediction performance

## Technical Details
**Files to create:**
- \`prediction/estimator.go\` - Core estimation logic
- \`prediction/strategies.go\` - Multiple estimation strategies

## Acceptance Criteria
- [ ] All estimation strategies implemented and tested
- [ ] Confidence scoring provides meaningful accuracy metrics
- [ ] Fallback logic handles tasks without profiles gracefully
- [ ] Prediction performance meets latency requirements
- [ ] Strategies can be combined for ensemble predictions
- [ ] Comprehensive test coverage for edge cases
- [ ] Performance benchmarks show acceptable overhead

## Implementation Notes
- Start with simple statistical strategies before ML
- Use cross-validation for confidence estimation
- Consider task similarity for better predictions
- Implement adaptive strategies that learn from prediction accuracy
- Add circuit breakers for unreliable predictions" \
"phase3,learning,prediction,estimation,enhancement" "3"

create_issue "18" "Integrate learning system with resource-aware scheduling" \
"## Description
Integrate the learning and profiling system with resource-aware scheduling to automatically use predictions for tasks without explicit resource requirements.

## Background
The learning system needs to be seamlessly integrated with the scheduler to provide automatic resource estimation and improve scheduling decisions.

## Requirements
- [ ] Integrate learning with scheduling decisions
- [ ] Add configuration for learning behavior
- [ ] Add metrics for prediction accuracy
- [ ] Add tests for learning convergence

## Technical Details
**Files to create:**
- \`scheduler/learning_scheduler.go\` - Learning-integrated scheduler

## Acceptance Criteria
- [ ] Learning system integrated with scheduling decisions
- [ ] Configuration controls learning behavior effectively
- [ ] Metrics track prediction accuracy over time
- [ ] Learning convergence tests pass consistently
- [ ] Fallback to default behavior when learning unavailable
- [ ] Performance impact of learning integration is minimal
- [ ] Comprehensive test coverage for learning scenarios

## Implementation Notes
- Make learning completely optional and configurable
- Use exponential backoff for failed predictions
- Consider A/B testing for learning strategies
- Implement learning rate adaptation based on accuracy
- Add learning health checks and circuit breakers" \
"phase3,learning,integration,scheduling,enhancement" "3"

# Phase 4 Issues
echo "📦 Creating Phase 4 Issues..."
create_issue "19" "Implement multi-dimensional bin packing algorithms" \
"## Description
Implement sophisticated bin packing algorithms for multi-dimensional resource scheduling to maximize resource utilization and minimize fragmentation.

## Background
Classic bin packing algorithms can significantly improve resource utilization compared to simple scheduling strategies, especially for heterogeneous workloads.

## Requirements
- [ ] Implement classic bin packing algorithms (FirstFit, BestFit, NextFit)
- [ ] Add multi-dimensional bin packing
- [ ] Add lookahead queue for better packing
- [ ] Benchmark against simple scheduling

## Technical Details
**Files to create:**
- \`scheduler/binpacking/types.go\` - Core bin packing structures
- \`scheduler/binpacking/algorithms.go\` - Algorithm implementations

## Acceptance Criteria
- [ ] All classic bin packing algorithms implemented
- [ ] Multi-dimensional packing handles resource constraints correctly
- [ ] Lookahead queue improves packing efficiency
- [ ] Performance benchmarks show significant improvement over simple scheduling
- [ ] Algorithms handle edge cases gracefully
- [ ] Comprehensive test coverage for all algorithms
- [ ] Configuration allows algorithm selection and tuning

## Implementation Notes
- Start with single-dimensional algorithms before multi-dimensional
- Use efficient data structures for bin representation
- Consider approximation algorithms for large problem sizes
- Implement adaptive algorithm selection based on workload
- Add metrics for packing efficiency and fragmentation" \
"phase4,algorithms,bin-packing,optimization,enhancement" "4"

create_issue "20" "Implement intelligent priority queue with dynamic scoring" \
"## Description
Implement an intelligent priority queue that dynamically re-scores tasks based on multiple factors to optimize scheduling decisions.

## Background
Static priority queues don't adapt to changing system conditions. Dynamic scoring provides better resource utilization and fairness.

## Requirements
- [ ] Implement priority queue with dynamic re-scoring
- [ ] Add multiple scoring strategies (ResourceEfficiency, WaitTime, Deadline)
- [ ] Add queue metrics (wait time, queue depth)
- [ ] Add tests for queue behavior under load

## Technical Details
**Files to create:**
- \`scheduler/queue/priority_queue.go\` - Core priority queue implementation
- \`scheduler/queue/scorers.go\` - Multiple scoring strategies

## Acceptance Criteria
- [ ] Priority queue supports dynamic re-scoring
- [ ] Multiple scoring strategies implemented and tested
- [ ] Queue metrics provide actionable insights
- [ ] Queue behavior is predictable under load
- [ ] Scoring strategies can be combined and weighted
- [ ] Performance overhead of re-scoring is acceptable
- [ ] Comprehensive test coverage for all scenarios

## Implementation Notes
- Use heap data structure for efficient priority operations
- Implement lazy re-scoring to minimize overhead
- Consider batch re-scoring for better performance
- Add circuit breakers for scoring failures
- Implement adaptive scoring based on system load" \
"phase4,algorithms,queue,priority,optimization,enhancement" "4"

create_issue "21" "Add defragmentation and placement optimization" \
"## Description
Implement resource fragmentation detection and task migration strategies to optimize resource utilization and prevent scheduling inefficiencies.

## Background
Over time, resource allocation can become fragmented, reducing scheduling efficiency. Defragmentation strategies can improve overall resource utilization.

## Requirements
- [ ] Implement fragmentation detection
- [ ] Add task migration for defragmentation
- [ ] Add placement optimization logic
- [ ] Add tests for fragmentation scenarios

## Technical Details
**Files to create:**
- \`scheduler/defragmentation.go\` - Fragmentation detection and resolution
- \`scheduler/placement.go\` - Optimal placement algorithms

## Acceptance Criteria
- [ ] Fragmentation detection provides accurate metrics
- [ ] Defragmentation strategies improve resource utilization
- [ ] Task migration handles edge cases gracefully
- [ ] Placement optimization reduces fragmentation
- [ ] Performance impact of defragmentation is minimal
- [ ] Comprehensive test coverage for all scenarios
- [ ] Configuration allows tuning of defragmentation behavior

## Implementation Notes
- Use conservative thresholds to avoid excessive migration
- Implement incremental defragmentation for large systems
- Consider migration cost vs. benefit analysis
- Add rollback mechanisms for failed migrations
- Implement defragmentation scheduling to minimize disruption" \
"phase4,algorithms,defragmentation,optimization,enhancement" "4"

# Phase 5 Issues
echo "📊 Creating Phase 5 Issues..."
create_issue "22" "Implement comprehensive metrics collection and export" \
"## Description
Implement a comprehensive metrics collection system that tracks scheduling performance, resource utilization, and system health with support for Prometheus and OpenTelemetry.

## Background
Production systems require detailed metrics to monitor performance, diagnose issues, and optimize resource utilization. This system provides the foundation for operational excellence.

## Requirements
- [ ] Implement comprehensive metrics collection
- [ ] Add Prometheus exporter
- [ ] Add OpenTelemetry support
- [ ] Add custom metrics API

## Technical Details
**Files to create:**
- \`metrics/collector.go\` - Core metrics collection
- \`metrics/prometheus.go\` - Prometheus integration

## Acceptance Criteria
- [ ] All critical metrics are collected and exposed
- [ ] Prometheus exporter provides standard metrics format
- [ ] OpenTelemetry integration works with existing observability stacks
- [ ] Custom metrics API allows application-specific metrics
- [ ] Metrics collection has minimal performance impact
- [ ] Comprehensive test coverage for all metric types
- [ ] Documentation explains all available metrics

## Implementation Notes
- Use atomic operations for thread-safe metric updates
- Implement metric batching for high-frequency updates
- Consider metric cardinality to prevent explosion
- Add metric validation and sanitization
- Implement metric lifecycle management" \
"phase5,observability,metrics,monitoring,enhancement" "5"

create_issue "23" "Create real-time debug dashboard with visualizations" \
"## Description
Create a real-time debug dashboard that provides visual insights into resource allocation, task scheduling, and system performance.

## Background
Visual representations of system state are essential for debugging complex scheduling issues and understanding resource utilization patterns.

## Requirements
- [ ] Implement HTTP dashboard server
- [ ] Add real-time WebSocket updates
- [ ] Create visualization endpoints (heatmaps, timelines, bin-packing view)
- [ ] Add basic HTML/JS frontend

## Technical Details
**Files to create:**
- \`dashboard/server.go\` - HTTP server and endpoints
- \`dashboard/visualizer.go\` - Visualization data generation

## Acceptance Criteria
- [ ] Dashboard server responds to HTTP requests
- [ ] WebSocket provides real-time updates
- [ ] Visualization endpoints generate valid JSON data
- [ ] Basic frontend displays key metrics
- [ ] Dashboard performance doesn't impact scheduling
- [ ] Authentication and access control implemented
- [ ] Responsive design works on different screen sizes

## Implementation Notes
- Use WebSocket for real-time updates
- Generate JSON data for frontend consumption
- Consider using existing charting libraries (D3.js, Chart.js)
- Implement rate limiting for dashboard endpoints
- Add dashboard configuration options" \
"phase5,observability,dashboard,visualization,enhancement" "5"

create_issue "24" "Add comprehensive task tracing and debugging capabilities" \
"## Description
Implement comprehensive task execution tracing and debugging capabilities to help diagnose scheduling issues and optimize performance.

## Background
Complex scheduling systems require detailed tracing to understand task lifecycle, identify bottlenecks, and debug resource allocation issues.

## Requirements
- [ ] Implement task execution tracing
- [ ] Add debug logging infrastructure
- [ ] Add slow task detection
- [ ] Add trace export functionality

## Technical Details
**Files to create:**
- \`debug/tracer.go\` - Task execution tracing
- \`debug/debugger.go\` - Debug mode configuration

## Acceptance Criteria
- [ ] Task tracing captures complete execution lifecycle
- [ ] Debug logging provides actionable information
- [ ] Slow task detection identifies performance issues
- [ ] Trace export supports multiple formats
- [ ] Tracing overhead is minimal in production
- [ ] Debug mode can be enabled/disabled at runtime
- [ ] Comprehensive test coverage for all tracing features

## Implementation Notes
- Use sampling for high-frequency tasks in production
- Implement trace buffering to handle high load
- Consider integration with existing tracing systems
- Add trace compression for long-running tasks
- Implement trace retention policies" \
"phase5,observability,tracing,debugging,enhancement" "5"

# Phase 6 Issues
echo "☁️ Creating Phase 6 Issues..."
create_issue "25" "Add automatic container resource limit detection" \
"## Description
Implement automatic detection of container resource limits from various container platforms to enable optimal resource-aware scheduling in containerized environments.

## Background
Containerized applications need to respect resource limits set by the orchestration platform. Automatic detection ensures the worker pool operates within these constraints.

## Requirements
- [ ] Implement cgroups v1/v2 detection
- [ ] Add Kubernetes resource detection
- [ ] Add Docker resource detection
- [ ] Add tests with container environments

## Technical Details
**Files to create:**
- \`cloud/container.go\` - Common container interface
- \`cloud/cgroups_linux.go\` - Linux cgroups implementation
- \`cloud/kubernetes.go\` - Kubernetes resource detection

## Acceptance Criteria
- [ ] cgroups v1/v2 detection works on Linux systems
- [ ] Kubernetes resource limits are detected automatically
- [ ] Docker resource limits are detected when available
- [ ] Graceful fallback when container detection fails
- [ ] Tests pass in containerized environments
- [ ] Performance impact of detection is minimal
- [ ] Configuration allows manual override of detected limits

## Implementation Notes
- Use build tags for platform-specific code
- Implement caching for resource limit detection
- Handle different cgroup versions gracefully
- Consider using Downward API for Kubernetes
- Add health checks for detection reliability" \
"phase6,cloud,containers,kubernetes,docker,enhancement" "6"

create_issue "26" "Integrate cost calculation and optimization reporting" \
"## Description
Integrate cost calculation and optimization reporting to help users understand the financial impact of their resource usage and identify optimization opportunities.

## Background
Cloud costs are a significant concern for many organizations. Understanding the cost implications of resource allocation helps optimize spending.

## Requirements
- [ ] Implement cost calculation logic
- [ ] Add provider-specific pricing (AWS, GCP, Azure)
- [ ] Generate cost optimization reports
- [ ] Add cost prediction for tasks

## Technical Details
**Files to create:**
- \`cloud/cost/calculator.go\` - Core cost calculation
- \`cloud/cost/providers.go\` - Provider-specific implementations

## Acceptance Criteria
- [ ] Cost calculation provides accurate estimates
- [ ] Multiple cloud providers are supported
- [ ] Cost reports identify optimization opportunities
- [ ] Cost prediction works for individual tasks
- [ ] Pricing data can be updated without code changes
- [ ] Cost calculations handle different resource types
- [ ] Comprehensive test coverage for cost scenarios

## Implementation Notes
- Use external pricing APIs when available
- Implement pricing data caching to reduce API calls
- Consider spot/preemptible instance pricing
- Add cost alerts and thresholds
- Implement cost optimization recommendations" \
"phase6,cloud,cost,optimization,reporting,enhancement" "6"

create_issue "27" "Create Kubernetes operator for WorkerPool management" \
"## Description
Create a Kubernetes operator that manages WorkerPool resources, enabling declarative configuration and integration with Kubernetes ecosystem.

## Background
Kubernetes operators provide a way to manage complex applications declaratively. A WorkerPool operator would enable seamless integration with Kubernetes workflows.

## Requirements
- [ ] Create Kubernetes CRD
- [ ] Implement operator controller
- [ ] Add HPA integration
- [ ] Add example deployments

## Technical Details
**Files to create:**
- \`operator/controller.go\` - Operator controller logic
- \`operator/crd.go\` - Custom Resource Definition

## Acceptance Criteria
- [ ] CRD is properly defined and validated
- [ ] Operator controller reconciles desired state
- [ ] HPA integration works with resource metrics
- [ ] Example deployments demonstrate usage
- [ ] Operator handles scaling and updates gracefully
- [ ] Comprehensive test coverage for operator logic
- [ ] Documentation explains operator usage

## Implementation Notes
- Use kubebuilder or operator-sdk for scaffolding
- Implement proper error handling and retry logic
- Consider using existing HPA patterns
- Add operator metrics and health checks
- Implement proper RBAC and security" \
"phase6,cloud,kubernetes,operator,crd,enhancement" "6"

# Phase 7 Issues
echo "🧪 Creating Phase 7 Issues..."
create_issue "28" "Add comprehensive testing suite with chaos testing" \
"## Description
Create a comprehensive testing suite that covers all resource-aware scheduling features, includes performance benchmarks, and implements chaos testing for production readiness.

## Background
Production systems require thorough testing to ensure reliability under various conditions. Chaos testing helps identify failure modes and resilience issues.

## Requirements
- [ ] Add integration tests for all features
- [ ] Add performance benchmarks
- [ ] Add chaos/stress tests
- [ ] Achieve >80% code coverage

## Technical Details
**Files to create:**
- \`tests/integration_test.go\` - Integration test suite
- \`tests/benchmark_test.go\` - Performance benchmarks
- \`tests/chaos_test.go\` - Chaos and stress testing

## Acceptance Criteria
- [ ] All resource-aware features have integration tests
- [ ] Performance benchmarks show acceptable overhead
- [ ] Chaos tests identify failure modes and recovery
- [ ] Code coverage exceeds 80% for new code
- [ ] Tests run reliably in CI/CD pipeline
- [ ] Performance regression tests catch degradations
- [ ] Chaos tests can be run in controlled environments

## Implementation Notes
- Use table-driven tests for comprehensive coverage
- Implement test fixtures for consistent test data
- Use test containers for integration testing
- Add performance regression detection
- Consider property-based testing for complex scenarios" \
"phase7,testing,chaos,integration,benchmarks,enhancement" "7"

create_issue "29" "Create comprehensive documentation and examples" \
"## Description
Create comprehensive documentation, code examples, and migration guides to help users adopt and effectively use the resource-aware scheduling features.

## Background
Good documentation is essential for user adoption and community growth. Examples and guides help users understand best practices and avoid common pitfalls.

## Requirements
- [ ] Write comprehensive documentation
- [ ] Add code examples for all features
- [ ] Create migration guide from other pools
- [ ] Add architecture diagrams

## Technical Details
**Files to create:**
- \`docs/getting-started.md\` - Quick start guide
- \`docs/resource-aware-scheduling.md\` - Feature documentation
- \`docs/api-reference.md\` - API documentation
- \`docs/best-practices.md\` - Usage guidelines
- \`docs/migration-guide.md\` - Migration from other pools

**Examples to create:**
- \`examples/basic/main.go\` - Basic usage
- \`examples/resource-limits/main.go\` - Resource limits
- \`examples/auto-learning/main.go\` - Auto-learning
- \`examples/kubernetes/deployment.yaml\` - K8s deployment
- \`examples/cost-optimization/main.go\` - Cost optimization

## Acceptance Criteria
- [ ] Documentation covers all features comprehensively
- [ ] Code examples demonstrate real-world usage
- [ ] Migration guide helps users transition smoothly
- [ ] Architecture diagrams explain system design
- [ ] Documentation is searchable and well-organized
- [ ] Examples run successfully and are tested
- [ ] Best practices guide prevents common mistakes

## Implementation Notes
- Use clear, concise language
- Include code snippets for all major features
- Add troubleshooting sections for common issues
- Consider video tutorials for complex concepts
- Implement documentation versioning" \
"phase7,documentation,examples,guides,enhancement" "7"

echo ""
echo "🎉 All 29 GitHub Issues created successfully!"
echo ""
echo "📋 Summary:"
echo "- Phase 1: 3 issues (Foundation)"
echo "- Phase 2: 3 issues (Basic Resource-Aware Scheduling)"
echo "- Phase 3: 3 issues (Learning and Profiling)"
echo "- Phase 4: 3 issues (Advanced Scheduling Algorithms)"
echo "- Phase 5: 3 issues (Observability and Monitoring)"
echo "- Phase 6: 3 issues (Cloud Integration)"
echo "- Phase 7: 2 issues (Testing and Documentation)"
echo ""
echo "🔗 View issues at: https://github.com/$REPO/issues"
echo "📊 View milestones at: https://github.com/$REPO/milestones"
echo ""
echo "✅ Ready for implementation!"
