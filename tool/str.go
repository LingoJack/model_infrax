package tool

import (
	"regexp"
	"strings"
)

// ToPascalCase 将字符串转换为 PascalCase（大驼峰）
func ToPascalCase(s string) string {
	// 处理下划线分隔的情况
	if strings.Contains(s, "_") {
		parts := strings.Split(s, "_")
		for i, part := range parts {
			if part != "" {
				parts[i] = strings.ToUpper(part[:1]) + part[1:]
			}
		}
		return strings.Join(parts, "")
	}

	// 处理驼峰命名的情况
	// 使用正则表达式找到所有大写字母的位置
	re := regexp.MustCompile(`([A-Z])`)
	result := re.ReplaceAllString(s, " $1")
	result = strings.TrimSpace(result)

	// 将首字母大写
	if len(result) > 0 {
		result = strings.ToUpper(result[:1]) + result[1:]
	}

	// 移除所有空格
	result = strings.ReplaceAll(result, " ", "")

	return result
}

// ToCamelCase 将字符串转换为 camelCase（小驼峰）
// 参数:
//   - s: 待转换的字符串（支持下划线分隔或已有驼峰格式）
//
// 返回:
//   - string: 转换后的小驼峰格式字符串
//
// 示例:
//   - "user_name" -> "userName"
//   - "UserName" -> "userName"
//   - "user_id" -> "userId"
func ToCamelCase(s string) string {
	// 先转换为大驼峰
	pascal := ToPascalCase(s)

	// 如果字符串为空，直接返回
	if len(pascal) == 0 {
		return pascal
	}

	// 将首字母转换为小写
	return strings.ToLower(pascal[:1]) + pascal[1:]
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
//   - TrimPrefix("[]byte", "[]") -> "byte"
func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}
