.PHONY: all build test test-coverage clean install lint fmt vet run help

# Variables
BINARY_NAME := mackerel-plugin-puma-v2
BINARY_PATH := bin/$(BINARY_NAME)
GO := go
GOFLAGS := -v
LDFLAGS := -ldflags "-s -w"
COVERAGE_FILE := coverage.out

# Default target
all: clean fmt vet lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_PATH) ./cmd/$(BINARY_NAME)
	@echo "Build complete: $(BINARY_PATH)"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test $(GOFLAGS) -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@echo "Coverage report generated: $(COVERAGE_FILE)"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GO) test $(GOFLAGS) -tags=integration ./test/integration/...

# Run E2E tests
test-e2e:
	@echo "Running E2E tests..."
	$(GO) test $(GOFLAGS) -tags=e2e ./test/e2e/...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/ $(COVERAGE_FILE) coverage.html
	@echo "Clean complete"

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install ./cmd/$(BINARY_NAME)
	@echo "Installation complete"

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Run the plugin (for development)
run: build
	$(BINARY_PATH) $(ARGS)

# Run with example stats endpoint
run-example: build
	$(BINARY_PATH) -socket=/tmp/puma.sock

# Show help
help:
	@echo "Available targets:"
	@echo "  all              - Run fmt, vet, lint, test, and build"
	@echo "  build            - Build the binary"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  test-integration - Run integration tests"
	@echo "  test-e2e         - Run E2E tests"
	@echo "  clean            - Remove build artifacts"
	@echo "  install          - Install the binary"
	@echo "  lint             - Run linter"
	@echo "  fmt              - Format code"
	@echo "  vet              - Run go vet"
	@echo "  run              - Run the plugin (use ARGS for arguments)"
	@echo "  run-example      - Run with example socket"
	@echo "  help             - Show this help message"