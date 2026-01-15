package app

import (
	"fmt"
	"log"

	"github.com/LingoJack/model_infrax/internal/config"
	"github.com/LingoJack/model_infrax/internal/generator"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/internal/parser"
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

func (app *App) Run() error {
	var schemas []model.Schema
	var err error
	switch app.Config.GenerateConfig.GenerateMode {
	case "database":
		log.Println("🚀 开始从数据库解析表结构...")
		var databaseParser *parser.DatabaseParser
		databaseParser, err = parser.NewDatabaseParser(app.Config)
		if err != nil {
			return fmt.Errorf("初始化数据库解析器失败: %w", err)
		}
		schemas, err = databaseParser.Parse()
		if err != nil {
			return err
		}
		log.Printf("✅ 数据库解析完成，共获取到 %d 个表", len(schemas))
		schemas = databaseParser.FilterTables(schemas)

	case "statement":
		schemas, err = parser.ParseCreateTableStatements()
		if err != nil {
			return err
		}
		schemas = parser.FilterTables(schemas)
	default:
		return fmt.Errorf("不支持的生成模式: %s，请使用 'database' 或 'statement'", app.Config.GenerateConfig.GenerateMode)
	}

	log.Printf("🔍 过滤后需要处理的表数量: %d", len(schemas))

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
