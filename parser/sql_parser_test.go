package parser

import (
	"fmt"
	"model_infrax/config"
	"model_infrax/tool"
	"os"
	"testing"
)

// loadSQL 加载SQL文件内容
func loadSQL() string {
	byts, err := os.ReadFile("/Users/jacklingo/dev/model_infrax/assert/database.sql")
	if err != nil {
		panic(fmt.Sprintf("读取SQL文件失败: %v", err))
	}
	return string(byts)
}

func TestParser_AllTables(t *testing.T) {
	configger, err := config.NewConfigger("/Users/jacklingo/dev/model_infrax/assert/application.yml")
	if err != nil {
		panic(err)
	}
	parser, err := NewParser(configger)
	if err != nil {
		panic(err)
	}
	tables, err := parser.AllTables()
	if err != nil {
		panic(err)
	}
	fmt.Println(tool.JsonifyIndent(tables))
}

func TestSql1(t *testing.T) {
	configger, err := config.NewConfigger("/Users/jacklingo/dev/model_infrax/assert/application.yml")
	if err != nil {
		panic(err)
	}
	parser, err := NewParser(configger)
	if err != nil {
		panic(err)
	}
	var tables []mysqlTable
	err = parser.db.Raw("show table status").Scan(&tables).Error
	if err != nil {
		panic(err)
	}

	fmt.Println(tool.JsonifyIndent(tables))
}

func TestSql2(t *testing.T) {
	configger, err := config.NewConfigger("/Users/jacklingo/dev/model_infrax/assert/application.yml")
	if err != nil {
		panic(err)
	}
	parser, err := NewParser(configger)
	if err != nil {
		panic(err)
	}
	var table map[string]interface{}
	err = parser.db.Raw(fmt.Sprintf("show table status like '%s'", "t_artifact")).Scan(&table).Error
	if err != nil {
		panic(err)
	}

	fmt.Println(tool.JsonifyIndent(table))
}

func TestSql3(t *testing.T) {
	configger, err := config.NewConfigger("/Users/jacklingo/dev/model_infrax/assert/application.yml")
	if err != nil {
		panic(err)
	}
	parser, err := NewParser(configger)
	if err != nil {
		panic(err)
	}
	var table []map[string]interface{}
	err = parser.db.Raw(fmt.Sprintf("show full fields from %s", "t_artifact")).Scan(&table).Error
	if err != nil {
		panic(err)
	}

	fmt.Println(tool.JsonifyIndent(table))
}

func TestSql4(t *testing.T) {
	configger, err := config.NewConfigger("/Users/jacklingo/dev/model_infrax/assert/application.yml")
	if err != nil {
		panic(err)
	}
	parser, err := NewParser(configger)
	if err != nil {
		panic(err)
	}
	var table []map[string]interface{}
	err = parser.db.Raw(fmt.Sprintf("show index from %s", "t_artifact")).Scan(&table).Error
	if err != nil {
		panic(err)
	}

	fmt.Println(tool.JsonifyIndent(table))
}
