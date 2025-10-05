# Variables
BINARY_NAME=tmux-session-launcher
OUT_DIR=./out
MAIN_PATH=./main.go

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(OUT_DIR)

# Run the program
.PHONY: run
run: build
	$(OUT_DIR)/$(BINARY_NAME)

# Run tests
.PHONY: test
test:
	go test ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -cover ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Vet code
.PHONY: vet
vet:
	go vet ./...

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Build for multiple platforms
.PHONY: build-all
build-all:
	@mkdir -p $(OUT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(OUT_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -o $(OUT_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -o $(OUT_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary to $(OUT_DIR)"
	@echo "  clean        - Remove build artifacts"
	@echo "  run          - Build and run the program"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  help         - Show this help"
