package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/internal/tool"
)

func GenerateVos(schemas []model.Schema) (err error) {
	logger.Infof("[GenerateVoOneByOne] 开始生成 VO, 共 %d 个表", len(schemas))
	for i, schema := range schemas {

		fileName := fmt.Sprintf("%s_vo.go", schema.Name)

		logger.Infof("[GenerateVoOneByOne] 正在生成第 %d/%d 个 VO: %s", i+1, len(schemas), schema.Name)

		// 从嵌入的文件系统中读取 VO 模板文件
		var tmplContent []byte
		tmplContent, err = fs.ReadFile(VoTemplatePath())
		if err != nil {
			logger.Errorf("[GenerateVO] 读取嵌入式 VO 模板文件失败, 模板路径: %s, 错误: %v", VoTemplatePath(), err)
			return fmt.Errorf("读取嵌入式 VO 模板文件失败: %w", err)
		}

		path := filepath.Join(outputPath, voPackage, fileName)
		err = GenerateVO(schema, string(tmplContent), path)
		if err != nil {
			logger.Errorf("[GenerateVoOneByOne] 生成 VO 失败, 表名: %s, 文件名: %s, 错误: %v", schema.Name, fileName, err)
			return err
		}
	}
	logger.Infof("[GenerateVoOneByOne] 所有 VO 生成完成, 共 %d 个", len(schemas))
	return nil
}

func GenerateVO(schema model.Schema, templateContent, outputFilePath string) (err error) {
	logger.Infof("[GenerateVO] 开始生成 VO, 文件名: %s", outputFilePath)

	code, err := GenerateVoCode(schema, templateContent)
	if err != nil {
		logger.Errorf("[GenerateVO] 生成 VO 代码失败, 错误: %v", err)
		return fmt.Errorf("生成 VO 代码失败: %w", err)
	}

	err = tool.WriteFileWithDir(outputFilePath, []byte(code), 0644)
	if err != nil {
		logger.Errorf("[GenerateVO] 写入输出文件失败, 文件路径: %s, 错误: %v", outputFilePath, err)
		return fmt.Errorf("写入输出文件失败: %w", err)
	}

	logger.Infof("[GenerateVO] 成功生成文件: %s", outputFilePath)
	return nil
}

func GenerateVoCode(schema model.Schema, templateContent string) (code string, err error) {
	tmpl, err := template.New("vo").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		logger.Errorf("[GenerateVO] 解析 VO 模板失败, 错误: %v", err)
		err = fmt.Errorf("解析 VO 模板失败: %w", err)
		return
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, TemplateData{
		DaoPackageName: getPackageName(daoPackage),
		PoPackageName:  getPackageName(poPackage),
		DtoPackageName: getPackageName(dtoPackage),
		VoPackageName:  getPackageName(voPackage),
		Schemas:        []model.Schema{schema},
	}); err != nil {
		logger.Errorf("[GenerateVO] 执行 VO 模板失败, 错误: %v", err)
		err = fmt.Errorf("执行 VO 模板失败: %w", err)
		return
	}

	code = buf.String()

	code = tool.FormatGoCode(code)

	return
}
