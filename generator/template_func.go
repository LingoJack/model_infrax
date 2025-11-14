package generator

import (
	"model_infrax/model"
	"model_infrax/tool"
	"strings"
)

// typeMapping 定义了数据库类型到 Go 类型的映射规则
type typeMapping struct {
	// 匹配函数：判断数据库类型是否匹配
	matcher func(dbType string) bool
	// 非空时的 Go 类型
	goType string
	// 可空时的 Go 类型
	nullableGoType string
}

// typeMappings 定义所有类型映射规则
// 规则顺序很重要，优先匹配更具体的规则
var typeMappings = []typeMapping{
	// 1. bigint unsigned - 无符号大整数
	{
		matcher: func(dbType string) bool {
			return strings.Contains(dbType, "bigint") && strings.Contains(dbType, "unsigned")
		},
		goType:         "uint64",
		nullableGoType: "*uint64",
	},
	// 2. bigint - 有符号大整数
	{
		matcher: func(dbType string) bool {
			return strings.Contains(dbType, "bigint")
		},
		goType:         "int64",
		nullableGoType: "*int64",
	},
	// 3. tinyint(1) - 布尔类型
	{
		matcher: func(dbType string) bool {
			return strings.Contains(dbType, "tinyint(1)")
		},
		goType:         "bool",
		nullableGoType: "*bool",
	},
	// 4. tinyint - 小整数
	{
		matcher: func(dbType string) bool {
			return strings.Contains(dbType, "tinyint")
		},
		goType:         "int8",
		nullableGoType: "*int8",
	},
	// 5. int unsigned - 无符号整数
	{
		matcher: func(dbType string) bool {
			return strings.HasPrefix(dbType, "int") && strings.Contains(dbType, "unsigned")
		},
		goType:         "uint",
		nullableGoType: "*uint",
	},
	// 6. int - 有符号整数
	{
		matcher: func(dbType string) bool {
			return strings.HasPrefix(dbType, "int")
		},
		goType:         "int",
		nullableGoType: "*int",
	},
	// 7. float - 单精度浮点数
	{
		matcher: func(dbType string) bool {
			return strings.Contains(dbType, "float")
		},
		goType:         "float32",
		nullableGoType: "*float32",
	},
	// 8. double/decimal - 双精度浮点数
	{
		matcher: func(dbType string) bool {
			return strings.Contains(dbType, "double") || strings.Contains(dbType, "decimal")
		},
		goType:         "float64",
		nullableGoType: "*float64",
	},
	// 9. datetime/timestamp - 日期时间类型
	{
		matcher: func(dbType string) bool {
			return strings.Contains(dbType, "datetime") || strings.Contains(dbType, "timestamp")
		},
		goType:         "time.Time",
		nullableGoType: "*time.Time",
	},
	// 10. date - 日期类型
	{
		matcher: func(dbType string) bool {
			return strings.Contains(dbType, "date")
		},
		goType:         "time.Time",
		nullableGoType: "*time.Time",
	},
	// 11. varchar/char/text/blob - 字符串类型
	{
		matcher: func(dbType string) bool {
			return strings.Contains(dbType, "varchar") ||
				strings.Contains(dbType, "char") ||
				strings.Contains(dbType, "text") ||
				strings.Contains(dbType, "blob")
		},
		goType:         "string",
		nullableGoType: "*string",
	},
	// 12. json - JSON 类型
	{
		matcher: func(dbType string) bool {
			return strings.Contains(dbType, "json")
		},
		goType:         "string",
		nullableGoType: "*string",
	},
}

// GetGoType 根据列的数据库类型返回对应的 Go 类型
// 支持可空类型自动转换为指针类型
//
// 参数:
//   - col: 列定义，包含数据库类型和是否可空等信息
//
// 返回:
//   - string: 对应的 Go 类型，如 "uint64", "*string", "time.Time" 等
//
// 示例:
//   - bigint unsigned + 非空 -> uint64
//   - varchar(128) + 可空 -> *string
//   - datetime + 非空 -> time.Time
func GetGoType(col model.Column) string {
	// 将类型转换为小写便于比较
	dbType := strings.ToLower(col.Type)

	// 遍历所有映射规则，找到第一个匹配的规则
	for _, mapping := range typeMappings {
		if mapping.matcher(dbType) {
			// 根据是否可空返回对应的 Go 类型
			if col.IsNullable {
				return mapping.nullableGoType
			}
			return mapping.goType
		}
	}

	// 默认映射为 string 类型
	if col.IsNullable {
		return "*string"
	}
	return "string"
}

func ToPascalCase(s string) string {
	return tool.ToPascalCase(s)
}
