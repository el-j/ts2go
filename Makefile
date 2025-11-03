.PHONY: build build-cli build-desktop test clean install help

# Default target
all: build

# Build CLI tool
build: build-cli

build-cli:
	GOWORK=off go build -o ts2go ./cmd/ts2go

# Build desktop app
build-desktop:
	cd packages/desktop && npm run build

# Run tests
test:
	cd tests && go test -v ./...

# Clean build artifacts
clean:
	rm -f ts2go
	cd packages/desktop && rm -rf dist node_modules

# Install CLI globally
install:
	GOWORK=off go install ./cmd/ts2go

# Development mode for desktop app
dev-desktop:
	cd packages/desktop && npm run dev

# Install desktop dependencies
install-desktop:
	cd packages/desktop && npm install

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