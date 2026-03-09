package transpiler

import (
	"strings"
	"testing"
)

func TestIfStatement(t *testing.T) {
	node := &ASTNode{
		Kind:       IfStatement,
		Expression: &ASTNode{Kind: Identifier, Text: "cond"},
		Children: []ASTNode{
			{Kind: Identifier, Text: "cond"}, // 0: condition fallback ignored if Expression is set
			{
				Kind: Block, // 1: then block
				Statements: []ASTNode{
					{Kind: ExpressionStatement, Expression: &ASTNode{Kind: Identifier, Text: "execute()"}},
				},
			},
		},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	expected := "if cond {\n\texecute()\n}"
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestIfElseStatement(t *testing.T) {
	node := &ASTNode{
		Kind:       IfStatement,
		Expression: &ASTNode{Kind: Identifier, Text: "cond"},
		Children: []ASTNode{
			{Kind: Identifier, Text: "cond"}, // 0: condition fallback
			{
				Kind: Block, // 1: then block
				Statements: []ASTNode{
					{Kind: ExpressionStatement, Expression: &ASTNode{Kind: Identifier, Text: "y()"}},
				},
			},
			{
				Kind: Block, // 2: else block
				Statements: []ASTNode{
					{Kind: ExpressionStatement, Expression: &ASTNode{Kind: Identifier, Text: "n()"}},
				},
			},
		},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	expected := "if cond {\n\ty()\n} else {\n\tn()\n}"
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestForStatement(t *testing.T) {
	node := &ASTNode{
		Kind: ForStatement,
		Initializer: &ASTNode{
			Kind: "VariableDeclarationList",
			Children: []ASTNode{
				{Kind: VariableDeclaration, Name: "i", Initializer: &ASTNode{Kind: NumericLiteral, Text: "0"}},
			},
		},
		Children: []ASTNode{
			{ // 0: condition
				Kind:     BinaryExpression,
				Operator: "LessThanToken",
				Children: []ASTNode{
					{Kind: Identifier, Text: "i"},
					{Kind: "LessThanToken"},
					{Kind: NumericLiteral, Text: "10"},
				},
			},
			{ // 1: incrementor
				Kind:     "PostfixUnaryExpression",
				Children: []ASTNode{{Kind: Identifier, Text: "i"}},
			},
			{ // 2: body block
				Kind: Block,
				Statements: []ASTNode{
					{Kind: ExpressionStatement, Expression: &ASTNode{Kind: Identifier, Text: "log()"}},
				},
			},
		},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	// Expected loop
	if !strings.Contains(result, "for") || !strings.Contains(result, "i < 10") {
		t.Errorf("Expected loop, got:\n%s", result)
	}
}

func TestForOfStatement(t *testing.T) {
	node := &ASTNode{
		Kind: ForOfStatement,
		Initializer: &ASTNode{
			Kind: "VariableDeclarationList",
			Children: []ASTNode{
				{Kind: VariableDeclaration, Name: "item"},
			},
		},
		Expression: &ASTNode{Kind: Identifier, Text: "items"},
		Body: &ASTNode{
			Kind: Block,
			Statements: []ASTNode{
				{Kind: ExpressionStatement, Expression: &ASTNode{Kind: Identifier, Text: "process(item)"}},
			},
		},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	expected := "for _, item := range items {\n\tprocess(item)\n}"
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestVariableDeclaration_const(t *testing.T) {
	node := &ASTNode{
		Kind: VariableStatement,
		Declarations: []ASTNode{
			{
				Kind:        VariableDeclaration,
				Name:        "MAX",
				Initializer: &ASTNode{Kind: NumericLiteral, Text: "100"},
			},
		},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	if !strings.Contains(result, "MAX := 100") {
		t.Errorf("Expected MAX := 100 declaration, got:\n%s", result)
	}
}

func TestVariableDeclaration_let(t *testing.T) {
	node := &ASTNode{
		Kind: VariableStatement,
		Declarations: []ASTNode{
			{
				Kind:        VariableDeclaration,
				Name:        "count",
				Type:        &ASTNode{Kind: "NumberKeyword"},
				Initializer: &ASTNode{Kind: NumericLiteral, Text: "0"},
			},
		},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	if !strings.Contains(result, "count := 0") {
		t.Errorf("Expected count let declaration, got:\n%s", result)
	}
}

func TestReturnStatement(t *testing.T) {
	node := &ASTNode{
		Kind:       ReturnStatement,
		Expression: &ASTNode{Kind: NumericLiteral, Text: "42"},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	if result != "return 42" {
		t.Errorf("Expected 'return 42', got '%s'", result)
	}
}

func TestBreakStatement(t *testing.T) {
	node := &ASTNode{Kind: BreakStatement}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	if result != "break" {
		t.Errorf("Expected 'break', got '%s'", result)
	}
}

func TestObjectDestructuring(t *testing.T) {
	node := &ASTNode{
		Kind: VariableStatement,
		Declarations: []ASTNode{
			{
				Kind: VariableDeclaration,
				NameNode: &ASTNode{
					Kind: "ObjectBindingPattern",
					Elements: []ASTNode{
						{Kind: "BindingElement", Name: "a"},
						{Kind: "BindingElement", Name: "b"},
						{
							Kind:     "BindingElement",
							Name:     "rest",
							Children: []ASTNode{{Kind: "DotDotDotToken"}},
						},
					},
				},
				Initializer: &ASTNode{Kind: Identifier, Text: "obj"},
			},
		},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	if !strings.Contains(result, "rest") || !strings.Contains(result, "delete") {
		t.Errorf("Expected rest destructuring to emit delete statements, got:\n%s", result)
	}
}
