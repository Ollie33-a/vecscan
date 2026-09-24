# VecScan Makefile
# Build: make all
# Install: sudo make install

BINARY_NAME=vecscan
VERSION=1.0.0
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: all build clean test install uninstall deps

# Default target: download deps and build
all: deps build

# Build depends on deps to ensure go.sum exists
build: deps
	@echo "Building VecScan..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/vecscan
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download
	@echo "Dependencies ready"

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f go.sum

test:
	$(GOTEST) -v ./...

install: build
	@echo "Installing VecScan to $(INSTALL_DIR)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "VecScan installed successfully!"
	@echo "Run 'vecscan --help' to get started."

uninstall:
	@echo "Removing VecScan..."
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "VecScan uninstalled."

# Cross compilation
build-linux:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/vecscan

build-arm:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/vecscan

# Development helpers
run: deps
	$(GOCMD) run ./cmd/vecscan

dev:
	air -c .air.toml

# Quick rebuild without checking deps (faster for development)
rebuild:
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/vecscan