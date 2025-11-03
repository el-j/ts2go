package transpiler

// ASTNode represents a simplified TypeScript AST node
type ASTNode struct {
	Kind            string    `json:"kind"`
	KindNumber      int       `json:"kindNumber"`
	Pos             int       `json:"pos"`
	End             int       `json:"end"`
	Text            string    `json:"text,omitempty"`
	Name            string    `json:"name,omitempty"`
	NameNode        *ASTNode  `json:"nameNode,omitempty"`  // For binding patterns (destructuring)
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
	Expression      *ASTNode  `json:"expression,omitempty"` // For template spans
	Literal         *ASTNode  `json:"literal,omitempty"`    // For template spans
	Head            *ASTNode  `json:"head,omitempty"`       // For template expressions
}

// TypeScriptKind maps to TypeScript SyntaxKind enum
const (
	// Declarations
	InterfaceDeclaration = "InterfaceDeclaration"
	TypeAliasDeclaration = "TypeAliasDeclaration"
	EnumDeclaration      = "EnumDeclaration"
	FunctionDeclaration  = "FunctionDeclaration"
	ArrowFunction        = "ArrowFunction" // Arrow function expressions
	ClassDeclaration     = "ClassDeclaration"
	VariableStatement    = "VariableStatement"
	FirstStatement       = "FirstStatement" // Alias for VariableStatement
	VariableDeclaration  = "VariableDeclaration"

	// Types
	StringKeyword  = "StringKeyword"
	NumberKeyword  = "NumberKeyword"
	BooleanKeyword = "BooleanKeyword"
	VoidKeyword    = "VoidKeyword"
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
	TrueKeyword                 = "TrueKeyword"
	FalseKeyword                = "FalseKeyword"
	NullKeyword                 = "NullKeyword"
	TemplateExpression          = "TemplateExpression"      // Template literal: `hello ${name}`
	TemplateLiteralTypeSpan     = "TemplateLiteralTypeSpan" // Part of template
	TemplateHead                = "TemplateHead"            // Start of template: `hello ${
	TemplateMiddle              = "TemplateMiddle"          // Middle of template: } world ${
	TemplateTail                = "TemplateTail"            // End of template: } !`
	BinaryExpression            = "BinaryExpression"
	CallExpression              = "CallExpression"
	TypeOfExpression            = "TypeOfExpression"
	DeleteExpression            = "DeleteExpression"
	PropertyAccessExpression    = "PropertyAccessExpression"
	ElementAccessExpression     = "ElementAccessExpression" // array[index]
	NewExpression               = "NewExpression"
	ConditionalExpression       = "ConditionalExpression"   // Ternary operator: condition ? true : false
	ObjectLiteralExpression     = "ObjectLiteralExpression" // { key: value }
	ArrayLiteralExpression      = "ArrayLiteralExpression"  // [1, 2, 3]
	SpreadElement               = "SpreadElement"           // ...args
	ObjectBindingPattern        = "ObjectBindingPattern"    // { x, y } = obj
	ArrayBindingPattern         = "ArrayBindingPattern"     // [a, b] = arr
	BindingElement              = "BindingElement"          // Element in binding pattern
	PrefixUnaryExpression       = "PrefixUnaryExpression"   // ++x, --x, !x
	PostfixUnaryExpression      = "PostfixUnaryExpression"  // x++, x--
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
	ForOfStatement      = "ForOfStatement"    // for (const item of array)
	ForInStatement      = "ForInStatement"    // for (const key in object)
	WhileStatement      = "WhileStatement"    // while (condition)
	DoStatement         = "DoStatement"       // do { } while (condition)
	SwitchStatement     = "SwitchStatement"   // switch (expr) { case: ... }
	CaseBlock           = "CaseBlock"         // { case 1: ... case 2: ... default: ... }
	CaseClause          = "CaseClause"        // case value: statements
	DefaultClause       = "DefaultClause"     // default: statements
	BreakStatement      = "BreakStatement"    // break;
	ContinueStatement   = "ContinueStatement" // continue;
	TryStatement        = "TryStatement"      // try { } catch (e) { } finally { }
	CatchClause         = "CatchClause"       // catch (error) { }
	ThrowStatement      = "ThrowStatement"    // throw new Error("msg");
	ExpressionStatement = "ExpressionStatement"

	// File structure
	SourceFile = "SourceFile"
)
