# Implementation Plan for Critical Code Generation Fixes

**Date:** November 7, 2025  
**Status:** READY TO IMPLEMENT  
**Target:** v0.6.0-beta Release  
**Timeline:** 2-3 weeks

---

## Executive Summary

This document outlines the implementation plan to fix critical code generation issues that prevent ts2go from producing compilable Go code. These fixes are **P0 blockers** that must be completed before any release.

**Goal:** Enable ts2go to generate Go code that compiles and runs correctly for supported TypeScript features.

---

## Phase 1: Code Generation Core Fixes (Week 1)

### 1.1 Fix Function Body Generation (Days 1-2)

**Problem:** Functions generate with empty `return` statements

**Files to Modify:**
- `internal/transpiler/codegen_functions.go`
- `internal/transpiler/codegen_expressions.go`

**Tasks:**
1. Investigate AST traversal for function bodies
2. Ensure all expression types are properly transpiled
3. Fix return statements to include actual values
4. Handle implicit returns in arrow functions

**Test Cases:**
```typescript
// Test 1: Simple function
function add(a: number, b: number): number {
    return a + b;
}

// Test 2: Arrow function with expression
const double = (x: number) => x * 2;

// Test 3: Arrow function with block
const greet = (name: string) => {
    return `Hello, ${name}!`;
};
```

**Expected Go Output:**
```go
func Add(a float64, b float64) float64 {
    return a + b
}

double := func(x float64) float64 {
    return x * 2
}

greet := func(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
```

**Validation:**
- Go code compiles without errors
- Functions return correct values
- Arrow functions work for both expression and block bodies

---

### 1.2 Fix Import Generation (Days 2-3)

**Problem:** Generated Go files missing required imports

**Files to Modify:**
- `internal/transpiler/codegen.go`
- `internal/transpiler/codegen_core.go`

**Tasks:**
1. Track all Go packages used in generated code
2. Automatically add import statements based on usage
3. Handle standard library imports (fmt, strings, etc.)
4. Handle runtime library imports

**Import Tracking:**
```go
type ImportTracker struct {
    imports map[string]bool
}

// Track when fmt.Println is used
func (t *ImportTracker) UseFmt() {
    t.imports["fmt"] = true
}

// Track when strings functions are used
func (t *ImportTracker) UseStrings() {
    t.imports["strings"] = true
}
```

**Test Cases:**
```typescript
// Should generate: import "fmt"
console.log("Hello");

// Should generate: import "strings"
const upper = name.toUpperCase();

// Should generate: import "github.com/el-j/ts2go/runtime/console"
console.error("Error:", err);
```

**Validation:**
- All required imports are present
- No unused imports
- Import paths are correct
- Code compiles without "undefined" errors

---

### 1.3 Fix Template Literal Transpilation (Day 3)

**Problem:** Template literals not fully converted to Go string formatting

**Files to Modify:**
- `internal/transpiler/codegen_expressions.go`

**Tasks:**
1. Parse template literal expressions
2. Convert to `fmt.Sprintf()` calls
3. Handle nested expressions
4. Escape special characters

**Implementation:**
```go
func (cg *CodeGenerator) transpileTemplateLiteral(node *TemplateNode) string {
    // "Hello, ${name}!" → fmt.Sprintf("Hello, %s!", name)
    format := ""
    args := []string{}
    
    for _, part := range node.Parts {
        if part.IsExpression {
            format += "%v"
            args = append(args, cg.transpileExpression(part.Expr))
        } else {
            format += part.Text
        }
    }
    
    if len(args) > 0 {
        return fmt.Sprintf("fmt.Sprintf(\"%s\", %s)", format, strings.Join(args, ", "))
    }
    return fmt.Sprintf("\"%s\"", format)
}
```

**Test Cases:**
```typescript
const name = "World";
const msg = `Hello, ${name}!`;
const multi = `Name: ${user.name}, Age: ${user.age}`;
```

**Expected Go Output:**
```go
name := "World"
msg := fmt.Sprintf("Hello, %s!", name)
multi := fmt.Sprintf("Name: %s, Age: %v", user.Name, user.Age)
```

**Validation:**
- Template literals compile correctly
- String interpolation works
- Special characters are escaped
- Nested expressions work

---

### 1.4 Fix Arrow Function Syntax (Days 4-5)

**Problem:** Arrow functions generate invalid Go syntax

**Files to Modify:**
- `internal/transpiler/codegen_functions.go`
- `internal/transpiler/codegen_expressions.go`

**Tasks:**
1. Fix anonymous function syntax
2. Handle expression bodies correctly
3. Handle block bodies correctly
4. Ensure proper variable scoping

**Current Issue:**
```typescript
const add = (a: number, b: number) => a + b;
```

Generates (BROKEN):
```go
add := func(a float64, b float64) float64 { return }
```

**Should Generate:**
```go
add := func(a float64, b float64) float64 {
    return a + b
}
```

**Implementation Strategy:**
```go
func (cg *CodeGenerator) transpileArrowFunction(node *ArrowFunctionNode) string {
    params := cg.transpileParameters(node.Parameters)
    returnType := cg.transpileType(node.ReturnType)
    
    var body string
    if node.IsExpression {
        // Expression body: (x) => x * 2
        expr := cg.transpileExpression(node.Body)
        body = fmt.Sprintf("{\n    return %s\n}", expr)
    } else {
        // Block body: (x) => { return x * 2; }
        body = cg.transpileBlock(node.Body)
    }
    
    return fmt.Sprintf("func(%s) %s %s", params, returnType, body)
}
```

**Test Cases:**
```typescript
// Expression body
const square = (x: number) => x * x;

// Block body
const greet = (name: string) => {
    const msg = `Hello, ${name}`;
    return msg;
};

// Multiple parameters
const calc = (a: number, b: number, op: string) => {
    if (op === "+") return a + b;
    if (op === "-") return a - b;
    return 0;
};
```

**Validation:**
- All arrow function types compile
- Return values are correct
- Variable scoping works
- Type inference works

---

## Phase 2: Additional Syntax Fixes (Week 1-2)

### 2.1 Fix Object Literal Syntax (Day 6)

**Problem:** Object literals generate invalid Go syntax

**Files to Modify:**
- `internal/transpiler/codegen_expressions.go`

**Tasks:**
1. Fix struct literal syntax
2. Handle map literals correctly
3. Fix array of objects syntax

**Test Cases:**
```typescript
const user = { name: "Alice", age: 30 };
const users: User[] = [
    { name: "Alice", age: 30 },
    { name: "Bob", age: 25 }
];
```

**Expected Go Output:**
```go
user := User{Name: "Alice", Age: 30}
users := []User{
    {Name: "Alice", Age: 30},
    {Name: "Bob", Age: 25},
}
```

**Validation:**
- Object literals compile
- Type inference works
- Arrays of objects work
- Nested objects work

---

### 2.2 Fix Class Inheritance (Day 7)

**Problem:** `super()` calls generate `/* unsupported expression */`

**Files to Modify:**
- `internal/transpiler/codegen_classes.go`

**Tasks:**
1. Fix `super()` constructor calls
2. Fix `super.method()` calls
3. Ensure proper initialization order

**Test Cases:**
```typescript
class Base {
    id: number;
    constructor(id: number) {
        this.id = id;
    }
}

class Derived extends Base {
    name: string;
    constructor(id: number, name: string) {
        super(id);
        this.name = name;
    }
}
```

**Expected Go Output:**
```go
type Base struct {
    Id float64
}

func NewBase(id float64) *Base {
    instance := &Base{}
    instance.Id = id
    return instance
}

type Derived struct {
    Base
    Name string
}

func NewDerived(id float64, name string) *Derived {
    instance := &Derived{}
    instance.Base = *NewBase(id)
    instance.Name = name
    return instance
}
```

**Validation:**
- Class inheritance compiles
- Super calls work
- Method inheritance works
- Field initialization correct

---

### 2.3 Fix Anonymous Function Syntax (Day 8)

**Problem:** Callback functions generate invalid syntax

**Files to Modify:**
- `internal/transpiler/codegen_expressions.go`
- `internal/transpiler/codegen_functions.go`

**Test Cases:**
```typescript
app.get("/", function(req, res) {
    res.send("Hello");
});

app.get("/user", (req, res) => {
    res.json({ user: "John" });
});
```

**Expected Go Output:**
```go
app.Get("/", func(req interface{}, res interface{}) {
    res.Send("Hello")
})

app.Get("/user", func(req interface{}, res interface{}) {
    res.Json(map[string]interface{}{"user": "John"})
})
```

**Validation:**
- Callback functions compile
- Parameter types correct
- Function bodies complete

---

## Phase 3: Testing & Validation (Week 2)

### 3.1 Add Compilation Tests (Days 9-10)

**Files to Create/Modify:**
- `tests/compilation_test.go` (NEW)
- `tests/integration_test.go` (MODIFY)

**Tasks:**
1. Create new test that verifies generated code compiles
2. Update existing integration tests
3. Test all example projects
4. Add to CI/CD pipeline

**Implementation:**
```go
func TestGeneratedCodeCompiles(t *testing.T) {
    testCases := []struct {
        name       string
        typescript string
    }{
        {
            name: "simple function",
            typescript: `
                function add(a: number, b: number): number {
                    return a + b;
                }
            `,
        },
        {
            name: "arrow function",
            typescript: `
                const double = (x: number) => x * 2;
            `,
        },
        // ... more test cases
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Generate Go code
            goCode := transpile(tc.typescript)
            
            // Write to temp file
            tmpFile := writeTempFile(t, goCode)
            defer os.Remove(tmpFile)
            
            // Try to compile
            cmd := exec.Command("go", "build", tmpFile)
            output, err := cmd.CombinedOutput()
            
            if err != nil {
                t.Errorf("Compilation failed:\n%s\n\nGenerated code:\n%s", 
                    output, goCode)
            }
        })
    }
}
```

**Test Coverage:**
- All basic types
- All control flow structures
- Functions (regular and arrow)
- Classes with inheritance
- Template literals
- Object literals
- Async/await patterns

**Validation:**
- All test cases pass
- Generated code compiles
- No compilation errors
- Tests run in CI/CD

---

### 3.2 Test Real-World Examples (Days 11-12)

**Examples to Test:**
1. `examples/real-world/express-hello/`
2. `examples/real-world/data-processing/`
3. `examples/real-world/commander-cli/`

**Tasks:**
1. Transpile each example
2. Verify compilation
3. Run generated binaries
4. Compare output with expected behavior
5. Document any limitations

**Testing Script:**
```bash
#!/bin/bash

for example in examples/real-world/*/; do
    echo "Testing: $example"
    
    # Transpile
    ./ts2go transpile "$example" --out /tmp/test-output
    
    # Build
    cd /tmp/test-output
    go mod tidy
    go build
    
    if [ $? -eq 0 ]; then
        echo "✅ $example: BUILD SUCCESS"
    else
        echo "❌ $example: BUILD FAILED"
        exit 1
    fi
done
```

**Validation:**
- All examples transpile
- All examples compile
- Binaries run successfully
- Output matches expectations (where applicable)

---

### 3.3 Update Integration Tests (Day 13)

**Files to Modify:**
- `tests/integration_test.go`

**Tasks:**
1. Add compilation verification to existing tests
2. Add new test cases for fixed issues
3. Ensure tests fail when code doesn't compile
4. Update test expectations

**Changes:**
```go
func TestTranspileSimple(t *testing.T) {
    // ... existing setup ...
    
    // Generate Go code
    goCode := transpile(input)
    
    // NEW: Verify it compiles
    if err := verifyCompilation(goCode); err != nil {
        t.Fatalf("Generated code does not compile: %v\n%s", err, goCode)
    }
    
    // ... existing assertions ...
}
```

**Validation:**
- All integration tests pass
- Tests verify compilation
- Test coverage increased
- No regression in existing functionality

---

## Phase 4: Documentation & Release Prep (Week 2-3)

### 4.1 Update Documentation (Days 14-15)

**Files to Update:**
- `CHANGELOG.md`
- `STATUS.md`
- `README.md`
- `KNOWN_ISSUES.md` (NEW)

**Tasks:**
1. Document all fixes in CHANGELOG
2. Update STATUS.md with accurate completion percentages
3. Create KNOWN_ISSUES.md with remaining limitations
4. Update README with what's now supported

**CHANGELOG Entry:**
```markdown
## [0.6.0-beta] - 2025-11-XX

### Fixed
- Function bodies now generate complete code
- Import statements automatically added
- Template literals properly transpiled to fmt.Sprintf
- Arrow functions generate valid Go syntax
- Object literals generate correct struct syntax
- Class inheritance super() calls now work
- Anonymous functions (callbacks) now compile

### Added
- Compilation tests to verify generated code
- Real-world example testing
- KNOWN_ISSUES.md documenting limitations

### Changed
- Integration tests now verify compilation
- Improved error messages for unsupported features
```

**KNOWN_ISSUES.md Structure:**
```markdown
# Known Issues and Limitations

## Not Yet Implemented
- Array destructuring
- Object destructuring
- Spread operator for objects
- Rest parameters
- Default parameters

## Known Bugs
- Complex generic types may not transpile correctly
- Some edge cases in async/await may not work

## Workarounds
...
```

**Validation:**
- Documentation is accurate
- Known issues clearly documented
- Examples are up to date
- Release notes complete

---

### 4.2 Performance Testing (Day 16)

**Tasks:**
1. Benchmark transpilation speed
2. Test memory usage
3. Test large project transpilation
4. Identify and fix performance bottlenecks

**Benchmarks to Create:**
```go
func BenchmarkTranspileSimple(b *testing.B) {
    input := readTestFile("simple.ts")
    for i := 0; i < b.N; i++ {
        transpile(input)
    }
}

func BenchmarkTranspileLarge(b *testing.B) {
    input := readTestFile("large-project.ts")
    for i := 0; i < b.N; i++ {
        transpile(input)
    }
}
```

**Validation:**
- Transpilation is reasonably fast
- Memory usage is acceptable
- No memory leaks
- Large projects work

---

### 4.3 Beta Release Preparation (Days 17-18)

**Tasks:**
1. Final code review
2. Run all tests
3. Update version to 0.6.0-beta
4. Build release artifacts
5. Test release binaries
6. Prepare release notes

**Pre-Release Checklist:**
- [ ] All P0 issues fixed
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Examples work
- [ ] Security issues addressed
- [ ] Version updated
- [ ] Release notes prepared
- [ ] Artifacts built and tested

**Release Artifacts:**
- CLI binaries (Linux, macOS, Windows)
- Desktop app (macOS, Linux, Windows)
- Docker image
- Source tarball
- Checksums

**Validation:**
- All artifacts build successfully
- Binaries work on target platforms
- Installation instructions accurate
- Release notes complete

---

## Implementation Schedule

### Week 1: Core Fixes
```
Day 1-2:  Function body generation
Day 3:    Import generation + Template literals
Day 4-5:  Arrow function syntax
Day 6:    Object literal syntax
Day 7:    Class inheritance
Day 8:    Anonymous functions
```

### Week 2: Testing & Validation
```
Day 9-10:  Compilation tests
Day 11-12: Real-world examples
Day 13:    Integration tests update
Day 14-15: Documentation
Day 16:    Performance testing
```

### Week 3: Release Prep
```
Day 17-18: Beta release preparation
Day 19:    Beta release
Day 20-21: Bug fixes from early testing
```

---

## Risk Management

### High Risk Items
1. **Function body generation** - Most complex, highest impact
2. **Import tracking** - Must be comprehensive
3. **Arrow functions** - Multiple edge cases

**Mitigation:**
- Start with high-risk items first
- Incremental testing after each change
- Keep changes small and focused
- Extensive test coverage

### Rollback Plan
If a fix breaks existing functionality:
1. Revert the specific commit
2. Add test case for the regression
3. Fix the issue
4. Re-run all tests before committing

---

## Success Metrics

### Code Quality
- [ ] All generated code compiles
- [ ] 95%+ of test cases pass
- [ ] No new compilation errors
- [ ] No regression in existing functionality

### Testing
- [ ] 100+ compilation test cases
- [ ] All real-world examples work
- [ ] Integration tests verify compilation
- [ ] CI/CD includes compilation checks

### Documentation
- [ ] CHANGELOG complete
- [ ] Known issues documented
- [ ] Examples updated
- [ ] Migration guide available

### Release
- [ ] Beta artifacts built
- [ ] Installation tested
- [ ] Community notified
- [ ] Feedback mechanism in place

---

## Dependencies

### Tools Required
- Go 1.24+
- Node.js 20+
- TypeScript compiler
- Git

### External Dependencies
- No new Go dependencies needed
- Desktop UI dependencies already updated

---

## Team Assignments

### Code Generation Fixes
- Primary: Core transpiler team
- Reviewer: Lead architect
- Estimated: 8-10 days

### Testing
- Primary: QA team
- Reviewer: Core transpiler team
- Estimated: 5-7 days

### Documentation
- Primary: Documentation team
- Reviewer: Product manager
- Estimated: 2-3 days

### Release Management
- Primary: DevOps team
- Reviewer: Project lead
- Estimated: 2-3 days

---

## Communication Plan

### Daily Standups
- Progress updates
- Blockers discussion
- Next steps planning

### Weekly Reviews
- Demo of completed features
- Test results review
- Documentation review

### Beta Release Announcement
- GitHub release
- Email to early adopters
- Social media announcement
- Discord/Slack notification

---

## Post-Beta Plan

### Week 4-5: Beta Feedback
- Monitor issues
- Fix critical bugs
- Gather user feedback
- Plan for next features

### Week 6-8: Modern JS Features
- Implement destructuring
- Implement spread/rest operators
- Implement default parameters

### Month 3-4: v1.0 Preparation
- Complete desktop UI
- Performance optimization
- Comprehensive testing
- Community feedback integration

---

## Appendix: Code Files Reference

### Files to Modify (Priority Order)

1. **P0 - Critical:**
   - `internal/transpiler/codegen_functions.go`
   - `internal/transpiler/codegen_expressions.go`
   - `internal/transpiler/codegen.go`
   - `internal/transpiler/codegen_core.go`

2. **P1 - High Priority:**
   - `internal/transpiler/codegen_classes.go`
   - `internal/transpiler/codegen_statements.go`
   - `tests/integration_test.go`

3. **P2 - Medium Priority:**
   - `tests/compilation_test.go` (NEW)
   - `CHANGELOG.md`
   - `STATUS.md`
   - `KNOWN_ISSUES.md` (NEW)

### Test Files to Create

1. `tests/compilation_test.go` - New compilation verification tests
2. `tests/fixtures/compilation/` - New test fixtures
3. `scripts/test-examples.sh` - Example testing script

---

**Document Version:** 1.0  
**Last Updated:** November 7, 2025  
**Next Review:** After Week 1 completion
