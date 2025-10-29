# Project Plan: "TS2Go" TypeScript-to-Go Transpiler

## 1. Project Mission

To create a toolchain that transpiles a subset of TypeScript modules into idiomatic, efficient, and "slim" Go modules. The initial goal is a "cold-compile" tool, not a hot-reload development server.

## 2. Core Challenges & "Translation Philosophy"

This is the most critical section. We must decide how to map "apples" to "oranges."

### Dynamic vs. Static Types:

- **TS:** `any`, `unknown`, union types, structural "duck" typing.
- **Go:** Strict, nominal, static typing. `interface{}` is the "slow" equivalent of `any`.
- **Decision:** We must disallow `any` and `unknown`. All types must be explicitly defined. Union types will be a major challenge, potentially requiring Go interface implementations or structs with multiple fields (e.g., `type StringOrNumber struct { StringVal string; NumVal float64; IsNum bool }`).

### Execution Model:

- **TS (Node.js):** Single-threaded, event-loop-based concurrency. Promises and async/await are the standard.
- **Go:** Multi-threaded, goroutine-based (CSP) concurrency. Channels and waitgroups are the standard.
- **Decision:** We must map async/await to Go's concurrency model.
  - An async function returning a `Promise<T>` could be transpiled to a Go function returning a `chan T` (or `chan Result{ Val T; Err error }`).
  - An `await` would become a blocking read from that channel: `val := <-myFunc()`.

### Data Structures:

- **TS:** `class` (syntactic sugar for prototypes), dynamic objects (properties can be added/removed).
- **Go:** `struct`. Fields are fixed at compile time.
- **Decision:**
  - TS `class` and `interface` will be transpiled to Go `struct`.
  - Dynamic objects (`{ [key: string]: any }`) must be transpiled to `map[string]interface{}`. This will be a significant performance/type-safety trade-off and should be discouraged.

### Standard Library & Ecosystem:

- **TS:** Relies on Node.js APIs (`fs`, `http`, `path`) and npm modules.
- **Go:** Has its own standard library (`os`, `net/http`, `path/filepath`) and Go modules.
- **Decision:** We cannot transpile `node_modules`. Instead, we must create a Go-based runtime library that provides an adapter layer.
  - TS `import { readFile } from 'fs'` will be mapped to `import "github.com/ts2go/runtime/fs"` and `fs.ReadFile(...)`.
  - This `ts2go/runtime` will be a Go module we must also build, which implements the Node.js API signatures (or a subset) using Go's standard library.

## 3. Phased Development Plan

### Phase 1: Research, Specification, & Foundational Tooling

*(Goal: Define the exact subset of TS we will support and set up the compiler architecture.)*

#### Specification Document (CRITICAL):

Write a `SPEC.md` that explicitly lists supported and unsupported TS features.

- **Supported (Example MVP):** `string`, `number`, `boolean`, `Array<T>`, `interface`, `type`, `function`, `if/else`, `for` loops, `return`, `console.log`.
- **Unsupported (Example MVP):** `class`, `any`, `unknown`, union types, `Promise`, `async/await`, `try/catch` (map to Go's error returns), dynamic properties, import of `node_modules`.

#### Compiler Architecture:

- **Parser:** Use the official TypeScript Compiler API (`tsc`). This gives us the Abstract Syntax Tree (AST) for free. We don't need to write our own parser.
- **Transformer:** This is the core of our tool. A Go (or TS) program that recursively walks the TS AST.
- **Emitter (Code Generator):** A module that generates Go code.
- **Recommendation:** Use Go's `go/ast` and `go/format` packages. This allows us to build a Go AST in memory and then format it correctly, which is much safer than string concatenation.

#### Setup Initial CLI Tool:

Create a simple CLI (e.g., using `cobra` in Go or `commander` in TS) that takes an input file.

```
ts2go --in ./src/index.ts --out ./go-dist/
```

This tool will:

- Run the TS Compiler API to get the AST of `index.ts`.
- Pipe this AST (as JSON) to stdout.

### Phase 2: MVP - Type & Data Structure Transpilation

*(Goal: Transpile only type definitions, not logic.)*

- **AST Walker:** Build the "Transformer" that walks the AST.
- **Type Emitter:** Implement handlers for TS AST nodes:
  - `ts.SyntaxKind.InterfaceDeclaration` → `type MyInterface struct { ... }`
  - `ts.SyntaxKind.TypeAliasDeclaration` → `type MyType = ...`
  - `ts.SyntaxKind.StringKeyword` → `string`
  - `ts.SyntaxKind.NumberKeyword` → `float64` (Safer default than `int`).
  - `ts.SyntaxKind.BooleanKeyword` → `bool`
  - `ts.SyntaxKind.ArrayType` → `[]MyType`
  - `ts.SyntaxKind.TypeLiteral` (for inline objects) → `struct { ... }`
  - `ts.SyntaxKind.PropertySignature` (for interface fields) → `MyField string \`json:"myField"\`` (auto-add JSON tags for JS compatibility).
- **CLI Update:** The tool now parses the AST and emits a `.go` file with Go type definitions.
- **Module Generation:** The tool also generates a `go.mod` file in the output directory (e.g., `module my-ts-project`).

### Phase 3: Core Logic Transpilation

*(Goal: Transpile basic functions, variables, and control flow.)*

- **Extend AST Walker:** Add handlers for logic nodes.
  - `ts.SyntaxKind.FunctionDeclaration` → `func MyFunction(...) ...` (Note: camelCase to PascalCase for public Go functions).
  - `ts.SyntaxKind.VariableStatement` → `var myVar = ...` or `myVar := ...`
  - `ts.SyntaxKind.BinaryExpression` (+, -, ===, !==) → +, -, ==, !=
  - `ts.SyntaxKind.IfStatement` → `if ... { ... }`
  - `ts.SyntaxKind.ForStatement` / `ts.SyntaxKind.ForOfStatement` → `for ... { ... }`
- **Error Handling:**
  - TS `throw new Error(...)` → `return nil, fmt.Errorf(...)`
  - TS `try/catch` → This is complex. An MVP might just disallow it and force explicit error returns from functions. A later version could transpile try blocks to `if err != nil { ... }` checks.
- **Function Calls:**
  - `console.log("hello")` → `fmt.Println("hello")` (This is a specific mapping we must implement).

### Phase 4: The ts2go Runtime Library

*(Goal: Support Node.js/Browser APIs by providing a Go equivalent.)*

- **Create a new Go Module:** `github.com/your-org/ts2go-runtime`.
- **Implement fs adapter:**
  - Create a `fs/fs.go` package.
  - Implement `func ReadFile(path string) (string, error)` that internally calls `os.ReadFile` and returns a string (not a Buffer).
  - Implement `func ReadFileSync(...)` (which will be identical to the async one, as all Go I/O is blocking by default, unless we use goroutines).
- **Implement console adapter:**
  - `console/console.go` with `func Log(...interface{})`, `func Error(...interface{})` that map to `fmt.Println` and `log.Println`.
- **Compiler Integration:**
  - Update the transpiler to detect `import { readFile } from 'fs'`.
  - This import should be rewritten to: `import "github.com/your-org/ts2go-runtime/fs"`
  - The call `readFile(...)` becomes `fs.ReadFile(...)`.

### Phase 5: Advanced Features (Post-MVP)

*(Goal: Support key TS features like classes, modules, and async.)*

- **Module Resolution:**
  - Handle `import` statements for other TS files in the project.
  - The tool must transpile the entire dependency tree of local files.
  - `import { MyType } from './types'` → Transpile `./types.ts` to `types.go` and add `import "my-ts-project/types"` and change the call to `types.MyType`.
- **Class Transpilation:**
  - `class User { ... }` → `type User struct { ... }`
  - `constructor(...)` → `func NewUser(...) *User { ... }`
  - `myMethod(...)` → `func (u *User) MyMethod(...) { ... }`
- **Async/Await (The Big One):**
  - As defined in Phase 1: `async function(...) Promise<T>` → `func(...) chan Result{T, error}`.
  - This will be a massive architectural challenge and will fundamentally change the "feel" of the code.

### Phase 6: Testing, Documentation & Release

- **Testing Strategy:**
  - **Unit Tests:** For the transpiler itself. "Given TS AST node X, assert Go AST node Y is produced."
  - **Integration Tests:** Create a directory of small, self-contained TS packages.
    - The test script will:
      - Run `ts2go` on the TS package.
      - `cd` into the output directory.
      - Run `go tidy` and `go build`.
      - Run the resulting binary and check stdout against an expected value.
- **Documentation:**
  - This is more important than the code.
  - The `SPEC.md` from Phase 1 is the most critical doc.
  - "Getting Started" guide.
  - "Limitations & Pitfalls" (e.g., "Why your dynamic objects are slow in Go").
- **Release:**
  - Publish to GitHub.
  - Provide `go install ...` instructions.
  - Publish binary releases for different platforms.

## 4. Key Pitfalls & Cautions

- **The "Slim and Nice" Fallacy:** Be prepared that transpiled JS/TS will not look like idiomatic Go. It will likely be "ugly" Go, full of `map[string]interface{}` and channel-based "promises." The goal is functional equivalence, not aesthetic equivalence.
- **npm Dependencies:** You cannot transpile the `node_modules` directory. It's too complex and often contains native C++ bindings. Any required library (e.g., lodash) would need to be either:
  - Rewritten in TS (and be part of your transpilation).
  - Rewritten in Go (and imported by the `ts2go-runtime`).
- **Performance:** The transpiled Go code, especially logic relying on `map[string]interface{}`, may be slower than the V8 JIT-compiled Node.js code. The "slim" benefit will be in binary size and memory usage, not necessarily CPU speed.