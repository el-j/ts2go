.PHONY: build build-cli build-desktop test clean install help

# Detect OS for binary name
ifeq ($(OS),Windows_NT)
    BINARY_NAME := ts2go.exe
    RM := del /Q
else
    BINARY_NAME := ts2go
    RM := rm -f
endif

# Default target
all: build

# Build CLI tool
build: build-cli

build-cli:
	go build -o $(BINARY_NAME) ./cmd/ts2go

# Build desktop app
build-desktop: build-cli
	@echo "Copying CLI binary to desktop-ui/src-tauri/bin..."
	@mkdir -p desktop-ui/src-tauri/bin
	@cp $(BINARY_NAME) desktop-ui/src-tauri/bin/ts2go-cli$(if $(findstring .exe,$(BINARY_NAME)),.exe,)
	cd desktop-ui && npm run tauri build

# Run tests
test:
	cd tests && go test -v ./...

# Clean build artifacts
clean:
	$(RM) $(BINARY_NAME)
	cd desktop-ui && rm -rf dist node_modules src-tauri/target

# Install CLI globally
install:
	go install ./cmd/ts2go

# Development mode for desktop app
dev-desktop: build-cli
	@echo "Copying CLI binary to desktop-ui/src-tauri/bin..."
	@mkdir -p desktop-ui/src-tauri/bin
	@cp $(BINARY_NAME) desktop-ui/src-tauri/bin/ts2go-cli$(if $(findstring .exe,$(BINARY_NAME)),.exe,)
	cd desktop-ui && npm run tauri dev

# Install desktop dependencies
install-desktop:
	cd desktop-ui && npm install

# Show help
help:
	@echo "TS2Go Monorepo Build Commands"
	@echo ""
	@echo "  make build          - Build CLI tool (default)"
	@echo "  make build-cli      - Build CLI tool"
	@echo "  make build-desktop  - Build desktop app"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install        - Install CLI globally"
	@echo "  make dev-desktop    - Run desktop app in development mode"
	@echo "  make install-desktop- Install desktop app dependencies"
	@echo ""