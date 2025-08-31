# Repository Setup Guide

This guide provides comprehensive instructions for setting up, developing, and contributing to the Worker Pool package repository.

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or later
- Git
- Make (optional, but recommended)

### Initial Setup
```bash
# Clone the repository
git clone https://github.com/go-foundations/workerpool.git
cd workerpool

# Install dependencies
make deps

# Run validation
make validate
```

## 📁 Repository Structure

```
workerpool/
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions CI/CD
├── benchmarks/
│   └── performance_test.go     # Performance benchmarks
├── examples/
│   ├── http_example/
│   │   └── main.go            # HTTP requests example
│   └── string_example/
│       └── main.go            # String processing example
├── strategies/                 # Execution strategy implementations
├── workerpool/
│   └── workerpool.go          # Core worker pool implementation
├── .gitignore                 # Git ignore patterns
├── .pre-commit-config.yaml    # Pre-commit hooks
├── CHANGELOG.md               # Version history
├── CONTRIBUTING.md            # Contribution guidelines
├── LICENSE                    # Project license
├── Makefile                   # Development automation
├── README.md                  # Project documentation
├── SECURITY.md                # Security policies
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── workerpool.go              # Main package file
└── workerpool_test.go         # Test suite
```

## 🛠️ Development Tools

### Makefile Commands

The project includes a comprehensive Makefile for common development tasks:

#### 🔧 Development Workflow
```bash
make dev              # Run complete development workflow
make fmt              # Format code
make fmt-check        # Check code formatting
make lint             # Run linter
make security         # Run security scan
make vet              # Run go vet
```

#### 🧪 Testing
```bash
make test             # Run all tests
make test-coverage    # Run tests with coverage report
make benchmark        # Run benchmarks
make test-specific TEST=TestName  # Run specific test
```

#### 🔨 Building
```bash
make build            # Build package and examples
make build-examples   # Build only examples
```

#### 📦 Dependencies
```bash
make deps             # Install dependencies
make deps-update      # Update dependencies
```

#### 🚀 CI/CD
```bash
make ci               # Run full CI pipeline
make validate         # Run full validation
make pre-commit       # Run pre-commit checks
```

#### 📊 Analysis
```bash
make perf             # Run performance tests
make docs             # Generate documentation
```

#### 🧹 Maintenance
```bash
make clean            # Clean build artifacts
make help             # Show help message
```

### Pre-commit Hooks

The project includes pre-commit hooks for automatic code quality checks:

```bash
# Install pre-commit
pip install pre-commit

# Install git hooks
pre-commit install

# Run manually
pre-commit run --all-files
```

**Available Hooks:**
- Code formatting (gofmt)
- Linting (golangci-lint)
- Security scanning (gosec)
- Go vet
- Import organization
- Unit tests
- Build verification

## 🔄 CI/CD Pipeline

### GitHub Actions

The repository includes a comprehensive CI/CD pipeline that runs on:
- Push to main/develop branches
- Pull requests to main/develop branches

#### Pipeline Stages

1. **Test and Validate**
   - Multi-version Go testing (1.21, 1.22)
   - Race condition detection
   - Code coverage reporting
   - Formatting checks
   - Linting
   - Security scanning
   - Build verification

2. **Build and Package**
   - Example compilation
   - Artifact upload

3. **Release** (main branch only)
   - Automated release creation
   - Version tagging

### Local CI Simulation

```bash
# Run full CI pipeline locally
make ci

# Run individual stages
make validate          # Full validation
make pre-commit       # Pre-commit checks
```

## 📝 Code Quality Standards

### Go Code Style
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Run `go vet` before committing
- Enable race detection in tests

### Testing Requirements
- All new code must include tests
- Tests should cover success and error cases
- Use table-driven tests for multiple scenarios
- Include benchmarks for performance-critical code
- Run tests with race detection enabled

### Documentation
- All exported types, functions, and methods must have GoDoc comments
- Include examples in documentation
- Update README.md for new features
- Maintain CHANGELOG.md for version changes

## 🔒 Security

### Security Tools
```bash
# Run security scan
make security

# Install gosec
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
```

### Security Guidelines
- Never commit sensitive information
- Run security scans before submitting PRs
- Follow secure coding practices
- Report vulnerabilities privately

## 📊 Performance

### Benchmarking
```bash
# Run benchmarks
make benchmark

# Run specific benchmark
go test -bench=BenchmarkWorkerPool -benchmem ./benchmarks
```

### Performance Monitoring
- Built-in metrics collection
- Execution time tracking
- Resource usage monitoring
- Performance regression detection

## 🚀 Release Process

### Versioning
- Follow [Semantic Versioning](https://semver.org/)
- Update version in relevant files
- Create release notes

### Release Checklist
- [ ] All tests pass
- [ ] Documentation updated
- [ ] Examples tested
- [ ] Benchmarks run
- [ ] Security scan completed
- [ ] Release notes prepared

## 🤝 Contributing

### Before Contributing
1. Read [CONTRIBUTING.md](CONTRIBUTING.md)
2. Check existing issues and PRs
3. Ensure development environment is set up

### Development Workflow
1. Fork the repository
2. Create a feature branch
3. Make changes following coding standards
4. Add/update tests
5. Run validation: `make validate`
6. Submit pull request

### PR Requirements
- Descriptive title and description
- Reference related issues
- Include test results
- Ensure CI checks pass
- Follow PR template

## 📞 Support and Communication

### Getting Help
- [GitHub Issues](https://github.com/go-foundations/workerpool/issues)
- [GitHub Discussions](https://github.com/go-foundations/workerpool/discussions)
- [GoDoc Documentation](https://pkg.go.dev/github.com/go-foundations/workerpool)

### Security Issues
- **DO NOT** create public issues for security vulnerabilities
- Use [GitHub Security Advisories](https://github.com/go-foundations/workerpool/security/advisories/new)
- Follow [SECURITY.md](SECURITY.md) guidelines

## 🔧 Troubleshooting

### Common Issues

#### Build Failures
```bash
# Clean and rebuild
make clean
make build

# Check dependencies
make deps
```

#### Test Failures
```bash
# Run tests with verbose output
make test

# Check specific test
make test-specific TEST=TestName
```

#### Linting Issues
```bash
# Auto-fix formatting
make fmt

# Check formatting
make fmt-check
```

#### Security Issues
```bash
# Run security scan
make security

# Install security tools
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
```

### Environment Issues

#### Go Version
```bash
# Check Go version
go version

# Should be 1.21 or later
```

#### Dependencies
```bash
# Clean module cache
go clean -modcache

# Reinstall dependencies
make deps
```

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Go Modules](https://golang.org/ref/mod)
- [Go Testing](https://golang.org/pkg/testing/)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Pre-commit](https://pre-commit.com/)

---

## 🎯 Next Steps

1. **Set up your development environment** using the Makefile
2. **Run the validation suite** to ensure everything works
3. **Explore the examples** to understand usage patterns
4. **Review the contributing guidelines** if you plan to contribute
5. **Set up pre-commit hooks** for automatic code quality checks

For questions or issues, please refer to the support channels listed above.

Happy coding! 🚀
