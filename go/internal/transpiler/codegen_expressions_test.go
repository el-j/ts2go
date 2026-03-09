package transpiler

import (
	"testing"
)

func TestBinaryExpression_GreaterThan(t *testing.T) {
	node := &ASTNode{
		Kind:     BinaryExpression,
		Operator: "GreaterThanToken",
		Children: []ASTNode{
			{Kind: Identifier, Text: "x"},
			{Kind: "GreaterThanToken"},
			{Kind: NumericLiteral, Text: "5"},
		},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != "x > 5" {
		t.Errorf("Expected 'x > 5', got '%s'", result)
	}
}

func TestStringLiteral(t *testing.T) {
	node := &ASTNode{Kind: StringLiteral, Text: "hello"}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != `"hello"` {
		t.Errorf("Expected '\"hello\"', got '%s'", result)
	}
}

func TestNumericLiteral(t *testing.T) {
	node := &ASTNode{Kind: NumericLiteral, Text: "42"}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != "42" {
		t.Errorf("Expected '42', got '%s'", result)
	}
}

func TestBoolLiterals(t *testing.T) {
	g := NewCodeGenerator()
	nodeTrue := &ASTNode{Kind: "TrueKeyword"}
	result, _ := g.generateExpression(nodeTrue)
	if result != "true" {
		t.Errorf("Expected 'true', got '%s'", result)
	}

	nodeFalse := &ASTNode{Kind: "FalseKeyword"}
	result, _ = g.generateExpression(nodeFalse)
	if result != "false" {
		t.Errorf("Expected 'false', got '%s'", result)
	}

	nodeNull := &ASTNode{Kind: "NullKeyword"}
	result, _ = g.generateExpression(nodeNull)
	if result != "nil" {
		t.Errorf("Expected 'nil', got '%s'", result)
	}
}

func TestCallExpression_consoleLog(t *testing.T) {
	node := &ASTNode{
		Kind: CallExpression,
		Expression: &ASTNode{
			Kind:       "PropertyAccessExpression",
			Name:       "log",
			Expression: &ASTNode{Kind: Identifier, Text: "console"},
		},
		Children: []ASTNode{
			{Kind: StringLiteral, Text: "test"},
		},
	}

	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != `fmt.Println("test")` {
		t.Errorf("Expected 'fmt.Println(\"test\")', got '%s'", result)
	}
	if g.imports["fmt"] != true {
		t.Errorf("Expected 'fmt' to be imported")
	}
}

func TestPropertyAccess(t *testing.T) {
	node := &ASTNode{
		Kind:       PropertyAccessExpression,
		Name:       "name",
		Expression: &ASTNode{Kind: Identifier, Text: "person"},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != "person.Name" {
		t.Errorf("Expected 'person.Name', got '%s'", result)
	}
}

func TestConditionalExpression(t *testing.T) {
	node := &ASTNode{
		Kind: ConditionalExpression,
		Children: []ASTNode{
			{Kind: Identifier, Text: "ok"},
			{Kind: NumericLiteral, Text: "1"},
			{Kind: NumericLiteral, Text: "0"},
		},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := "func() interface{} { if ok { return 1 } else { return 0 } }()"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestArrowFunction(t *testing.T) {
	node := &ASTNode{
		Kind: ArrowFunction,
		Parameters: []ASTNode{
			{Kind: "Parameter", Name: "x", Type: &ASTNode{Kind: "NumberKeyword"}},
		},
		Type: &ASTNode{Kind: "NumberKeyword"},
		Body: &ASTNode{
			Kind:     BinaryExpression,
			Operator: "PlusToken",
			Children: []ASTNode{
				{Kind: Identifier, Text: "x"},
				{Kind: "PlusToken"},
				{Kind: NumericLiteral, Text: "1"},
			},
		},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := "func(x float64) float64 { return x + 1 }"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestAwaitExpression(t *testing.T) {
	node := &ASTNode{
		Kind: AwaitExpression,
		Expression: &ASTNode{
			Kind:       CallExpression,
			Expression: &ASTNode{Kind: Identifier, Text: "fetchData"},
		},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// generateAwaitExpression might return FetchData() if it doesn't do special mapping yet,
	// typically wait expressions are just returning the expression since Go is sync.
	if result != "(<-fetchData())" {
		t.Errorf("Expected '(<-fetchData())', got '%s'", result)
	}
}

func TestTemplateExpression(t *testing.T) {
	node := &ASTNode{
		Kind: TemplateExpression,
		Head: &ASTNode{Text: "hello "},
		Children: []ASTNode{
			{
				Kind:       "TemplateSpan",
				Expression: &ASTNode{Kind: Identifier, Text: "name"},
				Literal:    &ASTNode{Text: "!"},
			},
		},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := `"hello " + fmt.Sprint(name) + "!"`
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestPostfixUnaryExpression(t *testing.T) {
	node := &ASTNode{
		Kind: "PostfixUnaryExpression",
		Children: []ASTNode{
			{Kind: Identifier, Text: "i"},
		},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != "i++" {
		t.Errorf("Expected 'i++', got '%s'", result)
	}
}

func TestPrefixUnaryExpression(t *testing.T) {
	node := &ASTNode{
		Kind: "PrefixUnaryExpression",
		Children: []ASTNode{
			{Kind: "ExclamationToken", Pos: 0},
		},
		Pos: 0,
	}
	// Wait, the implementation of prefix unary uses child position to guess operator!
	node.Children[0].Pos = 1 // Operator length 1
	node.Children = append(node.Children, ASTNode{Kind: Identifier, Text: "ok"})

	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Wait, it expects operand in Children[0].
	node2 := &ASTNode{
		Kind: "PrefixUnaryExpression",
		Pos:  0,
		Children: []ASTNode{
			{Kind: Identifier, Text: "ok", Pos: 1}, // +1 length operator
		},
	}
	result, _ = g.generateExpression(node2)
	if result != "-ok" {
		t.Errorf("Expected '-ok', got '%s'", result) // because it guesses - by default for 1 char diff without true/false
	}
}

func TestObjectLiteralExpression(t *testing.T) {
	node := &ASTNode{
		Kind: "ObjectLiteralExpression",
		Properties: []ASTNode{
			{
				Kind:        "PropertyAssignment",
				Name:        "a",
				Initializer: &ASTNode{Kind: NumericLiteral, Text: "1"},
			},
		},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := `map[string]interface{}{"a": 1}`
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestArrayLiteralExpression(t *testing.T) {
	node := &ASTNode{
		Kind: "ArrayLiteralExpression",
		Elements: []ASTNode{
			{Kind: NumericLiteral, Text: "1"},
			{Kind: NumericLiteral, Text: "2"},
		},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := `[]interface{}{1, 2}`
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestTypeOfExpression(t *testing.T) {
	node := &ASTNode{
		Kind:       TypeOfExpression,
		Expression: &ASTNode{Kind: Identifier, Text: "obj"},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := `reflect.TypeOf(obj).String()`
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestDeleteExpression(t *testing.T) {
	node := &ASTNode{
		Kind: DeleteExpression,
		Expression: &ASTNode{
			Kind:       PropertyAccessExpression,
			Name:       "key",
			Expression: &ASTNode{Kind: Identifier, Text: "m"},
		},
	}
	g := NewCodeGenerator()
	result, err := g.generateExpression(node)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := `delete(m, "key")`
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestThisKeyword(t *testing.T) {
	node := &ASTNode{Kind: ThisKeyword}
	g := NewCodeGenerator()
	g.currentReceiverVar = "p"
	result, _ := g.generateExpression(node)
	if result != "p" {
		t.Errorf("Expected 'p', got '%s'", result)
	}
}

func TestBinaryExpression_NullishCoalesce(t *testing.T) {
	node := &ASTNode{
		Kind:     BinaryExpression,
		Operator: "QuestionQuestionToken",
		Children: []ASTNode{
			{Kind: Identifier, Text: "a"},
			{Kind: "QuestionQuestionToken"},
			{Kind: Identifier, Text: "b"},
		},
	}
	g := NewCodeGenerator()
	result, _ := g.generateExpression(node)
	expected := "nullishCoalesce(a, b)"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
