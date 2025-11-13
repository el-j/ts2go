# TS2Go Specification

## Supported TypeScript Features

### Types
- `string` → `string`
- `number` → `float64`
- `boolean` → `bool`
- `Array<T>` → `[]T`
- `interface` → `type MyInterface struct { ... }`
- `type` aliases → `type MyType = ...`
- Inline object types → `struct { ... }`

### Unsupported Types
- `any`, `unknown` (disallowed)
- Union types (e.g., `string | number`)
- Dynamic objects (use `map[string]interface{}` with caution)

### Logic
- Functions: `function myFunc(...)` → `func MyFunc(...)`
- Variables: `let/const` → `var` or `:=`
- Control flow: `if/else`, `for` loops
- Binary operations: `+`, `-`, `==`, `!=`
- `console.log` → `fmt.Println`

### Unsupported Logic
- `class` (initially)
- `async/await`, `Promise`
- `try/catch` (use error returns)
- Dynamic property access

### Modules
- Local imports: `import { ... } from './file'` → transpile and import
- Node.js APIs: `import { readFile } from 'fs'` → `import "github.com/ts2go/runtime/fs"`

## Limitations
- No npm dependencies transpilation
- Performance trade-offs for dynamic features
- Code may not be idiomatic Go