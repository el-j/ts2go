package transpiler

// ASTNode represents a simplified TypeScript AST node
type ASTNode struct {
	Kind            string    `json:"kind"`
	KindNumber      int       `json:"kindNumber"`
	Pos             int       `json:"pos"`
	End             int       `json:"end"`
	Text            string    `json:"text,omitempty"`
	Name            string    `json:"name,omitempty"`
	Type            *ASTNode  `json:"type,omitempty"`
	Types           []ASTNode `json:"types,omitempty"`
	Elements        []ASTNode `json:"elements,omitempty"`
	Parameters      []ASTNode `json:"parameters,omitempty"`
	Members         []ASTNode `json:"members,omitempty"`
	HeritageClauses []ASTNode `json:"heritageClauses,omitempty"`
	Properties      []ASTNode `json:"properties,omitempty"`
	Statements      []ASTNode `json:"statements,omitempty"`
	Body            *ASTNode  `json:"body,omitempty"`
	Initializer     *ASTNode  `json:"initializer,omitempty"`
	Declarations    []ASTNode `json:"declarations,omitempty"`
	Children        []ASTNode `json:"children,omitempty"`
	Operator        string    `json:"operator,omitempty"`
	OperatorNumber  int       `json:"operatorNumber,omitempty"`
	QuestionDot     bool      `json:"questionDot,omitempty"`
}

// TypeScriptKind maps to TypeScript SyntaxKind enum
const (
	// Declarations
	InterfaceDeclaration = "InterfaceDeclaration"
	TypeAliasDeclaration = "TypeAliasDeclaration"
	EnumDeclaration      = "EnumDeclaration"
	FunctionDeclaration  = "FunctionDeclaration"
	ClassDeclaration     = "ClassDeclaration"
	VariableStatement    = "VariableStatement"
	VariableDeclaration  = "VariableDeclaration"

	// Types
	StringKeyword  = "StringKeyword"
	NumberKeyword  = "NumberKeyword"
	BooleanKeyword = "BooleanKeyword"
	ArrayType      = "ArrayType"
	TupleType      = "TupleType"
	TypeLiteral    = "TypeLiteral"
	TypeReference  = "TypeReference"
	UnionType      = "UnionType"
	LiteralType    = "LiteralType"

	// Properties
	PropertySignature   = "PropertySignature"
	PropertyDeclaration = "PropertyDeclaration"
	EnumMember          = "EnumMember"
	Parameter           = "Parameter"
	Constructor         = "Constructor"
	MethodDeclaration   = "MethodDeclaration"
	GetAccessor         = "GetAccessor"
	SetAccessor         = "SetAccessor"

	// Expressions
	Identifier                  = "Identifier"
	StringLiteral               = "StringLiteral"
	NumericLiteral              = "NumericLiteral"
	BinaryExpression            = "BinaryExpression"
	CallExpression              = "CallExpression"
	PropertyAccessExpression    = "PropertyAccessExpression"
	NewExpression               = "NewExpression"
	ThisKeyword                 = "ThisKeyword"
	SuperKeyword                = "SuperKeyword"
	ExpressionWithTypeArguments = "ExpressionWithTypeArguments"
	HeritageClause              = "HeritageClause"

	// Operators/Tokens
	QuestionDotToken      = "QuestionDotToken"
	QuestionQuestionToken = "QuestionQuestionToken"

	// Modifiers
	PublicKeyword  = "PublicKeyword"
	PrivateKeyword = "PrivateKeyword"
	StaticKeyword  = "StaticKeyword"

	// Statements
	Block               = "Block"
	ReturnStatement     = "ReturnStatement"
	IfStatement         = "IfStatement"
	ForStatement        = "ForStatement"
	ExpressionStatement = "ExpressionStatement"

	// File structure
	SourceFile = "SourceFile"
)
