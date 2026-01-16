package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/tool"
)

func GenerateTools() (err error) {
	logger.Infof("[GenerateAllTools]开始生成所有工具文件")

	entries, err := fs.ReadDir(ToolTemplateDir())
	if err != nil {
		logger.Errorf("[GenerateAllTools]读取嵌入式工具模板目录失败, 错误: %v", err)
		return fmt.Errorf("读取嵌入式工具模板目录失败: %w", err)
	}
	logger.Infof("[GenerateAllTools]找到 %d 个文件/目录", len(entries))

	generatedCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			logger.Infof("[GenerateAllTools]跳过目录: %s", entry.Name())
			continue
		}

		tplName := entry.Name()
		if !strings.HasSuffix(tplName, templatePathSuffix) {
			logger.Infof("[GenerateAllTools]跳过非模板文件: %s", tplName)
			continue
		}
		var byts []byte
		byts, err = fs.ReadFile(ToolTemplatePath(tplName))
		if err != nil {
			logger.Errorf("[GenerateAllTools]读取嵌入式工具模板文件失败, 模板路径: %s, 错误: %v", ToolTemplatePath(tplName), err)
			err = nil
			continue
		}

		goFileName := strings.TrimSuffix(tplName, templatePathSuffix) + ".go"

		logger.Infof("[GenerateAllTools]正在生成工具文件: %s", goFileName)

		path := filepath.Join(outputPath, toolPackage, goFileName)
		err = GenerateTool(string(byts), path)
		if err != nil {
			logger.Errorf("[GenerateAllTools]生成工具文件失败, 模板: %s, 输出: %s, 错误: %v", tplName, goFileName, err)
			err = nil
			continue
		}
		generatedCount++
	}

	logger.Infof("[GenerateAllTools]所有工具文件生成完成, 共生成 %d 个文件", generatedCount)
	return nil
}

func GenerateTool(templateContent, outputFilePath string) (err error) {
	logger.Infof("[GenerateTool]开始生成工具文件, 输出: %s", outputFilePath)

	code, err := GenerateToolCode(templateContent)
	if err != nil {
		logger.Errorf("[GenerateTool]生成工具代码失败, 错误: %v", err)
		return fmt.Errorf("生成工具代码失败: %w", err)
	}

	err = tool.WriteFileWithDir(outputFilePath, []byte(code), 0644)
	if err != nil {
		logger.Errorf("[GenerateTool]写入工具输出文件失败, 文件路径: %s, 错误: %v", outputFilePath, err)
		return fmt.Errorf("写入工具输出文件失败: %w", err)
	}

	logger.Infof("[GenerateTool]成功生成工具文件: %s", outputFilePath)
	return nil
}

func GenerateToolCode(templateContent string) (code string, err error) {
	code = tool.FormatGoCode(templateContent)
	return
}
