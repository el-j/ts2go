package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PackageInfo represents the contents of a package.json file
type PackageInfo struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Description      string            `json:"description"`
	Main             string            `json:"main"`
	Type             string            `json:"type"` // "module" or "commonjs"
	Scripts          map[string]string `json:"scripts"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	Author           string            `json:"author"`
	License          string            `json:"license"`
}

// ParsePackageJSON reads and parses a package.json file
func ParsePackageJSON(path string) (*PackageInfo, error) {
	// Resolve to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if path is a directory - if so, append package.json
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	if fileInfo.IsDir() {
		absPath = filepath.Join(absPath, "package.json")
	}

	// Read file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}

	// Parse JSON
	var pkg PackageInfo
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	// Set default type if not specified
	if pkg.Type == "" {
		pkg.Type = "commonjs"
	}

	return &pkg, nil
}

// GetAllDependencies returns all dependencies (including dev and peer)
func (p *PackageInfo) GetAllDependencies() map[string]string {
	all := make(map[string]string)

	// Regular dependencies
	for name, version := range p.Dependencies {
		all[name] = version
	}

	// Dev dependencies
	for name, version := range p.DevDependencies {
		all[name] = version
	}

	// Peer dependencies
	for name, version := range p.PeerDependencies {
		all[name] = version
	}

	return all
}

// GetProductionDependencies returns only runtime dependencies (not dev)
func (p *PackageInfo) GetProductionDependencies() map[string]string {
	prod := make(map[string]string)

	for name, version := range p.Dependencies {
		prod[name] = version
	}

	return prod
}

// IsScopedPackage checks if a package name is scoped (e.g., @org/package)
func IsScopedPackage(name string) bool {
	return len(name) > 0 && name[0] == '@'
}

// ExtractScope returns the scope from a scoped package name
// Example: "@babel/core" returns "babel"
func ExtractScope(name string) string {
	if !IsScopedPackage(name) {
		return ""
	}

	for i := 1; i < len(name); i++ {
		if name[i] == '/' {
			return name[1:i]
		}
	}

	return ""
}

// ExtractPackageName returns the package name without scope
// Example: "@babel/core" returns "core"
func ExtractPackageName(name string) string {
	if !IsScopedPackage(name) {
		return name
	}

	for i := 1; i < len(name); i++ {
		if name[i] == '/' {
			if i+1 < len(name) {
				return name[i+1:]
			}
			return ""
		}
	}

	return name
}

// IsBuiltinModule checks if a package is a Node.js built-in module
func IsBuiltinModule(name string) bool {
	builtins := map[string]bool{
		"assert":         true,
		"buffer":         true,
		"child_process":  true,
		"cluster":        true,
		"crypto":         true,
		"dgram":          true,
		"dns":            true,
		"domain":         true,
		"events":         true,
		"fs":             true,
		"http":           true,
		"https":          true,
		"net":            true,
		"os":             true,
		"path":           true,
		"punycode":       true,
		"querystring":    true,
		"readline":       true,
		"stream":         true,
		"string_decoder": true,
		"timers":         true,
		"tls":            true,
		"tty":            true,
		"url":            true,
		"util":           true,
		"v8":             true,
		"vm":             true,
		"zlib":           true,
		// Node 10+
		"async_hooks":    true,
		"perf_hooks":     true,
		"worker_threads": true,
	}

	// Remove "node:" prefix if present
	if len(name) > 5 && name[:5] == "node:" {
		name = name[5:]
	}

	return builtins[name]
}
