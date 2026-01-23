package generator

import (
	"path/filepath"
	"strings"
	"text/template"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/model"
)

// getPoPackage 获取 PO 包路径配置
// 必须在 conf.InitWithPath 之后调用
func getPoPackage() string {
	return conf.ValueStr("generate_option.package.po")
}

// getDaoPackage 获取 DAO 包路径配置
// 必须在 conf.InitWithPath 之后调用
func getDaoPackage() string {
	return conf.ValueStr("generate_option.package.dao")
}

// getDtoPackage 获取 DTO 包路径配置
// 必须在 conf.InitWithPath 之后调用
func getDtoPackage() string {
	return conf.ValueStr("generate_option.package.dto")
}

// getVoPackage 获取 VO 包路径配置
// 必须在 conf.InitWithPath 之后调用
func getVoPackage() string {
	return conf.ValueStr("generate_option.package.vo")
}

// getToolPackage 获取 Tool 包路径配置
// 必须在 conf.InitWithPath 之后调用
func getToolPackage() string {
	return conf.ValueStr("generate_option.package.tool")
}

// getOutputPath 获取输出路径配置
// 必须在 conf.InitWithPath 之后调用
func getOutputPath() string {
	return conf.ValueStr("generate_option.output_path")
}

var (
	funcMap = template.FuncMap{
		"ToPascalCase":     ToPascalCase,
		"ToCamelCase":      ToCamelCase,
		"ToSafeParamName":  ToSafeParamName,
		"TrimPointer":      TrimPointer,
		"GetGoType":        GetGoType,
		"ToJsonTag":        ToJsonTag,
		"contains":         strings.Contains,
		"IsSnakeCaseStyle": IsSnakeCaseStyle,
	}
)

// TemplateData 传递给模板的数据结构
type TemplateData struct {
	PoPackageName  string         // po 包名（从路径最后一段提取）
	DtoPackageName string         // dto 包名（从路径最后一段提取）
	VoPackageName  string         // vo 包名（从路径最后一段提取）
	DaoPackageName string         // dao 包名（从路径最后一段提取）
	Schemas        []model.Schema // 表结构列表
}

// getPackageName 从路径中提取包名（取路径的最后一段）
// 参数:
//   - path: 包路径，例如 "model/entity" 或 "dao"
//
// 返回:
//   - string: 包名，例如 "entity" 或 "dao"
//
// 示例:
//   - "model/entity" -> "entity"
//   - "dao" -> "dao"
//   - "model/query" -> "query"
func getPackageName(path string) string {
	// 使用 filepath.Base 获取路径的最后一段
	// 同时处理 Windows 和 Unix 风格的路径分隔符
	path = filepath.ToSlash(path) // 统一转换为 / 分隔符
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		panic("invalid path")
	}
	return parts[len(parts)-1]
}
