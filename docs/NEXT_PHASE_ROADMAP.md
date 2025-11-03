# Next Implementation Phase: Roadmap & Priority Plan

**Last Updated:** November 3, 2025  
**Based on:** ROADMAP.md, STATUS.md, and project analysis  
**Timeline:** 8-12 weeks to significantly improve coverage

---

## 🎯 Executive Summary

### Current State
- **Working:** Type system (interfaces, aliases, enums, unions, tuples)
- **Working:** Basic logic (functions, variables, expressions)
- **Working:** Classes with full OOP support
- **Working:** Dependency management and npm package mapping
- **Working:** Desktop UI (Phase 21 - 52% complete)
- **Coverage:** ~15-20% of typical TypeScript codebases

### Critical Gaps (Must Fix)
1. ❌ **Control Flow** - No if/else, no loops, no switch statements
2. ❌ **Modern JavaScript** - No arrow functions, no template literals, no destructuring
3. ❌ **Async/Await** - No Promise support
4. ❌ **Error Handling** - No try/catch support

---

## 📋 Phase 22: Control Flow Statements (CRITICAL - Week 1-2)

**Priority:** HIGHEST  
**Impact:** Moves coverage from 20% to 50%+  
**Status:** Not Started

### Why This is Critical
Without control flow, you cannot transpile even the simplest real-world application. This is the biggest blocker to adoption.

### Implementation Tasks

#### 1. If/Else Statements
```typescript
// TypeScript
if (x > 0) {
    console.log("positive");
} else if (x < 0) {
    console.log("negative");
} else {
    console.log("zero");
}
```
```go
// Go
if x > 0 {
    fmt.Println("positive")
} else if x < 0 {
    fmt.Println("negative")
} else {
    fmt.Println("zero")
}
```

**Tasks:**
- [ ] Add AST handler for `ts.SyntaxKind.IfStatement`
- [ ] Handle condition expressions
- [ ] Handle else/else-if chains
- [ ] Add tests for nested if statements
- [ ] Add tests for if with complex conditions

#### 2. For Loops
```typescript
// For loop
for (let i = 0; i < 10; i++) {
    console.log(i);
}

// For-of loop
for (const item of array) {
    console.log(item);
}

// For-in loop
for (const key in object) {
    console.log(key, object[key]);
}
```
```go
// For loop
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// For-range loop (for-of)
for _, item := range array {
    fmt.Println(item)
}

// For-range loop (for-in)
for key := range object {
    fmt.Println(key, object[key])
}
```

**Tasks:**
- [ ] Add handler for `ts.SyntaxKind.ForStatement`
- [ ] Add handler for `ts.SyntaxKind.ForOfStatement`
- [ ] Add handler for `ts.SyntaxKind.ForInStatement`
- [ ] Handle loop initialization, condition, increment
- [ ] Add tests for nested loops
- [ ] Add tests for break/continue

#### 3. While Loops
```typescript
while (condition) {
    // body
}

do {
    // body
} while (condition);
```
```go
for condition {
    // body
}

// Do-while requires special handling
for {
    // body
    if !condition {
        break
    }
}
```

**Tasks:**
- [ ] Add handler for `ts.SyntaxKind.WhileStatement`
- [ ] Add handler for `ts.SyntaxKind.DoStatement`
- [ ] Add tests for various while conditions

#### 4. Switch Statements
```typescript
switch (value) {
    case 1:
        console.log("one");
        break;
    case 2:
        console.log("two");
        break;
    default:
        console.log("other");
}
```
```go
switch value {
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("two")
default:
    fmt.Println("other")
}
```

**Tasks:**
- [ ] Add handler for `ts.SyntaxKind.SwitchStatement`
- [ ] Handle case clauses
- [ ] Handle default clause
- [ ] Handle fall-through behavior differences
- [ ] Add tests for switch with various types

**Deliverables:**
- ✅ All control flow AST handlers implemented
- ✅ 100+ tests covering control flow scenarios
- ✅ Documentation with before/after examples
- ✅ Updated STATUS.md showing new coverage

---

## 📋 Phase 23: Modern JavaScript Syntax (HIGH - Week 3-4)

**Priority:** HIGH  
**Impact:** Moves coverage from 50% to 70%+  
**Status:** Not Started

### Implementation Tasks

#### 1. Arrow Functions
```typescript
const add = (a: number, b: number): number => a + b;
const greet = (name: string) => {
    console.log("Hello " + name);
};
```
```go
add := func(a float64, b float64) float64 { return a + b }
greet := func(name string) {
    fmt.Println("Hello " + name)
}
```

**Tasks:**
- [ ] Add handler for `ts.SyntaxKind.ArrowFunction`
- [ ] Handle implicit returns
- [ ] Handle explicit block bodies
- [ ] Add tests for various arrow function forms

#### 2. Template Literals
```typescript
const name = "World";
const message = `Hello, ${name}!`;
const multiline = `
    Line 1
    Line 2
`;
```
```go
name := "World"
message := fmt.Sprintf("Hello, %s!", name)
multiline := "\n    Line 1\n    Line 2\n"
```

**Tasks:**
- [ ] Add handler for `ts.SyntaxKind.TemplateExpression`
- [ ] Convert to `fmt.Sprintf`
- [ ] Handle nested expressions
- [ ] Handle multiline strings
- [ ] Add tests for complex templates

#### 3. Destructuring
```typescript
// Array destructuring
const [a, b, ...rest] = array;

// Object destructuring
const { name, age } = person;
const { x: newX, y: newY } = point;
```
```go
// Array destructuring
a := array[0]
b := array[1]
rest := array[2:]

// Object destructuring
name := person.Name
age := person.Age
newX := point.X
newY := point.Y
```

**Tasks:**
- [ ] Add handler for array destructuring
- [ ] Add handler for object destructuring
- [ ] Handle rest/spread operators
- [ ] Handle nested destructuring
- [ ] Handle default values
- [ ] Add comprehensive tests

#### 4. Spread/Rest Operators
```typescript
const arr2 = [...arr1, 4, 5];
const obj2 = { ...obj1, extra: "value" };
function sum(...numbers: number[]) { }
```
```go
arr2 := append(append([]int{}, arr1...), 4, 5)
// Object spread requires careful struct handling
func sum(numbers ...float64) { }
```

**Tasks:**
- [ ] Add handler for spread in arrays
- [ ] Add handler for spread in objects
- [ ] Add handler for rest parameters
- [ ] Add tests for various spread/rest scenarios

**Deliverables:**
- ✅ Modern syntax support implemented
- ✅ 80+ tests covering modern syntax
- ✅ Documentation updated
- ✅ Coverage reaches 70%+

---

## 📋 Phase 24: Async/Await & Promises (HIGH - Week 5-6)

**Priority:** HIGH  
**Impact:** Required for most Node.js applications  
**Status:** Not Started

### Challenge
Go doesn't have promises natively. We need to map async/await to goroutines and channels.

### Proposed Strategy

#### Option A: Goroutines + Channels (Recommended)
```typescript
async function fetchData(): Promise<string> {
    const result = await http.get("/api/data");
    return result.data;
}
```
```go
func FetchData() <-chan Result[string] {
    ch := make(chan Result[string], 1)
    go func() {
        result := http.Get("/api/data")
        if result.Err != nil {
            ch <- Result[string]{Err: result.Err}
        } else {
            ch <- Result[string]{Val: result.Data, Err: nil}
        }
        close(ch)
    }()
    return ch
}

// Usage
resultCh := FetchData()
result := <-resultCh
if result.Err != nil {
    // handle error
}
```

**Tasks:**
- [ ] Define `Result[T]` generic type for async returns
- [ ] Add handler for `async` function declarations
- [ ] Convert `Promise<T>` to `<-chan Result[T]`
- [ ] Add handler for `await` expressions
- [ ] Convert `await expr` to `result := <-expr; check result.Err`
- [ ] Add tests for async/await patterns
- [ ] Document async/await mapping strategy

#### Option B: Synchronous (Simpler, but loses concurrency)
```typescript
async function fetchData(): Promise<string> {
    return http.get("/api/data");
}
```
```go
func FetchData() (string, error) {
    return http.Get("/api/data")
}
```

**Tasks:**
- [ ] Convert `async` functions to regular functions
- [ ] Convert `Promise<T>` to `(T, error)`
- [ ] Convert `await` to direct calls with error checking
- [ ] Add tests
- [ ] Document limitations

**Recommendation:** Start with Option B for MVP, add Option A later for better concurrency.

**Deliverables:**
- ✅ Async/await support (at least basic)
- ✅ 50+ tests for async patterns
- ✅ Documentation explaining mapping
- ✅ Examples showing before/after

---

## 📋 Phase 25: Error Handling (MEDIUM - Week 7-8)

**Priority:** MEDIUM  
**Impact:** Required for production code  
**Status:** Not Started

### Implementation Tasks

#### Try/Catch/Finally
```typescript
try {
    const result = riskyOperation();
    console.log(result);
} catch (error) {
    console.error("Error:", error);
} finally {
    cleanup();
}
```
```go
result, err := riskyOperation()
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
} else {
    fmt.Println(result)
}
cleanup()
```

**Tasks:**
- [ ] Add handler for `ts.SyntaxKind.TryStatement`
- [ ] Convert try block to regular code with error checking
- [ ] Convert catch to error handling
- [ ] Convert finally to defer statements
- [ ] Add tests for try/catch/finally
- [ ] Handle nested try/catch

#### Throw Statements
```typescript
throw new Error("Something went wrong");
```
```go
return nil, fmt.Errorf("something went wrong")
```

**Tasks:**
- [ ] Add handler for `ts.SyntaxKind.ThrowStatement`
- [ ] Convert to early return with error
- [ ] Add tests for various throw scenarios

**Deliverables:**
- ✅ Try/catch/finally support
- ✅ Throw statement support
- ✅ 40+ tests for error handling
- ✅ Documentation with examples
- ✅ Coverage reaches 80%+

---

## 📋 Phase 26: Desktop UI Completion (MEDIUM - Week 9-10)

**Priority:** MEDIUM  
**Current:** 52% complete  
**Status:** In Progress

### Remaining Work (from PHASE21_TODO.md)

#### Week 3: File Tree & Advanced Features (0% Complete)
- [ ] File tree component with icons
- [ ] Multi-file transpilation
- [ ] Real-time file watcher
- [ ] Settings panel
- [ ] Terminal/logs panel
- [ ] Rust backend commands integration

#### Week 4: Polish & Release (0% Complete)
- [ ] Error handling improvements
- [ ] Performance optimization
- [ ] Documentation
- [ ] Tauri multi-platform builds
- [ ] Release automation

**Deliverables:**
- ✅ Desktop UI 100% complete
- ✅ Multi-platform builds (Windows, macOS, Linux)
- ✅ Published releases on GitHub

---

## 📋 Phase 27: Production Readiness (LOW - Week 11-12)

**Priority:** LOW  
**Status:** Not Started

### Tasks
- [ ] Performance benchmarks
- [ ] Memory profiling and optimization
- [ ] Comprehensive integration tests
- [ ] Documentation improvements
- [ ] Migration guides
- [ ] Video tutorials
- [ ] Community examples
- [ ] Plugin system (future)

---

## 🎯 Success Metrics

### Coverage Goals
- **After Phase 22 (Control Flow):** 50% of backend TypeScript code
- **After Phase 23 (Modern JS):** 70% of backend TypeScript code
- **After Phase 24 (Async):** 80% of backend TypeScript code
- **After Phase 25 (Errors):** 85% of backend TypeScript code

### Quality Goals
- ✅ All tests passing
- ✅ 80%+ code coverage
- ✅ Zero critical bugs
- ✅ Documentation complete
- ✅ CI/CD fully automated

---

## 🚀 Getting Started with Next Phase

### For Phase 22 (Control Flow)
1. Read the control flow implementation guide
2. Start with if/else statements (simplest)
3. Add comprehensive tests
4. Move to for loops
5. Add while loops
6. Finish with switch statements

### Development Workflow
1. Create feature branch: `git checkout -b feature/phase22-control-flow`
2. Implement AST handlers in `internal/transpiler/`
3. Add tests in `tests/`
4. Update documentation in `docs/`
5. Run full test suite: `make test`
6. Create PR to `develop` branch
7. After review, merge to `develop`
8. After stabilization, merge to `main`

---

## 📚 Resources

### Key Files
- `internal/transpiler/transpiler.go` - Main transpiler logic
- `internal/transpiler/emitter.go` - Go code generation
- `tests/` - Integration tests
- `docs/EXAMPLES.md` - Before/after examples

### External Resources
- [TypeScript AST Viewer](https://ts-ast-viewer.com/)
- [Go AST Documentation](https://pkg.go.dev/go/ast)
- [TypeScript Compiler API](https://github.com/microsoft/TypeScript/wiki/Using-the-Compiler-API)

---

## 🤝 Contributing

If you want to contribute to any of these phases:
1. Check the priority and status
2. Look for open issues or create one
3. Discuss approach in the issue
4. Create feature branch and implement
5. Submit PR with tests and documentation

---

**Questions?** Open an issue or discussion on GitHub.
