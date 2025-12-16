{{- /* 生成 Vo 结构体的模板 */ -}}
package {{ .VoPackageName }}

import (
	"encoding/json"
	"time"
)

{{- range $schema := .Schemas }}

// {{ $schema.Name | ToPascalCase }}Vo {{ $schema.Comment }} 视图对象
type {{ $schema.Name | ToPascalCase }}Vo struct {
{{- range $schema.Columns }}
	{{ .ColumnName | ToPascalCase }} {{ . | GetGoType }} `json:"{{ .ColumnName }},omitempty"` // {{ .Comment }}
{{- end }}
}

// Jsonify 将结构体序列化为 JSON 字符串（紧凑格式）
// 返回:
//   - string: JSON 字符串，如果序列化失败则返回错误信息的 JSON
func (t *{{ $schema.Name | ToPascalCase }}Vo) Jsonify() string {
	byts, err := json.Marshal(t)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(byts)
}

// JsonifyIndent 将结构体序列化为格式化的 JSON 字符串（带缩进）
// 返回:
//   - string: 格式化的 JSON 字符串，如果序列化失败则返回错误信息的 JSON
func (t *{{ $schema.Name | ToPascalCase }}Vo) JsonifyIndent() string {
	byts, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(byts)
}

{{- end }}