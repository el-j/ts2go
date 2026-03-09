# NPM to Go Package Mapper

This module provides a comprehensive mapping database from npm packages to their Go equivalents, along with tools to query and utilize these mappings.

## Overview

The mapper consists of:

1. **Mapping Database** (`mappings/npm-to-go.yaml`) - 49 package mappings including:
   - Node.js built-in modules (fs, path, http, crypto, etc.)
   - Popular npm packages (axios, lodash, express, uuid, etc.)
   - Status indicators (supported, partial, unsupported)
   - API-level mappings for method calls

2. **Go Loader** (`internal/mapper/loader.go`) - Programmatic access to mappings
3. **CLI Tool** (`internal/mapper/cmd/mapping-cli`) - Query tool for developers

## Statistics

- **Total Packages:** 49
- **Supported:** 32 (65.3%)
- **Partial Support:** 12 (24.5%)
- **Unsupported:** 5 (10.2%)

**By Type:**
- Runtime: 3 packages (custom runtime wrappers)
- Stdlib: 17 packages (Go standard library)
- Equivalent: 21 packages (third-party Go libraries)
- Framework: 3 packages (web frameworks)
- Unsupported: 5 packages (frontend/build tools)

## Usage

### CLI Tool

Query the mapping database from the command line:

```bash
# Show summary statistics
go run internal/mapper/cmd/mapping-cli/main.go summary

# Look up a specific package
go run internal/mapper/cmd/mapping-cli/main.go lookup -package axios

# List all runtime mappings
go run internal/mapper/cmd/mapping-cli/main.go list -type runtime

# List all stdlib mappings
go run internal/mapper/cmd/mapping-cli/main.go list -type stdlib

# List all equivalent packages
go run internal/mapper/cmd/mapping-cli/main.go list -type equivalent
```

### Programmatic Access

Use the Go API in your code:

```go
import "github.com/yourusername/ts2go/internal/mapper"

// Load mappings from default location
db, err := mapper.LoadDefaultMappings()
if err != nil {
    log.Fatal(err)
}

// Look up a package
mapping, err := db.GetMapping("axios")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("axios → %s\n", mapping.Go)
// Output: axios → github.com/go-resty/resty/v2

// Check if package is supported
if mapping.IsSupported() {
    fmt.Println("Package is supported!")
}

// Get API-level mapping
if goAPI, ok := mapping.GetAPIMapping("axios.get"); ok {
    fmt.Printf("axios.get → %s\n", goAPI)
    // Output: axios.get → resty.R().Get
}

// List all runtime packages
runtimePkgs := db.GetRuntimeMappings()
for _, pkg := range runtimePkgs {
    fmt.Printf("%s → %s\n", pkg.Npm, pkg.Go)
}
```

## Mapping Structure

Each mapping includes:

```yaml
- npm: "axios"                              # NPM package name
  go: "github.com/go-resty/resty/v2"       # Go package import path
  type: "equivalent"                        # runtime|stdlib|equivalent|framework|unsupported
  status: "supported"                       # supported|partial|unsupported
  complexity: "medium"                      # simple|medium|high
  description: "Promise based HTTP client"  # Brief description
  api_mappings:                            # Method-level mappings
    "axios.get": "resty.R().Get"
    "axios.post": "resty.R().Post"
  notes: "Optional notes about differences" # Implementation notes
  example: |                                # Code example
    // TypeScript
    const res = await axios.get('https://api.com/users');
    // Go
    res, err := resty.New().R().Get("https://api.com/users")
```

## Examples

### Node.js Built-ins

```
fs     → github.com/yourusername/ts2go/runtime/fs  (runtime)
path   → github.com/yourusername/ts2go/runtime/path (runtime)
os     → os                                         (stdlib)
http   → net/http                                   (stdlib)
crypto → crypto                                     (stdlib)
url    → net/url                                    (stdlib)
```

### Popular Packages

```
axios          → github.com/go-resty/resty/v2        (HTTP client)
lodash         → github.com/samber/lo                (utilities)
express        → github.com/gin-gonic/gin            (web framework)
uuid           → github.com/google/uuid              (UUID generation)
bcrypt         → golang.org/x/crypto/bcrypt          (password hashing)
jsonwebtoken   → github.com/golang-jwt/jwt/v5        (JWT tokens)
moment         → time                                 (date/time)
```

### API Mappings

**Axios HTTP Client:**
```
axios.get      → resty.R().Get
axios.post     → resty.R().Post
axios.put      → resty.R().Put
axios.delete   → resty.R().Delete
axios.patch    → resty.R().Patch
```

**Lodash Utilities:**
```
_.map      → lo.Map
_.filter   → lo.Filter
_.find     → lo.Find
_.reduce   → lo.Reduce
_.groupBy  → lo.GroupBy
_.sortBy   → lo.SortBy
```

**Express Framework:**
```
express()     → gin.Default()
app.get       → router.GET
app.post      → router.POST
app.listen    → router.Run
```

## Testing

Run the test suite:

```bash
cd internal/mapper
go test -v
```

Tests include:
- Unit tests for loader functionality (12 tests)
- Integration tests with actual mapping file (5 test suites)
- Verification of all 49 package mappings
- API mapping validation

## Future Enhancements

Planned additions:
1. More npm packages (targeting 100+ mappings)
2. Advanced transformation rules for complex APIs
3. Custom mapping overrides via config file
4. Automatic Go package installation
5. Version compatibility matrix
6. Community-contributed mappings

## Contributing

To add a new mapping:

1. Edit `mappings/npm-to-go.yaml`
2. Add the mapping with required fields (npm, go, type, status, complexity, description)
3. Include API mappings if applicable
4. Add integration test in `integration_test.go`
5. Run tests: `go test -v`

## License

Part of the ts2go TypeScript to Go transpiler project.
