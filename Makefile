.PHONY: help build test clean install examples

help:
	@echo "Available commands:"
	@echo "  make install   - Install dependencies"
	@echo "  make build     - Build the project"
	@echo "  make test      - Run tests"
	@echo "  make examples  - Run example programs"
	@echo "  make clean     - Clean build artifacts"

install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

build:
	@echo "Building..."
	go build ./...

test:
	@echo "Running tests..."
	go test -v ./tests/...

examples:
	@echo "Building examples..."
	go build -o bin/basic examples/basic/main.go
	go build -o bin/advanced examples/advanced/main.go
	go build -o bin/auth examples/auth/main.go
	go build -o bin/realtime examples/realtime/main.go

clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"
