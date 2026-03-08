---
id: TASK-010
title: 'Add transpiler code generation unit tests'
status: in-progress
priority: high
dependencies: ['TASK-005']
estimated_effort: L
---

## Goal

Dramatically increase transpiler coverage by adding focused unit tests for each `codegen_*.go` file, targeting the many currently-untested expression/statement generators.

## Context

Current state: `go/internal/transpiler/` tests exist only in `errors_test.go` and `debug_test.go` (debug must be deleted by TASK-005). Coverage of `codegen_expressions.go`, `codegen_statements.go`, `codegen_classes.go`, `codegen_types.go` is essentially 0%.

The `CodeGenerator` has a public `Generate(*ASTNode) (string, error)` method which allows black-box testing by constructing AST nodes inline without involving the Node.js parser.

## Scope

### Files to create

- `go/internal/transpiler/codegen_expressions_test.go`
- `go/internal/transpiler/codegen_statements_test.go`
- `go/internal/transpiler/codegen_classes_test.go`
- `go/internal/transpiler/codegen_types_test.go`

## Implementation

For each test file, construct `ASTNode` structs directly and call `g.generateExpression()` or `g.generateStatement()`:

**codegen_expressions_test.go** (20+ tests):

- `TestBinaryExpression_GreaterThan`: node with `BinaryExpression`, children `[x, GreaterThanToken, 5]`, expect `"x > 5"`
- `TestStringLiteral`, `TestNumericLiteral`, `TestBoolLiterals`
- `TestCallExpression_consoleLog`: expect `fmt.Println`
- `TestPropertyAccess`, `TestConditionalExpression` (ternary)
- `TestArrowFunction`, `TestAwaitExpression`

**codegen_statements_test.go** (15+ tests):

- `TestIfStatement`, `TestIfElseStatement`, `TestForStatement`, `TestForOfStatement`
- `TestVariableDeclaration_const`, `TestVariableDeclaration_let`
- `TestReturnStatement`, `TestBreakStatement`
- `TestObjectDestructuring` with rest pattern (TASK-003 fix)

**codegen_classes_test.go** (10+ tests):

- `TestClassDeclaration`, `TestClassWithConstructor`
- `TestPublicMethod`, `TestPrivateMethod`
- `TestGetAccessor`, `TestSetAccessor`

**codegen_types_test.go** (8+ tests):

- `TestInterface_toStruct`, `TestTypeAlias`
- `TestEnumDeclaration`, `TestUnionType`

## Acceptance Criteria

- [ ] `go test ./internal/transpiler/... -cover` exits 0
- [ ] Coverage for `codegen_expressions.go` ≥ 70%
- [ ] Coverage for `codegen_statements.go` ≥ 70%
- [ ] At least 50 new test cases total
- [ ] No `any` types introduced
