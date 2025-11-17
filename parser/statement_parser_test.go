package parser

import (
	"fmt"
	"log"
	"model_infrax/config"
	"model_infrax/tool"
	"os"
	"testing"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/parser/test_driver"
)

type colX struct {
	colNames []string
}

func (v *colX) Enter(in ast.Node) (ast.Node, bool) {
	if name, ok := in.(*ast.ColumnName); ok {
		v.colNames = append(v.colNames, name.Name.O)
	}
	return in, false
}

func (v *colX) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func extract(rootNode *ast.StmtNode) []string {
	v := &colX{}
	(*rootNode).Accept(v)
	return v.colNames
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, _, err := p.ParseSQL(sql)
	if err != nil {
		return nil, err
	}

	return &stmtNodes[0], nil
}

func TestTidbPaseNode(t *testing.T) {
	p := parser.New()

	byts, err := os.ReadFile("/Users/jacklingo/dev/model_infrax/assert/database.sql")
	if err != nil {
		panic(err)
	}

	statement := string(byts)

	stmtNodes, _, err := p.ParseSQL(statement)
	if err != nil {
		panic(err)
	}

	log.Printf("stmtNodes: %s", tool.JsonifyIndent(stmtNodes))
}

func TestTiDBParser(t *testing.T) {
	byts, err := os.ReadFile("/Users/jacklingo/dev/model_infrax/assert/database.sql")
	if err != nil {
		panic(err)
	}

	statement := string(byts)

	astNode, err := parse(statement)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return
	}
	fmt.Printf("%s\n", tool.JsonifyIndent(extract(astNode)))
}

func TestStatementParser_Parse(t *testing.T) {
	sqlFilePath := "/Users/jacklingo/dev/model_infrax/assert/database.sql"

	configger, err := config.NewConfigger("/Users/jacklingo/dev/model_infrax/assert/application.yml")
	if err != nil {
		panic(err)
	}

	p, err := NewStatementParser(configger, sqlFilePath)
	if err != nil {
		panic(err)
	}

	schemas, err := p.Parse()
	if err != nil {
		panic(err)
	}

	log.Printf("schemas: %s", tool.JsonifyIndent(schemas))
}
