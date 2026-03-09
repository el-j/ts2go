package transpiler

import (
	"strings"
	"testing"
)

func TestClassDeclaration(t *testing.T) {
	node := &ASTNode{
		Kind: ClassDeclaration,
		Name: "Person",
		Members: []ASTNode{
			{
				Kind: PropertyDeclaration,
				Name: "name",
				Type: &ASTNode{Kind: "StringKeyword"},
			},
		},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	if !strings.Contains(result, "type Person struct") || !strings.Contains(result, "Name string") {
		t.Errorf("Expected struct definition, got:\n%s", result)
	}
}

func TestClassWithConstructor(t *testing.T) {
	node := &ASTNode{
		Kind: ClassDeclaration,
		Name: "User",
		Members: []ASTNode{
			{
				Kind: "Constructor",
				Parameters: []ASTNode{
					{Kind: "Parameter", Name: "id", Type: &ASTNode{Kind: "NumberKeyword"}},
				},
				Body: &ASTNode{
					Kind: Block,
					Statements: []ASTNode{
						{
							Kind: ExpressionStatement,
							Expression: &ASTNode{
								Kind:     BinaryExpression,
								Operator: "FirstAssignment",
								Children: []ASTNode{
									{Kind: PropertyAccessExpression, Name: "id", Expression: &ASTNode{Kind: ThisKeyword}},
									{Kind: "FirstAssignment"},
									{Kind: Identifier, Text: "id"},
								},
							},
						},
					},
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
	if !strings.Contains(result, "func NewUser(id float64) *User") {
		t.Errorf("Expected NewUser constructor, got:\n%s", result)
	}
}

func TestPublicMethod(t *testing.T) {
	node := &ASTNode{
		Kind: ClassDeclaration,
		Name: "Greeter",
		Members: []ASTNode{
			{
				Kind:       MethodDeclaration,
				Name:       "greet",
				Parameters: []ASTNode{},
				Type:       &ASTNode{Kind: "StringKeyword"},
				Body: &ASTNode{
					Kind: Block,
					Statements: []ASTNode{
						{Kind: ReturnStatement, Expression: &ASTNode{Kind: StringLiteral, Text: "Hello"}},
					},
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
	if !strings.Contains(result, "func (g *Greeter) Greet() string") {
		t.Errorf("Expected Greet method, got:\n%s", result)
	}
}

func TestPrivateMethod(t *testing.T) {
	node := &ASTNode{
		Kind: ClassDeclaration,
		Name: "Helper",
		Members: []ASTNode{
			{
				Kind:       MethodDeclaration,
				Name:       "doWork",
				Modifiers:  []ASTNode{{Kind: "PrivateKeyword"}},
				Parameters: []ASTNode{},
				Body:       &ASTNode{Kind: Block},
			},
		},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	// The current transpiler implementation makes everything PascalCase in Go
	if !strings.Contains(result, "func (h *Helper) DoWork()") {
		t.Errorf("Expected DoWork private method, got:\n%s", result)
	}
}

func TestGetAccessor(t *testing.T) {
	node := &ASTNode{
		Kind: ClassDeclaration,
		Name: "Box",
		Members: []ASTNode{
			{
				Kind: GetAccessor,
				Name: "value",
				Body: &ASTNode{
					Kind: Block,
					Statements: []ASTNode{
						{Kind: ReturnStatement, Expression: &ASTNode{Kind: NumericLiteral, Text: "1"}},
					},
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
	if !strings.Contains(result, "func (b *Box) GetValue()") {
		t.Errorf("Expected GetValue getter, got:\n%s", result)
	}
}

func TestSetAccessor(t *testing.T) {
	node := &ASTNode{
		Kind: ClassDeclaration,
		Name: "Box",
		Members: []ASTNode{
			{
				Kind:       SetAccessor,
				Name:       "value",
				Parameters: []ASTNode{{Kind: "Parameter", Name: "val", Type: &ASTNode{Kind: "NumberKeyword"}}},
				Body:       &ASTNode{Kind: Block},
			},
		},
	}
	g := NewCodeGenerator()
	err := g.generateStatement(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	result := strings.TrimSpace(g.output.String())
	if !strings.Contains(result, "func (b *Box) SetValue(val float64)") {
		t.Errorf("Expected SetValue setter, got:\n%s", result)
	}
}
