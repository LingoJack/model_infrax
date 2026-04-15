package generator

import (
	"strings"

	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/pkg/tool"
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

// ToCamelCase 将字符串转换为 camelCase（小驼峰）
// 用于模板中生成参数名等需要小驼峰格式的场景
func ToCamelCase(s string) string {
	return tool.ToCamelCase(s)
}

// ToSafeParamName 将字符串转换为安全的参数名（避免 Go 关键字冲突）
// 用于模板中生成函数参数名，如果是 Go 关键字则添加后缀
func ToSafeParamName(s string) string {
	return tool.ToSafeParamName(s)
}

// TrimPrefix 去除字符串的前缀
// 用于模板中去除类型前面的指针符号等
// 参数:
//   - s: 原始字符串
//   - prefix: 要去除的前缀
//
// 返回:
//   - string: 去除前缀后的字符串
//
// 示例:
//   - TrimPrefix("*string", "*") -> "string"
func TrimPrefix(s, prefix string) string {
	return tool.TrimPrefix(s, prefix)
}

// TrimPointer 去除类型前面的指针符号 *
// 专门用于模板中将指针类型转换为值类型
// 参数:
//   - s: 原始类型字符串（如 "*string", "*int64"）
//
// 返回:
//   - string: 去除指针符号后的类型（如 "string", "int64"）
//
// 示例:
//   - TrimPointer("*string") -> "string"
//   - TrimPointer("*int64") -> "int64"
//   - TrimPointer("string") -> "string" (如果没有指针符号，返回原值)
func TrimPointer(s string) string {
	return tool.TrimPrefix(s, "*")
}

// ToJsonTag 根据列定义列表的整体命名风格生成 JSON tag
// 用于模板中为带后缀的字段生成合适的 JSON tag
// 参数:
//   - fieldName: 当前字段名
//   - suffix: 后缀（如 "List", "Fuzzy", "Start", "End"）
//   - columns: 列定义列表，用于判断整体命名风格
//
// 返回:
//   - string: 生成的 JSON tag
//
// 示例:
//   - ToJsonTag("user_id", "List", columns) -> "user_id_list" (蛇形风格)
//   - ToJsonTag("userId", "List", columns) -> "userIdList" (驼峰风格)
func ToJsonTag(fieldName, suffix string, columns []model.Column) string {
	// 提取所有字段名
	fieldNames := make([]string, len(columns))
	for i, col := range columns {
		fieldNames[i] = col.ColumnName
	}
	// 调用 tool 包的函数生成 JSON tag
	return tool.ToJsonTag(fieldName, suffix, fieldNames)
}

// HasTimeColumnInSchemas 判断 schemas 中是否存在 time.Time / *time.Time 类型的列
// 用于模板中决定是否需要 import "time"
func HasTimeColumnInSchemas(schemas []model.Schema) bool {
	for _, schema := range schemas {
		for _, col := range schema.Columns {
			goType := GetGoType(col)
			if strings.Contains(goType, "time.Time") {
				return true
			}
		}
	}
	return false
}

// IsSnakeCaseStyle 判断列定义列表的整体命名风格是否为蛇形命名
// 用于模板中判断应该生成蛇形还是驼峰风格的 JSON tag
// 参数:
//   - columns: 列定义列表
//
// 返回:
//   - bool: true 表示蛇形命名风格，false 表示驼峰命名风格
//
// 示例:
//   - IsSnakeCaseStyle(columns) -> true (如果字段包含 user_id, created_at 等)
//   - IsSnakeCaseStyle(columns) -> false (如果字段包含 userId, createdAt 等)
func IsSnakeCaseStyle(columns []model.Column) bool {
	// 提取所有字段名
	fieldNames := make([]string, len(columns))
	for i, col := range columns {
		fieldNames[i] = col.ColumnName
	}
	// 调用 tool 包的函数判断命名风格
	return tool.IsSnakeCaseStyle(fieldNames)
}
