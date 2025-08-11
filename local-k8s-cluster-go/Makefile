# Makefile for k3s-manager

# Variables
BINARY_NAME=k3s-manager
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

# Default target
.DEFAULT_GOAL := help

# Build targets
.PHONY: build
build: ## Build the binary for current platform
	go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd

.PHONY: build-all
build-all: ## Build binaries for all platforms
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 ./cmd
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-arm64 ./cmd
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 ./cmd
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-arm64 ./cmd
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe ./cmd
	@echo "All builds completed!"

.PHONY: install
install: ## Install the binary to $GOPATH/bin
	go install ${LDFLAGS} ./cmd

# Development targets
.PHONY: run
run: ## Run the application
	go run ./cmd

.PHONY: test
test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

.PHONY: test-coverage
test-coverage: test ## Run tests and show coverage
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run --timeout=5m

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	goimports -w .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

# Dependency management
.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod verify

.PHONY: deps-update
deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

.PHONY: deps-clean
deps-clean: ## Clean module cache
	go clean -modcache

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t ${BINARY_NAME}:${VERSION} .
	docker tag ${BINARY_NAME}:${VERSION} ${BINARY_NAME}:latest

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run --rm -it ${BINARY_NAME}:latest --help

# Release targets
.PHONY: release-dry
release-dry: ## Dry run of release process
	@echo "Would create release for version: ${VERSION}"
	@echo "Build time: ${BUILD_TIME}"

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)