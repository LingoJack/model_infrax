package generator

import (
	"bytes"
	"fmt"
	"log"

	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/internal/tool"
)

var (
	poPackage   = conf.ValueStr("generate_option.package.po")
	daoPackage  = conf.ValueStr("generate_option.package.dao")
	dtoPackage  = conf.ValueStr("generate_option.package.dto")
	voPackage   = conf.ValueStr("generate_option.package.vo")
	toolPackage = conf.ValueStr("generate_option.package.tool")
	outputPath  = conf.ValueStr("generate_option.output_path")
)

// TemplateData 传递给模板的数据结构
type TemplateData struct {
	PoPackageName  string         // po 包名（从路径最后一段提取）
	DtoPackageName string         // dto 包名（从路径最后一段提取）
	VoPackageName  string         // vo 包名（从路径最后一段提取）
	DaoPackageName string         // dao 包名（从路径最后一段提取）
	Schemas        []model.Schema // 表结构列表
}

func GenerateModelOneByOne(schemas []model.Schema) (err error) {
	for _, schema := range schemas {
		fileName := fmt.Sprintf("%s.go", schema.Name)
		err = GenerateModelAllInOne([]model.Schema{schema}, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateModelAllInOne(schemas []model.Schema, outputFileName string) (err error) {
	tmplContent, err := fs.ReadFile(FrameworkPrefix() + "po.template")
	if err != nil {
		return fmt.Errorf("读取嵌入式模板文件失败: %w", err)
	}

	// 创建模板并注册函数
	tmpl, err := template.New("model").Funcs(template.FuncMap{
		"ToPascalCase":    ToPascalCase,
		"ToCamelCase":     ToCamelCase,
		"ToSafeParamName": ToSafeParamName,
		"TrimPointer":     TrimPointer,
		"GetGoType":       GetGoType,
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(daoPackage),
		PoPackageName:  getPackageName(poPackage),
		DtoPackageName: getPackageName(dtoPackage),
		Schemas:        schemas,
	}

	// 先将模板执行结果写入缓冲区
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	filePath := filepath.Join(outputPath, poPackage, outputFileName)
	err = tool.WriteFileWithDir(filePath, []byte(tool.FormatGoCode(buf.String())), 0644)
	if err != nil {
		return fmt.Errorf("写入输出文件失败: %w", err)
	}

	log.Printf("成功生成文件: %s\n", filePath)
	return nil
}

func GenerateDtoOneByOne(schemas []model.Schema) (err error) {
	for _, schema := range schemas {
		fileName := fmt.Sprintf("%s_dto.go", schema.Name)
		err = GenerateDTO([]model.Schema{schema}, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateDTO(schemas []model.Schema, outputFileName string) (err error) {
	tmplContent, err := fs.ReadFile(FrameworkPrefix() + "dto.template")
	if err != nil {
		return fmt.Errorf("读取嵌入式 DTO 模板文件失败: %w", err)
	}

	// 创建模板并注册函数
	tmpl, err := template.New("dto").Funcs(template.FuncMap{
		"ToPascalCase":    ToPascalCase,
		"ToCamelCase":     ToCamelCase,
		"ToSafeParamName": ToSafeParamName,
		"TrimPointer":     TrimPointer,
		"GetGoType":       GetGoType,
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("解析 DTO 模板失败: %w", err)
	}

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(daoPackage),
		PoPackageName:  getPackageName(poPackage),
		DtoPackageName: getPackageName(dtoPackage),
		Schemas:        schemas,
	}

	// 先将模板执行结果写入缓冲区
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return fmt.Errorf("执行 DTO 模板失败: %w", err)
	}

	filePath := filepath.Join(outputPath, dtoPackage, outputFileName)
	err = tool.WriteFileWithDir(filePath, []byte(tool.FormatGoCode(buf.String())), 0644)
	if err != nil {
		return fmt.Errorf("写入 DTO 输出文件失败: %w", err)
	}

	log.Printf("成功生成 DTO 文件: %s\n", filePath)
	return nil
}

// GenerateTool 生成工具文件
// 参数:
//   - templateFileName: 模板文件名（如 "ptr.template"）
//   - outputFileName: 输出文件名（如 "ptr.go"）
//
// 返回:
//   - error: 生成过程中的错误
func GenerateTool(templateFileName, outputFileName string) (err error) {
	// 构建嵌入式模板文件路径
	templatePath := filepath.Join(templatePathPrefix+"tool", templateFileName)

	// 从嵌入的文件系统中读取模板文件
	tmplContent, err := fs.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("读取嵌入式工具模板文件失败: %w", err)
	}

	// 创建模板并注册函数
	tmpl, err := template.New("tool").Funcs(template.FuncMap{
		"ToPascalCase":    ToPascalCase,
		"ToCamelCase":     ToCamelCase,
		"ToSafeParamName": ToSafeParamName,
		"TrimPointer":     TrimPointer,
		"GetGoType":       GetGoType,
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("解析工具模板失败: %w", err)
	}

	// 先将模板执行结果写入缓冲区
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		return fmt.Errorf("执行工具模板失败: %w", err)
	}

	// 创建输出文件并写入格式化后的代码
	filePath := filepath.Join(outputPath, toolPackage, outputFileName)
	err = tool.WriteFileWithDir(filePath, []byte(tool.FormatGoCode(buf.String())), 0644)
	if err != nil {
		return fmt.Errorf("写入工具输出文件失败: %w", err)
	}

	log.Printf("成功生成工具文件: %s\n", filePath)
	return nil
}

// GenerateDaoOneByOne 根据模板生成 DAO 代码，每个表生成一个文件
// 参数:
//   - schemas: 表结构列表
//
// 返回:
//   - error: 生成过程中的错误
func GenerateDaoOneByOne(schemas []model.Schema) (err error) {
	for _, schema := range schemas {
		fileName := fmt.Sprintf("%s_dao.go", schema.Name)
		err = GenerateDAO([]model.Schema{schema}, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateDAO 生成 DAO 文件
// 参数:
//   - schemas: 表结构列表
//   - outputFileName: 输出文件名
//
// 返回:
//   - error: 生成过程中的错误
func GenerateDAO(schemas []model.Schema, outputFileName string) (err error) {
	// 从嵌入的文件系统中读取 DAO 模板文件
	tmplContent, err := fs.ReadFile(FrameworkPrefix() + "dao.template")
	if err != nil {
		return fmt.Errorf("读取嵌入式 DAO 模板文件失败: %w", err)
	}

	// 创建模板并注册函数
	tmpl, err := template.New("dao").Funcs(template.FuncMap{
		"ToPascalCase":    ToPascalCase,
		"ToCamelCase":     ToCamelCase,
		"ToSafeParamName": ToSafeParamName,
		"TrimPointer":     TrimPointer,
		"GetGoType":       GetGoType,
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("解析 DAO 模板失败: %w", err)
	}

	// 从配置中获取输出路径（已在配置解析时展开 ~ 符号）
	o := filepath.Join(outputPath, daoPackage)

	// 确保输出目录存在
	if err = os.MkdirAll(o, 0755); err != nil {
		return fmt.Errorf("创建 DAO 输出目录失败: %w", err)
	}

	// 生成文件路径
	filePath := filepath.Join(o, outputFileName)

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(daoPackage),
		PoPackageName:  getPackageName(poPackage),
		DtoPackageName: getPackageName(dtoPackage),
		Schemas:        schemas,
	}

	// 先将模板执行结果写入缓冲区
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return fmt.Errorf("执行 DAO 模板失败: %w", err)
	}

	// 创建输出文件并写入格式化后的代码
	err = os.WriteFile(filePath, []byte(tool.FormatGoCode(buf.String())), 0644)
	if err != nil {
		return fmt.Errorf("写入 DAO 输出文件失败: %w", err)
	}

	log.Printf("成功生成 DAO 文件: %s\n", filePath)
	return nil
}

// GenerateAllTools 生成所有工具文件
// 返回:
//   - error: 生成过程中的错误
func GenerateAllTools() (err error) {
	// 从嵌入的文件系统中读取工具模板目录
	entries, err := fs.ReadDir(templatePathPrefix + "tool")
	if err != nil {
		return fmt.Errorf("读取嵌入式工具模板目录失败: %w", err)
	}

	// 遍历所有模板文件
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 只处理 .template 文件
		templateFileName := entry.Name()
		if !strings.HasSuffix(templateFileName, ".template") {
			continue
		}

		// 生成输出文件名（将 .template 替换为 .go）
		outputFileName := strings.TrimSuffix(templateFileName, ".template") + ".go"

		// 生成工具文件
		err = GenerateTool(templateFileName, outputFileName)
		if err != nil {
			return fmt.Errorf("生成工具文件 %s 失败: %w", outputFileName, err)
		}
	}

	return nil
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
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "model" // 默认返回 model
}

// GenerateVO 生成 VO 文件
// 参数:
//   - schemas: 表结构列表
//   - outputFileName: 输出文件名
//
// 返回:
//   - error: 生成过程中的错误
func GenerateVO(schemas []model.Schema, outputFileName string) (err error) {
	// 从嵌入的文件系统中读取 VO 模板文件
	tmplContent, err := fs.ReadFile(FrameworkPrefix() + "vo.template")
	if err != nil {
		return fmt.Errorf("读取嵌入式 VO 模板文件失败: %w", err)
	}

	// 创建模板并注册函数
	tmpl, err := template.New("vo").Funcs(template.FuncMap{
		"ToPascalCase":    ToPascalCase,
		"ToCamelCase":     ToCamelCase,
		"ToSafeParamName": ToSafeParamName,
		"TrimPointer":     TrimPointer,
		"GetGoType":       GetGoType,
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("解析 VO 模板失败: %w", err)
	}

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(daoPackage),
		PoPackageName:  getPackageName(poPackage),
		DtoPackageName: getPackageName(dtoPackage),
		VoPackageName:  getPackageName(voPackage),
		Schemas:        schemas,
	}

	// 先将模板执行结果写入缓冲区
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	filePath := filepath.Join(outputPath, voPackage, outputFileName)
	err = tool.WriteFileWithDir(filePath, []byte(tool.FormatGoCode(buf.String())), 0644)
	if err != nil {
		return fmt.Errorf("写入输出文件失败: %w", err)
	}

	log.Printf("成功生成文件: %s\n", filePath)
	return nil
}

// GenerateVoOneByOne 根据模板生成 VO 代码，每个表生成一个文件
// 参数:
//   - schemas: 表结构列表
//
// 返回:
//   - error: 生成过程中的错误
func GenerateVoOneByOne(schemas []model.Schema) (err error) {
	for _, schema := range schemas {
		fileName := fmt.Sprintf("%s_vo.go", schema.Name)
		err = GenerateVO([]model.Schema{schema}, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}
