.PHONY: build build-all test test-cover test-race vet lint clean install uninstall completions dev help all

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -ldflags "-s -w \
	-X github.com/Gu1llaum-3/tinymonitor/cmd.Version=$(VERSION) \
	-X github.com/Gu1llaum-3/tinymonitor/cmd.Commit=$(COMMIT) \
	-X github.com/Gu1llaum-3/tinymonitor/cmd.BuildDate=$(DATE)"

# Binary name
BINARY := tinymonitor

# Default target
all: vet test build

# Build for current platform
build:
	go build $(LDFLAGS) -o $(BINARY) .

# Build for all platforms (Linux + macOS)
build-all: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-x86_64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-x86_64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 .

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run tests with race detector
test-race:
	go test -race -v ./...

# Run go vet
vet:
	go vet ./...

# Lint
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install with:" && \
		echo "  brew install golangci-lint  (macOS)" && \
		echo "  or visit: https://golangci-lint.run/usage/install/" && \
		exit 1)
	golangci-lint run

# Clean build artifacts
clean:
	rm -f $(BINARY)
	rm -rf dist/
	rm -f coverage.out coverage.html

# Install to /usr/local/bin
install: build
	sudo cp $(BINARY) /usr/local/bin/$(BINARY)
	sudo chmod +x /usr/local/bin/$(BINARY)

# Uninstall
uninstall:
	sudo rm -f /usr/local/bin/$(BINARY)

# Generate shell completions
completions: build
	mkdir -p completions
	./$(BINARY) completion bash > completions/$(BINARY).bash
	./$(BINARY) completion zsh > completions/_$(BINARY)
	./$(BINARY) completion fish > completions/$(BINARY).fish

# Development: build and run
dev: build
	./$(BINARY)

# Show help
help:
	@echo "Available targets:"
	@echo "  build       - Build for current platform"
	@echo "  build-all   - Build for all platforms (Linux + macOS)"
	@echo "  test        - Run tests"
	@echo "  test-cover  - Run tests with coverage report"
	@echo "  test-race   - Run tests with race detector"
	@echo "  vet         - Run go vet"
	@echo "  lint        - Run golangci-lint"
	@echo "  clean       - Clean build artifacts"
	@echo "  install     - Install to /usr/local/bin"
	@echo "  uninstall   - Remove from /usr/local/bin"
	@echo "  completions - Generate shell completions"
	@echo "  dev         - Build and run"
	@echo "  help        - Show this help"
