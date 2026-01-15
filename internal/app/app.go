package app

import (
	"context"
	"fmt"
	"log"

	"github.com/LingoJack/model_infrax/internal/config"
	"github.com/LingoJack/model_infrax/internal/constant"
	"github.com/LingoJack/model_infrax/internal/generator"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/internal/parser/database_parser"
	"github.com/LingoJack/model_infrax/internal/parser/statement_parser"
)

type App struct {
	Config    *config.Configger
	Generator *generator.Generator
}

func NewApp(cfg *config.Configger, g *generator.Generator) *App {
	return &App{
		Config:    cfg,
		Generator: g,
	}
}

func (app *App) Run(ctx context.Context) error {
	var schemas []model.Schema
	var err error
	switch app.Config.GenerateConfig.GenerateMode {
	case constant.GenerateModeDatabase:
		schemas, err = database_parser.Parse(ctx)
		if err != nil {
			return err
		}
		schemas = database_parser.Filter(schemas)
	case constant.GenerateModeStatement:
		schemas, err = statement_parser.Parse()
		if err != nil {
			return err
		}
		schemas = statement_parser.Filter(schemas)
	default:
		return fmt.Errorf("不支持的生成模式: %s，请使用 '%s' 或 '%s'", app.Config.GenerateConfig.GenerateMode, constant.GenerateModeDatabase, constant.GenerateModeStatement)
	}

	if len(schemas) == 0 {
		log.Println("⚠️ 没有找到需要处理的表，请检查配置文件中的表过滤规则")
		return nil
	}

	log.Println("🏗️ 开始生成 Model 代码...")
	if app.Config.GenerateOption.ModelAllInOneFile {
		err = app.Generator.GenerateModelAllInOne(schemas, app.Config.GenerateOption.ModelAllInOneFileName)
	} else {
		err = app.Generator.GenerateModelToEachFile(schemas)
	}
	if err != nil {
		return fmt.Errorf("生成Model代码失败: %w", err)
	}
	log.Println("✅ Model 代码生成完成")

	log.Println("📝 开始生成 DTO 代码...")
	err = app.Generator.GenerateDTOOneByOne(schemas)
	if err != nil {
		return fmt.Errorf("生成DTO代码失败: %w", err)
	}
	log.Println("✅ DTO 代码生成完成")

	log.Println("👁️ 开始生成 VO 代码...")
	err = app.Generator.GenerateVoToEachFile(schemas)
	if err != nil {
		return fmt.Errorf("生成VO代码失败: %w", err)
	}
	log.Println("✅ VO 代码生成完成")

	log.Println("🗄️ 开始生成 DAO 代码...")
	err = app.Generator.GenerateDAOOneByOne(schemas)
	if err != nil {
		return fmt.Errorf("生成DAO代码失败: %w", err)
	}
	log.Println("✅ DAO 代码生成完成")

	log.Println("🛠️ 开始生成 Tool 工具代码...")
	err = app.Generator.GenerateAllTools()
	if err != nil {
		return fmt.Errorf("生成Tool代码失败: %w", err)
	}

	log.Println("🎉 所有代码生成完成！")
	log.Printf("📊 生成统计: %d个表 -> Model + DTO + VO + DAO + Tools", len(schemas))

	return nil
}
