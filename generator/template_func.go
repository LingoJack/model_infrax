package generator

import (
	"model_infrax/model"
	"regexp"
	"strings"
)

// ToPascalCase 将字符串转换为 PascalCase（大驼峰）
// 例如: t_artifact -> TArtifact, artifactId -> ArtifactID
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

// GetGoType 根据列信息返回对应的 Go 类型
func GetGoType(col model.Column) string {
	// 根据列名推断类型
	columnName := strings.ToLower(col.ColumnName)

	// 时间类型
	if strings.Contains(columnName, "time") || strings.Contains(columnName, "date") {
		if col.IsNullable {
			return "*time.Time"
		}
		return "time.Time"
	}

	// ID 类型通常是 uint64
	if strings.HasSuffix(columnName, "id") || columnName == "id" {
		if col.IsNullable {
			return "*uint64"
		}
		return "uint64"
	}

	// 步骤、数量等整数类型
	if strings.Contains(columnName, "step") || strings.Contains(columnName, "count") ||
		strings.Contains(columnName, "num") || strings.Contains(columnName, "status") {
		if col.IsNullable {
			return "*int"
		}
		return "int"
	}

	// 默认为字符串类型
	if col.IsNullable {
		return "*string"
	}
	return "string"
}

// GetMySQLType 根据列信息返回对应的 MySQL 类型
func GetMySQLType(col model.Column) string {

	// TODO 这里太刻意了

	columnName := strings.ToLower(col.ColumnName)

	// 时间类型
	if strings.Contains(columnName, "time") || strings.Contains(columnName, "date") {
		return "datetime"
	}

	// ID 类型
	if strings.HasSuffix(columnName, "id") || columnName == "id" {
		if col.IsAutoIncrement {
			return "bigint(20) unsigned"
		}
		return "bigint(20)"
	}

	// 步骤、状态等整数类型
	if strings.Contains(columnName, "step") || strings.Contains(columnName, "status") {
		return "int(11)"
	}

	// 内容字段使用 text
	if columnName == "content" || strings.Contains(columnName, "content") {
		return "text"
	}

	// 默认为 varchar
	return "varchar(128)"
}
