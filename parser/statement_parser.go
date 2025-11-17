package parser

import (
	"fmt"
	"log"
	"model_infrax/config"
	"model_infrax/model"
	"os"
	"strings"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/test_driver"
	"github.com/samber/lo"
)

type StatementParser struct {
	configger  *config.Configger
	statements []string
}

// NewStatementParser 创建SQL语句解析器
// 从配置文件中读取SQL文件路径，解析SQL文件内容
func NewStatementParser(cfg *config.Configger) (*StatementParser, error) {
	// 从配置中获取SQL文件路径
	sqlFilePath := cfg.GenerateConfig.SqlFilePath
	if sqlFilePath == "" {
		return nil, fmt.Errorf("statement模式下必须配置sql_file_path")
	}

	// 读取SQL文件内容
	byts, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return nil, fmt.Errorf("读取SQL文件失败 [%s]: %w", sqlFilePath, err)
	}

	// 按分号分割SQL语句
	statements := strings.Split(string(byts), ";")

	log.Printf("📄 成功加载SQL文件: %s, 共 %d 条语句", sqlFilePath, len(statements))

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
				// 提取列注释：从ValueExpr的Datum.b字段中读取UTF-8编码的字节数组
				if option.Expr != nil {
					// 尝试类型断言为test_driver.ValueExpr
					if valueExpr, ok := option.Expr.(*test_driver.ValueExpr); ok {
						// Datum.b 存储的是UTF-8编码的字节数组
						if len(valueExpr.Datum.GetBytes()) > 0 {
							column.Comment = string(valueExpr.Datum.GetBytes())
						}
					}
				}
				// 如果Expr方式没取到，尝试StrValue（兼容处理）
				if column.Comment == "" && option.StrValue != "" {
					column.Comment = option.StrValue
				}
			case ast.ColumnOptionDefaultValue:
				// 提取默认值：需要区分ValueExpr（字符串/数值）和FuncCallExpr（函数如CURRENT_TIMESTAMP）
				if option.Expr != nil {
					var defaultVal string

					// 处理ValueExpr类型（字符串或数值默认值）
					if valueExpr, ok := option.Expr.(*test_driver.ValueExpr); ok {
						// 如果Datum.b为空字节数组，表示空字符串''
						if len(valueExpr.Datum.GetBytes()) == 0 {
							defaultVal = ""
						} else {
							// 否则转换字节数组为字符串
							defaultVal = string(valueExpr.Datum.GetBytes())
						}
					} else if funcExpr, ok := option.Expr.(*ast.FuncCallExpr); ok {
						// 处理FuncCallExpr类型（如CURRENT_TIMESTAMP）
						defaultVal = funcExpr.FnName.O
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
