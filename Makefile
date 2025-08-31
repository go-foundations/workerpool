# Worker Pool Package Makefile
# Provides targets for testing, validation, building, and CI pipeline support

# Variables
BINARY_NAME=workerpool
MAIN_PATH=./examples
BENCHMARK_PATH=./benchmarks
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html
LINT_TIMEOUT=5m

# Go commands
GO=go
GOFMT=gofmt
GOIMPORTS=goimports
GOLINT=golangci-lint
GOSEC=gosec

# Default target
.PHONY: all
all: clean test build

# Development targets
.PHONY: dev
dev: fmt lint test build

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BINARY_NAME)
	@rm -rf $(COVERAGE_FILE)
	@rm -rf $(COVERAGE_HTML)
	@rm -rf examples/*/$(BINARY_NAME)
	@go clean -cache -testcache -modcache

# Format code
.PHONY: fmt
fmt:
	@echo "🎨 Formatting code..."
	@$(GOFMT) -s -w .
	@$(GO) mod tidy

# Check formatting
.PHONY: fmt-check
fmt-check:
	@echo "🔍 Checking code formatting..."
	@if [ "$$($(GOFMT) -l . | wc -l)" -gt 0 ]; then \
		echo "❌ Code is not formatted. Run 'make fmt' to fix."; \
		$(GOFMT) -l .; \
		exit 1; \
	fi
	@echo "✅ Code formatting is correct"

# Lint code
.PHONY: lint
lint:
	@echo "🔍 Linting code..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run --timeout=$(LINT_TIMEOUT); \
	else \
		echo "⚠️  golangci-lint not found, installing..." \
		&& curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2 \
		&& $(GOLINT) run --timeout=$(LINT_TIMEOUT); \
	fi

# Security scan
.PHONY: security
security:
	@echo "🔒 Running security scan..."
	@if command -v $(GOSEC) >/dev/null 2>&1; then \
		$(GOSEC) ./... || echo "⚠️  Security scan completed with warnings"; \
	else \
		echo "⚠️  gosec not found, installing..." \
		&& go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest \
		&& $(GOSEC) ./... || echo "⚠️  Security scan completed with warnings"; \
	fi

# Run go vet
.PHONY: vet
vet:
	@echo "🔍 Running go vet..."
	@$(GO) vet ./...

# Run all tests
.PHONY: test
test: vet
	@echo "🧪 Running all tests..."
	@$(GO) test -v -race -coverprofile=$(COVERAGE_FILE) ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage: test
	@echo "📊 Generating coverage report..."
	@$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "📈 Coverage report generated: $(COVERAGE_HTML)"
	@$(GO) tool cover -func=$(COVERAGE_FILE)

# Run benchmarks
.PHONY: benchmark
benchmark:
	@echo "⚡ Running benchmarks..."
	@$(GO) test -bench=. -benchmem ./benchmarks

# Run specific test
.PHONY: test-specific
test-specific:
	@echo "🧪 Running specific test: $(TEST)"
	@$(GO) test -v -run $(TEST) .

# Build examples
.PHONY: build-examples
build-examples:
	@echo "🔨 Building examples..."
	@$(GO) build -o examples/http_example/$(BINARY_NAME) ./examples/http_example
	@$(GO) build -o examples/string_example/$(BINARY_NAME) ./examples/string_example
	@echo "✅ Examples built successfully"

# Build all
.PHONY: build
build: build-examples
	@echo "🔨 Building package..."
	@$(GO) build ./...
	@echo "✅ Package built successfully"

# Install dependencies
.PHONY: deps
deps:
	@echo "📦 Installing dependencies..."
	@$(GO) mod download
	@$(GO) mod verify

# Update dependencies
.PHONY: deps-update
deps-update:
	@echo "📦 Updating dependencies..."
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@$(GO) mod verify

# Validate package
.PHONY: validate
validate: fmt-check security vet test build
	@echo "✅ Package validation completed successfully!"

# CI pipeline target
.PHONY: ci
ci: clean deps validate
	@echo "🚀 CI pipeline completed successfully!"

# Pre-commit hook
.PHONY: pre-commit
pre-commit: fmt lint test
	@echo "✅ Pre-commit checks passed!"

# Performance testing
.PHONY: perf
perf: benchmark
	@echo "📊 Performance testing completed"

# Documentation
.PHONY: docs
docs:
	@echo "📚 Generating documentation..."
	@$(GO) doc -all ./...
	@echo "📖 Documentation generated"

# Help target
.PHONY: help
help:
	@echo "Worker Pool Package - Available Targets:"
	@echo ""
	@echo "🔧 Development:"
	@echo "  dev          - Run development workflow (fmt, lint, test, build)"
	@echo "  fmt          - Format code"
	@echo "  fmt-check    - Check code formatting"
	@echo "  lint         - Run linter"
	@echo "  security     - Run security scan"
	@echo "  vet          - Run go vet"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test         - Run all tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  benchmark    - Run benchmarks"
	@echo "  test-specific- Run specific test (TEST=TestName)"
	@echo ""
	@echo "🔨 Building:"
	@echo "  build        - Build package and examples"
	@echo "  build-examples- Build only examples"
	@echo ""
	@echo "📦 Dependencies:"
	@echo "  deps         - Install dependencies"
	@echo "  deps-update  - Update dependencies"
	@echo ""
	@echo "🚀 CI/CD:"
	@echo "  ci           - Run full CI pipeline"
	@echo "  validate     - Run full validation"
	@echo "  pre-commit   - Run pre-commit checks"
	@echo ""
	@echo "📊 Analysis:"
	@echo "  perf         - Run performance tests"
	@echo "  docs         - Generate documentation"
	@echo ""
	@echo "🧹 Maintenance:"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help message"

# Default target
.DEFAULT_GOAL := help
