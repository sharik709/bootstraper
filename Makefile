.PHONY: build test lint clean

GOFLAGS := -ldflags="-s -w"
VERSION := $(shell grep -E "Version = " cmd/root.go | cut -d '"' -f 2)

all: build

build:
	@echo "Building bootstraper $(VERSION)..."
	@go build $(GOFLAGS) -o dist/bt

clean:
	@echo "Cleaning..."
	@rm -rf dist
	@go clean

test:
	@echo "Running tests..."
	@go test -v ./...

coverage:
	@echo "Generating test coverage..."
	@go test -cover -coverprofile=coverage.txt ./...
	@go tool cover -html=coverage.txt

lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

install: build
	@echo "Installing bootstraper to $$GOPATH/bin/bt"
	@cp dist/bt $$GOPATH/bin/bt

release:
	@echo "Building release artifacts..."
	@mkdir -p dist/release
	@GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o dist/release/bt-darwin-amd64
	@GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o dist/release/bt-darwin-arm64
	@GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o dist/release/bt-linux-amd64
	@GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o dist/release/bt-linux-arm64
	@GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o dist/release/bt-windows-amd64.exe

help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make test       - Run tests"
	@echo "  make coverage   - Generate test coverage"
	@echo "  make lint       - Run linter"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make install    - Install to GOPATH/bin"
	@echo "  make release    - Build release artifacts"
