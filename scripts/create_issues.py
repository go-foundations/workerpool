#!/usr/bin/env python3
"""
Auto-create GitHub Issues from documentation files
Prerequisites: 
1. Install PyGithub: pip install PyGithub
2. Set GITHUB_TOKEN environment variable
3. Run from repository root
"""

import os
import sys
import re
from github import Github
from pathlib import Path

# Configuration
REPO_NAME = "go-foundations/workerpool"
MILESTONE_PREFIX = "Phase"

def create_milestone(github, repo, phase, title, description):
    """Create a milestone if it doesn't exist"""
    try:
        milestone = repo.create_milestone(
            title=title,
            description=description,
            state="open"
        )
        print(f"✅ Created milestone: {title}")
        return milestone
    except Exception as e:
        print(f"⚠️  Milestone {title} already exists or error: {e}")
        # Try to find existing milestone
        for m in repo.get_milestones():
            if m.title == title:
                return m
        return None

def create_issue(github, repo, issue_num, title, body, labels, milestone):
    """Create a GitHub issue"""
    try:
        issue = repo.create_issue(
            title=title,
            body=body,
            labels=labels,
            milestone=milestone
        )
        print(f"✅ Created Issue #{issue_num}: {title}")
        return issue
    except Exception as e:
        print(f"❌ Failed to create Issue #{issue_num}: {e}")
        return None

def main():
    # Check for GitHub token
    token = os.getenv('GITHUB_TOKEN')
    if not token:
        print("❌ Please set GITHUB_TOKEN environment variable")
        print("   export GITHUB_TOKEN=your_github_token")
        sys.exit(1)
    
    # Initialize GitHub client
    github = Github(token)
    repo = github.get_repo(REPO_NAME)
    
    print(f"🚀 Creating GitHub Issues for Resource-Aware Scheduling Implementation")
    print(f"Repository: {REPO_NAME}")
    print("")
    
    # Create milestones
    milestones = {}
    milestone_data = {
        1: "Foundation - Resource types, monitoring, task API",
        2: "Basic Resource-Aware Scheduling - Resource tracking, admission control",
        3: "Learning and Profiling - Profiling, prediction, auto-learning",
        4: "Advanced Scheduling Algorithms - Bin packing, optimization",
        5: "Observability and Monitoring - Metrics, dashboard, tracing",
        6: "Cloud Integration - Containers, cost, Kubernetes",
        7: "Testing and Documentation - Comprehensive testing, documentation"
    }
    
    print("📋 Creating milestones...")
    for phase, description in milestone_data.items():
        title = f"{MILESTONE_PREFIX} {phase}"
        milestone = create_milestone(github, repo, phase, title, description)
        if milestone:
            milestones[phase] = milestone
    
    print("✅ All milestones created")
    print("")
    
    # Phase 1 Issues
    print("🌱 Creating Phase 1 Issues...")
    
    # Issue 10
    create_issue(github, repo, 10, 
        "Add resource type definitions and monitoring interfaces",
        """## Description
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
- `resources/types.go` - Core resource type definitions
- `resources/monitor.go` - Monitoring interfaces

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
- Use context.Context for monitoring lifecycle""",
        ["phase1", "foundation", "resources", "types", "enhancement"],
        milestones[1]
    )
    
    # Issue 11
    create_issue(github, repo, 11,
        "Implement cross-platform system resource monitoring",
        """## Description
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
- `resources/system_monitor.go` - Common interface and logic
- `resources/system_monitor_linux.go` - Linux implementation using /proc
- `resources/system_monitor_darwin.go` - macOS implementation using syscall
- `resources/system_monitor_windows.go` - Windows implementation

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
- Consider using existing libraries like gopsutil for complex platforms""",
        ["phase1", "foundation", "monitoring", "cross-platform", "enhancement"],
        milestones[1]
    )
    
    # Issue 12
    create_issue(github, repo, 12,
        "Extend Task struct with resource requirements and builder API",
        """## Description
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
- `task.go` - Extend existing Task struct
- `task_profile.go` - Add task profiling support

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
- Use pointer for Resources to distinguish "no requirements" from "zero requirements"
- Validate resource values are positive and within reasonable bounds
- Consider adding convenience methods for common resource patterns
- Maintain existing Task constructors for backward compatibility""",
        ["phase1", "foundation", "task", "api", "builder", "enhancement"],
        milestones[1]
    )
    
    # Phase 2 Issues
    print("⚡ Creating Phase 2 Issues...")
    
    # Issue 13
    create_issue(github, repo, 13,
        "Implement resource tracking and admission control system",
        """## Description
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
- `scheduler/resource_tracker.go` - Resource allocation tracking
- `scheduler/admission_controller.go` - Task admission decisions

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
- Handle partial resource allocation gracefully""",
        ["phase2", "scheduling", "resources", "admission-control", "enhancement"],
        milestones[2]
    )
    
    # Issue 14
    create_issue(github, repo, 14,
        "Implement core resource-aware scheduling with multiple strategies",
        """## Description
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
- `scheduler/resource_scheduler.go` - Core scheduling logic
- `scheduler/strategies.go` - Scheduling strategy implementations

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
- Implement graceful degradation when resources limited""",
        ["phase2", "scheduling", "algorithms", "strategies", "enhancement"],
        milestones[2]
    )
    
    # Issue 15
    create_issue(github, repo, 15,
        "Integrate resource scheduler with main worker pool",
        """## Description
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
- `pool.go` - Integrate resource scheduler

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
- Add configuration validation with helpful error messages""",
        ["phase2", "integration", "pool", "backward-compatibility", "enhancement"],
        milestones[2]
    )
    
    print("")
    print("🎉 Phase 1 and 2 Issues created successfully!")
    print("")
    print("📋 Summary:")
    print("- Phase 1: 3 issues (Foundation)")
    print("- Phase 2: 3 issues (Basic Resource-Aware Scheduling)")
    print("")
    print("🔗 View issues at: https://github.com/{}/issues".format(REPO_NAME))
    print("📊 View milestones at: https://github.com/{}/milestones".format(REPO_NAME))
    print("")
    print("💡 To create remaining issues, run this script again with additional phases")
    print("✅ Ready for implementation!")

if __name__ == "__main__":
    main()
