.PHONY: build install test clean tidy help

BINARY_NAME=evm-profiler

# Default target
all: build

## 🏗️ Build targets

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) main.go

install: ## Install binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME) to $(GOPATH)/bin..."
	go install .

tidy: ## Tidy up dependencies
	@echo "Tidying up go.mod..."
	go mod tidy

test: ## Run unit tests
	@echo "Running tests..."
	go test ./... -v

clean: ## Remove build artifacts
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
