package tool

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// ToPascalCase 将字符串转换为 PascalCase（大驼峰）
// 使用 strcase 库实现，支持多种命名格式转换
// 参数:
//   - s: 待转换的字符串（支持下划线分隔、短横线分隔或已有驼峰格式）
//
// 返回:
//   - string: 转换后的大驼峰格式字符串
//
// 示例:
//   - "user_name" -> "UserName"
//   - "userName" -> "UserName"
//   - "user-name" -> "UserName"
func ToPascalCase(s string) string {
	return strcase.ToCamel(s)
}

// ToCamelCase 将字符串转换为 camelCase（小驼峰）
// 使用 strcase 库实现，支持多种命名格式转换
// 参数:
//   - s: 待转换的字符串（支持下划线分隔、短横线分隔或已有驼峰格式）
//
// 返回:
//   - string: 转换后的小驼峰格式字符串
//
// 示例:
//   - "user_name" -> "userName"
//   - "UserName" -> "userName"
//   - "user_id" -> "userId"
//   - "user-name" -> "userName"
func ToCamelCase(s string) string {
	return strcase.ToLowerCamel(s)
}

// goKeywords Go 语言关键字集合
var goKeywords = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"defer":       true,
	"else":        true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
}

// ToSafeParamName 将字符串转换为安全的参数名（避免 Go 关键字冲突）
// 参数:
//   - s: 待转换的字符串
//
// 返回:
//   - string: 安全的参数名，如果是 Go 关键字则添加 "Val" 后缀
//
// 示例:
//   - "interface" -> "interfaceVal"
//   - "type" -> "typeVal"
//   - "user_name" -> "userName"
func ToSafeParamName(s string) string {
	// 先转换为小驼峰
	camelCase := ToCamelCase(s)

	// 检查是否为 Go 关键字
	if goKeywords[camelCase] {
		return camelCase + "Val"
	}

	return camelCase
}

// TrimPrefix 去除字符串的前缀
// 参数:
//   - s: 原始字符串
//   - prefix: 要去除的前缀
//
// 返回:
//   - string: 去除前缀后的字符串
//
// 示例:
//   - TrimPrefix("*string", "*") -> "string"
//   - TrimPrefix("[] byte", "[] ") -> "byte"
func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// ToJsonTag 根据字段列表的整体命名风格生成 JSON tag
// 通过判断所有字段的命名风格（而非单个字段）来决定后缀的拼接方式
// 如果整体是蛇形命名风格，则后缀也转换为蛇形并用下划线连接
// 如果整体是驼峰命名风格，则后缀直接拼接保持驼峰风格
// 参数:
//   - fieldName: 当前字段名
//   - suffix: 后缀（如 "List", "Fuzzy", "Start", "End"）
//   - allFieldNames: 所有字段名列表，用于判断整体命名风格
//
// 返回:
//   - string: 生成的 JSON tag
//
// 示例:
//   - ToJsonTag("user_id", "List", []string{"id", "user_id", "created_at"}) -> "user_id_list"
//   - ToJsonTag("userId", "List", []string{"id", "userId", "createdAt"}) -> "userIdList"
//   - ToJsonTag("id", "List", []string{"id", "user_id", "created_at"}) -> "id_list"
//   - ToJsonTag("id", "List", []string{"id", "userId", "createdAt"}) -> "idList"
func ToJsonTag(fieldName, suffix string, allFieldNames []string) string {
	// 判断整体命名风格
	if IsSnakeCaseStyle(allFieldNames) {
		// 蛇形命名：将后缀转换为蛇形并用下划线连接
		return fieldName + "_" + strcase.ToSnake(suffix)
	}
	// 驼峰命名：直接拼接后缀（保持驼峰风格）
	return fieldName + suffix
}

// IsSnakeCaseStyle 判断字段列表的整体命名风格是否为蛇形命名
// 遍历所有字段名（排除 "id" 这种特殊字段），如果发现有字段包含下划线，则认为是蛇形命名风格
// 参数:
//   - fieldNames: 字段名列表
//
// 返回:
//   - bool: true 表示蛇形命名风格，false 表示驼峰命名风格
//
// 示例:
//   - IsSnakeCaseStyle([]string{"id", "user_id", "created_at"}) -> true
//   - IsSnakeCaseStyle([]string{"id", "userId", "createdAt"}) -> false
//   - IsSnakeCaseStyle([]string{"id"}) -> false (默认驼峰)
func IsSnakeCaseStyle(fieldNames []string) bool {
	// 遍历所有字段名
	for _, fieldName := range fieldNames {
		// 排除 id 这种特殊字段，因为它无法判断命名风格
		if fieldName == "id" {
			continue
		}
		// 如果发现有字段包含下划线，说明是蛇形命名
		if strings.Contains(fieldName, "_") {
			return true
		}
	}
	// 如果没有发现下划线，默认为驼峰命名
	return false
}
