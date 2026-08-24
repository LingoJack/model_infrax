{{- /* 生成 GORM 模型结构体的模板 */ -}}
package {{ .PoPackageName }}

import (
	"encoding/json"
{{- if HasTimeColumnInSchemas .Schemas }}
	"time"
{{- end }}
)

{{- range $schema := .Schemas }}

// {{ $schema.Name | ToPascalCase }} {{ $schema.Comment }}
type {{ $schema.Name | ToPascalCase }} struct {
{{- range $schema.Columns }}
	{{ .ColumnName | ToPascalCase }} {{ . | GetGoType }} `gorm:"{{ BuildGormTag . $.GormTagStyle }}" json:"{{ .ColumnName }}"`
{{- end }}
}

{{ if ne $.CommentStyle "none" }}// TableName 返回表名
{{- end }}
func (t *{{ $schema.Name | ToPascalCase }}) TableName() string {
	return "{{ $schema.Name }}"
}

{{- if $.Methods.PoJsonify }}

{{ if ne $.CommentStyle "none" }}// Jsonify 将结构体序列化为 JSON 字符串（紧凑格式）
{{- if eq $.CommentStyle "full" }}
// 返回:
//   - string: JSON 字符串，如果序列化失败则返回错误信息的 JSON
{{- end }}
{{- end }}
func (t *{{ $schema.Name | ToPascalCase }}) Jsonify() string {
	byts, err := json.Marshal(t)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(byts)
}

{{ if ne $.CommentStyle "none" }}// JsonifyIndent 将结构体序列化为格式化的 JSON 字符串（带缩进）
{{- if eq $.CommentStyle "full" }}
// 返回:
//   - string: 格式化的 JSON 字符串，如果序列化失败则返回错误信息的 JSON
{{- end }}
{{- end }}
func (t *{{ $schema.Name | ToPascalCase }}) JsonifyIndent() string {
	byts, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(byts)
}

{{- end }}

{{- if $.Methods.PoBuilder }}

{{ if ne $.CommentStyle "none" }}// {{ $schema.Name | ToPascalCase }}Builder 用于构建 {{ $schema.Name | ToPascalCase }} 实例的 Builder
{{- end }}
type {{ $schema.Name | ToPascalCase }}Builder struct {
	instance *{{ $schema.Name | ToPascalCase }}
}

{{ if ne $.CommentStyle "none" }}// New{{ $schema.Name | ToPascalCase }}Builder 创建一个新的 {{ $schema.Name | ToPascalCase }}Builder 实例
{{- if eq $.CommentStyle "full" }}
// 返回:
//   - *{{ $schema.Name | ToPascalCase }}Builder: Builder 实例，用于链式调用
{{- end }}
{{- end }}
func New{{ $schema.Name | ToPascalCase }}Builder() *{{ $schema.Name | ToPascalCase }}Builder {
	return &{{ $schema.Name | ToPascalCase }}Builder{
		instance: &{{ $schema.Name | ToPascalCase }}{},
	}
}

{{- range $column := $schema.Columns }}
{{- if not $column.IsAutoIncrement }}

{{ if ne $.CommentStyle "none" }}// With{{ $column.ColumnName | ToPascalCase }} 设置 {{ $column.ColumnName }} 字段
{{- if eq $.CommentStyle "full" }}
// 参数:
//   - {{ $column.ColumnName | ToSafeParamName }}: {{ $column.Comment }}
// 返回:
//   - *{{ $schema.Name | ToPascalCase }}Builder: 返回 Builder 实例，支持链式调用
{{- end }}
{{- end }}
func (b *{{ $schema.Name | ToPascalCase }}Builder) With{{ $column.ColumnName | ToPascalCase }}({{ $column.ColumnName | ToSafeParamName }} {{ $column | GetGoType }}) *{{ $schema.Name | ToPascalCase }}Builder {
	b.instance.{{ $column.ColumnName | ToPascalCase }} = {{ $column.ColumnName | ToSafeParamName }}
	return b
}
{{- if $column.IsNullable }}

{{ if ne $.CommentStyle "none" }}// With{{ $column.ColumnName | ToPascalCase }}Value 设置 {{ $column.ColumnName }} 字段（便捷方法，自动转换为指针）
{{- if eq $.CommentStyle "full" }}
// 参数:
//   - {{ $column.ColumnName | ToSafeParamName }}: {{ $column.Comment }}
// 返回:
//   - *{{ $schema.Name | ToPascalCase }}Builder: 返回 Builder 实例，支持链式调用
{{- end }}
{{- end }}
func (b *{{ $schema.Name | ToPascalCase }}Builder) With{{ $column.ColumnName | ToPascalCase }}Value({{ $column.ColumnName | ToSafeParamName }} {{ $column | GetGoType | TrimPointer }}) *{{ $schema.Name | ToPascalCase }}Builder {
	b.instance.{{ $column.ColumnName | ToPascalCase }} = &{{ $column.ColumnName | ToSafeParamName }}
	return b
}
{{- end }}
{{- end }}
{{- end }}

{{ if ne $.CommentStyle "none" }}// Build 构建并返回 {{ $schema.Name | ToPascalCase }} 实例
{{- if eq $.CommentStyle "full" }}
// 返回:
//   - *{{ $schema.Name | ToPascalCase }}: 构建完成的实例
{{- end }}
{{- end }}
func (b *{{ $schema.Name | ToPascalCase }}Builder) Build() *{{ $schema.Name | ToPascalCase }} {
	return b.instance
}

{{- end }}

{{- end }}
