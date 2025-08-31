# Contributing to Worker Pool

Thank you for your interest in contributing to the Worker Pool package! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for using the Makefile)

### Setting up the development environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/workerpool.git
   cd workerpool
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/go-foundations/workerpool.git
   ```
4. Install dependencies:
   ```bash
   make deps
   # or manually:
   go mod download
   ```

## 🧪 Development Workflow

### Using the Makefile

The project includes a comprehensive Makefile for common development tasks:

```bash
# Run the full development workflow
make dev

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make benchmark

# Format code
make fmt

# Check formatting
make fmt-check

# Run linter
make lint

# Security scan
make security

# Build examples
make build

# Full validation
make validate
```

### Manual Commands

If you prefer not to use the Makefile:

```bash
# Format code
go fmt ./...

# Run tests
go test -v -race ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
go test -bench=. -benchmem ./benchmarks

# Build examples
go build ./examples/...
```

## 📝 Code Style and Standards

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for code formatting
- Run `go vet` before committing
- Ensure all tests pass with `-race` flag

### Naming Conventions

- Use descriptive names for variables, functions, and types
- Follow Go naming conventions (e.g., `WorkerPool`, `NewWithConfig`)
- Use camelCase for exported identifiers
- Use snake_case for file names

### Documentation

- All exported types, functions, and methods must have GoDoc comments
- Include examples in documentation where appropriate
- Update README.md for new features or breaking changes

## 🧪 Testing

### Test Requirements

- All new code must include tests
- Tests should cover both success and error cases
- Use table-driven tests for multiple test cases
- Include benchmarks for performance-critical code
- Run tests with race detection enabled

### Running Tests

```bash
# Run all tests
make test

# Run specific test
make test-specific TEST=TestWorkerPool

# Run tests with coverage
make test-coverage

# Run benchmarks
make benchmark
```

## 🔒 Security

### Security Guidelines

- Never commit sensitive information (API keys, passwords, etc.)
- Run security scans before submitting PRs: `make security`
- Follow secure coding practices
- Report security vulnerabilities privately

## 📦 Dependencies

### Adding Dependencies

- Only add dependencies that are absolutely necessary
- Use Go modules for dependency management
- Ensure dependencies are actively maintained
- Document why a dependency is needed

### Updating Dependencies

```bash
# Update all dependencies
make deps-update

# Update specific dependency
go get -u github.com/example/package
go mod tidy
```

## 🔄 Pull Request Process

### Before Submitting a PR

1. Ensure all tests pass: `make test`
2. Run the full validation: `make validate`
3. Update documentation if needed
4. Add or update tests for new functionality
5. Ensure code follows the project's style guidelines

### PR Guidelines

- Use descriptive PR titles
- Include a detailed description of changes
- Reference related issues
- Include test results and coverage information
- Ensure CI checks pass

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] All tests pass
- [ ] New tests added
- [ ] Coverage maintained/improved

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
```

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

## 📞 Getting Help

- Open an issue for bugs or feature requests
- Use GitHub Discussions for questions
- Check existing issues and PRs
- Review the documentation

## 📄 License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing to the Worker Pool package! 🎉
