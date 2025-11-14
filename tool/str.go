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
