package tool

import (
	"fmt"

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

// ToSnakeCase 将字符串转换为 snake_case（下划线分隔）
// 使用 strcase 库实现，支持多种命名格式转换
// 参数:
//   - s: 待转换的字符串（支持大驼峰、小驼峰、短横线分隔或已有驼峰格式）
//
// 返回:
//   - string: 转换后的小驼峰格式字符串
//
// 示例:
//   - "userName" -> "user_name"
//   - "UserName" -> "user_name"
//   - "userId" -> "user_id"
//   - "user-name" -> "user_name"
func ToSnakeCase(s string) string {
	return strcase.ToSnake(s)
}

// InterfaceToString 将 interface{} 类型安全地转换为 string
// 参数:
//   - v: 任意类型的值
//
// 返回:
//   - string: 转换后的字符串，如果是 nil 则返回空字符串
//
// 示例:
//   - "hello" -> "hello"
//   - 123 -> "123"
//   - nil -> ""
func InterfaceToString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		return fmt.Sprintf("%v", v)
	}
}