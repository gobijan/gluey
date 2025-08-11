# Gluey Makefile
# https://github.com/gobijan/gluey

# Variables
BINARY_NAME := gluey
MAIN_PATH := ./cmd/gluey
BUILD_DIR := ./build
COVERAGE_DIR := ./coverage
EXAMPLE_DIR := ./example_app

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOVET := $(GOCMD) vet
GOLINT := golangci-lint

# Build variables
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

# Phony targets
.PHONY: all help build install clean test test-verbose test-short test-race test-coverage \
        test-coverage-html test-unit test-integration bench lint fmt vet check deps tidy \
        update-deps tools example-app demo docs serve-docs release-dry release-snapshot ci

## help: Display this help message
help:
	@echo "$(BLUE)Gluey Makefile$(NC)"
	@echo "$(BLUE)===============$(NC)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@grep -E '^## ' Makefile | sed 's/## /  /' | column -t -s ':' | sed 's/^/  /'
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make build       # Build the gluey CLI"
	@echo "  make test        # Run all tests"
	@echo "  make check       # Run all quality checks"
	@echo "  make example-app # Create and test an example app"

## all: Run tests and build
all: test build

## build: Build the gluey CLI tool
build:
	@echo "$(BLUE)Building gluey...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## install: Install gluey to $$GOPATH/bin
install:
	@echo "$(BLUE)Installing gluey...$(NC)"
	@$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Installed to $(GOPATH)/bin/$(BINARY_NAME)$(NC)"

## clean: Clean build artifacts and test cache
clean:
	@echo "$(BLUE)Cleaning...$(NC)"
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR) $(COVERAGE_DIR) $(EXAMPLE_DIR)
	@$(GOTEST) -cache clean
	@echo "$(GREEN)✓ Clean complete$(NC)"

## test: Run all tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	@$(GOTEST) -v ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(BLUE)Running tests (verbose)...$(NC)"
	@$(GOTEST) -v -count=1 ./...

## test-short: Run short tests only
test-short:
	@echo "$(BLUE)Running short tests...$(NC)"
	@$(GOTEST) -short ./...
	@echo "$(GREEN)✓ Short tests passed$(NC)"

## test-race: Run tests with race detector
test-race:
	@echo "$(BLUE)Running tests with race detector...$(NC)"
	@$(GOTEST) -race ./...
	@echo "$(GREEN)✓ No race conditions detected$(NC)"

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.out
	@echo "$(GREEN)✓ Coverage report generated$(NC)"

## test-coverage-html: Generate HTML coverage report
test-coverage-html: test-coverage
	@echo "$(BLUE)Generating HTML coverage report...$(NC)"
	@$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "$(GREEN)✓ HTML report: $(COVERAGE_DIR)/coverage.html$(NC)"
	@echo "$(YELLOW)Opening in browser...$(NC)"
	@open $(COVERAGE_DIR)/coverage.html 2>/dev/null || xdg-open $(COVERAGE_DIR)/coverage.html 2>/dev/null || echo "Please open $(COVERAGE_DIR)/coverage.html manually"

## test-unit: Run unit tests only
test-unit:
	@echo "$(BLUE)Running unit tests...$(NC)"
	@$(GOTEST) -v -tags=unit ./...
	@echo "$(GREEN)✓ Unit tests passed$(NC)"

## test-integration: Run integration tests
test-integration:
	@echo "$(BLUE)Running integration tests...$(NC)"
	@$(GOTEST) -v -tags=integration ./...
	@echo "$(GREEN)✓ Integration tests passed$(NC)"

## bench: Run benchmarks
bench:
	@echo "$(BLUE)Running benchmarks...$(NC)"
	@$(GOTEST) -bench=. -benchmem ./...

## lint: Run golangci-lint
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run --no-config ./...; \
		echo "$(GREEN)✓ Linting complete$(NC)"; \
	else \
		echo "$(RED)✗ golangci-lint not installed. Run 'make tools' to install$(NC)"; \
		exit 1; \
	fi

## fmt: Format code with gofmt
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@$(GOFMT) -s -w .
	@echo "$(GREEN)✓ Code formatted$(NC)"

## vet: Run go vet
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	@$(GOVET) ./...
	@echo "$(GREEN)✓ Vet complete$(NC)"

## check: Run all quality checks (fmt, vet, lint)
check: fmt vet lint
	@echo "$(GREEN)✓ All quality checks passed$(NC)"

## deps: Download dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@$(GOMOD) download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

## tidy: Run go mod tidy
tidy:
	@echo "$(BLUE)Tidying dependencies...$(NC)"
	@$(GOMOD) tidy
	@echo "$(GREEN)✓ Dependencies tidied$(NC)"

## update-deps: Update all dependencies
update-deps:
	@echo "$(BLUE)Updating dependencies...$(NC)"
	@$(GOGET) -u ./...
	@$(GOMOD) tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

## tools: Install development tools
tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	@if ! command -v $(GOLINT) >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin; \
	fi
	@echo "$(GREEN)✓ Tools installed$(NC)"

## example-app: Create and run an example application
example-app: build
	@echo "$(BLUE)Creating example application...$(NC)"
	@rm -rf $(EXAMPLE_DIR)
	@$(BUILD_DIR)/$(BINARY_NAME) new example_app
	@cd $(EXAMPLE_DIR) && \
		$(GOMOD) init example_app && \
		$(GOMOD) edit -replace github.com/gobijan/gluey=../ && \
		$(GOMOD) tidy && \
		echo "$(YELLOW)Generating interfaces...$(NC)" && \
		$(GOCMD) run -mod=mod ../$(BUILD_DIR)/$(BINARY_NAME) gen && \
		echo "$(YELLOW)Generating examples...$(NC)" && \
		$(GOCMD) run -mod=mod ../$(BUILD_DIR)/$(BINARY_NAME) example
	@echo "$(GREEN)✓ Example app created in $(EXAMPLE_DIR)$(NC)"
	@echo "$(YELLOW)To run the example:$(NC)"
	@echo "  cd $(EXAMPLE_DIR) && go run main.go"

## demo: Run a quick demo
demo: example-app
	@echo "$(BLUE)Running demo...$(NC)"
	@echo "$(YELLOW)Example app structure:$(NC)"
	@tree $(EXAMPLE_DIR) -I 'go.mod|go.sum' 2>/dev/null || find $(EXAMPLE_DIR) -type f -name "*.go" -o -name "*.html" | head -20
	@echo ""
	@echo "$(GREEN)✓ Demo complete$(NC)"
	@echo "$(YELLOW)Start the app with: cd $(EXAMPLE_DIR) && go run main.go$(NC)"

## docs: Generate documentation
docs:
	@echo "$(BLUE)Generating documentation...$(NC)"
	@$(GOCMD) doc -all > docs/API.md
	@echo "$(GREEN)✓ Documentation generated$(NC)"

## serve-docs: Serve documentation locally
serve-docs:
	@echo "$(BLUE)Serving documentation...$(NC)"
	@echo "$(YELLOW)Documentation server starting at http://localhost:6060$(NC)"
	@godoc -http=:6060

## release-dry: Dry run of release process
release-dry:
	@echo "$(BLUE)Dry run of release process...$(NC)"
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --skip-publish --rm-dist; \
		echo "$(GREEN)✓ Dry run complete$(NC)"; \
	else \
		echo "$(RED)✗ goreleaser not installed$(NC)"; \
		echo "Install with: brew install goreleaser/tap/goreleaser"; \
		exit 1; \
	fi

## release-snapshot: Create snapshot release
release-snapshot:
	@echo "$(BLUE)Creating snapshot release...$(NC)"
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --rm-dist; \
		echo "$(GREEN)✓ Snapshot created in ./dist$(NC)"; \
	else \
		echo "$(RED)✗ goreleaser not installed$(NC)"; \
		echo "Install with: brew install goreleaser/tap/goreleaser"; \
		exit 1; \
	fi

## ci: Run full CI pipeline locally
ci: clean deps check test-race test-coverage build
	@echo "$(GREEN)✓ CI pipeline complete$(NC)"

# Development shortcuts
.PHONY: t b c r

## t: Shortcut for test
t: test

## b: Shortcut for build
b: build

## c: Shortcut for check
c: check

## r: Shortcut for run (build and run example)
r: example-app
	@cd $(EXAMPLE_DIR) && $(GOCMD) run main.go