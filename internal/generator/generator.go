package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/internal/tool"
)

// artifactKind 描述一类生成物（PO/DTO/VO/DAO）的生成参数
type artifactKind struct {
	label        string                        // 日志标识，如 "PO"/"DTO"
	templatePath func() string                 // 嵌入模板路径
	packageDir   func() string                 // 输出子目录（配置项 getter）
	fileName     func(tableName string) string // 输出文件名
}

var (
	poKind  = artifactKind{"PO", poTemplatePath, getPoPackage, func(t string) string { return t + ".go" }}
	dtoKind = artifactKind{"DTO", dtoTemplatePath, getDtoPackage, func(t string) string { return t + "_dto.go" }}
	voKind  = artifactKind{"VO", voTemplatePath, getVoPackage, func(t string) string { return t + "_vo.go" }}
	daoKind = artifactKind{"DAO", daoTemplatePath, getDaoPackage, func(t string) string { return t + "_dao.go" }}
)

// GeneratePos 逐表生成 PO 实体代码
func GeneratePos(schemas []model.Schema) error {
	return generateArtifacts(schemas, poKind)
}

// GenerateDtos 逐表生成 DTO 代码
func GenerateDtos(schemas []model.Schema) error {
	return generateArtifacts(schemas, dtoKind)
}

// GenerateVos 逐表生成 VO 代码
func GenerateVos(schemas []model.Schema) error {
	return generateArtifacts(schemas, voKind)
}

// GenerateDaos 逐表生成 DAO 代码
func GenerateDaos(schemas []model.Schema) error {
	return generateArtifacts(schemas, daoKind)
}

// generateArtifacts 读取模板后逐表渲染、格式化并写出文件
// 任一表生成失败即中断返回
func generateArtifacts(schemas []model.Schema, kind artifactKind) error {
	tmplContent, err := fs.ReadFile(kind.templatePath())
	if err != nil {
		logger.Errorf("[generateArtifacts] 读取嵌入模板文件失败, 模板路径: %s, 错误: %v", kind.templatePath(), err)
		return fmt.Errorf("读取嵌入式模板文件失败: %w", err)
	}

	logger.Infof("[generateArtifacts] 开始生成 %s, 共 %d 个表", kind.label, len(schemas))
	for i, schema := range schemas {
		fileName := kind.fileName(schema.Name)
		logger.Infof("[generateArtifacts] 正在生成第 %d/%d 个 %s: %s", i+1, len(schemas), kind.label, schema.Name)

		outputPath := filepath.Join(getOutputPath(), kind.packageDir(), fileName)
		if err := generateArtifactFile(schema, string(tmplContent), kind.label, outputPath); err != nil {
			logger.Errorf("[generateArtifacts] 生成 %s 失败, 表名: %s, 文件名: %s, 错误: %v", kind.label, schema.Name, fileName, err)
			return err
		}
	}
	logger.Infof("[generateArtifacts] 所有 %s 生成完成, 共 %d 个", kind.label, len(schemas))
	return nil
}

// generateArtifactFile 渲染单表模板并写出文件
func generateArtifactFile(schema model.Schema, templateContent, label, outputFilePath string) error {
	logger.Infof("[generateArtifactFile] 开始生成 %s, 文件名: %s", label, outputFilePath)

	code, err := renderArtifactCode(schema, templateContent, label)
	if err != nil {
		logger.Errorf("[generateArtifactFile] 生成 %s 代码失败, 错误: %v", label, err)
		return fmt.Errorf("生成 %s 代码失败: %w", label, err)
	}

	if err = tool.WriteFileWithDir(outputFilePath, []byte(code), 0644); err != nil {
		logger.Errorf("[generateArtifactFile] 写入输出文件失败, 文件路径: %s, 错误: %v", outputFilePath, err)
		return fmt.Errorf("写入输出文件失败: %w", err)
	}

	logger.Infof("[generateArtifactFile] 成功生成 %s 文件: %s", label, outputFilePath)
	return nil
}

// renderArtifactCode 解析并执行模板，然后格式化 Go 代码
func renderArtifactCode(schema model.Schema, templateContent, name string) (code string, err error) {
	tmpl, err := template.New(name).Funcs(funcMap).Parse(templateContent)
	if err != nil {
		logger.Errorf("[renderArtifactCode] 解析模板失败, 错误: %v", err)
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, newTemplateData([]model.Schema{schema})); err != nil {
		logger.Errorf("[renderArtifactCode] 执行模板失败, 错误: %v", err)
		return "", fmt.Errorf("执行模板失败: %w", err)
	}

	return tool.FormatGoCode(buf.String()), nil
}

// GenerateTools 将嵌入式工具模板目录整包输出到目标项目
// 工具模板是静态代码（无模板变量），仅做格式化；单个文件失败跳过不中断
func GenerateTools() (err error) {
	logger.Infof("[GenerateTools] 开始生成所有工具文件")

	entries, err := fs.ReadDir(toolTemplateDir())
	if err != nil {
		logger.Errorf("[GenerateTools] 读取嵌入式工具模板目录失败, 错误: %v", err)
		return fmt.Errorf("读取嵌入式工具模板目录失败: %w", err)
	}
	logger.Infof("[GenerateTools] 找到 %d 个文件/目录", len(entries))

	generatedCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			logger.Infof("[GenerateTools] 跳过目录: %s", entry.Name())
			continue
		}

		tplName := entry.Name()
		if !strings.HasSuffix(tplName, templatePathSuffix) {
			logger.Infof("[GenerateTools] 跳过非模板文件: %s", tplName)
			continue
		}
		byts, err := fs.ReadFile(toolTemplatePath(tplName))
		if err != nil {
			logger.Errorf("[GenerateTools] 读取嵌入式工具模板文件失败, 模板路径: %s, 错误: %v", toolTemplatePath(tplName), err)
			continue
		}

		goFileName := strings.TrimSuffix(tplName, templatePathSuffix) + ".go"
		logger.Infof("[GenerateTools] 正在生成工具文件: %s", goFileName)

		outputPath := filepath.Join(getOutputPath(), getToolPackage(), goFileName)
		if err := tool.WriteFileWithDir(outputPath, []byte(tool.FormatGoCode(string(byts))), 0644); err != nil {
			logger.Errorf("[GenerateTools] 生成工具文件失败, 模板: %s, 输出: %s, 错误: %v", tplName, outputPath, err)
			continue
		}
		generatedCount++
	}

	logger.Infof("[GenerateTools] 所有工具文件生成完成, 共生成 %d 个文件", generatedCount)
	return nil
}
