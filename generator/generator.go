package generator

import (
	"fmt"
	"model_infrax/model"
	"os"
	"path/filepath"
	"text/template"
)

// Generator 代码生成器
type Generator struct {
	templatePath string // 模板文件路径
	outputPath   string // 输出路径
}

// NewGenerator 创建新的生成器实例
func NewGenerator(templatePath, outputPath string) *Generator {
	return &Generator{
		templatePath: templatePath,
		outputPath:   outputPath,
	}
}

// GenerateOneByOne 根据模板生成代码
func (g *Generator) GenerateOneByOne(schemas []model.Schema) (err error) {
	for _, schema := range schemas {
		fileName := fmt.Sprintf("%s.go", schema.Name)
		err = g.Generate([]model.Schema{schema}, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// Generate 生成所有表到一个文件
func (g *Generator) Generate(schemas []model.Schema, outputFileName string) (err error) {
	// 读取模板文件
	tmplContent, err := os.ReadFile(g.templatePath)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %w", err)
	}

	// 创建模板并注册函数
	tmpl, err := template.New("model").Funcs(template.FuncMap{
		"ToPascalCase": ToPascalCase,
		"GetGoType":    GetGoType,
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	// 确保输出目录存在
	if err = os.MkdirAll(g.outputPath, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 生成文件路径
	filePath := filepath.Join(g.outputPath, outputFileName)

	// 创建输出文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer file.Close()

	// 执行模板
	err = tmpl.Execute(file, schemas)
	if err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	fmt.Printf("成功生成文件: %s\n", filePath)
	return nil
}
