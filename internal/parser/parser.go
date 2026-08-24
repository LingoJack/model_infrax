package parser

import (
	"fmt"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/samber/lo"

	"github.com/LingoJack/model_infrax/internal/model"
)

// 两种解析模式共用的表过滤配置
var (
	filterAll        bool
	filterTableNames []string
)

// loadFilterConfig 从配置文件读取表过滤配置
// 必须在 conf.InitWithPath 之后调用
func loadFilterConfig() error {
	filterAll = conf.ValueBool("generate_config.all_tables")
	var err error
	filterTableNames, err = conf.ValueStrSlice("generate_config.table_names")
	if err != nil {
		return fmt.Errorf("[loadFilterConfig] 读取表名配置失败: %w", err)
	}
	return nil
}

// FilterTables 根据配置过滤表（all_tables 为 true 时返回全部）
func FilterTables(schemas []model.Schema) (filtered []model.Schema) {
	if filterAll {
		return schemas
	}
	return lo.Filter(schemas, func(schema model.Schema, _ int) bool {
		return lo.Contains(filterTableNames, schema.Name)
	})
}
