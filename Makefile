.PHONY: build test vet clean all release

BINARY_NAME=tinymonitor
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

all: vet test build

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BINARY_NAME) ./cmd/tinymonitor

test:
	go test -v ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*

release: clean
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BINARY_NAME)-linux-amd64 ./cmd/tinymonitor
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BINARY_NAME)-linux-arm64 ./cmd/tinymonitor
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BINARY_NAME)-darwin-amd64 ./cmd/tinymonitor
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BINARY_NAME)-darwin-arm64 ./cmd/tinymonitor
