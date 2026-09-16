# Makefile for the importable Go library.
# Targets verify packages; they do not run or package an application binary.

SHELL := /bin/bash
.DEFAULT_GOAL := help

.PHONY: all
all: test vet build

.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: build
build: ## Verify compilation of all packages
	go build -v ./...

.PHONY: test
test: ## Run unit and integration tests with race detection
	go test -race -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with race detection and HTML coverage report
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: vet
vet: ## Run go vet analysis
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (or go vet if golangci-lint not installed)
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run || (echo "golangci-lint not found in PATH; running go vet ./..." && go vet ./...)

.PHONY: tidy
tidy: ## Tidy and verify Go module dependencies
	go mod tidy
	go mod verify

.PHONY: clean
clean: ## Clean temporary test coverage and artifact files
	rm -rf bin tmp coverage.out coverage.html profile.out
	@echo "Cleaned build and test artifacts"

