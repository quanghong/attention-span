BINARY_NAME=attention-span
BUILD_DIR=bin
VERSION=1.0.0

.PHONY: all build build-mac test clean install run help

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/attention-span
	@echo "✓ Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

build-mac:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/attention-span
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/attention-span
	@echo "✓ macOS binaries built:"
	@echo "  • $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 (Apple Silicon M1/M2/M3/M4)"
	@echo "  • $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 (Intel Mac)"

test:
	go test -v ./...

install: build
	@if [ -w /usr/local/bin ]; then \
		install -m 755 $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME); \
		echo "✓ Installed to /usr/local/bin/$(BINARY_NAME)"; \
	else \
		mkdir -p $(HOME)/.local/bin; \
		install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME); \
		echo "✓ Installed to $(HOME)/.local/bin/$(BINARY_NAME)"; \
		echo "  Ensure $(HOME)/.local/bin is in your PATH"; \
	fi

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)
	@echo "✓ Cleaned $(BUILD_DIR)"

help:
	@echo "Attention Span Tracker Makefile targets:"
	@echo "  make build      - Compile binary for current platform"
	@echo "  make build-mac  - Cross-compile for macOS (arm64 & amd64)"
	@echo "  make test       - Run all unit tests"
	@echo "  make install    - Install binary to /usr/local/bin or ~/.local/bin"
	@echo "  make run        - Build and run attention-span"
	@echo "  make clean      - Clean build artifacts"
