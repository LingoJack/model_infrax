package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/logger"
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

var (
	funcMap = template.FuncMap{
		"ToPascalCase":    ToPascalCase,
		"ToCamelCase":     ToCamelCase,
		"ToSafeParamName": ToSafeParamName,
		"TrimPointer":     TrimPointer,
		"GetGoType":       GetGoType,
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

func GenerateModels(schemas []model.Schema) (err error) {
	logger.Infof("[GenerateModelOneByOne]开始生成模型, 共 %d 个表", len(schemas))
	for i, schema := range schemas {
		fileName := fmt.Sprintf("%s.go", schema.Name)
		logger.Infof("[GenerateModelOneByOne]正在生成第 %d/%d 个模型: %s", i+1, len(schemas), schema.Name)
		err = GenerateModel([]model.Schema{schema}, outputPath, fileName)
		if err != nil {
			logger.Errorf("[GenerateModelOneByOne]生成模型失败, 表名: %s, 文件名: %s, 错误: %v", schema.Name, fileName, err)
			return err
		}
	}
	logger.Infof("[GenerateModelOneByOne]所有模型生成完成, 共 %d 个", len(schemas))
	return nil
}

func GenerateModel(schemas []model.Schema, outputDir, outputFileName string) (err error) {
	logger.Infof("[GenerateModel]开始生成模型, 文件名: %s, 表数量: %d", outputFileName, len(schemas))

	templatePath := PoTemplatePath()
	logger.Infof("[GenerateModel]读取模板文件: %s", templatePath)
	tmplContent, err := fs.ReadFile(templatePath)
	if err != nil {
		logger.Errorf("[GenerateModel]读取嵌入模板文件失败, 模板路径: %s, 错误: %v", templatePath, err)
		return fmt.Errorf("读取嵌入式模板文件失败: %w", err)
	}
	logger.Infof("[GenerateModel]模板文件读取成功, 大小: %d 字节", len(tmplContent))

	// 创建模板并注册函数
	logger.Infof("[GenerateModel]开始解析模板")
	tmpl, err := template.New("model").Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		logger.Errorf("[GenerateModel]解析模板失败, 错误: %v", err)
		return fmt.Errorf("解析模板失败: %w", err)
	}
	logger.Infof("[GenerateModel]模板解析成功")

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(daoPackage),
		PoPackageName:  getPackageName(poPackage),
		DtoPackageName: getPackageName(dtoPackage),
		Schemas:        schemas,
	}
	logger.Infof("[GenerateModel]模板数据准备完成, PoPackage: %s, DaoPackage: %s, DtoPackage: %s",
		templateData.PoPackageName, templateData.DaoPackageName, templateData.DtoPackageName)

	// 先将模板执行结果写入缓冲区
	logger.Infof("[GenerateModel]开始执行模板")
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		logger.Errorf("[GenerateModel]执行模板失败, 错误: %v", err)
		return fmt.Errorf("执行模板失败: %w", err)
	}
	logger.Infof("[GenerateModel]模板执行成功, 生成代码大小: %d 字节", buf.Len())

	filePath := filepath.Join(outputDir, poPackage, outputFileName)
	logger.Infof("[GenerateModel]开始写入文件: %s", filePath)
	err = tool.WriteFileWithDir(filePath, []byte(tool.FormatGoCode(buf.String())), 0644)
	if err != nil {
		logger.Errorf("[GenerateModel]写入输出文件失败, 文件路径: %s, 错误: %v", filePath, err)
		return fmt.Errorf("写入输出文件失败: %w", err)
	}

	logger.Infof("[GenerateModel]成功生成模型文件: %s", filePath)
	return nil
}

func GenerateDtos(schemas []model.Schema) (err error) {
	logger.Infof("[GenerateDtoOneByOne]开始生成 DTO, 共 %d 个表", len(schemas))
	for i, schema := range schemas {
		fileName := fmt.Sprintf("%s_dto.go", schema.Name)
		logger.Infof("[GenerateDtoOneByOne]正在生成第 %d/%d 个 DTO: %s", i+1, len(schemas), schema.Name)
		err = GenerateDto([]model.Schema{schema}, outputPath, fileName)
		if err != nil {
			logger.Errorf("[GenerateDtoOneByOne]生成 DTO 失败, 表名: %s, 文件名: %s, 错误: %v", schema.Name, fileName, err)
			return err
		}
	}
	logger.Infof("[GenerateDtoOneByOne]所有 DTO 生成完成, 共 %d 个", len(schemas))
	return nil
}

func GenerateDto(schemas []model.Schema, outputDir, outputFileName string) (err error) {
	logger.Infof("[GenerateDTO]开始生成 DTO, 文件名: %s, 表数量: %d", outputFileName, len(schemas))

	templatePath := DtoTemplatePath()
	logger.Infof("[GenerateDTO]读取模板文件: %s", templatePath)
	tmplContent, err := fs.ReadFile(templatePath)
	if err != nil {
		logger.Errorf("[GenerateDTO]读取嵌入式 DTO 模板文件失败, 模板路径: %s, 错误: %v", templatePath, err)
		return fmt.Errorf("读取嵌入式 DTO 模板文件失败: %w", err)
	}
	logger.Infof("[GenerateDTO]模板文件读取成功, 大小: %d 字节", len(tmplContent))

	// 创建模板并注册函数
	logger.Infof("[GenerateDTO]开始解析模板")
	tmpl, err := template.New("dto").Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		logger.Errorf("[GenerateDTO]解析 DTO 模板失败, 错误: %v", err)
		return fmt.Errorf("解析 DTO 模板失败: %w", err)
	}
	logger.Infof("[GenerateDTO]模板解析成功")

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(daoPackage),
		PoPackageName:  getPackageName(poPackage),
		DtoPackageName: getPackageName(dtoPackage),
		Schemas:        schemas,
	}
	logger.Infof("[GenerateDTO]模板数据准备完成, DtoPackage: %s, PoPackage: %s, DaoPackage: %s",
		templateData.DtoPackageName, templateData.PoPackageName, templateData.DaoPackageName)

	// 先将模板执行结果写入缓冲区
	logger.Infof("[GenerateDTO]开始执行模板")
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		logger.Errorf("[GenerateDTO]执行 DTO 模板失败, 错误: %v", err)
		return fmt.Errorf("执行 DTO 模板失败: %w", err)
	}
	logger.Infof("[GenerateDTO]模板执行成功, 生成代码大小: %d 字节", buf.Len())

	filePath := filepath.Join(outputDir, dtoPackage, outputFileName)
	logger.Infof("[GenerateDTO]开始写入文件: %s", filePath)
	err = tool.WriteFileWithDir(filePath, []byte(tool.FormatGoCode(buf.String())), 0644)
	if err != nil {
		logger.Errorf("[GenerateDTO]写入 DTO 输出文件失败, 文件路径: %s, 错误: %v", filePath, err)
		return fmt.Errorf("写入 DTO 输出文件失败: %w", err)
	}

	logger.Infof("[GenerateDTO]成功生成 DTO 文件: %s", filePath)
	return nil
}

func GenerateDaos(schemas []model.Schema) (err error) {
	logger.Infof("[GenerateDaoOneByOne]开始生成 DAO, 共 %d 个表", len(schemas))
	for i, schema := range schemas {
		fileName := fmt.Sprintf("%s_dao.go", schema.Name)
		logger.Infof("[GenerateDaoOneByOne]正在生成第 %d/%d 个 DAO: %s", i+1, len(schemas), schema.Name)
		err = GenerateDao([]model.Schema{schema}, outputPath, fileName)
		if err != nil {
			logger.Errorf("[GenerateDaoOneByOne]生成 DAO 失败, 表名: %s, 文件名: %s, 错误: %v", schema.Name, fileName, err)
			return err
		}
	}
	logger.Infof("[GenerateDaoOneByOne]所有 DAO 生成完成, 共 %d 个", len(schemas))
	return nil
}

func GenerateDao(schemas []model.Schema, outputDir, outputFileName string) (err error) {
	logger.Infof("[GenerateDao]开始生成 DAO, 文件名: %s, 表数量: %d", outputFileName, len(schemas))

	// 从嵌入的文件系统中读取 DAO 模板文件
	templatePath := DaoTemplatePath()
	logger.Infof("[GenerateDao]读取模板文件: %s", templatePath)
	tmplContent, err := fs.ReadFile(templatePath)
	if err != nil {
		logger.Errorf("[GenerateDao]读取嵌入式 DAO 模板文件失败, 模板路径: %s, 错误: %v", templatePath, err)
		return fmt.Errorf("读取嵌入式 DAO 模板文件失败: %w", err)
	}
	logger.Infof("[GenerateDao]模板文件读取成功, 大小: %d 字节", len(tmplContent))

	// 创建模板并注册函数
	logger.Infof("[GenerateDao]开始解析模板")
	tmpl, err := template.New("dao").Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		logger.Errorf("[GenerateDao]解析 DAO 模板失败, 错误: %v", err)
		return fmt.Errorf("解析 DAO 模板失败: %w", err)
	}
	logger.Infof("[GenerateDao]模板解析成功")

	// 从配置中获取输出路径（已在配置解析时展开 ~ 符号）
	o := filepath.Join(outputDir, daoPackage)
	logger.Infof("[GenerateDao]输出目录: %s", o)

	// 确保输出目录存在
	if err = os.MkdirAll(o, 0755); err != nil {
		logger.Errorf("[GenerateDao]创建 DAO 输出目录失败, 目录: %s, 错误: %v", o, err)
		return fmt.Errorf("创建 DAO 输出目录失败: %w", err)
	}
	logger.Infof("[GenerateDao]输出目录创建成功")

	// 生成文件路径
	filePath := filepath.Join(o, outputFileName)

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(daoPackage),
		PoPackageName:  getPackageName(poPackage),
		DtoPackageName: getPackageName(dtoPackage),
		Schemas:        schemas,
	}
	logger.Infof("[GenerateDao]模板数据准备完成, DaoPackage: %s, PoPackage: %s, DtoPackage: %s",
		templateData.DaoPackageName, templateData.PoPackageName, templateData.DtoPackageName)

	// 先将模板执行结果写入缓冲区
	logger.Infof("[GenerateDao]开始执行模板")
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		logger.Errorf("[GenerateDao]执行 DAO 模板失败, 错误: %v", err)
		return fmt.Errorf("执行 DAO 模板失败: %w", err)
	}
	logger.Infof("[GenerateDao]模板执行成功, 生成代码大小: %d 字节", buf.Len())

	// 创建输出文件并写入格式化后的代码
	logger.Infof("[GenerateDao]开始写入文件: %s", filePath)
	err = os.WriteFile(filePath, []byte(tool.FormatGoCode(buf.String())), 0644)
	if err != nil {
		logger.Errorf("[GenerateDao]写入 DAO 输出文件失败, 文件路径: %s, 错误: %v", filePath, err)
		return fmt.Errorf("写入 DAO 输出文件失败: %w", err)
	}

	logger.Infof("[GenerateDao]成功生成 DAO 文件: %s", filePath)
	return nil
}

func GenerateVos(schemas []model.Schema) (err error) {
	logger.Infof("[GenerateVoOneByOne]开始生成 VO, 共 %d 个表", len(schemas))
	for i, schema := range schemas {
		fileName := fmt.Sprintf("%s_vo.go", schema.Name)
		logger.Infof("[GenerateVoOneByOne]正在生成第 %d/%d 个 VO: %s", i+1, len(schemas), schema.Name)
		err = GenerateVO([]model.Schema{schema}, outputPath, fileName)
		if err != nil {
			logger.Errorf("[GenerateVoOneByOne]生成 VO 失败, 表名: %s, 文件名: %s, 错误: %v", schema.Name, fileName, err)
			return err
		}
	}
	logger.Infof("[GenerateVoOneByOne]所有 VO 生成完成, 共 %d 个", len(schemas))
	return nil
}

func GenerateVO(schemas []model.Schema, outputDir, outputFileName string) (err error) {
	logger.Infof("[GenerateVO]开始生成 VO, 文件名: %s, 表数量: %d", outputFileName, len(schemas))

	// 从嵌入的文件系统中读取 VO 模板文件
	tmplContent, err := fs.ReadFile(VoTemplatePath())
	if err != nil {
		logger.Errorf("[GenerateVO]读取嵌入式 VO 模板文件失败, 模板路径: %s, 错误: %v", VoTemplatePath(), err)
		return fmt.Errorf("读取嵌入式 VO 模板文件失败: %w", err)
	}
	logger.Infof("[GenerateVO]模板文件读取成功, 大小: %d 字节", len(tmplContent))

	// 创建模板并注册函数
	logger.Infof("[GenerateVO]开始解析模板")
	tmpl, err := template.New("vo").Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		logger.Errorf("[GenerateVO]解析 VO 模板失败, 错误: %v", err)
		return fmt.Errorf("解析 VO 模板失败: %w", err)
	}
	logger.Infof("[GenerateVO]模板解析成功")

	// 准备模板数据，包含包名和表结构
	templateData := TemplateData{
		DaoPackageName: getPackageName(daoPackage),
		PoPackageName:  getPackageName(poPackage),
		DtoPackageName: getPackageName(dtoPackage),
		VoPackageName:  getPackageName(voPackage),
		Schemas:        schemas,
	}
	logger.Infof("[GenerateVO]模板数据准备完成, VoPackage: %s, PoPackage: %s, DtoPackage: %s, DaoPackage: %s",
		templateData.VoPackageName, templateData.PoPackageName, templateData.DtoPackageName, templateData.DaoPackageName)

	// 先将模板执行结果写入缓冲区
	logger.Infof("[GenerateVO]开始执行模板")
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		logger.Errorf("[GenerateVO]执行模板失败, 错误: %v", err)
		return fmt.Errorf("执行模板失败: %w", err)
	}
	logger.Infof("[GenerateVO]模板执行成功, 生成代码大小: %d 字节", buf.Len())

	filePath := filepath.Join(outputDir, voPackage, outputFileName)
	logger.Infof("[GenerateVO]开始写入文件: %s", filePath)
	err = tool.WriteFileWithDir(filePath, []byte(tool.FormatGoCode(buf.String())), 0644)
	if err != nil {
		logger.Errorf("[GenerateVO]写入输出文件失败, 文件路径: %s, 错误: %v", filePath, err)
		return fmt.Errorf("写入输出文件失败: %w", err)
	}

	logger.Infof("[GenerateVO]成功生成文件: %s", filePath)
	return nil
}

func GenerateTools() (err error) {
	logger.Infof("[GenerateAllTools]开始生成所有工具文件")

	// 从嵌入的文件系统中读取工具模板目录
	toolDir := ToolTemplateDir()
	logger.Infof("[GenerateAllTools]读取工具模板目录: %s", toolDir)
	entries, err := fs.ReadDir(toolDir)
	if err != nil {
		logger.Errorf("[GenerateAllTools]读取嵌入式工具模板目录失败, 目录: %s, 错误: %v", toolDir, err)
		return fmt.Errorf("读取嵌入式工具模板目录失败: %w", err)
	}
	logger.Infof("[GenerateAllTools]找到 %d 个文件/目录", len(entries))

	// 遍历所有模板文件
	generatedCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			logger.Infof("[GenerateAllTools]跳过目录: %s", entry.Name())
			continue
		}

		// 只处理 templatePathSuffix 文件
		tplFileName := entry.Name()
		if !strings.HasSuffix(tplFileName, templatePathSuffix) {
			logger.Infof("[GenerateAllTools]跳过非模板文件: %s", tplFileName)
			continue
		}

		goFileName := strings.TrimSuffix(tplFileName, templatePathSuffix) + ".go"
		logger.Infof("[GenerateAllTools]准备生成工具文件 [%d]: %s -> %s", generatedCount+1, tplFileName, goFileName)

		err = GenerateTool(ToolTemplatePath(tplFileName), outputPath, goFileName)
		if err != nil {
			logger.Errorf("[GenerateAllTools]生成工具文件失败, 模板: %s, 输出: %s, 错误: %v", tplFileName, goFileName, err)
			return fmt.Errorf("生成工具文件 %s 失败: %w", goFileName, err)
		}
		generatedCount++
	}

	logger.Infof("[GenerateAllTools]所有工具文件生成完成, 共生成 %d 个文件", generatedCount)
	return nil
}

func GenerateTool(templatePath, outputDir, outputFileName string) (err error) {
	logger.Infof("[GenerateTool]开始生成工具文件, 模板: %s, 输出: %s", templatePath, outputFileName)

	// 从嵌入的文件系统中读取模板文件
	tmplContent, err := fs.ReadFile(templatePath)
	if err != nil {
		logger.Errorf("[GenerateTool]读取嵌入式工具模板文件失败, 模板路径: %s, 错误: %v", templatePath, err)
		return fmt.Errorf("读取嵌入式工具模板文件失败: %w", err)
	}
	logger.Infof("[GenerateTool]模板文件读取成功, 大小: %d 字节", len(tmplContent))

	// 创建模板并注册函数
	logger.Infof("[GenerateTool]开始解析模板")
	tmpl, err := template.New("tool").Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		logger.Errorf("[GenerateTool]解析工具模板失败, 错误: %v", err)
		return fmt.Errorf("解析工具模板失败: %w", err)
	}
	logger.Infof("[GenerateTool]模板解析成功")

	// 先将模板执行结果写入缓冲区
	logger.Infof("[GenerateTool]开始执行模板")
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		logger.Errorf("[GenerateTool]执行工具模板失败, 错误: %v", err)
		return fmt.Errorf("执行工具模板失败: %w", err)
	}
	logger.Infof("[GenerateTool]模板执行成功, 生成代码大小: %d 字节", buf.Len())

	// 创建输出文件并写入格式化后的代码
	filePath := filepath.Join(outputDir, toolPackage, outputFileName)
	logger.Infof("[GenerateTool]开始写入文件: %s", filePath)
	err = tool.WriteFileWithDir(filePath, []byte(tool.FormatGoCode(buf.String())), 0644)
	if err != nil {
		logger.Errorf("[GenerateTool]写入工具输出文件失败, 文件路径: %s, 错误: %v", filePath, err)
		return fmt.Errorf("写入工具输出文件失败: %w", err)
	}

	logger.Infof("[GenerateTool]成功生成工具文件: %s", filePath)
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
