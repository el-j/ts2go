package transpiler

// ASTNode represents a simplified TypeScript AST node
type ASTNode struct {
	Kind         string    `json:"kind"`
	KindNumber   int       `json:"kindNumber"`
	Pos          int       `json:"pos"`
	End          int       `json:"end"`
	Text         string    `json:"text,omitempty"`
	Name         string    `json:"name,omitempty"`
	Type         *ASTNode  `json:"type,omitempty"`
	Parameters   []ASTNode `json:"parameters,omitempty"`
	Members      []ASTNode `json:"members,omitempty"`
	Properties   []ASTNode `json:"properties,omitempty"`
	Elements     []ASTNode `json:"elements,omitempty"`
	Statements   []ASTNode `json:"statements,omitempty"`
	Body         *ASTNode  `json:"body,omitempty"`
	Initializer  *ASTNode  `json:"initializer,omitempty"`
	Declarations []ASTNode `json:"declarations,omitempty"`
	Children     []ASTNode `json:"children,omitempty"`
}

// TypeScriptKind maps to TypeScript SyntaxKind enum
const (
	// Declarations
	InterfaceDeclaration = "InterfaceDeclaration"
	TypeAliasDeclaration = "TypeAliasDeclaration"
	FunctionDeclaration  = "FunctionDeclaration"
	VariableStatement    = "VariableStatement"
	VariableDeclaration  = "VariableDeclaration"

	// Types
	StringKeyword  = "StringKeyword"
	NumberKeyword  = "NumberKeyword"
	BooleanKeyword = "BooleanKeyword"
	ArrayType      = "ArrayType"
	TypeLiteral    = "TypeLiteral"
	TypeReference  = "TypeReference"

	// Properties
	PropertySignature   = "PropertySignature"
	PropertyDeclaration = "PropertyDeclaration"
	Parameter           = "Parameter"

	// Expressions
	Identifier       = "Identifier"
	StringLiteral    = "StringLiteral"
	NumericLiteral   = "NumericLiteral"
	BinaryExpression = "BinaryExpression"
	CallExpression   = "CallExpression"

	// Statements
	Block               = "Block"
	ReturnStatement     = "ReturnStatement"
	IfStatement         = "IfStatement"
	ForStatement        = "ForStatement"
	ExpressionStatement = "ExpressionStatement"

	// File structure
	SourceFile = "SourceFile"
)
