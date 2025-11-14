.PHONY: build build-cli build-desktop build-desktop-windows build-desktop-linux build-desktop-all build-desktop-debug test clean install help dev-desktop install-desktop

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
	@echo "🧹 Cleaning and preparing bundling directory..."
	@rm -rf desktop-ui/src-tauri/bin
	@mkdir -p desktop-ui/src-tauri/bin
	@echo "Copying CLI binary to desktop-ui/src-tauri/bin..."
	@cp $(BINARY_NAME) desktop-ui/src-tauri/bin/ts2go-cli$(if $(findstring .exe,$(BINARY_NAME)),.exe,)
	@echo "Copying parser dependencies to desktop-ui/src-tauri/bin/parser..."
	@mkdir -p desktop-ui/src-tauri/bin/parser
	@cp -R internal/transpiler/parser/* desktop-ui/src-tauri/bin/parser/
	@echo "Copying package mappings to desktop-ui/src-tauri/bin/mappings..."
	@mkdir -p desktop-ui/src-tauri/bin/mappings
	@cp -R mappings/* desktop-ui/src-tauri/bin/mappings/
	@echo "Bundling Node.js runtime for $(PLATFORM)..."
	@if command -v node >/dev/null 2>&1; then \
		if [ "$(PLATFORM)" = "windows" ]; then \
			cp "$$(command -v node)" desktop-ui/src-tauri/bin/node.exe; \
			chmod +x desktop-ui/src-tauri/bin/node.exe; \
			echo "  ✓ Node.js runtime bundled for Windows ($$(node --version))"; \
		else \
			cp "$$(command -v node)" desktop-ui/src-tauri/bin/node; \
			chmod +x desktop-ui/src-tauri/bin/node; \
			echo "  ✓ Node.js runtime bundled for $(PLATFORM) ($$(node --version))"; \
		fi; \
	else \
		echo "  ⚠️  Warning: Node.js not found in PATH. App will require system Node.js."; \
	fi
	@echo "Building Tauri desktop application..."
	@echo "Note: DMG creation may fail (bundle_dmg.sh issue), but we'll create it manually if needed."
	cd desktop-ui && npm ci && (npm run tauri:build || echo "⚠️  Tauri build had errors, checking if .app was created...")
	@echo ""
	@echo "🔍 Verifying build artifacts..."
	@if [ "$(PLATFORM)" = "macos" ]; then \
		APP_PATH="desktop-ui/src-tauri/target/release/bundle/macos/TS2Go Desktop.app"; \
		if [ ! -d "$$APP_PATH" ]; then \
			echo "❌ ERROR: .app bundle not found"; \
			exit 1; \
		fi; \
		APP_SIZE=$$(du -sm "$$APP_PATH" | cut -f1); \
		echo "�📦 .app size: $${APP_SIZE}MB"; \
		if [ $$APP_SIZE -lt 100 ]; then \
			echo "❌ ERROR: .app too small ($${APP_SIZE}MB, expected >100MB)"; \
			exit 1; \
		fi; \
		if [ $$APP_SIZE -gt 180 ]; then \
			echo "⚠️  WARNING: .app larger than expected ($${APP_SIZE}MB, expected ~125-150MB)"; \
			echo "  This might indicate duplicate binaries."; \
		fi; \
		if [ ! -f "$$APP_PATH/Contents/Resources/bin/ts2go-cli" ]; then \
			echo "❌ ERROR: ts2go-cli missing from .app"; \
			exit 1; \
		fi; \
		NODE_COUNT=$$(find "$$APP_PATH/Contents/Resources/bin" -name "node*" -type f 2>/dev/null | grep -v node_modules | wc -l | tr -d ' '); \
		if [ $$NODE_COUNT -eq 0 ]; then \
			echo "❌ ERROR: Node.js runtime missing from .app"; \
			exit 1; \
		elif [ $$NODE_COUNT -gt 1 ]; then \
			echo "⚠️  WARNING: Multiple Node.js binaries found ($$NODE_COUNT)"; \
			find "$$APP_PATH/Contents/Resources/bin" -name "node*" -type f 2>/dev/null | grep -v node_modules; \
		fi; \
		echo "✅ .app bundle verified"; \
	fi
	@echo ""
	@echo "📦 Copying release to project root..."
	@mkdir -p $(RELEASE_DIR)/stable/$(PLATFORM)
	@if [ "$(PLATFORM)" = "macos" ]; then \
		rm -rf "$(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).app"; \
		if [ -d "desktop-ui/src-tauri/target/release/bundle/macos" ]; then \
			cp -R "desktop-ui/src-tauri/target/release/bundle/macos/"*.app "$(RELEASE_DIR)/stable/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		if [ -d "desktop-ui/src-tauri/target/release/bundle/dmg" ]; then \
			DMG_FILE=$$(find "desktop-ui/src-tauri/target/release/bundle/dmg" -name "*.dmg" -size +10M 2>/dev/null | head -1); \
			if [ -n "$$DMG_FILE" ]; then \
				cp "$$DMG_FILE" "$(RELEASE_DIR)/stable/$(PLATFORM)/" 2>/dev/null || true; \
				echo "✅ DMG copied from Tauri build"; \
			else \
				echo "⚠️  No valid DMG from Tauri, creating manually with hdiutil..."; \
				APP_TO_PACKAGE=$$(find "$(RELEASE_DIR)/stable/$(PLATFORM)" -name "*.app" | head -1); \
				if [ -n "$$APP_TO_PACKAGE" ]; then \
					hdiutil create -volname "TS2Go Desktop" \
						-srcfolder "$$APP_TO_PACKAGE" \
						-ov -format UDZO "$(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).dmg" || echo "⚠️  hdiutil failed, DMG not created"; \
					if [ -f "$(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).dmg" ]; then \
						echo "✅ DMG created manually with hdiutil"; \
					fi; \
				fi; \
			fi; \
		else \
			echo "⚠️  DMG directory not found, creating manually..."; \
			APP_TO_PACKAGE=$$(find "$(RELEASE_DIR)/stable/$(PLATFORM)" -name "*.app" | head -1); \
			if [ -n "$$APP_TO_PACKAGE" ]; then \
				hdiutil create -volname "TS2Go Desktop" \
					-srcfolder "$$APP_TO_PACKAGE" \
					-ov -format UDZO "$(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).dmg" || echo "⚠️  hdiutil failed, DMG not created"; \
				if [ -f "$(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).dmg" ]; then \
					echo "✅ DMG created manually with hdiutil"; \
				fi; \
			fi; \
		fi; \
		echo "✅ Build complete!"; \
		echo "📍 Release location: $(RELEASE_DIR)/stable/$(PLATFORM)/"; \
	elif [ "$(PLATFORM)" = "windows" ]; then \
		if [ -d "desktop-ui/src-tauri/target/release/bundle/msi" ]; then \
			cp desktop-ui/src-tauri/target/release/bundle/msi/*.msi "$(RELEASE_DIR)/stable/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		if [ -d "desktop-ui/src-tauri/target/release/bundle/nsis" ]; then \
			cp desktop-ui/src-tauri/target/release/bundle/nsis/*.exe "$(RELEASE_DIR)/stable/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		echo "✅ Build complete!"; \
		echo "📍 Release location: $(RELEASE_DIR)/stable/$(PLATFORM)/"; \
	else \
		if [ -d "desktop-ui/src-tauri/target/release/bundle/appimage" ]; then \
			cp desktop-ui/src-tauri/target/release/bundle/appimage/*.AppImage "$(RELEASE_DIR)/stable/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		if [ -d "desktop-ui/src-tauri/target/release/bundle/deb" ]; then \
			cp desktop-ui/src-tauri/target/release/bundle/deb/*.deb "$(RELEASE_DIR)/stable/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
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
	@echo "Copying package mappings to desktop-ui/src-tauri/bin/mappings..."
	@mkdir -p desktop-ui/src-tauri/bin/mappings
	@cp -R mappings/* desktop-ui/src-tauri/bin/mappings/
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
	@echo "Copying package mappings to desktop-ui/src-tauri/bin/mappings..."
	@mkdir -p desktop-ui/src-tauri/bin/mappings
	@cp -R mappings/* desktop-ui/src-tauri/bin/mappings/
	@echo "Building Tauri desktop application (DEBUG mode)..."
	cd desktop-ui && npm run tauri:build -- --debug
	@echo ""
	@echo "📦 Copying debug release to project root..."
	@mkdir -p $(RELEASE_DIR)/debug/$(PLATFORM)
	@if [ "$(PLATFORM)" = "macos" ]; then \
		rm -rf "$(RELEASE_DIR)/debug/$(PLATFORM)/TS2Go-Desktop-v$(VERSION)-debug.app"; \
		if [ -d "desktop-ui/src-tauri/target/debug/bundle/macos" ]; then \
			cp -R "desktop-ui/src-tauri/target/debug/bundle/macos/"*.app "$(RELEASE_DIR)/debug/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		if [ -d "desktop-ui/src-tauri/target/debug/bundle/dmg" ]; then \
			cp desktop-ui/src-tauri/target/debug/bundle/dmg/*.dmg "$(RELEASE_DIR)/debug/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		echo "✅ Debug build complete!"; \
		echo "📍 Debug release location: $(RELEASE_DIR)/debug/$(PLATFORM)/"; \
	elif [ "$(PLATFORM)" = "windows" ]; then \
		if [ -d "desktop-ui/src-tauri/target/debug/bundle/msi" ]; then \
			cp desktop-ui/src-tauri/target/debug/bundle/msi/*.msi "$(RELEASE_DIR)/debug/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		if [ -d "desktop-ui/src-tauri/target/debug/bundle/nsis" ]; then \
			cp desktop-ui/src-tauri/target/debug/bundle/nsis/*.exe "$(RELEASE_DIR)/debug/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		echo "✅ Debug build complete!"; \
		echo "📍 Debug release location: $(RELEASE_DIR)/debug/$(PLATFORM)/"; \
	else \
		if [ -d "desktop-ui/src-tauri/target/debug/bundle/appimage" ]; then \
			cp desktop-ui/src-tauri/target/debug/bundle/appimage/*.AppImage "$(RELEASE_DIR)/debug/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		if [ -d "desktop-ui/src-tauri/target/debug/bundle/deb" ]; then \
			cp desktop-ui/src-tauri/target/debug/bundle/deb/*.deb "$(RELEASE_DIR)/debug/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		echo "✅ Debug build complete!"; \
		echo "📍 Debug release location: $(RELEASE_DIR)/debug/$(PLATFORM)/"; \
	fi
	@echo ""

# Build for Windows (requires Windows or cross-compilation setup)
build-desktop-windows:
	@echo "⚠️  Windows builds from macOS require additional setup."
	@echo "Please choose one of the following options:"
	@echo ""
	@echo "Option 1: Build natively on Windows"
	@echo "  - Clone repo on Windows machine"
	@echo "  - Run: make build-desktop"
	@echo ""
	@echo "Option 2: Use GitHub Actions (recommended)"
	@echo "  - Push to GitHub"
	@echo "  - CI will build all platforms automatically"
	@echo ""
	@echo "Option 3: Install cargo-xwin (experimental)"
	@echo "  - cargo install cargo-xwin"
	@echo "  - brew install --cask wine-stable"
	@echo "  - Then retry this command"
	@echo ""
	@if command -v cargo-xwin >/dev/null 2>&1; then \
		echo "✓ cargo-xwin detected, attempting build..."; \
		$(MAKE) build-desktop-windows-xwin; \
	else \
		exit 1; \
	fi

# Internal target using cargo-xwin
build-desktop-windows-xwin: build-cli
	@echo "Building TS2Go Desktop for Windows using cargo-xwin..."
	@mkdir -p desktop-ui/src-tauri/bin
	@echo "Copying CLI binary to desktop-ui/src-tauri/bin..."
	@cp $(BINARY_NAME) desktop-ui/src-tauri/bin/ts2go-cli$(if $(findstring .exe,$(BINARY_NAME)),.exe,)
	@echo "Copying parser dependencies..."
	@mkdir -p desktop-ui/src-tauri/bin/parser
	@cp -R internal/transpiler/parser/* desktop-ui/src-tauri/bin/parser/
	@echo "Copying package mappings..."
	@mkdir -p desktop-ui/src-tauri/bin/mappings
	@cp -R mappings/* desktop-ui/src-tauri/bin/mappings/
	@echo "Note: Using host Node.js. For true Windows build, download Windows Node.js binary."
	@if command -v node >/dev/null 2>&1; then \
		cp "$$(command -v node)" desktop-ui/src-tauri/bin/node 2>/dev/null || true; \
	fi
	@echo "Building Tauri application for Windows with cargo-xwin..."
	cd desktop-ui && cargo xwin build --release --target x86_64-pc-windows-msvc
	@echo "📦 Copying Windows release to project root..."
	@mkdir -p $(RELEASE_DIR)/stable/windows
	@if [ -d "desktop-ui/src-tauri/target/x86_64-pc-windows-msvc/release/bundle/msi" ]; then \
		cp desktop-ui/src-tauri/target/x86_64-pc-windows-msvc/release/bundle/msi/*.msi "$(RELEASE_DIR)/stable/windows/" 2>/dev/null || true; \
		echo "✅ Windows build complete!"; \
		echo "📍 Release location: $(RELEASE_DIR)/stable/windows/"; \
	else \
		echo "⚠️  Windows bundle not found. Check build output."; \
	fi
	@echo ""

# Build for Linux (requires Linux or cross-compilation setup)
build-desktop-linux:
	@echo "Building TS2Go Desktop for Linux..."
	@if command -v cross >/dev/null 2>&1; then \
		echo "✓ Using cross for containerized build"; \
		$(MAKE) build-desktop-linux-cross; \
	elif [ "$$(uname -s)" = "Linux" ]; then \
		echo "✓ Building natively on Linux"; \
		$(MAKE) build-desktop-linux-native; \
	else \
		echo "⚠️  Linux builds from macOS require Docker-based cross-compilation."; \
		echo ""; \
		echo "Option 1: Install cross (recommended)"; \
		echo "  - cargo install cross"; \
		echo "  - Ensure Docker is running"; \
		echo "  - Then retry this command"; \
		echo ""; \
		echo "Option 2: Build natively on Linux"; \
		echo "  - Clone repo on Linux machine"; \
		echo "  - Run: make build-desktop"; \
		echo ""; \
		echo "Option 3: Use GitHub Actions"; \
		echo "  - Push to GitHub"; \
		echo "  - CI will build all platforms"; \
		echo ""; \
		exit 1; \
	fi

# Internal target using cross (Docker-based)
build-desktop-linux-cross: build-cli
	@echo "Building TS2Go Desktop for Linux using cross..."
	@mkdir -p desktop-ui/src-tauri/bin
	@cp $(BINARY_NAME) desktop-ui/src-tauri/bin/ts2go-cli
	@mkdir -p desktop-ui/src-tauri/bin/parser
	@cp -R internal/transpiler/parser/* desktop-ui/src-tauri/bin/parser/
	@mkdir -p desktop-ui/src-tauri/bin/mappings
	@cp -R mappings/* desktop-ui/src-tauri/bin/mappings/
	@echo "Building with cross..."
	cd desktop-ui/src-tauri && cross build --release --target x86_64-unknown-linux-gnu
	@echo "📦 Copying Linux release to project root..."
	@mkdir -p $(RELEASE_DIR)/stable/linux
	@if [ -d "desktop-ui/src-tauri/target/x86_64-unknown-linux-gnu/release/bundle" ]; then \
		if [ -d "desktop-ui/src-tauri/target/x86_64-unknown-linux-gnu/release/bundle/appimage" ]; then \
			cp desktop-ui/src-tauri/target/x86_64-unknown-linux-gnu/release/bundle/appimage/*.AppImage "$(RELEASE_DIR)/stable/linux/" 2>/dev/null || true; \
		fi; \
		if [ -d "desktop-ui/src-tauri/target/x86_64-unknown-linux-gnu/release/bundle/deb" ]; then \
			cp desktop-ui/src-tauri/target/x86_64-unknown-linux-gnu/release/bundle/deb/*.deb "$(RELEASE_DIR)/stable/linux/" 2>/dev/null || true; \
		fi; \
		echo "✅ Linux build complete!"; \
		echo "📍 Release location: $(RELEASE_DIR)/stable/linux/"; \
	else \
		echo "⚠️  Linux bundle not found. Check build output."; \
	fi
	@echo ""

# Internal target for native Linux build
build-desktop-linux-native: build-cli
	@echo "Building TS2Go Desktop for Linux (native)..."
	@mkdir -p desktop-ui/src-tauri/bin
	@cp $(BINARY_NAME) desktop-ui/src-tauri/bin/ts2go-cli
	@mkdir -p desktop-ui/src-tauri/bin/parser
	@cp -R internal/transpiler/parser/* desktop-ui/src-tauri/bin/parser/
	@mkdir -p desktop-ui/src-tauri/bin/mappings
	@cp -R mappings/* desktop-ui/src-tauri/bin/mappings/
	@if command -v node >/dev/null 2>&1; then \
		cp "$$(command -v node)" desktop-ui/src-tauri/bin/node; \
		chmod +x desktop-ui/src-tauri/bin/node; \
	fi
	cd desktop-ui && npm run tauri:build
	@mkdir -p $(RELEASE_DIR)/stable/linux
	@if [ -d "desktop-ui/src-tauri/target/release/bundle/appimage" ]; then \
		cp desktop-ui/src-tauri/target/release/bundle/appimage/*.AppImage "$(RELEASE_DIR)/stable/linux/" 2>/dev/null || true; \
	fi
	@if [ -d "desktop-ui/src-tauri/target/release/bundle/deb" ]; then \
		cp desktop-ui/src-tauri/target/release/bundle/deb/*.deb "$(RELEASE_DIR)/stable/linux/" 2>/dev/null || true; \
	fi
	@echo "✅ Linux build complete!"
	@echo "📍 Release location: $(RELEASE_DIR)/stable/linux/"
	@echo ""

# Build for all platforms (requires appropriate toolchains)
build-desktop-all:
	@echo "🏗️  Building TS2Go Desktop for all platforms..."
	@echo ""
	@echo "Building macOS..."
	@$(MAKE) build-desktop
	@echo ""
	@echo "Attempting Windows build..."
	@-$(MAKE) build-desktop-windows || echo "⚠️  Windows build skipped (requires setup)"
	@echo ""
	@echo "Attempting Linux build..."
	@-$(MAKE) build-desktop-linux || echo "⚠️  Linux build skipped (requires setup)"
	@echo ""
	@echo "📦 Build summary:"
	@ls -lh $(RELEASE_DIR)/stable/ 2>/dev/null || echo "No builds completed"

# Install desktop dependencies
install-desktop:
	cd desktop-ui && npm install

# Show help
help:
	@echo "TS2Go Monorepo Build Commands"
	@echo ""
	@echo "CLI Builds:"
	@echo "  make build                  - Build CLI tool (default)"
	@echo "  make build-cli              - Build CLI tool"
	@echo "  make install                - Install CLI globally"
	@echo ""
	@echo "Desktop App Builds:"
	@echo "  make build-desktop          - Build desktop app for current platform (stable/release)"
	@echo "  make build-desktop-debug    - Build desktop app for current platform (debug/unoptimized)"
	@echo "  make build-desktop-windows  - Build desktop app for Windows (x86_64)"
	@echo "  make build-desktop-linux    - Build desktop app for Linux (x86_64)"
	@echo "  make build-desktop-all      - Build desktop app for all platforms"
	@echo "  make dev-desktop            - Run desktop app in development mode"
	@echo "  make install-desktop        - Install desktop app dependencies"
	@echo ""
	@echo "Testing & Maintenance:"
	@echo "  make test                   - Run tests"
	@echo "  make clean                  - Clean build artifacts"
	@echo ""
	@echo "Output locations:"
	@echo "  Stable builds:  release/stable/{macos,windows,linux}/"
	@echo "  Debug builds:   release/debug/{platform}/"
	@echo ""
	@echo "Cross-compilation notes:"
	@echo "  - Windows builds from macOS/Linux require: rustup target add x86_64-pc-windows-msvc"
	@echo "  - Linux builds from macOS/Windows require: rustup target add x86_64-unknown-linux-gnu"
	@echo "  - Some targets may require additional system dependencies"
	@echo ""