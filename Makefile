# Define the name of your Go binary
BINARY_NAME := mold

# Define the main package path for your CLI application
# Adjust this path if your main package is located elsewhere, e.g., ./cmd/mycli
MAIN_PACKAGE := ./cmd/cli

# Find all Go files in the current directory and its subdirectories, excluding vendor
GO_FILES := $(shell find . -type f -name "*.go" ! -path "./vendor/*")

.PHONY: all lint fmt run build clean help

# Default target: runs lint and format
all: lint fmt

# Target to run golangci-lint for code quality checks
# Install golangci-lint if you haven't already: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.58.1
lint:
	@echo "Running golangci-lint..."
	@if ! command -v go tool golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Please install it: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi
	go tool golangci-lint run ./...

# Target to format Go source code using go fmt
fmt:
	@echo "Formatting Go files..."
	go tool golines -m 120 -w .
	go tool goimports-reviser -imports-order "std,company,project,general" ./...

# Target to run your Go CLI application directly
# Assumes your main function is in a file within the path defined by MAIN_PACKAGE
run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_PACKAGE)

# Target to build the Go CLI application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Target to clean up build artifacts
clean:
	@echo "Cleaning up build artifacts..."
	rm -f $(BINARY_NAME)

# Help target to display available commands
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all       Runs lint and format (default)"
	@echo "  lint      Runs golangci-lint for code quality checks"
	@echo "  fmt       Formats Go source code using go fmt"
	@echo "  run       Runs the Go CLI application"
	@echo "  build     Builds the Go CLI application executable"
	@echo "  clean     Removes the built executable"
	@echo "  help      Displays this help message"
