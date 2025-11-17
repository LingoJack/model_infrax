package main

import (
	"log"
	"model_infrax/config"
	"model_infrax/generator"
	"model_infrax/parser"

	flag "github.com/spf13/pflag"
)

type App struct {
	Config    *config.Configger
	Parser    *parser.DatabaseParser
	Generator *generator.Generator
}

func NewApp(cfg *config.Configger, p *parser.DatabaseParser, g *generator.Generator) *App {
	return &App{
		Config:    cfg,
		Parser:    p,
		Generator: g,
	}
}

// Run 运行应用程序，执行完整的代码生成流程
// 流程包括：
// 1. 解析数据库表结构
// 2. 根据配置过滤需要处理的表
// 3. 生成Model实体类代码
// 4. 生成DTO数据传输对象代码
// 5. 生成DAO数据访问对象代码
// 6. 生成Tool工具类代码
func (a *App) Run() error {
	log.Println("🚀 开始解析数据库...")

	// 获取数据库中所有表的结构信息
	schemas, err := a.Parser.Parse()
	if err != nil {
		return err
	}

	log.Printf("✅ 解析完成，共获取到 %d 个表", len(schemas))

	// 根据配置文件中的表名过滤规则，筛选需要生成代码的表
	schemas = a.Parser.FilterTables(schemas)
	log.Printf("🔍 过滤后需要处理的表数量: %d", len(schemas))

	log.Println("🏗️ 开始生成 Model 代码...")

	// 生成Model实体类代码
	// 根据配置决定是生成到一个文件还是分别生成
	if a.Config.GenerateOption.ModelAllInOneFile {
		err = a.Generator.GenerateModel(schemas, a.Config.GenerateOption.ModelAllInOneFileName)
	} else {
		err = a.Generator.GenerateModelOneByOne(schemas)
	}
	if err != nil {
		return err
	}

	log.Println("✅ Model 代码生成完成")

	log.Println("📝 开始生成 DTO 代码...")

	// 生成DTO数据传输对象代码，用于API接口的数据交换
	err = a.Generator.GenerateDTOOneByOne(schemas)
	if err != nil {
		return err
	}

	log.Println("✅ DTO 代码生成完成")

	log.Println("🗄️ 开始生成 DAO 代码...")

	// 生成DAO数据访问对象代码，提供数据库操作方法
	err = a.Generator.GenerateDAOOneByOne(schemas)
	if err != nil {
		return err
	}

	log.Println("✅ DAO 代码生成完成")

	log.Println("🛠️ 开始生成 Tool 工具代码...")

	// 生成Tool工具类代码，提供通用的辅助功能
	err = a.Generator.GenerateAllTools()
	if err != nil {
		return err
	}

	log.Println("🎉 所有代码生成完成！")

	return nil
}

func main() {

	configPath := flag.StringP("config", "c", "./assert/application.yml", "配置文件路径")

	flag.Parse()

	app, err := InitializeApp(*configPath)
	if err != nil {
		log.Fatalf("初始化应用失败: %v", err)
	}

	if err = app.Run(); err != nil {
		log.Fatalf("运行应用失败: %v", err)
	}
}
