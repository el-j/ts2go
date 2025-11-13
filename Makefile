.PHONY: build build-cli build-desktop test clean install help

# Get version from VERSION file
VERSION := $(shell cat VERSION | tr -d '\n')

# Detect OS for binary name
ifeq ($(OS),Windows_NT)
    BINARY_NAME := ts2go.exe
    RM := del /Q
    RELEASE_DIR := release
    PLATFORM := windows
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        PLATFORM := macos
        RELEASE_DIR := release
    else
        PLATFORM := linux
        RELEASE_DIR := release
    endif
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
	@echo "Copying parser dependencies to desktop-ui/src-tauri/bin/parser..."
	@mkdir -p desktop-ui/src-tauri/bin/parser
	@cp -R internal/transpiler/parser/* desktop-ui/src-tauri/bin/parser/
	@echo "Bundling Node.js runtime..."
	@if command -v node >/dev/null 2>&1; then \
		cp "$$(command -v node)" desktop-ui/src-tauri/bin/node; \
		chmod +x desktop-ui/src-tauri/bin/node; \
		echo "  ✓ Node.js runtime bundled ($$(node --version))"; \
	else \
		echo "  ⚠️  Warning: Node.js not found in PATH. App will require system Node.js."; \
	fi
	@echo "Building Tauri desktop application..."
	cd desktop-ui && npm run tauri:build
	@echo ""
	@echo "📦 Copying release to project root..."
	@mkdir -p $(RELEASE_DIR)/stable/$(PLATFORM)
	@if [ "$(PLATFORM)" = "macos" ]; then \
		rm -rf "$(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).app"; \
		cp -R "desktop-ui/src-tauri/target/release/bundle/macos/TS2Go Desktop.app" "$(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).app"; \
		echo "✅ Build complete!"; \
		echo "📍 Release location: $(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).app"; \
	elif [ "$(PLATFORM)" = "windows" ]; then \
		cp -R desktop-ui/src-tauri/target/release/bundle/msi/* "$(RELEASE_DIR)/stable/$(PLATFORM)/"; \
		echo "✅ Build complete!"; \
		echo "📍 Release location: $(RELEASE_DIR)/stable/$(PLATFORM)/"; \
	else \
		cp -R desktop-ui/src-tauri/target/release/bundle/appimage/* "$(RELEASE_DIR)/stable/$(PLATFORM)/"; \
		echo "✅ Build complete!"; \
		echo "📍 Release location: $(RELEASE_DIR)/stable/$(PLATFORM)/"; \
	fi
	@echo ""

# Run tests
test:
	cd tests && go test -v ./...

# Clean build artifacts
clean:
	$(RM) $(BINARY_NAME)
	cd desktop-ui && rm -rf dist node_modules src-tauri/target
	rm -rf $(RELEASE_DIR)

# Install CLI globally
install:
	go install ./cmd/ts2go

# Development mode for desktop app
dev-desktop: build-cli
	@echo "Copying CLI binary to desktop-ui/src-tauri/bin..."
	@mkdir -p desktop-ui/src-tauri/bin
	@cp $(BINARY_NAME) desktop-ui/src-tauri/bin/ts2go-cli$(if $(findstring .exe,$(BINARY_NAME)),.exe,)
	@echo "Copying parser dependencies to desktop-ui/src-tauri/bin/parser..."
	@mkdir -p desktop-ui/src-tauri/bin/parser
	@cp -R internal/transpiler/parser/* desktop-ui/src-tauri/bin/parser/
	@echo "Bundling Node.js runtime..."
	@if command -v node >/dev/null 2>&1; then \
		cp "$$(command -v node)" desktop-ui/src-tauri/bin/node; \
		chmod +x desktop-ui/src-tauri/bin/node; \
		echo "  ✓ Node.js runtime bundled ($$(node --version))"; \
	fi
	cd desktop-ui && npm run tauri:dev

# Build desktop app in debug mode (unoptimized, with symbols)
build-desktop-debug: build-cli
	@echo "Copying CLI binary to desktop-ui/src-tauri/bin..."
	@mkdir -p desktop-ui/src-tauri/bin
	@cp $(BINARY_NAME) desktop-ui/src-tauri/bin/ts2go-cli$(if $(findstring .exe,$(BINARY_NAME)),.exe,)
	@echo "Copying parser dependencies to desktop-ui/src-tauri/bin/parser..."
	@mkdir -p desktop-ui/src-tauri/bin/parser
	@cp -R internal/transpiler/parser/* desktop-ui/src-tauri/bin/parser/
	@echo "Bundling Node.js runtime..."
	@if command -v node >/dev/null 2>&1; then \
		cp "$$(command -v node)" desktop-ui/src-tauri/bin/node; \
		chmod +x desktop-ui/src-tauri/bin/node; \
		echo "  ✓ Node.js runtime bundled ($$(node --version))"; \
	fi
	@cp -R internal/transpiler/parser/* desktop-ui/src-tauri/bin/parser/
	@echo "Building Tauri desktop application (DEBUG mode)..."
	cd desktop-ui && npm run tauri:build -- --debug
	@echo ""
	@echo "📦 Copying debug release to project root..."
	@mkdir -p $(RELEASE_DIR)/debug/$(PLATFORM)
	@if [ "$(PLATFORM)" = "macos" ]; then \
		rm -rf "$(RELEASE_DIR)/debug/$(PLATFORM)/TS2Go-Desktop-v$(VERSION)-debug.app"; \
		cp -R "desktop-ui/src-tauri/target/debug/bundle/macos/TS2Go Desktop.app" "$(RELEASE_DIR)/debug/$(PLATFORM)/TS2Go-Desktop-v$(VERSION)-debug.app"; \
		echo "✅ Debug build complete!"; \
		echo "📍 Debug release location: $(RELEASE_DIR)/debug/$(PLATFORM)/TS2Go-Desktop-v$(VERSION)-debug.app"; \
	elif [ "$(PLATFORM)" = "windows" ]; then \
		cp -R desktop-ui/src-tauri/target/debug/bundle/msi/* "$(RELEASE_DIR)/debug/$(PLATFORM)/"; \
		echo "✅ Debug build complete!"; \
		echo "📍 Debug release location: $(RELEASE_DIR)/debug/$(PLATFORM)/"; \
	else \
		cp -R desktop-ui/src-tauri/target/debug/bundle/appimage/* "$(RELEASE_DIR)/debug/$(PLATFORM)/"; \
		echo "✅ Debug build complete!"; \
		echo "📍 Debug release location: $(RELEASE_DIR)/debug/$(PLATFORM)/"; \
	fi
	@echo ""

# Install desktop dependencies
install-desktop:
	cd desktop-ui && npm install

# Show help
help:
	@echo "TS2Go Monorepo Build Commands"
	@echo ""
	@echo "  make build             - Build CLI tool (default)"
	@echo "  make build-cli         - Build CLI tool"
	@echo "  make build-desktop     - Build desktop app (stable/release)"
	@echo "  make build-desktop-debug - Build desktop app (debug/unoptimized)"
	@echo "  make test              - Run tests"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make install           - Install CLI globally"
	@echo "  make dev-desktop       - Run desktop app in development mode"
	@echo "  make install-desktop   - Install desktop app dependencies"
	@echo ""
	@echo "Output locations:"
	@echo "  Stable builds:  release/stable/{platform}/"
	@echo "  Debug builds:   release/debug/{platform}/"
	@echo ""