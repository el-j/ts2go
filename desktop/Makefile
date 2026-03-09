# Desktop Application Makefile
# Builds and runs the Tauri + Vue.js desktop application

.PHONY: help dev build test clean install

help:
	@echo "Desktop Application Build System"
	@echo ""
	@echo "Targets:"
	@echo "  dev       Run in development mode"
	@echo "  build     Build production application"
	@echo "  test      Run tests"
	@echo "  clean     Clean build artifacts"
	@echo "  install   Install dependencies"

# Development mode
dev:
	@echo "🟢 Starting development server..."
	@cd tauri && cargo tauri dev

# Production build
build:
	@echo "🟢 Building production application..."
	@cd ui && npm install
	@cd ui && npm run build
	@cd tauri && cargo tauri build
	@echo "✅ Build complete → tauri/target/release/"

# Run tests
test:
	@echo "🧪 Running tests..."
	@cd ui && npm test

# Run tests with coverage
test-coverage:
	@cd ui && npm run test:coverage

# Clean artifacts
clean:
	@echo "🧹 Cleaning..."
	@cd ui && rm -rf dist/ node_modules/
	@cd tauri && cargo clean
	@echo "✅ Cleaned"

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	@cd ui && npm install
	@cd tauri && cargo fetch
	@echo "✅ Dependencies installed"

# Lint code
lint:
	@cd ui && npm run lint || echo "⚠️  Linting not configured"

# Format code
fmt:
	@cd ui && npm run format || echo "⚠️  Formatting not configured"

# Type check
type-check:
	@cd ui && npm run type-check

.DEFAULT_GOAL := help
