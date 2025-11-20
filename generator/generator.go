package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"log"

	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/LingoJack/model_infrax/config"
	"github.com/LingoJack/model_infrax/model"
)

// Generator 代码生成器
type Generator struct {
	modelTemplatePath string            // 模板文件路径（嵌入式路径）
	daoTemplatePath   string            // dao文件路径（嵌入式路径）
	dtoTemplatePath   string            // dto文件路径（嵌入式路径）
	toolTemplateDir   string            // tool文件路径（嵌入式路径）
	configger         *config.Configger // 配置对象
}

// TemplateData 传递给模板的数据结构
type TemplateData struct {
	PoPackageName  string         // po 包名（从路径最后一段提取）
	DtoPackageName string         // dto 包名（从路径最后一段提取）
	DaoPackageName string         // dao 包名（从路径最后一段提取）
	Schemas        []model.Schema // 表结构列表
}

// NewGenerator 创建新的生成器实例
// 参数:
//   - cfg: 配置对象，用于获取模板路径和输出路径等配置信息
//
// 返回:
//   - *Generator: 生成器实例
func NewGenerator(cfg *config.Configger) *Generator {
	// 使用嵌入式模板路径（相对于 embed.FS 的根路径）
	// 由于 embed.go 在 generator 目录下，需要使用 ../assert/template/ 前缀
	generator := Generator{
		modelTemplatePath: templatePathPrefix + "po.template",
		daoTemplatePath:   templatePathPrefix + "dao.template",
		dtoTemplatePath:   templatePathPrefix + "dto.template",
		toolTemplateDir:   templatePathPrefix + "tools",
		configger:         cfg,
	}

	// 如果使用 itea-go 框架，切换到对应的模板路径
	if cfg.GenerateOption.UseFramework == "itea-go" {
		generator = Generator{
			modelTemplatePath: templatePathPrefix + "itea-go/po.template",
			daoTemplatePath:   templatePathPrefix + "itea-go/dao.template",
			dtoTemplatePath:   templatePathPrefix + "itea-go/dto.template",
			toolTemplateDir:   templatePathPrefix + "tools",
			configger:         cfg,
		}
	}

	return &generator
}

// GenerateModelOneByOne 根据模板生成代码，每个表生成一个文件
// 参数:
//   - schemas: 表结构列表
//
// 返回:
//   - error: 生成过程中的错误
func (g *Generator) GenerateModelOneByOne(schemas []model.Schema) (err error) {
	for _, schema := range schemas {
		fileName := fmt.Sprintf("%s.go", schema.Name)
		err = g.GenerateModel([]model.Schema{schema}, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateModel 生成所有表到一个文件
// 参数:
//   - schemas: 表结构列表
//   - outputFileName: 输出文件名
//
// 返回:
//   - error: 生成过程中的错误
func (g *Generator) GenerateModel(schemas []model.Schema, outputFileName string) (err error) {

	// 从嵌入的文件系统中读取模板文件
	tmplContent, err := templateFS.ReadFile(g.modelTemplatePath)
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

	// 从配置中获取输出路径（已在配置解析时展开 ~ 符号）
	outputPath := filepath.Join(g.configger.GenerateOption.OutputPath, g.configger.GenerateOption.Package.PoPackage)

	// 确保输出目录存在
	if err = os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 生成文件路径
	filePath := filepath.Join(outputPath, outputFileName)

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(g.configger.GenerateOption.Package.DaoPackage),
		PoPackageName:  getPackageName(g.configger.GenerateOption.Package.PoPackage),
		DtoPackageName: getPackageName(g.configger.GenerateOption.Package.DtoPackage),
		Schemas:        schemas,
	}

	// 先将模板执行结果写入缓冲区
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	// 使用 go/format 格式化代码
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		// 如果格式化失败，记录警告但仍然写入未格式化的代码
		log.Printf("警告: 格式化代码失败: %v，将写入未格式化的代码\n", err)
		formattedCode = buf.Bytes()
	}

	// 创建输出文件并写入格式化后的代码
	err = os.WriteFile(filePath, formattedCode, 0644)
	if err != nil {
		return fmt.Errorf("写入输出文件失败: %w", err)
	}

	log.Printf("成功生成文件: %s\n", filePath)
	return nil
}

// GenerateDTOOneByOne 根据模板生成 DTO 代码，每个表生成一个文件
// 参数:
//   - schemas: 表结构列表
//
// 返回:
//   - error: 生成过程中的错误
func (g *Generator) GenerateDTOOneByOne(schemas []model.Schema) (err error) {
	for _, schema := range schemas {
		fileName := fmt.Sprintf("%s_dto.go", schema.Name)
		err = g.GenerateDTO([]model.Schema{schema}, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateDTO 生成 DTO 文件
// 参数:
//   - schemas: 表结构列表
//   - outputFileName: 输出文件名
//
// 返回:
//   - error: 生成过程中的错误
func (g *Generator) GenerateDTO(schemas []model.Schema, outputFileName string) (err error) {
	// 从嵌入的文件系统中读取 DTO 模板文件
	tmplContent, err := templateFS.ReadFile(g.dtoTemplatePath)
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

	// 从配置中获取输出路径（已在配置解析时展开 ~ 符号）
	outputPath := filepath.Join(g.configger.GenerateOption.OutputPath, g.configger.GenerateOption.Package.DtoPackage)

	// 确保输出目录存在
	if err = os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("创建 DTO 输出目录失败: %w", err)
	}

	// 生成文件路径
	filePath := filepath.Join(outputPath, outputFileName)

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(g.configger.GenerateOption.Package.DaoPackage),
		PoPackageName:  getPackageName(g.configger.GenerateOption.Package.PoPackage),
		DtoPackageName: getPackageName(g.configger.GenerateOption.Package.DtoPackage),
		Schemas:        schemas,
	}

	// 先将模板执行结果写入缓冲区
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return fmt.Errorf("执行 DTO 模板失败: %w", err)
	}

	// 使用 go/format 格式化代码
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		// 如果格式化失败，记录警告但仍然写入未格式化的代码
		log.Printf("警告: 格式化 DTO 代码失败: %v，将写入未格式化的代码\n", err)
		formattedCode = buf.Bytes()
	}

	// 创建输出文件并写入格式化后的代码
	err = os.WriteFile(filePath, formattedCode, 0644)
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
func (g *Generator) GenerateTool(templateFileName, outputFileName string) (err error) {
	// 构建嵌入式模板文件路径
	templatePath := filepath.Join(g.toolTemplateDir, templateFileName)

	// 从嵌入的文件系统中读取模板文件
	tmplContent, err := templateFS.ReadFile(templatePath)
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

	// 从配置中获取输出路径（已在配置解析时展开 ~ 符号）
	outputPath := filepath.Join(g.configger.GenerateOption.OutputPath, g.configger.GenerateOption.Package.ToolPackage)

	// 确保输出目录存在
	if err = os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("创建工具输出目录失败: %w", err)
	}

	// 生成文件路径
	filePath := filepath.Join(outputPath, outputFileName)

	// 先将模板执行结果写入缓冲区
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		return fmt.Errorf("执行工具模板失败: %w", err)
	}

	// 使用 go/format 格式化代码
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		// 如果格式化失败，记录警告但仍然写入未格式化的代码
		log.Printf("警告: 格式化工具代码失败: %v，将写入未格式化的代码\n", err)
		formattedCode = buf.Bytes()
	}

	// 创建输出文件并写入格式化后的代码
	err = os.WriteFile(filePath, formattedCode, 0644)
	if err != nil {
		return fmt.Errorf("写入工具输出文件失败: %w", err)
	}

	log.Printf("成功生成工具文件: %s\n", filePath)
	return nil
}

// GenerateDAOOneByOne 根据模板生成 DAO 代码，每个表生成一个文件
// 参数:
//   - schemas: 表结构列表
//
// 返回:
//   - error: 生成过程中的错误
func (g *Generator) GenerateDAOOneByOne(schemas []model.Schema) (err error) {
	for _, schema := range schemas {
		fileName := fmt.Sprintf("%s_dao.go", schema.Name)
		err = g.GenerateDAO([]model.Schema{schema}, fileName)
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
func (g *Generator) GenerateDAO(schemas []model.Schema, outputFileName string) (err error) {
	// 从嵌入的文件系统中读取 DAO 模板文件
	tmplContent, err := templateFS.ReadFile(g.daoTemplatePath)
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
	outputPath := filepath.Join(g.configger.GenerateOption.OutputPath, g.configger.GenerateOption.Package.DaoPackage)

	// 确保输出目录存在
	if err = os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("创建 DAO 输出目录失败: %w", err)
	}

	// 生成文件路径
	filePath := filepath.Join(outputPath, outputFileName)

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(g.configger.GenerateOption.Package.DaoPackage),
		PoPackageName:  getPackageName(g.configger.GenerateOption.Package.PoPackage),
		DtoPackageName: getPackageName(g.configger.GenerateOption.Package.DtoPackage),
		Schemas:        schemas,
	}

	// 先将模板执行结果写入缓冲区
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return fmt.Errorf("执行 DAO 模板失败: %w", err)
	}

	// 使用 go/format 格式化代码
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		// 如果格式化失败，记录警告但仍然写入未格式化的代码
		log.Printf("警告: 格式化 DAO 代码失败: %v，将写入未格式化的代码\n", err)
		formattedCode = buf.Bytes()
	}

	// 创建输出文件并写入格式化后的代码
	err = os.WriteFile(filePath, formattedCode, 0644)
	if err != nil {
		return fmt.Errorf("写入 DAO 输出文件失败: %w", err)
	}

	log.Printf("成功生成 DAO 文件: %s\n", filePath)
	return nil
}

// GenerateAllTools 生成所有工具文件
// 返回:
//   - error: 生成过程中的错误
func (g *Generator) GenerateAllTools() (err error) {
	// 从嵌入的文件系统中读取工具模板目录
	entries, err := templateFS.ReadDir(g.toolTemplateDir)
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
		err = g.GenerateTool(templateFileName, outputFileName)
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
