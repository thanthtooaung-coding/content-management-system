.PHONY: all build run test clean fmt vet lint tidy docker-build docker-run help install-tools check security bench profile dev watch

## Application Configuration
APP_NAME ?= server
DOCKER_IMAGE ?= $(APP_NAME):latest
PORT ?= 8080
GO_VERSION ?= 1.24
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

## Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

## Default target
all: install-tools fmt vet lint tidy test build

## Build targets
## build: Build the application binary
build:
	@echo "$(CYAN)🔨 Building application...$(RESET)"
	@mkdir -p bin/
	@go build -ldflags "-X main.version=$(shell git describe --tags --always --dirty)" -o bin/$(APP_NAME) ./cmd
	@echo "$(GREEN)✅ Build complete: bin/$(APP_NAME)$(RESET)"

## run: Run the application
run:
	@echo "$(CYAN)🚀 Running application...$(RESET)"
	@go run ./cmd

## dev: Run application with live reload (requires air)
dev:
	@echo "$(CYAN)🔥 Starting development server with live reload...$(RESET)"
	@air -c .air.toml || go run ./cmd

## watch: Watch for file changes and run tests
watch:
	@echo "$(CYAN)👀 Watching for changes...$(RESET)"
	@fswatch -o . | xargs -n1 -I{} make test

## Test targets
## test: Run all tests
test:
	@echo "$(CYAN)🧪 Running tests...$(RESET)"
	@go test -race -v ./...

## test-short: Run tests with -short flag
test-short:
	@echo "$(CYAN)🧪 Running short tests...$(RESET)"
	@go test -short -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(CYAN)🧪 Running tests with coverage...$(RESET)"
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(RESET)"

## bench: Run benchmarks
bench:
	@echo "$(CYAN)📊 Running benchmarks...$(RESET)"
	@go test -bench=. -benchmem ./...

## Code quality targets
## fmt: Format Go code
fmt:
	@echo "$(CYAN)📝 Formatting code...$(RESET)"
	@gofmt -s -w $(GO_FILES)
	@goimports -w $(GO_FILES)
	@echo "$(GREEN)✅ Code formatted$(RESET)"

## vet: Run go vet
vet:
	@echo "$(CYAN)🔍 Running go vet...$(RESET)"
	@go vet ./...
	@echo "$(GREEN)✅ Vet passed$(RESET)"

## lint: Run linters
lint:
	@echo "$(CYAN)🔍 Running linters...$(RESET)"
	@golangci-lint run --timeout=5m
	@echo "$(GREEN)✅ Linting passed$(RESET)"

## lint-fix: Run linters with auto-fix
lint-fix:
	@echo "$(CYAN)🔧 Running linters with auto-fix...$(RESET)"
	@golangci-lint run --fix --timeout=5m

## security: Run security checks
security:
	@echo "$(CYAN)🔒 Running security checks...$(RESET)"
	@gosec ./...
	@nancy sleuth
	@echo "$(GREEN)✅ Security checks passed$(RESET)"

## check: Run all checks (fmt, vet, lint, security, test)
check: fmt vet lint security test
	@echo "$(GREEN)✅ All checks passed$(RESET)"

## Dependencies
## tidy: Tidy go modules
tidy:
	@echo "$(CYAN)📦 Tidying go modules...$(RESET)"
	@go mod tidy
	@go mod verify
	@echo "$(GREEN)✅ Modules tidied$(RESET)"

## deps-update: Update all dependencies
deps-update:
	@echo "$(CYAN)📦 Updating dependencies...$(RESET)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies updated$(RESET)"

## deps-download: Download dependencies
deps-download:
	@echo "$(CYAN)📦 Downloading dependencies...$(RESET)"
	@go mod download
	@echo "$(GREEN)✅ Dependencies downloaded$(RESET)"

## Docker targets
## docker-build: Build Docker image
docker-build:
	@echo "$(CYAN)🐳 Building Docker image...$(RESET)"
	@docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)✅ Docker image built: $(DOCKER_IMAGE)$(RESET)"

## docker-run: Run Docker container
docker-run:
	@echo "$(CYAN)🐳 Running Docker container...$(RESET)"
	@docker run -p $(PORT):$(PORT) --rm --name $(APP_NAME) $(DOCKER_IMAGE)

## docker-run-bg: Run Docker container in background
docker-run-bg:
	@echo "$(CYAN)🐳 Running Docker container in background...$(RESET)"
	@docker run -d -p $(PORT):$(PORT) --name $(APP_NAME) $(DOCKER_IMAGE)

## docker-stop: Stop Docker container
docker-stop:
	@echo "$(CYAN)🐳 Stopping Docker container...$(RESET)"
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true

## docker-push: Push Docker image to registry
docker-push: docker-build
	@echo "$(CYAN)🐳 Pushing Docker image...$(RESET)"
	@docker push $(DOCKER_IMAGE)

## Tool installation
## install-tools: Install required development tools
install-tools:
	@echo "$(CYAN)🛠️  Installing development tools...$(RESET)"
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install github.com/sonatypeoss/nancy@latest
	@go install github.com/cosmtrek/air@latest
	@echo "$(GREEN)✅ Tools installed$(RESET)"

## Profiling targets
## profile-cpu: Run CPU profiling
profile-cpu:
	@echo "$(CYAN)📊 Running CPU profiling...$(RESET)"
	@go test -cpuprofile=cpu.prof -bench=. ./...
	@go tool pprof cpu.prof

## profile-mem: Run memory profiling
profile-mem:
	@echo "$(CYAN)📊 Running memory profiling...$(RESET)"
	@go test -memprofile=mem.prof -bench=. ./...
	@go tool pprof mem.prof

## Cleanup targets
## clean: Clean build artifacts and caches
clean:
	@echo "$(CYAN)🧹 Cleaning...$(RESET)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f cpu.prof mem.prof
	@go clean -cache -testcache -modcache
	@docker system prune -f
	@echo "$(GREEN)✅ Cleanup complete$(RESET)"

## clean-docker: Clean Docker images and containers
clean-docker:
	@echo "$(CYAN)🧹 Cleaning Docker...$(RESET)"
	@docker stop $(APP_NAME) 2>/dev/null || true
	@docker rm $(APP_NAME) 2>/dev/null || true
	@docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
	@docker system prune -f

## Info targets
## version: Show Go version and project info
version:
	@echo "$(CYAN)ℹ️  Version Information$(RESET)"
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse --short HEAD)"
	@echo "Git branch: $(shell git branch --show-current)"
	@echo "Build time: $(shell date)"

## env: Show environment variables
env:
	@echo "$(CYAN)ℹ️  Environment Variables$(RESET)"
	@echo "APP_NAME: $(APP_NAME)"
	@echo "DOCKER_IMAGE: $(DOCKER_IMAGE)"
	@echo "PORT: $(PORT)"
	@echo "GO_VERSION: $(GO_VERSION)"
	@echo "GOPATH: $(GOPATH)"
	@echo "GOROOT: $(GOROOT)"

## help: Show this help message
help:
	@echo "$(CYAN)Available targets:$(RESET)"
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/##//' | sort
	@echo ""
	@echo "$(YELLOW)Example usage:$(RESET)"
	@echo "  make all          # Run full build pipeline"
	@echo "  make dev          # Start development server"
	@echo "  make check        # Run all checks"
	@echo "  make docker-build # Build Docker image"