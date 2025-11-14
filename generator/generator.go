package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"model_infrax/config"
	"model_infrax/model"
	"os"
	"path/filepath"
	"text/template"
)

// Generator 代码生成器
type Generator struct {
	modelTemplatePath string            // 模板文件路径
	configger         *config.Configger // 配置对象
}

// NewGenerator 创建新的生成器实例
// 参数:
//   - cfg: 配置对象，用于获取模板路径和输出路径等配置信息
//
// 返回:
//   - *Generator: 生成器实例
func NewGenerator(cfg *config.Configger) *Generator {
	return &Generator{
		modelTemplatePath: "./assert/template/model.template",
		configger:         cfg,
	}
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
	// 读取模板文件
	tmplContent, err := os.ReadFile(g.modelTemplatePath)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %w", err)
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
	outputPath := g.configger.GenerateOption.OutputPath

	// 确保输出目录存在
	if err = os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 生成文件路径
	filePath := filepath.Join(outputPath, outputFileName)

	// 先将模板执行结果写入缓冲区
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, schemas)
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