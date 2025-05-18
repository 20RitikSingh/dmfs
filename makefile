# Makefile for dmfs - Production Ready

# Variables
BINARY_NAME=dmfs
BIN_DIR=bin
SRC=main.go

.PHONY: all build run test clean help

all: build

build:
	@echo "[INFO] Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(SRC)
	@echo "[INFO] Build complete."

run: build
	@echo "[INFO] Running $(BINARY_NAME)..."
	@./$(BIN_DIR)/$(BINARY_NAME)

test:
	@echo "[INFO] Running tests..."
	@go test -v ./...

clean:
	@echo "[INFO] Cleaning up..."
	@rm -f $(BIN_DIR)/$(BINARY_NAME)
	@echo "[INFO] Clean complete."

help:
	@echo "Available targets:"
	@echo "  build   - Build the project binary."
	@echo "  run     - Build and run the project."
	@echo "  test    - Run all tests."
	@echo "  clean   - Remove built binaries."
	@echo "  help    - Show this help message."

