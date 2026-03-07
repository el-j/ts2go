package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/el-j/ts2go/internal/transpiler"
)

func main() {
	data, _ := ioutil.ReadFile("/tmp/ast.json")
	var ast transpiler.ASTNode
	json.Unmarshal(data, &ast)
	funcDecl := ast.Statements[1]
	ifStmt := funcDecl.Body.Statements[2]

	conditionNode := ifStmt.Expression
	if conditionNode == nil && len(ifStmt.Children) > 0 {
		conditionNode = &ifStmt.Children[0]
	}

	b, _ := json.MarshalIndent(conditionNode, "", "  ")
	fmt.Println(string(b))
}
