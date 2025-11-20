package parser

import (
	"fmt"
	"log"

	"os"
	"testing"

	"github.com/LingoJack/model_infrax/config"
	"github.com/LingoJack/model_infrax/tool"
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

	byts, err := os.ReadFile("/Users/jacklingo/dev/model_infrax/assets/schema.sql")
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
	byts, err := os.ReadFile("/Users/jacklingo/dev/model_infrax/assets/schema.sql")
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
	// 加载配置文件
	configger, err := config.NewConfigger("/Users/jacklingo/dev/model_infrax/assets/application.yml")
	if err != nil {
		panic(err)
	}

	// 创建StatementParser，SQL文件路径从配置中读取
	p, err := NewStatementParser(configger)
	if err != nil {
		panic(err)
	}

	// 解析SQL文件
	schemas, err := p.Parse()
	if err != nil {
		panic(err)
	}

	log.Printf("schemas: %s", tool.JsonifyIndent(schemas))
}

// TestDebugAST 调试AST结构，查看列的Comment和Default值的实际存储方式
func TestDebugAST(t *testing.T) {
	p := parser.New()

	byts, err := os.ReadFile("/Users/jacklingo/dev/model_infrax/assets/schema.sql")
	if err != nil {
		panic(err)
	}

	statement := string(byts)

	stmtNodes, _, err := p.ParseSQL(statement)
	if err != nil {
		panic(err)
	}

	createTableStmt, ok := stmtNodes[0].(*ast.CreateTableStmt)
	if !ok {
		panic("不是CREATE TABLE语句")
	}

	log.Printf("表名: %s", createTableStmt.Table.Name.O)
	log.Printf("表选项数量: %d", len(createTableStmt.Options))

	// 打印表选项
	for i, option := range createTableStmt.Options {
		log.Printf("表选项[%d] Tp=%d, StrValue=%s", i, option.Tp, option.StrValue)
	}

	// 打印列信息
	for i, col := range createTableStmt.Cols {
		log.Printf("\n=== 列[%d]: %s ===", i, col.Name.Name.O)
		log.Printf("类型: %s", col.Tp.String())
		log.Printf("选项数量: %d", len(col.Options))

		for j, option := range col.Options {
			log.Printf("  选项[%d] Tp=%d", j, option.Tp)
			log.Printf("    StrValue: '%s'", option.StrValue)
			if option.Expr != nil {
				log.Printf("    Expr类型: %T", option.Expr)
				log.Printf("    Expr.Text(): '%s'", option.Expr.Text())
				// 打印完整的Expr结构，看看里面有什么字段
				log.Printf("    Expr完整结构: %#v", option.Expr)
			}
		}
	}
}
