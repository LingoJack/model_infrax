package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/pkg/logger"
	"github.com/LingoJack/model_infrax/pkg/tool"
)

func GenerateModels(schemas []model.Schema) (err error) {
	logger.Infof("[GenerateModelOneByOne] 开始生成模型, 共 %d 个表", len(schemas))
	for i, schema := range schemas {

		fileName := fmt.Sprintf("%s.go", schema.Name)

		logger.Infof("[GenerateModelOneByOne] 正在生成第 %d/%d 个模型: %s", i+1, len(schemas), schema.Name)

		var tmplContent []byte
		tmplContent, err = fs.ReadFile(PoTemplatePath())
		if err != nil {
			logger.Errorf("[GenerateModel] 读取嵌入模板文件失败, 模板路径: %s, 错误: %v", PoTemplatePath(), err)
			return fmt.Errorf("读取嵌入式模板文件失败: %w", err)
		}

		path := filepath.Join(outputPath, poPackage, fileName)
		err = GenerateModel(schema, string(tmplContent), path)
		if err != nil {
			logger.Errorf("[GenerateModelOneByOne] 生成模型失败, 表名: %s, 文件名: %s, 错误: %v", schema.Name, fileName, err)
			return err
		}
	}
	logger.Infof("[GenerateModelOneByOne] 所有模型生成完成, 共 %d 个", len(schemas))
	return nil
}

func GenerateModel(schema model.Schema, templateContent, outputFilePath string) (err error) {
	logger.Infof("[GenerateModel] 开始生成模型, 文件名: %s", outputFilePath)

	code, err := GenerateModelCode(schema, templateContent)
	if err != nil {
		logger.Errorf("[GenerateModel] 生成模型代码失败, 错误: %v", err)
		return fmt.Errorf("生成模型代码失败: %w", err)
	}

	if err = tool.WriteFileWithDir(outputFilePath, []byte(code), 0644); err != nil {
		logger.Errorf("[GenerateModel] 写入输出文件失败, 文件路径: %s, 错误: %v", outputFilePath, err)
		return fmt.Errorf("写入输出文件失败: %w", err)
	}

	logger.Infof("[GenerateModel] 成功生成模型文件: %s", outputFilePath)
	return nil
}

func GenerateModelCode(schema model.Schema, templateContent string) (code string, err error) {
	tmpl, err := template.New("model").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		logger.Errorf("[GenerateModel] 解析模板失败, 错误: %v", err)
		err = fmt.Errorf("解析模板失败: %w", err)
		return
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, TemplateData{
		DaoPackageName: getPackageName(daoPackage),
		PoPackageName:  getPackageName(poPackage),
		DtoPackageName: getPackageName(dtoPackage),
		Schemas:        []model.Schema{schema},
	}); err != nil {
		logger.Errorf("[GenerateModel] 执行模板失败, 错误: %v", err)
		err = fmt.Errorf("执行模板失败: %w", err)
		return
	}

	code = buf.String()

	code = tool.FormatGoCode(code)

	return
}
