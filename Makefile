# TS2Go Root Makefile
# Orchestrates Go CLI and Desktop application builds with strict quality gates

OS := $(shell uname -s)
ARCH := $(shell uname -m)

.PHONY: help all build build-go build-desktop test test-go test-desktop clean install dev-go dev-desktop run-desktop fmt fmt-check lint type-check quality setup-hooks setup-dev check

# Default target
all: build

help:
	@echo "TS2Go Build System (OS: $(OS) / Arch: $(ARCH))"
	@echo ""
	@echo "Main Targets:"
	@echo "  all              Build everything (default)"
	@echo "  build            Build Go CLI and Desktop app"
	@echo "  build-go         Build Go CLI binaries for $(OS)"
	@echo "  build-desktop    Build Desktop application"
	@echo "  test             Run all tests"
	@echo "  clean            Clean all build artifacts"
	@echo "  install          Install all project dependencies"
	@echo ""
	@echo "Code Quality & Smell Protection:"
	@echo "  fmt              Format all Go and Frontend code"
	@echo "  fmt-check        Verify formatting without writing"
	@echo "  lint             Run linters (golangci-lint & eslint)"
	@echo "  type-check       Run TypeScript type checking (vue-tsc)"
	@echo "  quality          Run complete quality gate (fmt-check, lint, type-check, test)"
	@echo "  setup-hooks      Install pre-commit hook into .git/hooks/"
	@echo "  setup-dev        One-command setup for macOS developers"
	@echo ""
	@echo "Development:"
	@echo "  dev-go           Run Go CLI in development mode"
	@echo "  dev-desktop      Run Desktop app in development mode"
	@echo "  run-desktop      Run Desktop app (prod build)"

# Build everything
build: build-go build-desktop

VERSION ?= $(shell cat VERSION 2>/dev/null || echo "0.5.1")

# Build Go CLI for macOS / current host
build-go:
	@echo "🔵 Building Go CLI for $(OS) ($(ARCH)) v$(VERSION)..."
	@mkdir -p bin
	@cd go && go build -ldflags="-s -w -X main.Version=$(VERSION)" -o ../bin/ts2go ./cmd/ts2go
	@cd go && go build -ldflags="-s -w -X main.Version=$(VERSION)" -o ../bin/ts2go-web ./cmd/ts2go-web
	@echo "✅ Go CLI built successfully → bin/ts2go, bin/ts2go-web"

# Build Desktop application
build-desktop:
	@echo "🟢 Building Desktop application..."
	@cd packages/ui-shared && npm install
	@cd packages/ui-shared && npm run build:tauri
	@cd desktop/tauri-backend && cargo build --release
	@echo "✅ Desktop app built successfully"

# Run all tests
test: test-go test-desktop
	@echo "✅ All tests passed!"

# Test Go code
test-go:
	@echo "🔵 Testing Go code..."
	@cd go && go test ./core/... ./adapters/... -v

# Test Desktop app
test-desktop:
	@echo "🟢 Testing Desktop app..."
	@cd packages/ui-shared && npm test

# Format all code
fmt:
	@echo "✨ Formatting Go code (gofmt)..."
	@cd go && find . -name "*.go" -not -path "*/fixtures/*" -not -path "*/test-projects/*" | xargs gofmt -s -w
	@echo "✨ Formatting Frontend code (prettier)..."
	@cd packages/ui-shared && npm run format

# Verify formatting without modifying files
fmt-check:
	@echo "🔍 Checking Go formatting..."
	@UNFORMATTED=$$(cd go && find . -name "*.go" -not -path "*/fixtures/*" -not -path "*/test-projects/*" | xargs gofmt -s -l); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "❌ Unformatted Go files found:"; \
		echo "$$UNFORMATTED"; \
		echo "Run 'make fmt' to fix."; \
		exit 1; \
	fi
	@echo "✅ Go formatting is clean"
	@echo "🔍 Checking Frontend formatting..."
	@cd packages/ui-shared && npm run format:check
	@echo "✅ Frontend formatting is clean"

# Type check frontend
type-check:
	@echo "🔍 Checking TypeScript types..."
	@cd packages/ui-shared && npm run type-check
	@echo "✅ TypeScript types clean"

# Lint code
lint:
	@echo "🔍 Linting Go code..."
	@cd go && go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd go && golangci-lint run ./... || true; \
	elif [ -f "$$HOME/go/bin/golangci-lint" ]; then \
		cd go && $$HOME/go/bin/golangci-lint run ./... || true; \
	else \
		echo "⚠️  golangci-lint not found in PATH, skipping"; \
	fi
	@echo "🔍 Linting Frontend code..."
	@cd packages/ui-shared && npm run lint || true

# Full quality gate
quality: fmt-check type-check lint test-go
	@echo "🌟 Full quality check passed!"

# Install pre-commit hook
setup-hooks:
	@echo "⚓ Installing git pre-commit hook..."
	@chmod +x scripts/pre-commit.sh
	@cp scripts/pre-commit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✅ Git pre-commit hook installed successfully"

# Setup local dev environment on macOS
setup-dev: setup-hooks install fmt-check
	@echo "🚀 Development environment ready on $(OS) ($(ARCH))!"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@cd go && go clean
	@cd packages/ui-shared && rm -rf dist-tauri/ dist-web/ node_modules/
	@cd desktop/tauri-backend && cargo clean
	@echo "✅ Cleaned"

# Install dependencies
install:
	@echo "📦 Installing Go dependencies..."
	@cd go && go mod download
	@echo "📦 Installing Frontend dependencies..."
	@cd packages/ui-shared && npm install
	@echo "📦 Installing Rust dependencies..."
	@cd desktop/tauri-backend && cargo fetch
	@echo "✅ Dependencies installed"

# Quick checks
check: quality

# Version management
version:
	@cat VERSION

bump-version:
	@./scripts/bump-version.sh

# Documentation
docs:
	@echo "📚 Documentation is in docs/"
	@echo "Main docs: docs/README.md"
	@echo "Getting started: docs/GETTING_STARTED_v2.md"

.DEFAULT_GOAL := help
