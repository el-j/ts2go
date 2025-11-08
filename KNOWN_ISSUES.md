# Known Issues and Limitations

**Last Updated:** November 8, 2025  
**Version:** 0.5.1

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

## Modern JavaScript Features (Not Yet Implemented)

### 6. Array/Object Destructuring

**Status:** ❌ Not Implemented  
**Priority:** P1 - High  
**Target Fix:** v0.7.0-beta

**Missing Features:**
```typescript
// Array destructuring
const [a, b, c] = [1, 2, 3];
const [first, ...rest] = array;

// Object destructuring
const {name, age} = person;
const {x, y, ...rest} = point;
```

**Workaround:**
```typescript
// Use explicit assignments
const a = array[0];
const b = array[1];
const c = array[2];

const name = person.name;
const age = person.age;
```

**Tracking:** See ROADMAP_TO_1.0.0.md Phase 2

---

### 7. Spread Operator for Objects

**Status:** ❌ Not Implemented  
**Priority:** P1 - High  
**Target Fix:** v0.7.0-beta

**Missing Features:**
```typescript
const merged = {...obj1, ...obj2, extra: value};
```

**Workaround:**
```typescript
// Use Object.assign or manual property copying
const merged = Object.assign({}, obj1, obj2, { extra: value });
```

**Tracking:** See ROADMAP_TO_1.0.0.md Phase 2

---

### 8. Rest Parameters

**Status:** ❌ Not Implemented  
**Priority:** P1 - High  
**Target Fix:** v0.7.0-beta

**Missing Features:**
```typescript
function sum(...numbers: number[]): number {
    return numbers.reduce((a, b) => a + b, 0);
}
```

**Workaround:**
```typescript
// Use explicit array parameter
function sum(numbers: number[]): number {
    return numbers.reduce((a, b) => a + b, 0);
}
```

**Tracking:** See ROADMAP_TO_1.0.0.md Phase 2

---

### 9. Default Parameters

**Status:** ❌ Not Implemented  
**Priority:** P1 - High  
**Target Fix:** v0.7.0-beta

**Missing Features:**
```typescript
function greet(name: string = "World"): string {
    return `Hello, ${name}!`;
}
```

**Workaround:**
```typescript
// Use function overloading or explicit checks
function greet(name?: string): string {
    const actualName = name || "World";
    return `Hello, ${actualName}!`;
}
```

**Tracking:** See ROADMAP_TO_1.0.0.md Phase 2

---

### 10. Computed Property Names

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

### 12. Desktop UI Incomplete (52% Complete)

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

### 13. Large Project Transpilation Performance

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

### 14. Desktop UI Dependency Vulnerabilities

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

### 15. Complex Generic Types

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

### 16. Conditional Types

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

### 17. Incomplete Node.js API Coverage

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

### 18. Compilation Tests Not Comprehensive

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
