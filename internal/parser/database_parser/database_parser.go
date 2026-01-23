package database_parser

import (
	"context"
	"fmt"

	"strings"

	"github.com/LingoJack/model_infrax/internal/infra/db_infra"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/pkg/conf"
	"github.com/LingoJack/model_infrax/pkg/tool"
	"github.com/samber/lo"
)

var (
	err        error
	all        bool
	tableNames []string
)

func init() {
	all = conf.ValueBool("generate_config.all_tables")
	tableNames, err = conf.ValueStrSlice("generate_config.table_names")
	if err != nil {
		panic(err)
	}
}

// Parse 获取所有表名
func Parse(ctx context.Context) (schemas []model.Schema, err error) {
	var tables []mysqlTable
	err = db_infra.ExecSql(ctx, &tables, "show table status")
	if err != nil {
		return nil, fmt.Errorf("查询数据库表失败: %w", err)
	}

	// 使用 lo.Map 遍历所有表并构建 Schema
	schemas = lo.Map(tables, func(table mysqlTable, index int) model.Schema {
		tableName := table.Name
		tableComment := table.Comment

		// 查询表的所有字段信息
		var fields []mysqlField
		err = db_infra.ExecSql(ctx, &fields, "show full fields from ?", tableName)
		if err != nil {
			panic(err)
		}

		// 构建列信息和列名到列的映射
		var columns []model.Column
		name2Column := make(map[string]model.Column) // 初始化 map
		lo.ForEach(fields, func(field mysqlField, index int) {
			column := model.Column{
				ColumnName:      field.Field,
				Collate:         tool.Stringify(field.Collation), // Stringify 已经处理了 nil 指针
				Comment:         field.Comment,
				Type:            field.Type,
				Default:         field.Default, // 设置默认值
				IsAutoIncrement: strings.Contains(field.Extra, "auto_increment"),
				IsNullable:      field.Null == "YES",
			}
			columns = append(columns, column)
			name2Column[column.ColumnName] = column
		})

		// 查询表的所有索引信息
		var mysqlIndexes []mysqlIndex
		err = db_infra.ExecSql(ctx, &mysqlIndexes, "show index from ?", tableName)
		if err != nil {
			panic(err)
		}

		// 构建索引名到列的映射
		indexName2Columns := make(map[string][]model.Column)
		lo.ForEach(mysqlIndexes, func(index mysqlIndex, i int) {
			indexName2Columns[index.KeyName] = append(indexName2Columns[index.KeyName], name2Column[index.ColumnName])
		})

		// 构建索引名到索引对象的映射，用于后续查找
		inexName2Index := make(map[string]model.Index)

		var primaryKey model.Index
		var indexes []model.Index

		// 将 map 转换为 Index 切片
		indexes = lo.MapToSlice(indexName2Columns, func(key string, value []model.Column) model.Index {
			idx := model.Index{
				IndexName: key,
				Columns:   value,
			}
			// 识别主键
			if key == "PRIMARY" {
				primaryKey = idx
			}
			inexName2Index[key] = idx
			return idx
		})

		// 提取唯一索引
		var uniqueIndexes []model.Index
		lo.ForEach(mysqlIndexes, func(index mysqlIndex, i int) {
			if index.NonUnique == 0 {
				uniqueIndexes = append(uniqueIndexes, inexName2Index[index.KeyName])
			}
		})

		// 找到所有索引的列
		var indexedColumns []model.Column
		lo.ForEach(indexes, func(index model.Index, i int) {
			indexedColumns = append(indexedColumns, index.Columns...)
		})

		// 找到所有唯一索引的列
		var uniqueIndexedColumns []model.Column
		lo.ForEach(uniqueIndexes, func(index model.Index, i int) {
			uniqueIndexedColumns = append(uniqueIndexedColumns, index.Columns...)
		})

		// 找到所有主键的列
		var primaryKeyColumns []model.Column
		lo.ForEach(primaryKey.Columns, func(column model.Column, i int) {
			primaryKeyColumns = append(primaryKeyColumns, column)
		})

		// 标记索引列：直接在 columns 切片中更新
		for i := range columns {
			columnName := columns[i].ColumnName

			// 检查是否为索引列
			for _, indexedCol := range indexedColumns {
				if indexedCol.ColumnName == columnName {
					columns[i].IsIndexed = true
					break
				}
			}

			// 检查是否为唯一索引列
			for _, uniqueCol := range uniqueIndexedColumns {
				if uniqueCol.ColumnName == columnName {
					columns[i].IsUnique = true
					break
				}
			}

			// 检查是否为主键列
			for _, pkCol := range primaryKeyColumns {
				if pkCol.ColumnName == columnName {
					columns[i].IsPrimaryKey = true
					break
				}
			}
		}

		// 构建 Schema 对象
		return model.Schema{
			Name:        tableName,
			Comment:     tableComment,
			Columns:     columns,
			PrimaryKey:  primaryKey,
			Indexes:     indexes,
			UniqueIndex: uniqueIndexes,
		}
	})

	return
}

// Filter 根据配置文件过滤表
func Filter(schemas []model.Schema) (filtered []model.Schema) {
	if all {
		filtered = schemas
		return
	}
	filtered = lo.Filter(schemas, func(schema model.Schema, index int) bool {
		return lo.Contains(tableNames, schema.Name)
	})
	return
}
