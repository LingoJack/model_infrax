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

func GenerateDaos(schemas []model.Schema) (err error) {
	logger.Infof("[GenerateDaoOneByOne] 开始生成 DAO, 共 %d 个表", len(schemas))
	for i, schema := range schemas {

		fileName := fmt.Sprintf("%s_dao.go", schema.Name)

		logger.Infof("[GenerateDaoOneByOne] 正在生成第 %d/%d 个 DAO: %s", i+1, len(schemas), schema.Name)

		var tmplContent []byte
		tmplContent, err = fs.ReadFile(DaoTemplatePath())
		if err != nil {
			logger.Errorf("[GenerateDao] 读取嵌入式 DAO 模板文件失败, 模板路径: %s, 错误: %v", DaoTemplatePath(), err)
			return fmt.Errorf("读取嵌入式 DAO 模板文件失败: %w", err)
		}

	path := filepath.Join(getOutputPath(), getDaoPackage(), fileName)
	err = GenerateDao(schema, string(tmplContent), path)
		if err != nil {
			logger.Errorf("[GenerateDaoOneByOne] 生成 DAO 失败, 表名: %s, 文件名: %s, 错误: %v", schema.Name, fileName, err)
			return err
		}
	}
	logger.Infof("[GenerateDaoOneByOne] 所有 DAO 生成完成, 共 %d 个", len(schemas))
	return nil
}

func GenerateDao(schema model.Schema, templateContent, outputFilePath string) (err error) {
	logger.Infof("[GenerateDao] 开始生成 DAO, 文件名: %s", outputFilePath)

	code, err := GenerateDaoCode(schema, templateContent)
	if err != nil {
		logger.Errorf("[GenerateDao] 生成 DAO 代码失败, 错误: %v", err)
		return fmt.Errorf("生成 DAO 代码失败: %w", err)
	}

	err = tool.WriteFileWithDir(outputFilePath, []byte(code), 0644)
	if err != nil {
		logger.Errorf("[GenerateDao] 写入 DAO 输出文件失败, 文件路径: %s, 错误: %v", outputFilePath, err)
		return fmt.Errorf("写入 DAO 输出文件失败: %w", err)
	}

	logger.Infof("[GenerateDao] 成功生成 DAO 文件: %s", outputFilePath)
	return nil
}

func GenerateDaoCode(schema model.Schema, templateContent string) (code string, err error) {
	tmpl, err := template.New("dao").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		logger.Errorf("[GenerateDao] 解析 DAO 模板失败, 错误: %v", err)
		err = fmt.Errorf("解析 DAO 模板失败: %w", err)
		return
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, TemplateData{
		DaoPackageName: getPackageName(getDaoPackage()),
		PoPackageName:  getPackageName(getPoPackage()),
		DtoPackageName: getPackageName(getDtoPackage()),
		Schemas:        []model.Schema{schema},
	}); err != nil {
		logger.Errorf("[GenerateDao] 执行 DAO 模板失败, 错误: %v", err)
		err = fmt.Errorf("执行 DAO 模板失败: %w", err)
		return
	}

	code = buf.String()

	code = tool.FormatGoCode(code)

	return
}
