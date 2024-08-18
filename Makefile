# Name of the binary
BINARY_NAME=ayo

# Output directory for the build
BUILD_DIR=build

# Coverage file name
COVERAGE_FILE=$(BUILD_DIR)/coverage.out

# Default target
.PHONY: all
all: build

# Build the Go project
.PHONY: build
build: ## Build the Go project
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ayo

# Run the Go project
.PHONY: run
run: build ## Build and run the Go project
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean up the build artifacts
.PHONY: clean
clean: ## Clean up build artifacts
	@if [ -d $(BUILD_DIR) ]; then rm -r $(BUILD_DIR); fi

# Run tests with coverage
.PHONY: test
test: ## Run tests with coverage
	@mkdir -p $(BUILD_DIR)
	go test ./... -coverprofile=$(COVERAGE_FILE)
	go tool cover -func=$(COVERAGE_FILE)

# Format the code
.PHONY: fmt
fmt: ## Format the code
	go fmt ./...

# Lint the code
.PHONY: lint
lint: ## Lint the code
	golangci-lint run

# Display help
.PHONY: help
help: ## Show help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

