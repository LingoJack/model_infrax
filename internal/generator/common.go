package generator

import (
	"path/filepath"
	"strings"
	"text/template"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/internal/tool"
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
		"ToPascalCase":           ToPascalCase,
		"ToCamelCase":            ToCamelCase,
		"ToSafeParamName":        ToSafeParamName,
		"TrimPointer":            TrimPointer,
		"GetGoType":              GetGoType,
		"ToJsonTag":              ToJsonTag,
		"contains":               strings.Contains,
		"IsSnakeCaseStyle":       IsSnakeCaseStyle,
		"HasTimeColumnInSchemas": HasTimeColumnInSchemas,
	}
)

// TemplateData 传递给模板的数据结构
type TemplateData struct {
	PoPackageName  string         // po 包名（从路径最后一段提取）
	DtoPackageName string         // dto 包名（从路径最后一段提取）
	VoPackageName  string         // vo 包名（从路径最后一段提取）
	DaoPackageName string         // dao 包名（从路径最后一段提取）
	Schemas        []model.Schema // 表结构列表

	// 完整 import 路径，用于跨包引用（如 dao 引用 po/dto）
	// 格式: {module}/{output_path}/{package}，例如 github.com/foo/app/target/jen/model/entity
	// 若 go.mod 读取失败则为空字符串，模板应对其做 nil 判断
	PoImportPath  string
	DtoImportPath string
	VoImportPath  string
	DaoImportPath string
}

// newTemplateData 构建 TemplateData，自动从 go.mod 推导完整 import 路径
func newTemplateData(schemas []model.Schema) TemplateData {
	modulePath, err := tool.ReadModulePath()
	if err != nil {
		logger.Infof("[newTemplateData] 读取 go.mod 失败，跨包 import path 将为空: %v", err)
	}

	buildImportPath := func(pkg string) string {
		if modulePath == "" || pkg == "" {
			return ""
		}
		outputPath := filepath.ToSlash(getOutputPath())
		pkg = filepath.ToSlash(pkg)
		return modulePath + "/" + outputPath + "/" + pkg
	}

	return TemplateData{
		DaoPackageName: getPackageName(getDaoPackage()),
		PoPackageName:  getPackageName(getPoPackage()),
		DtoPackageName: getPackageName(getDtoPackage()),
		VoPackageName:  getPackageName(getVoPackage()),
		Schemas:        schemas,
		PoImportPath:   buildImportPath(getPoPackage()),
		DtoImportPath:  buildImportPath(getDtoPackage()),
		VoImportPath:   buildImportPath(getVoPackage()),
		DaoImportPath:  buildImportPath(getDaoPackage()),
	}
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
