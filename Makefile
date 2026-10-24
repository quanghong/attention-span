BINARY_NAME=attention-span
BUILD_DIR=bin
INSTALL_DIR="${HOME}/.local/bin"
VERSION=1.0.0

.PHONY: all install test run help

all: help

install: clean
	@./install.sh

test:
	go test -v ./...

run:
	$(INSTALL_DIR)/$(BINARY_NAME)

clean:
	go mod tidy
	rm -rf $(BUILD_DIR)
	@echo "✓ Cleaned $(BUILD_DIR)"

help:
	@echo "Attention Span Tracker Makefile targets:"
	@echo "  make install       - Install binary for Linux AMD64 and Darwin ARM64 platforms"
	@echo "  make test          - Run all unit tests"
	@echo "  make run           - Run attention-span"
