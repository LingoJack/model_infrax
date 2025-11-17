package parser

import (
	"fmt"
	"log"
	"model_infrax/config"
	"model_infrax/model"
	"model_infrax/tool"
	"os"
	"strings"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/parser/test_driver"
	"github.com/samber/lo"
)

type StatementParser struct {
	configger  *config.Configger
	statements []string
}

func NewStatementParser(cfg *config.Configger, sqlFilePath string) (*StatementParser, error) {
	path := tool.EscapeHomeDir(sqlFilePath)

	byts, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	statements := strings.Split(string(byts), ";")

	return &StatementParser{
		configger:  cfg,
		statements: statements,
	}, nil
}

func (p *StatementParser) Parse() (schemas []model.Schema, err error) {
	for _, statement := range p.statements {
		// 跳过空语句
		trimmed := strings.TrimSpace(statement)
		if trimmed == "" {
			continue
		}

		log.Printf("⌛️ parsing statement: %s", statement)
		var schema model.Schema
		schema, err = p.parseStatement(statement)
		if err != nil {
			return nil, fmt.Errorf("解析语句失败: %w", err)
		}
		schemas = append(schemas, schema)
	}
	return
}

func (p *StatementParser) FilterTables(schemas []model.Schema) (filtered []model.Schema) {
	if p.configger.GenerateConfig.AllTables {
		filtered = schemas
		return
	}
	filtered = lo.Filter(schemas, func(schema model.Schema, index int) bool {
		return lo.Contains(p.configger.GenerateConfig.TableNames, schema.Name)
	})
	return
}

// parseStatement 解析单个CREATE TABLE语句，提取表结构信息
func (p *StatementParser) parseStatement(statement string) (schema model.Schema, err error) {
	// 创建TiDB parser实例
	tidbParser := parser.New()

	// 解析SQL语句
	stmtNodes, _, err := tidbParser.ParseSQL(statement)
	if err != nil {
		return schema, fmt.Errorf("SQL解析失败: %w", err)
	}

	// 确保至少有一个语句节点
	if len(stmtNodes) == 0 {
		return schema, fmt.Errorf("未找到有效的SQL语句")
	}

	// 获取第一个语句节点
	stmtNode := stmtNodes[0]

	// 类型断言为CREATE TABLE语句
	createTableStmt, ok := stmtNode.(*ast.CreateTableStmt)
	if !ok {
		return schema, fmt.Errorf("不是CREATE TABLE语句")
	}

	// 提取表名
	schema.Name = createTableStmt.Table.Name.O

	// 提取表注释
	for _, option := range createTableStmt.Options {
		if option.Tp == ast.TableOptionComment {
			schema.Comment = option.StrValue
			break
		}
	}

	// 用于存储列名到列的映射，方便后续索引处理
	columnMap := make(map[string]*model.Column)

	// 提取列信息
	for _, col := range createTableStmt.Cols {
		column := model.Column{
			ColumnName: col.Name.Name.O,
			Type:       col.Tp.String(),
		}

		// 提取列的各种属性
		for _, option := range col.Options {
			switch option.Tp {
			case ast.ColumnOptionComment:
				// 提取列注释：优先使用Expr，如果Expr为nil则使用StrValue
				if option.Expr != nil {
					column.Comment = option.Expr.Text()
				}
				// 如果Expr方式没取到，尝试StrValue
				if column.Comment == "" && option.StrValue != "" {
					column.Comment = option.StrValue
				}
			case ast.ColumnOptionDefaultValue:
				// 提取默认值：使用Expr.Text()获取原始文本
				if option.Expr != nil {
					defaultVal := option.Expr.Text()
					// 去除可能的引号
					if len(defaultVal) >= 2 && defaultVal[0] == '\'' && defaultVal[len(defaultVal)-1] == '\'' {
						defaultVal = defaultVal[1 : len(defaultVal)-1]
					}
					column.Default = &defaultVal
				}
			case ast.ColumnOptionAutoIncrement:
				// 标记自增列
				column.IsAutoIncrement = true
			case ast.ColumnOptionNull:
				// 标记允许NULL
				column.IsNullable = true
			case ast.ColumnOptionNotNull:
				// 标记不允许NULL
				column.IsNullable = false
			case ast.ColumnOptionPrimaryKey:
				// 标记主键
				column.IsPrimaryKey = true
			case ast.ColumnOptionUniqKey:
				// 标记唯一键
				column.IsUnique = true
			}
		}

		// 提取字符集校对规则
		if col.Tp.GetCollate() != "" {
			column.Collate = col.Tp.GetCollate()
		}

		schema.Columns = append(schema.Columns, column)
		columnMap[column.ColumnName] = &schema.Columns[len(schema.Columns)-1]
	}
	for _, constraint := range createTableStmt.Constraints {
		switch constraint.Tp {
		case ast.ConstraintPrimaryKey:
			// 处理主键
			var pkColumns []model.Column
			for _, indexCol := range constraint.Keys {
				colName := indexCol.Column.Name.O
				if col, exists := columnMap[colName]; exists {
					col.IsPrimaryKey = true
					col.IsIndexed = true
					pkColumns = append(pkColumns, *col)
				}
			}
			schema.PrimaryKey = model.Index{
				IndexName: "PRIMARY",
				Columns:   pkColumns,
			}

		case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
			// 处理唯一索引
			var uniqueColumns []model.Column
			for _, indexCol := range constraint.Keys {
				colName := indexCol.Column.Name.O
				if col, exists := columnMap[colName]; exists {
					col.IsUnique = true
					col.IsIndexed = true
					uniqueColumns = append(uniqueColumns, *col)
				}
			}
			indexName := constraint.Name
			if indexName == "" {
				// 如果没有指定索引名，使用列名组合
				indexName = "uk_" + strings.Join(lo.Map(uniqueColumns, func(c model.Column, _ int) string {
					return c.ColumnName
				}), "_")
			}
			schema.UniqueIndex = append(schema.UniqueIndex, model.Index{
				IndexName: indexName,
				Columns:   uniqueColumns,
			})

		case ast.ConstraintKey, ast.ConstraintIndex:
			// 处理普通索引
			var indexColumns []model.Column
			for _, indexCol := range constraint.Keys {
				colName := indexCol.Column.Name.O
				if col, exists := columnMap[colName]; exists {
					col.IsIndexed = true
					indexColumns = append(indexColumns, *col)
				}
			}
			indexName := constraint.Name
			if indexName == "" {
				// 如果没有指定索引名，使用列名组合
				indexName = "idx_" + strings.Join(lo.Map(indexColumns, func(c model.Column, _ int) string {
					return c.ColumnName
				}), "_")
			}
			schema.Indexes = append(schema.Indexes, model.Index{
				IndexName: indexName,
				Columns:   indexColumns,
			})
		default:

		}
	}

	log.Printf("✅ 成功解析表: %s, 列数: %d, 索引数: %d", schema.Name, len(schema.Columns), len(schema.Indexes))
	return schema, nil
}
