package model

import (
	"encoding/json"
	"fmt"
)

// Column 数据库列的元数据信息
type Column struct {
	ColumnName      string  // 列名
	Collate         string  // 字符集校对规则
	Comment         string  // 列注释
	Type            string  // 列类型
	Default         *string // 默认值（可能为null）
	IsAutoIncrement bool    // 是否自增
	IsNullable      bool    // 是否允许为NULL
	IsIndexed       bool    // 是否有索引
	IsUnique        bool    // 是否唯一索引
	IsPrimaryKey    bool    // 是否主键
}

func (f Column) Json() string {
	byts, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}

func (f Column) JsonIndent() string {
	byts, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}