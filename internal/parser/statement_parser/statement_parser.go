package statement_parser

import (
	"fmt"
	"strings"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/pkg/tool"
	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/test_driver"
	"github.com/samber/lo"
)

var (
	err        error
	sqlPath    string
	tableNames []string
	all        bool
)

func init() {
	all = conf.ValueBool("generate_config.all_tables")
	sqlPath = tool.EscapeHomeDir(conf.ValueStr("generate_config.sql_file_path"))
	tableNames, err = conf.ValueStrSlice("generate_config.table_names")
	if err != nil {
		panic(err)
	}
}

func SqlStatements() (statements []string, err error) {
	if !tool.IsValidFilePath(sqlPath) {
		err = fmt.Errorf("SQL文件路径无效: %s", sqlPath)
		return
	}
	statements = strings.Split(tool.MustReadText(sqlPath), ";")
	return
}

func Parse(statements []string) (schemas []model.Schema, err error) {
	for _, statement := range statements {
		if len(strings.TrimSpace(statement)) == 0 {
			continue
		}
		var schema model.Schema
		schema, err = ParseSqlStatement(statement)
		if err != nil {
			return nil, fmt.Errorf("解析语句失败: %w", err)
		}
		schemas = append(schemas, schema)
	}
	return
}

func Filter(schemas []model.Schema) (filtered []model.Schema) {
	if all {
		filtered = schemas
		return
	}
	logger.Infof("neededTableNames: %v", tableNames)
	filtered = lo.Filter(schemas, func(schema model.Schema, index int) bool {
		return lo.Contains(tableNames, schema.Name)
	})
	return
}

// ParseSqlStatement 解析单个CREATE TABLE语句，提取表结构信息
func ParseSqlStatement(statement string) (schema model.Schema, err error) {
	// 创建TiDB parser实例
	tidbParser := parser.New()

	// 解析 SQL 语句
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

	// 用于存储列名到列索引的映射，方便后续索引处理
	// 注意：不能使用指针映射，因为 slice 扩容会导致指针失效
	columnIndexMap := make(map[string]int)

	// 提取列信息
	for _, col := range createTableStmt.Cols {
		column := model.Column{
			ColumnName: col.Name.Name.O,
			Type:       col.Tp.String(),
			IsNullable: true, // MySQL默认列是可以为NULL的，除非显式声明NOT NULL
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
		// 存储列名到索引的映射，避免使用指针（slice 扩容会导致指针失效）
		columnIndexMap[column.ColumnName] = len(schema.Columns) - 1
	}
	for _, constraint := range createTableStmt.Constraints {
		switch constraint.Tp {
		case ast.ConstraintPrimaryKey:
			// 处理主键
			var pkColumns []model.Column
			for _, indexCol := range constraint.Keys {
				colName := indexCol.Column.Name.O
				if colIdx, exists := columnIndexMap[colName]; exists {
					// 通过索引直接修改 schema.Columns 中的列属性
					schema.Columns[colIdx].IsPrimaryKey = true
					schema.Columns[colIdx].IsIndexed = true
					pkColumns = append(pkColumns, schema.Columns[colIdx])
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
				if colIdx, exists := columnIndexMap[colName]; exists {
					// 通过索引直接修改 schema.Columns 中的列属性
					schema.Columns[colIdx].IsUnique = true
					schema.Columns[colIdx].IsIndexed = true
					uniqueColumns = append(uniqueColumns, schema.Columns[colIdx])
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
				if colIdx, exists := columnIndexMap[colName]; exists {
					// 通过索引直接修改 schema.Columns 中的列属性
					schema.Columns[colIdx].IsIndexed = true
					indexColumns = append(indexColumns, schema.Columns[colIdx])
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

	logger.Infof("✅ 成功解析表: %s, 列数: %d, 索引数: %d", schema.Name, len(schema.Columns), len(schema.Indexes))
	return schema, nil
}
