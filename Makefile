# Variables
BINARY_NAME=goup
MAIN_PATH=cmd/goup
BUILD_DIR=build
INSTALL_DIR=$(HOME)/go/bin

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOCLEAN=$(GOCMD) clean

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-v $(LDFLAGS)

.PHONY: all build clean test coverage install uninstall run help deps lint mocks

# Default target
all: clean deps mocks test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf internal/mocks/*.go
	@echo "Clean completed"

# Generate mocks
mocks:
	@echo "Generating mocks..."
	@mkdir -p internal/mocks
	@if ! command -v mockgen >/dev/null 2>&1; then \
		echo "mockgen not found. Installing..."; \
		go install go.uber.org/mock/mockgen@latest; \
	fi
	@chmod +x ./scripts/generate_mocks.sh
	@PATH="$$PATH:$$HOME/go/bin:$$(go env GOPATH)/bin" ./scripts/generate_mocks.sh

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install the application
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installation completed: $(INSTALL_DIR)/$(BINARY_NAME)"

# Uninstall the application
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstallation completed"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@if which golangci-lint >/dev/null 2>&1; then \
		echo "Found golangci-lint, running..."; \
		golangci-lint run; \
	elif [ -f "$(shell go env GOPATH)/bin/golangci-lint" ]; then \
		echo "Found golangci-lint in GOPATH, running..."; \
		$(shell go env GOPATH)/bin/golangci-lint run; \
	elif [ -f "$(HOME)/go/bin/golangci-lint" ]; then \
		echo "Found golangci-lint in ~/go/bin, running..."; \
		$(HOME)/go/bin/golangci-lint run; \
	else \
		echo "❌ golangci-lint not found!"; \
		echo "📦 Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		echo "✅ Installation completed, running lint..."; \
		$(shell go env GOPATH)/bin/golangci-lint run; \
	fi

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(MAIN_PATH)
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(MAIN_PATH)
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(MAIN_PATH)
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(MAIN_PATH)
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(MAIN_PATH)
	
	@echo "Multi-platform build completed"

# Development workflow
dev: clean deps mocks test build
	@echo "Development build completed"

# Quick build without tests
quick:
	@echo "Quick build..."
	$(GOBUILD) -o $(BINARY_NAME) ./$(MAIN_PATH)

# Help
help:
	@echo "Available targets:"
	@echo "  all        - Run clean, deps, mocks, test, and build"
	@echo "  build      - Build the application"
	@echo "  clean      - Clean build artifacts and mocks"
	@echo "  test       - Run tests"
	@echo "  coverage   - Run tests with coverage report"
	@echo "  mocks      - Generate mock files"
	@echo "  install    - Install the application to ~/go/bin"
	@echo "  uninstall  - Remove the application from ~/go/bin"
	@echo "  run        - Build and run the application (use ARGS='--help' for options)"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  lint       - Lint code (requires golangci-lint)"
	@echo "  bench      - Run benchmarks"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  dev        - Development workflow (clean, deps, mocks, test, build)"
	@echo "  quick      - Quick build without tests"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make run ARGS='--help'"
	@echo "  make run ARGS='--select --dry-run'"
	@echo "  make test"
	@echo "  make coverage"
	@echo "  make mocks"