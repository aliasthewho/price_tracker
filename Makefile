.PHONY: help lint test test-race cover build clean pre-commit setup-hooks

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=price-tracker

# Default target
help: ## Show this help message
	@echo "\n\033[1mEMMSA Price Tracker - Makefile Help\033[0m\n"
	@echo "\033[1mAvailable targets:\033[0m"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
	@echo ""

# --- Development ---

setup: tools ## Install all development dependencies

setup-hooks: ## Set up Git hooks
	@echo "Setting up Git hooks..."
	@chmod +x scripts/setup-hooks.sh
	@./scripts/setup-hooks.sh

pre-commit: lint test-race ## Run checks before committing
	@echo "\n\033[32m✓ All pre-commit checks passed!\033[0m\n"

# Install development dependencies
tools: ## Install development tools
	@echo "Installing development tools..."
	# Install golangci-lint
	if ! command -v golangci-lint &> /dev/null; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.0; \
	fi
	# Install pre-commit
	if ! command -v pre-commit &> /dev/null; then \
		pip install pre-commit; \
	fi
	@echo "\033[32m✓ Development tools installed\033[0m"

# Run golangci-lint
lint: ## Run linters
	@echo "\n\033[1m🔍 Running linters...\033[0m"
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		$(MAKE) tools; \
	fi
	@if ! golangci-lint run --timeout 5m; then \
		echo "\n\033[31m✗ Linting failed. Fix the issues and try again.\033[0m"; \
		exit 1; \
	else \
		echo "\n\033[32m✓ Linting passed!\033[0m"; \
	fi

# Run tests
test: ## Run tests
	@echo "\n\033[1m🧪 Running tests...\033[0m"
	@$(GOTEST) -v -cover ./...

# Run tests with race detector
test-race: ## Run tests with race detector
	@echo "\n\033[1m🔍 Running tests with race detector...\033[0m"
	@$(GOTEST) -race -v ./...

# Run tests with coverage
cover: ## Run tests with coverage report
	@echo "\n\033[1m📊 Running tests with coverage...\033[0m"
	@$(GOTEST) -coverprofile=coverage.out -covermode=atomic ./...
	@echo "\n\033[1m📈 Coverage Report:\033[0m"
	@go tool cover -func=coverage.out | grep total:
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "\n✅ Coverage report generated: \033[4mcoverage.html\033[0m\n"

# Build the application
build: ## Build the application
	@echo "\n\033[1m🔨 Building $(BINARY_NAME)...\033[0m"
	@mkdir -p bin
	@$(GOBUILD) -o bin/$(BINARY_NAME) ./cmd/price-tracker
	@echo "\n✅ Build successful: bin/$(BINARY_NAME)\n"

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "\n\033[1m🧹 Cleaning...\033[0m"
	@$(GOCLEAN)
	@rm -rf bin/ coverage.out coverage.html
	@echo "✅ Clean complete\n"

.DEFAULT_GOAL := help
