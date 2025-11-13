package parser

import (
	"fmt"
	"model_infrax/config"
	"model_infrax/model"
	"model_infrax/tool"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/samber/lo"
)

// Parser SQL解析器，用于解析数据库表结构
type Parser struct {
	db        *gorm.DB
	configger *config.Configger
}

const (
	// 获取单个表的描述
	sqlShowTableStatusForSingleTable = "show table status like '%s'"

	// 查看所有字段
	sqlShowTableFields = "show full fields from %s"

	// 查看所有索引
	sqlShowTableIndexes = "show full indexes from %s"
)

func NewParser(cfg *config.Configger) (*Parser, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.GenerateConfig.Username,
		cfg.GenerateConfig.Password,
		cfg.GenerateConfig.Host,
		cfg.GenerateConfig.Port,
		cfg.GenerateConfig.DatabaseName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	return &Parser{
		db:        db,
		configger: cfg,
	}, nil
}

type mysqlTable struct {
	Name          string  `json:"Name"`            // 表名
	Engine        string  `json:"Engine"`          // 存储引擎
	Version       int     `json:"Version"`         // 版本号
	RowFormat     string  `json:"Row_format"`      // 行格式
	Rows          int64   `json:"Rows"`            // 行数（估算值）
	AvgRowLength  int64   `json:"Avg_row_length"`  // 平均行长度
	DataLength    int64   `json:"Data_length"`     // 数据长度
	MaxDataLength int64   `json:"Max_data_length"` // 最大数据长度
	IndexLength   int64   `json:"Index_length"`    // 索引长度
	DataFree      int64   `json:"Data_free"`       // 空闲空间
	AutoIncrement *int64  `json:"Auto_increment"`  // 自增值（可能为null）
	CreateTime    *string `json:"Create_time"`     // 创建时间（可能为null）
	UpdateTime    *string `json:"Update_time"`     // 更新时间（可能为null）
	CheckTime     *string `json:"Check_time"`      // 检查时间（可能为null）
	Collation     string  `json:"Collation"`       // 字符集校对规则
	Checksum      *int64  `json:"Checksum"`        // 校验和（可能为null）
	CreateOptions string  `json:"Create_options"`  // 创建选项
	Comment       string  `json:"Comment"`         // 表注释
}

type mysqlField struct {
	Field      string  `json:"Field"`      // 字段名
	Type       string  `json:"Type"`       // 字段类型
	Collation  *string `json:"Collation"`  // 字符集校对规则（可能为null）
	Null       string  `json:"Null"`       // 是否允许为NULL（YES/NO）
	Key        string  `json:"Key"`        // 键类型（PRI/UNI/MUL等）
	Default    *string `json:"Default"`    // 默认值（可能为null）
	Extra      string  `json:"Extra"`      // 额外信息（如auto_increment）
	Privileges string  `json:"Privileges"` // 权限信息
	Comment    string  `json:"Comment"`    // 字段注释
}

// mysqlIndex MySQL索引信息结构体
type mysqlIndex struct {
	Table        string  `json:"Table"`         // 表名
	NonUnique    int     `json:"Non_unique"`    // 是否非唯一索引（0=唯一索引，1=非唯一索引）
	KeyName      string  `json:"Key_name"`      // 索引名称
	SeqInIndex   int     `json:"Seq_in_index"`  // 字段在索引中的序号（从1开始）
	ColumnName   string  `json:"Column_name"`   // 列名
	Collation    *string `json:"Collation"`     // 排序方式（A=升序，D=降序，NULL=未排序）
	Cardinality  *int64  `json:"Cardinality"`   // 索引中唯一值的数量估算（可能为null）
	SubPart      *int    `json:"Sub_part"`      // 索引前缀长度（可能为null）
	Packed       *string `json:"Packed"`        // 关键字如何被压缩（可能为null）
	Null         string  `json:"Null"`          // 列是否可以包含NULL值
	IndexType    string  `json:"Index_type"`    // 索引类型（BTREE/HASH/FULLTEXT/SPATIAL）
	Comment      string  `json:"Comment"`       // 索引注释
	IndexComment string  `json:"Index_comment"` // 索引注释（创建索引时的COMMENT）
	Visible      string  `json:"Visible"`       // 索引是否可见（YES/NO）
	Expression   *string `json:"Expression"`    // 表达式索引的表达式（可能为null）
}

// AllTables 获取所有表名
func (p *Parser) AllTables() (schemas []model.Schema, err error) {
	// 执行 sqlShowTableStatus 拿到所有返回值
	var tables []mysqlTable
	if err = p.db.Raw("show table status").Scan(&tables).Error; err != nil {
		return nil, fmt.Errorf("查询数据库表失败: %w", err)
	}

	// 使用 lo.Map 遍历所有表并构建 Schema
	schemas = lo.Map(tables, func(table mysqlTable, index int) model.Schema {
		tableName := table.Name
		tableComment := table.Comment

		// 查询表的所有字段信息
		var fields []mysqlField
		err = p.db.Raw(fmt.Sprintf("show full fields from `%s`", tableName)).Scan(&fields).Error
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
				IsAutoIncrement: strings.Contains(field.Extra, "auto_increment"),
				IsNullable:      field.Null == "YES",
			}
			columns = append(columns, column)
			name2Column[column.ColumnName] = column
		})

		// 查询表的所有索引信息
		var mysqlIndexes []mysqlIndex
		err = p.db.Raw(fmt.Sprintf("show index from `%s`", tableName)).Scan(&mysqlIndexes).Error
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
				IndexName: model.IndexName(key),
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

		// 构建 Schema 对象
		schema := model.Schema{
			Name:        tableName,
			Comment:     tableComment,
			Columns:     columns,
			PrimaryKey:  primaryKey,
			Indexes:     indexes,
			UniqueIndex: uniqueIndexes,
		}

		return schema
	})

	return
}

func (p *Parser) FilterTables(schemas []model.Schema) (filtered []model.Schema) {
	if p.configger.GenerateConfig.AllTables {
		filtered = schemas
		return
	}
	filtered = lo.Filter(schemas, func(schema model.Schema, index int) bool {
		return lo.Contains(p.configger.GenerateConfig.TableNames, schema.Name)
	})
	return
}
