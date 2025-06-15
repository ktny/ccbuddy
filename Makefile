.PHONY: build test lint fmt clean install dev

# Build the binary
build:
	go build -o ccbuddy ./cmd/ccbuddy

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint the code
lint:
	golangci-lint run

# Format the code
fmt:
	go fmt ./...

# Clean build artifacts
clean:
	rm -f ccbuddy coverage.out coverage.html
	go clean

# Install the binary
install:
	go install ./cmd/ccbuddy

# Development mode (run with file watching)
dev: build
	./ccbuddy

# Check if dependencies are up to date
deps:
	go mod tidy
	go mod verify

# Initialize development environment
init-dev:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest