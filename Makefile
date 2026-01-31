.PHONY: build test vet clean all release

BINARY_NAME=tinymonitor
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-s -w \
	-X github.com/Gu1llaum-3/tinymonitor/cmd.Version=$(VERSION) \
	-X github.com/Gu1llaum-3/tinymonitor/cmd.Commit=$(COMMIT) \
	-X github.com/Gu1llaum-3/tinymonitor/cmd.BuildDate=$(BUILD_DATE)

all: vet test build

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

test:
	go test -v ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*

release: clean
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-arm64 .
