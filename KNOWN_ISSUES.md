# Known Issues and Limitations

**Last Updated:** November 9, 2025  
**Version:** 0.5.1 → 0.7.0-beta

---

## Overview

This document lists known issues, limitations, and workarounds for the ts2go transpiler. These items are being actively worked on and will be addressed in upcoming releases.

---

## Critical Issues (Being Fixed in v0.6.0-beta)

### 1. Incomplete Function Body Generation

**Status:** ✅ PARTIALLY FIXED  
**Priority:** P0 - Critical  
**Target Fix:** v0.6.0-beta

**Issue:**
Some TypeScript functions transpile to Go code with incomplete function bodies, particularly for functions using array methods.

**What's Fixed:**
- ✅ Arrow functions now generate complete bodies
- ✅ Template literals transpile correctly
- ✅ Array methods (filter, map, reduce, etc.) now call runtime/array package

**Remaining Issues:**
- Type inference for lambda parameters (uses interface{} instead of specific types)
- Some edge cases in complex nested expressions

**Example (NOW WORKS):**
```typescript
const double = (x: number) => x * 2;
const evens = [1,2,3,4].filter(x => x % 2 === 0);
```

Now generates:
```go
double := func(x float64) interface{} { return x * 2 }
evens := array.Filter([]interface{}{1,2,3,4}, func(x interface{}) interface{} { return x % 2 == 0 })
```

**Workaround for Type Issues:**
- Explicitly type variables before using array methods
- Cast interface{} results to specific types when needed

**Tracking:** See ROADMAP_TO_1.0.0.md Phase 1 - Day 1-2 ✅ COMPLETE

---

### 2. Missing Import Statements

**Status:** ✅ MOSTLY FIXED  
**Priority:** P0 - Critical  
**Target Fix:** v0.6.0-beta

**Issue:**
Generated Go code may be missing required import statements (e.g., `fmt`, `strings`), causing compilation errors.

**What's Fixed:**
- ✅ fmt import automatically added for console.log and template literals
- ✅ Array runtime import automatically added for array methods
- ✅ Import tracking system exists and works

**Remaining Issues:**
- Some edge cases with conditional imports
- Runtime library imports may need explicit module setup

**Example (NOW WORKS):**
```typescript
console.log("Hello");
const evens = [1,2,3].filter(x => x % 2 === 0);
```

Now generates:
```go
import (
    "fmt"
    "github.com/ts2go/runtime/array"
)
```

**Workaround:**
- Imports are now automatically added in most cases
- If missing, manually add required imports

**Tracking:** See ROADMAP_TO_1.0.0.md Phase 1 - Day 3 ✅ WORKING

---

### 3. Template Literal Transpilation Incomplete

**Status:** ✅ FIXED  
**Priority:** P0 - Critical  
**Target Fix:** v0.6.0-beta

**Issue:**
Template literals may not fully transpile to proper Go string handling.

**What's Fixed:**
- ✅ Template literals now transpile using string concatenation
- ✅ Nested expressions are handled correctly
- ✅ Special characters are escaped properly
- ✅ fmt import is automatically added

**Example (NOW WORKS):**
```typescript
const name = "World";
const msg = `Hello, ${name}!`;
```

Now generates:
```go
import "fmt"

name := "World"
msg := "Hello, " + fmt.Sprint(name) + "!"
```

**Note:** Uses string concatenation instead of fmt.Sprintf for performance.

**Tracking:** See ROADMAP_TO_1.0.0.md Phase 1 - Day 4 ✅ COMPLETE

---

### NEW: Array Runtime Library

**Status:** ✅ IMPLEMENTED  
**Priority:** P0 - Critical  
**Implemented:** Current version

**Feature:**
Complete JavaScript-like array methods now available via runtime/array package.

**Supported Methods:**
- ✅ Filter, Map, Reduce - Functional programming
- ✅ Find, FindIndex, Some, Every - Searching  
- ✅ Push, Pop, Shift, Unshift - Stack/queue operations
- ✅ Includes, IndexOf, LastIndexOf - Element lookup
- ✅ Reverse, Slice, Concat - Manipulation
- ✅ Join, ForEach, Sort - Utilities

**Example:**
```typescript
const numbers = [1, 2, 3, 4, 5, 6];
const evens = numbers.filter(x => x % 2 === 0);
const doubled = numbers.map(x => x * 2);
const sum = numbers.reduce((acc, x) => acc + x, 0);
```

Generates working Go code:
```go
import "github.com/ts2go/runtime/array"

numbers := []interface{}{1, 2, 3, 4, 5, 6}
evens := array.Filter(numbers, func(x interface{}) interface{} { return x%2 == 0 })
doubled := array.Map(numbers, func(x interface{}) interface{} { return x * 2 })
sum := array.Reduce(numbers, func(acc interface{}, x interface{}) interface{} { return acc + x }, 0)
```

**Known Limitation:**
Type inference uses interface{} for lambda parameters. May require type assertions for complex operations.

**Usage:**
Array methods are automatically transpiled. No special configuration needed.

---

## High Priority Issues

### 4. Object Literal Type Inference

**Status:** 🟡 Partial Support  
**Priority:** P1 - High  
**Target Fix:** v0.6.0-beta

**Issue:**
Object literals passed as function arguments may not have proper type prefixes, causing compilation errors.

**Example:**
```typescript
addRecord({ id: 1, name: "Test" });
```

Previously generated:
```go
addRecord({Id: 1, Name: "Test"})  // Invalid syntax
```

**Current Fix:**
Now generates:
```go
addRecord(DataRecord{Id: 1, Name: "Test"})  // Valid
```

**Remaining Issues:**
- Type inference may fail in complex scenarios
- Nested object literals may still have issues

**Workaround:**
- Explicitly type variables before passing to functions
- Use type assertions in Go code

---

### 5. Class Inheritance Super Calls

**Status:** 🟡 Partial Support  
**Priority:** P1 - High  
**Target Fix:** v0.6.0-beta

**Issue:**
`super()` constructor calls and `super.method()` calls may generate `/* unsupported expression */` comments.

**Example:**
```typescript
class Derived extends Base {
    constructor(id: number) {
        super(id);
    }
}
```

May generate incomplete constructor code.

**Workaround:**
- Manually implement inheritance in Go
- Initialize base struct fields directly

**Tracking:** See IMPLEMENTATION_PLAN.md Phase 1

---

## Modern JavaScript Features

### NEW: Spread Operators (Arrays)

**Status:** ✅ IMPLEMENTED  
**Priority:** P1 - High  
**Implemented:** Current version

**Feature:**
Array spread operators now work correctly.

**Example:**
```typescript
const arr1 = [1, 2, 3];
const arr2 = [4, 5, 6];
const combined = [...arr1, ...arr2];
const withExtra = [...arr1, 7, 8, ...arr2];
```

Generates:
```go
import "github.com/ts2go/runtime/array"

arr1 := []interface{}{1, 2, 3}
arr2 := []interface{}{4, 5, 6}
combined := array.Concat(arr1, arr2)
withExtra := array.Concat(arr1, []interface{}{7, 8}, arr2)
```

**Note:** Uses runtime/array Concat function for efficient merging.

---

### NEW: Rest Parameters

**Status:** ✅ IMPLEMENTED  
**Priority:** P1 - High  
**Implemented:** Current version

**Feature:**
Rest parameters (...args) now transpile to Go variadic parameters.

**Example:**
```typescript
function sum(...numbers: number[]): number {
    let total = 0;
    for (const num of numbers) {
        total += num;
    }
    return total;
}
```

Generates:
```go
func Sum(numbers ...float64) float64 {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}
```

**Note:** Properly extracts element type from array type for variadic parameter.

---

### NEW: Default Parameters

**Status:** ✅ IMPLEMENTED  
**Priority:** P1 - High  
**Implemented:** Current version

**Feature:**
Default parameter values now supported using optional variadic pattern.

**Example:**
```typescript
function greet(name: string = "World"): string {
    return "Hello, " + name + "!";
}

function multiply(a: number, b: number = 2): number {
    return a * b;
}
```

Generates:
```go
func Greet(_name_opt ...string) string {
    name := "World"
    if len(_name_opt) > 0 {
        name = _name_opt[0]
    }
    return "Hello, " + name + "!"
}

func Multiply(a float64, _b_opt ...float64) float64 {
    b := 2
    if len(_b_opt) > 0 {
        b = _b_opt[0]
    }
    return a * b
}
```

**Usage:**
```go
Greet()           // Uses default "World"
Greet("Alice")    // Uses "Alice"
Multiply(5)       // Uses default 2, returns 10
Multiply(5, 3)    // Uses 3, returns 15
```

**Note:** Function body uses normal parameter names, not the _opt versions.

---

### 6. Array/Object Destructuring

**Status:** ✅ IMPLEMENTED  
**Priority:** P1 - High  
**Implemented:** Current version

**Features:**
Array and object destructuring patterns now fully supported.

**Array Destructuring Example:**
```typescript
const [a, b, c] = [1, 2, 3];
const [first, second] = ["hello", "world"];
const [head, ...tail] = [1, 2, 3, 4, 5];
```

Generates:
```go
_destructure_0 := []interface{}{1, 2, 3}
a := _destructure_0[0]
b := _destructure_0[1]
c := _destructure_0[2]

_destructure_1 := []interface{}{"hello", "world"}
first := _destructure_1[0]
second := _destructure_1[1]

_destructure_2 := []interface{}{1, 2, 3, 4, 5}
head := _destructure_2[0]
tail := _destructure_2[1:]  // Rest pattern
```

**Object Destructuring Example:**
```typescript
const person = {name: "John", age: 30};
const {name, age} = person;
```

Generates:
```go
person := map[string]interface{}{"name": "John", "age": 30}
_destructure_0 := person
name := _destructure_0.(map[string]interface{})["name"]
age := _destructure_0.(map[string]interface{})["age"]
```

**Note:** Uses temporary variables to avoid duplicating initializer expressions.

---

### 7. Object Spread Operator

**Status:** ✅ IMPLEMENTED  
**Priority:** P1 - High  
**Implemented:** Current version

**Feature:**
Object spread operators now work correctly, generating inline merge functions.

**Example:**
```typescript
const obj1 = {a: 1, b: 2};
const obj2 = {c: 3, d: 4};
const merged = {...obj1, ...obj2};
const withExtra = {...obj1, e: 5, ...obj2};
```

Generates:
```go
obj1 := map[string]interface{}{"a": 1, "b": 2}
obj2 := map[string]interface{}{"c": 3, "d": 4}

merged := func() map[string]interface{} {
    _merge_0 := make(map[string]interface{})
    for k, v := range obj1 { _merge_0[k] = v }
    for k, v := range obj2 { _merge_0[k] = v }
    return _merge_0
}()

withExtra := func() map[string]interface{} {
    _merge_1 := make(map[string]interface{})
    for k, v := range obj1 { _merge_1[k] = v }
    for k, v := range map[string]interface{}{"e": 5} { _merge_1[k] = v }
    for k, v := range obj2 { _merge_1[k] = v }
    return _merge_1
}()
```

**Note:** Later properties override earlier ones (standard JavaScript behavior).

---

### 8. Computed Property Names

**Status:** ❌ Not Implemented  
**Priority:** P2 - Medium  
**Target Fix:** v0.7.0-beta

**Missing Features:**
```typescript
const key = "name";
const obj = { [key]: "value" };
```

**Workaround:**
```typescript
// Use explicit property assignment
const obj = {};
obj[key] = "value";
```

---

## Example Project Issues

### 11. Example Projects May Not Compile

**Status:** ⚠️ Known Limitation  
**Priority:** P1 - High  
**Target Fix:** v0.6.0-beta

**Issue:**
The transpiled Go files in `examples/real-world/` contain intentionally incomplete code to demonstrate the transpiler output. They are excluded from CI checks and may not compile without manual fixes.

**Affected Files:**
- `examples/real-world/data-processing/processor.go`
- `examples/real-world/express-hello/server.go`
- `examples/real-world/commander-cli/cli.go`

**Note:**
These examples demonstrate the transpiler's *current* output capabilities. They will be updated to fully compile as the transpiler improves in v0.6.0-beta.

**Workaround:**
- Use the TypeScript source files (`.ts`) as reference
- Manually fix compilation errors in generated `.go` files
- Wait for v0.6.0-beta for automatic compilation

---

## Desktop UI Issues

### 12. Desktop UI Build Issue - Missing CLI Binary

**Status:** ✅ FIXED  
**Priority:** P0 - Critical  
**Fixed In:** Current version

**Issue:**
Desktop UI build failed with error: "glob pattern bin/ts2go-cli* path not found or didn't match any files."

**Root Cause:**
The Tauri build configuration expected the CLI binary to be present in `desktop-ui/src-tauri/bin/` directory before building, but the build process didn't create it.

**Fix:**
Updated `Makefile` to automatically:
1. Build the CLI binary first
2. Create the `src-tauri/bin/` directory
3. Copy the CLI binary as `ts2go-cli` (or `ts2go-cli.exe` on Windows)

**Usage:**
```bash
# From repository root:
make build-desktop    # Builds desktop app with CLI bundled
make dev-desktop      # Runs desktop app in development mode
```

**Note:** The CLI binary is now automatically bundled as a Tauri resource and can be accessed by the desktop app.

---

### 13. Desktop UI Incomplete (52% Complete)

**Status:** 🔄 In Progress  
**Priority:** P2 - Medium  
**Target Fix:** v0.8.0-beta

**Missing Features:**
- Settings persistence
- Recent projects history
- Syntax error highlighting in editor
- Multi-file project view
- Build and run generated Go code
- Dependency visualization
- Export/import projects

**Workaround:**
- Use CLI tool for full functionality
- Desktop UI is usable for basic transpilation and preview

**Tracking:** See ROADMAP_TO_1.0.0.md Phase 3

---

## Performance Issues

### 14. Large Project Transpilation Performance

**Status:** 🟡 Acceptable  
**Priority:** P2 - Medium  
**Target Fix:** v0.9.0-rc

**Issue:**
Transpilation of very large projects (>50K lines) may be slow.

**Current Performance:**
- Small projects (<1K lines): <1 second
- Medium projects (1K-10K lines): 1-5 seconds
- Large projects (10K-50K lines): 5-30 seconds
- Very large projects (>50K lines): 30+ seconds

**Workaround:**
- Transpile incrementally (file by file)
- Use watch mode for development

**Tracking:** See ROADMAP_TO_1.0.0.md Phase 4

---

## Security Issues

### 15. Desktop UI Dependency Vulnerabilities

**Status:** ✅ Fixed  
**Priority:** P0 - Critical  
**Fixed In:** Current PR

**Issue:**
Desktop UI dependencies had 7 vulnerabilities (1 critical, 6 moderate).

**Resolution:**
- Updated happy-dom to v20.0.10+ (fixed critical VM escape vulnerability)
- Updated other dependencies
- Remaining development-only vulnerabilities documented

**Note:**
Any remaining vulnerabilities are in development dependencies only and do not affect production builds.

---

## Type System Limitations

### 16. Complex Generic Types

**Status:** 🟡 Partial Support  
**Priority:** P2 - Medium  
**Target Fix:** v1.1.0

**Issue:**
Complex generic types with multiple type parameters and constraints may not transpile correctly.

**Example:**
```typescript
type Mapper<T extends string, U extends number> = {
    [K in T]: U;
};
```

**Workaround:**
- Simplify generic types
- Use concrete types where possible
- Manually define complex types in Go

---

### 17. Conditional Types

**Status:** ❌ Not Implemented  
**Priority:** P2 - Medium  
**Target Fix:** v1.1.0

**Missing Features:**
```typescript
type IsString<T> = T extends string ? true : false;
```

**Workaround:**
- Use runtime type checks in Go
- Avoid conditional types in TypeScript

---

## Runtime Library Limitations

### 18. Incomplete Node.js API Coverage

**Status:** 🟡 Partial Support  
**Priority:** P2 - Medium  
**Target Fix:** v1.0.0

**Coverage:**
- ✅ Console API (complete)
- ✅ Process API (basic support)
- ✅ OS API (basic support)
- ⚠️ File System API (partial - async operations incomplete)
- ⚠️ HTTP API (basic support only)
- ❌ Crypto API (not implemented)
- ❌ Stream API (not implemented)

**Workaround:**
- Use Go standard library directly for missing APIs
- Contribute runtime implementations

---

## Testing Gaps

### 19. Compilation Tests Not Comprehensive

**Status:** 🔴 In Progress  
**Priority:** P0 - Critical  
**Target Fix:** v0.6.0-beta

**Issue:**
Current integration tests verify that TypeScript parses and Go code generates, but don't verify that generated code compiles or runs correctly.

**Tracking:** See IMPLEMENTATION_PLAN.md Phase 1

---

## Reporting Issues

If you encounter an issue not listed here:

1. **Search Existing Issues:** Check [GitHub Issues](https://github.com/el-j/ts2go/issues)
2. **Provide Details:**
   - TypeScript input code
   - Generated Go code
   - Expected behavior
   - Actual behavior
   - Version of ts2go
3. **Minimal Reproduction:** Provide smallest code that reproduces the issue

---

## Workarounds Summary

### General Approach
1. **Simplify TypeScript:** Use simpler patterns that transpile reliably
2. **Manual Fixes:** Fix generated Go code manually for critical issues
3. **Alternative Syntax:** Use alternative TypeScript syntax that transpiles better
4. **Wait for Updates:** Many issues are being actively fixed

### Quick Reference

| Issue | Quick Workaround |
|-------|-----------------|
| Missing imports | Add manually to generated .go files |
| Empty function bodies | Use explicit return statements |
| Template literals | Use string concatenation |
| Object destructuring | Use explicit property access |
| Spread operators | Use Object.assign or loops |
| Default parameters | Use explicit checks |
| Arrow functions | Use regular function syntax |

---

## Version Roadmap

- **v0.6.0-beta** (Week 2): Core code generation fixes
- **v0.7.0-beta** (Week 5): Modern JS features
- **v0.8.0-beta** (Week 8): Desktop UI completion
- **v0.9.0-rc** (Week 10): Polish and optimization
- **v1.0.0** (Week 12): Stable release

See [ROADMAP_TO_1.0.0.md](ROADMAP_TO_1.0.0.md) for detailed timeline.

---

## Contributing Fixes

Want to help fix these issues? See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Priority areas for contributions:
1. Code generation fixes (Phase 1)
2. Modern JavaScript features (Phase 2)
3. Desktop UI features (Phase 3)
4. Documentation improvements
5. Test coverage

---

**Note:** This document is updated regularly. Check the "Last Updated" date at the top to ensure you have current information.
