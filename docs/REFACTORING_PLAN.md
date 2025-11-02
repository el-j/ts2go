# Refactoring Phase: Code Organization

**Date:** November 2, 2025  
**Status:** IN PROGRESS  
**Reason:** codegen.go is 1943 lines - too large to maintain

---

## Problem

**Current State:**
- `codegen.go`: 1943 lines, 51+ methods
- All code generation logic in one file
- Hard to navigate, test, and extend
- Adding new features (loops, ternary, async) will make it worse

**Impact:**
- Slows development
- Makes code review difficult  
- Increases bug risk
- Harder for new contributors

---

## Refactoring Strategy

### Phase 1: Extract to Separate Files (Same Package)
Split into focused files while keeping everything in `package transpiler`:

```
internal/transpiler/
├── ast.go (existing - AST node definitions)
├── transpiler.go (existing - public API)
├── codegen.go (REFACTORED - main orchestrator, ~150 lines)
├── codegen_core.go (NEW - Generate() entry point, imports, helpers)
├── codegen_declarations.go (NEW - interfaces, type aliases, enums)
├── codegen_classes.go (NEW - class generation, constructors, methods)
├── codegen_functions.go (NEW - function declarations)
├── codegen_statements.go (NEW - if/else, return, variable, blocks)
├── codegen_expressions.go (NEW - binary, call, property access)
├── codegen_types.go (NEW - type generation and mapping)
├── codegen_helpers.go (NEW - utility functions)
└── codegen_loops.go (NEW - for/while/switch - to be added)
```

### Phase 2: Future Package Structure (Optional)
If needed later, can extract to sub-packages:

```
internal/transpiler/
├── generator/
│   ├── generator.go (main CodeGenerator struct)
│   ├── declarations.go
│   ├── classes.go
│   ├── functions.go
│   ├── statements.go
│   ├── expressions.go
│   └── types.go
```

But **Phase 1 is sufficient** for now - same-package split is simpler.

---

## File Breakdown

### codegen.go (~150 lines)
**Purpose:** Main CodeGenerator struct definition
```go
- type CodeGenerator struct { ... }
- func NewCodeGenerator()
- func NewCodeGeneratorWithModule()
- Basic initialization
```

### codegen_core.go (~200 lines)
**Purpose:** Entry point and orchestration
```go
- func Generate() - main entry point
- Package/import generation
- Helper function injection
- Module system integration
```

### codegen_declarations.go (~300 lines)
**Purpose:** Top-level declarations
```go
- generateInterface()
- generateTypeAlias()
- generateEnum()
- generateUnionType()
- Enum helpers (numeric, string)
```

### codegen_classes.go (~500 lines)
**Purpose:** Class and OOP features
```go
- generateClass()
- generateConstructor()
- generateConstructorWithBase()
- generateMethod()
- generateGetter/Setter()
- generateStaticMethod()
- generateMethodBody()
- Class helpers (super calls, etc.)
```

### codegen_functions.go (~250 lines)
**Purpose:** Function generation
```go
- generateFunction()
- generateFunctionBody()
- Parameter generation
- Return type handling
```

### codegen_statements.go (~250 lines)
**Purpose:** Statement generation (INCLUDES OUR NEW IF/ELSE!)
```go
- generateStatement() - dispatcher
- generateVariableStatement()
- generateExpressionStatement()
- generateReturnStatement()
- generateIfStatement() ← OUR NEW CODE
- generateStatementBlock()
```

### codegen_expressions.go (~400 lines)
**Purpose:** Expression generation
```go
- generateExpression() - dispatcher
- generateBinaryExpression()
- generateCallExpression()
- generatePropertyAccessExpression()
- generateNewExpression()
- generateArrayLiteralExpression()
- generateObjectLiteralExpression()
- Literal handlers
```

### codegen_types.go (~300 lines)
**Purpose:** Type system
```go
- generateType()
- generateUnionType()
- generateTupleType()
- Type mapping (TS → Go)
```

### codegen_helpers.go (~200 lines)
**Purpose:** Utilities (ALREADY CREATED!)
```go
- writeLine()
- toPascalCase()
- toCamelCase()
- hasModifier(), isPrivate(), isStatic()
- collectImports()
- needsOptionalAccess()
- needsNullishCoalesce()
- Module system helpers
```

### codegen_loops.go (~300 lines - TO BE ADDED)
**Purpose:** Loop and control flow (FUTURE)
```go
- generateForStatement()
- generateForOfStatement()
- generateForInStatement()
- generateWhileStatement()
- generateSwitchStatement()
- generateBreakStatement()
- generateContinueStatement()
```

---

## Implementation Steps

### Step 1: ✅ Create Helper File (DONE)
- Created `codegen_helpers.go` with utility functions
- ~250 lines extracted

### Step 2: ⏳ Create Statements File (IN PROGRESS)
- Create `codegen_statements.go` with statement generation
- INCLUDES our new if/else logic
- Need to remove from original codegen.go

### Step 3: Create Type File
- Extract all type generation methods
- Move to `codegen_types.go`

### Step 4: Create Expression File
- Extract all expression generation methods
- Move to `codegen_expressions.go`

### Step 5: Create Declarations File
- Extract interface, type alias, enum generation
- Move to `codegen_declarations.go`

### Step 6: Create Classes File
- Extract all class-related methods (largest section)
- Move to `codegen_classes.go`

### Step 7: Create Functions File
- Extract function generation
- Move to `codegen_functions.go`

### Step 8: Create Core File
- Extract Generate() entry point
- Move to `codegen_core.go`

### Step 9: Slim Down Original
- Keep only CodeGenerator struct definition
- Keep NewCodeGenerator()
- ~150 lines total

### Step 10: Test Everything
- Run all existing tests
- Verify nothing breaks
- Fix any issues

---

## Testing Plan

### Verification Steps:
1. ✅ Backup original: `codegen.go.backup`
2. ⏳ Create new files with methods
3. ⏳ Build: `go build ./internal/transpiler`
4. ⏳ Run existing tests: `go test ./internal/transpiler`
5. ⏳ Test if-statement: `./ts2go convert tests/fixtures/if-test.ts`
6. ⏳ Run full integration test: `go test ./tests`

### Success Criteria:
- ✅ All files compile without errors
- ✅ All existing tests pass
- ✅ If-statement test still works
- ✅ No functionality lost
- ✅ Code is more maintainable

---

## Benefits After Refactoring

### For Development:
- ✅ Each file is 150-500 lines (manageable)
- ✅ Clear separation of concerns
- ✅ Easy to find specific functionality
- ✅ Easier to add new features
- ✅ Better for code review

### For Testing:
- ✅ Can test each area independently
- ✅ Easier to write unit tests
- ✅ Faster to locate bugs
- ✅ Better test organization

### For Future Features:
- ✅ New `codegen_loops.go` for loop features
- ✅ New `codegen_async.go` for async/await (future)
- ✅ New `codegen_operators.go` for ternary, etc. (future)
- ✅ Each feature in its own focused file

---

## Timeline

**Original Estimate:** 2-4 hours
**Current Progress:** 
- ✅ Step 1: Helper file (30 min)
- ⏳ Step 2: Statements file (in progress)
- ⏳ Steps 3-10: Remaining files

**Target Completion:** Today (Nov 2, 2025)

---

## Next Steps After Refactoring

Once refactoring is complete, continue with Phase 15:
1. ✅ If/else statements (DONE)
2. ⏳ Ternary operator - add to codegen_expressions.go
3. ⏳ For loops - add to new codegen_loops.go
4. ⏳ While loops - add to codegen_loops.go
5. ⏳ Switch - add to codegen_loops.go

**Impact:** Much easier to add these features with organized code!

---

**Status:** IN PROGRESS - Currently extracting statements to separate file
**Next:** Extract types, expressions, declarations, classes, functions, core
