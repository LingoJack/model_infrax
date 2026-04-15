package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/pkg/tool"
)

func GenerateDtos(schemas []model.Schema) (err error) {
	logger.Infof("[GenerateDtoOneByOne] 开始生成 DTO, 共 %d 个表", len(schemas))
	for i, schema := range schemas {

		fileName := fmt.Sprintf("%s_dto.go", schema.Name)

		logger.Infof("[GenerateDtoOneByOne] 正在生成第 %d/%d 个 DTO: %s", i+1, len(schemas), schema.Name)

		var tmplContent []byte
		tmplContent, err = fs.ReadFile(DtoTemplatePath())
		if err != nil {
			logger.Errorf("[GenerateDto] 读取嵌入模板文件失败, 模板路径: %s, 错误: %v", DtoTemplatePath(), err)
			return fmt.Errorf("读取嵌入式模板文件失败: %w", err)
		}

	path := filepath.Join(getOutputPath(), getDtoPackage(), fileName)
	err = GenerateDto(schema, string(tmplContent), path)
		if err != nil {
			logger.Errorf("[GenerateDtoOneByOne] 生成 DTO 失败, 表名: %s, 文件名: %s, 错误: %v", schema.Name, fileName, err)
			return err
		}
	}
	logger.Infof("[GenerateDtoOneByOne] 所有 DTO 生成完成, 共 %d 个", len(schemas))
	return nil
}

func GenerateDto(schema model.Schema, templateContent, outputFilePath string) (err error) {
	logger.Infof("[GenerateDTO] 开始生成 DTO, 文件名: %s", outputFilePath)

	code, err := GenerateDtoCode(schema, templateContent)
	if err != nil {
		logger.Errorf("[GenerateDTO] 生成 DTO 代码失败, 错误: %v", err)
		return fmt.Errorf("生成 DTO 代码失败: %w", err)
	}

	err = tool.WriteFileWithDir(outputFilePath, []byte(code), 0644)
	if err != nil {
		logger.Errorf("[GenerateDTO] 写入 DTO 输出文件失败, 文件路径: %s, 错误: %v", outputFilePath, err)
		return fmt.Errorf("写入 DTO 输出文件失败: %w", err)
	}

	logger.Infof("[GenerateDTO] 成功生成 DTO 文件: %s", outputFilePath)
	return nil
}

func GenerateDtoCode(schema model.Schema, templateContent string) (code string, err error) {
	tmpl, err := template.New("dto").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		logger.Errorf("[GenerateDTO] 解析 DTO 模板失败, 错误: %v", err)
		err = fmt.Errorf("解析 DTO 模板失败: %w", err)
		return
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, newTemplateData([]model.Schema{schema})); err != nil {
		logger.Errorf("[GenerateDTO] 执行 DTO 模板失败, 错误: %v", err)
		err = fmt.Errorf("执行 DTO 模板失败: %w", err)
		return
	}
	code = buf.String()

	code = tool.FormatGoCode(code)

	return
}
