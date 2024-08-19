# Name of the binary
BINARY_NAME=ayo

# Output directory for the build
BUILD_DIR=dist

# Coverage file name
COVERAGE_FILE=$(BUILD_DIR)/coverage.out

project_version := $(shell git describe --tags --always)
git_commit := $(shell git rev-parse --verify HEAD)
date := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
author := $(USER)

LDFLAGS := -s -w
LDFLAGS += -X 'github.com/banst/ayo/cmd/ayo/config.Version=$(project_version)'
LDFLAGS += -X 'github.com/banst/ayo/cmd/ayo/config.Commit=$(git_commit)'
LDFLAGS += -X 'github.com/banst/ayo/cmd/ayo/config.Date=$(date)'
LDFLAGS += -X 'github.com/banst/ayo/cmd/ayo/config.BuiltBy=$(author)'

# Default target
.PHONY: all
all: build

# Build the Go project
.PHONY: build
build-go: ## Build the Go project
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ayo

# Build the Documentation
.PHONY: build-docs
build-docs: ## Build the Documentation
	mdbook build -d ../$(BUILD_DIR)/docs docs/

# Build the project
.PHONY: build
build: build-go build-docs ## Build the project

# Run the Go project
.PHONY: run
run: build-go ## Build and run the Go project
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

