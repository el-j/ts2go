package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/el-j/ts2go/internal/transpiler"
)

func main() {
	data, err := os.ReadFile("/tmp/ast.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading AST file: %v\n", err)
		return
	}

	var ast transpiler.ASTNode
	if err := json.Unmarshal(data, &ast); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling AST: %v\n", err)
		return
	}

	if len(ast.Statements) < 2 || ast.Statements[1].Body == nil || len(ast.Statements[1].Body.Statements) < 3 {
		return
	}

	funcDecl := ast.Statements[1]
	ifStmt := funcDecl.Body.Statements[2]

	conditionNode := ifStmt.Expression
	if conditionNode == nil && len(ifStmt.Children) > 0 {
		conditionNode = &ifStmt.Children[0]
	}

	b, _ := json.MarshalIndent(conditionNode, "", "  ")
	fmt.Println(string(b))
}
