{{- /* 生成 GORM 模型结构体的模板 */ -}}
package {{ .PoPackageName }}

import (
	"encoding/json"
	"time"
)

{{- range $schema := .Schemas }}

// {{ $schema.Name | ToPascalCase }} {{ $schema.Comment }}
type {{ $schema.Name | ToPascalCase }} struct {
{{- range $schema.Columns }}
	{{ .ColumnName | ToPascalCase }} {{ . | GetGoType }} `gorm:"column:{{ .ColumnName }};type:{{ .Type }};{{ if .IsAutoIncrement }}primaryKey;autoIncrement;{{ end }}{{ if .Default }}default:{{ .Default }};{{ end }}comment:{{ .Comment }};{{ if not .IsNullable }}not null{{ end }}" json:"{{ .ColumnName }}"`
{{- end }}
}

// TableName 返回表名
func (t *{{ $schema.Name | ToPascalCase }}) TableName() string {
	return "{{ $schema.Name }}"
}

// Jsonify 将结构体序列化为 JSON 字符串（紧凑格式）
// 返回:
//   - string: JSON 字符串，如果序列化失败则返回错误信息的 JSON
func (t *{{ $schema.Name | ToPascalCase }}) Jsonify() string {
	byts, err := json.Marshal(t)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(byts)
}

// JsonifyIndent 将结构体序列化为格式化的 JSON 字符串（带缩进）
// 返回:
//   - string: 格式化的 JSON 字符串，如果序列化失败则返回错误信息的 JSON
func (t *{{ $schema.Name | ToPascalCase }}) JsonifyIndent() string {
	byts, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(byts)
}


// {{ $schema.Name | ToPascalCase }}Builder 用于构建 {{ $schema.Name | ToPascalCase }} 实例的 Builder
type {{ $schema.Name | ToPascalCase }}Builder struct {
	instance *{{ $schema.Name | ToPascalCase }}
}

// New{{ $schema.Name | ToPascalCase }}Builder 创建一个新的 {{ $schema.Name | ToPascalCase }}Builder 实例
// 返回:
//   - *{{ $schema.Name | ToPascalCase }}Builder: Builder 实例，用于链式调用
func New{{ $schema.Name | ToPascalCase }}Builder() *{{ $schema.Name | ToPascalCase }}Builder {
	return &{{ $schema.Name | ToPascalCase }}Builder{
		instance: &{{ $schema.Name | ToPascalCase }}{},
	}
}

{{- range $column := $schema.Columns }}
{{- if not $column.IsAutoIncrement }}

// With{{ $column.ColumnName | ToPascalCase }} 设置 {{ $column.ColumnName }} 字段
// 参数:
//   - {{ $column.ColumnName | ToSafeParamName }}: {{ $column.Comment }}
// 返回:
//   - *{{ $schema.Name | ToPascalCase }}Builder: 返回 Builder 实例，支持链式调用
func (b *{{ $schema.Name | ToPascalCase }}Builder) With{{ $column.ColumnName | ToPascalCase }}({{ $column.ColumnName | ToSafeParamName }} {{ $column | GetGoType }}) *{{ $schema.Name | ToPascalCase }}Builder {
	b.instance.{{ $column.ColumnName | ToPascalCase }} = {{ $column.ColumnName | ToSafeParamName }}
	return b
}
{{- if $column.IsNullable }}

// With{{ $column.ColumnName | ToPascalCase }}Value 设置 {{ $column.ColumnName }} 字段（便捷方法，自动转换为指针）
// 参数:
//   - {{ $column.ColumnName | ToSafeParamName }}: {{ $column.Comment }}
// 返回:
//   - *{{ $schema.Name | ToPascalCase }}Builder: 返回 Builder 实例，支持链式调用
func (b *{{ $schema.Name | ToPascalCase }}Builder) With{{ $column.ColumnName | ToPascalCase }}Value({{ $column.ColumnName | ToSafeParamName }} {{ $column | GetGoType | TrimPointer }}) *{{ $schema.Name | ToPascalCase }}Builder {
	b.instance.{{ $column.ColumnName | ToPascalCase }} = &{{ $column.ColumnName | ToSafeParamName }}
	return b
}
{{- end }}
{{- end }}
{{- end }}

// Build 构建并返回 {{ $schema.Name | ToPascalCase }} 实例
// 返回:
//   - *{{ $schema.Name | ToPascalCase }}: 构建完成的实例
func (b *{{ $schema.Name | ToPascalCase }}Builder) Build() *{{ $schema.Name | ToPascalCase }} {
	return b.instance
}

{{- end }}