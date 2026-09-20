package transpiler

import (
	"fmt"
)

// generateStatement generates code for statements
func (g *CodeGenerator) generateStatement(node *ASTNode) error {
	switch node.Kind {
	case InterfaceDeclaration:
		return g.generateInterface(node)
	case TypeAliasDeclaration:
		return g.generateTypeAlias(node)
	case EnumDeclaration:
		return g.generateEnum(node)
	case ClassDeclaration:
		return g.generateClass(node)
	case FunctionDeclaration:
		return g.generateFunction(node)
	case VariableStatement, "FirstStatement":
		return g.generateVariableStatement(node)
	case ExpressionStatement:
		return g.generateExpressionStatement(node)
	case ReturnStatement:
		return g.generateReturnStatement(node)
	case IfStatement:
		return g.generateIfStatement(node)
	case ForStatement:
		return g.generateForStatement(node)
	case ForOfStatement:
		return g.generateForOfStatement(node)
	case ForInStatement:
		return g.generateForInStatement(node)
	case WhileStatement:
		return g.generateWhileStatement(node)
	case BreakStatement:
		return g.generateBreakStatement(node)
	case ContinueStatement:
		return g.generateContinueStatement(node)
	case SwitchStatement:
		return g.generateSwitchStatement(node)
	case TryStatement:
		return g.generateTryStatement(node)
	case ThrowStatement:
		return g.generateThrowStatement(node)
	case Block:
		return g.generateBlock(node)
	case "EmptyStatement":
		return nil
	case "ImportDeclaration", "ExportDeclaration", "ExportAssignment":
		return nil
	default:
		return UnsupportedFeatureError("", 0, 0, fmt.Sprintf("unsupported statement: %s", node.Kind))
	}
}

// generateExpressionStatement generates code for expression statements
func (g *CodeGenerator) generateExpressionStatement(node *ASTNode) error {
	var exprNode *ASTNode
	if node.Expression != nil {
		exprNode = node.Expression
	} else if len(node.Children) > 0 {
		exprNode = &node.Children[0]
	}

	if exprNode != nil {
		expr, err := g.generateExpression(exprNode)
		if err != nil {
			return err
		}
		g.writeLine(expr)
	}
	return nil
}
