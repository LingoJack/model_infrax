package main

import (
	"fmt"

	"github.com/LingoJack/model_infrax/config"
	"github.com/LingoJack/model_infrax/generator"
	"github.com/LingoJack/model_infrax/model"
	"github.com/LingoJack/model_infrax/parser"

	"log"

	flag "github.com/spf13/pflag"
)

type App struct {
	Config    *config.Configger
	Generator *generator.Generator
}

// NewApp 创建应用实例
// 注意：DatabaseParser 和 StatementParser 不再作为依赖注入，而是在 Run 方法中根据模式动态创建
// 这样可以避免 statement 模式下不必要的数据库连接
func NewApp(cfg *config.Configger, g *generator.Generator) *App {
	return &App{
		Config:    cfg,
		Generator: g,
	}
}

// NewAppFromBuilder 从配置构建器创建应用实例
// 这是使用 Go 代码配置的便捷方法
// 示例:
//
//	app := NewAppFromBuilder(
//	    config.NewBuilder().
//	        DatabaseMode("localhost", 3306, "mydb", "root", "password").
//	        AllTables().
//	        OutputPath("./output"),
//	)
func NewAppFromBuilder(builder *config.ConfiggerBuilder) (*App, error) {
	cfg, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("构建配置失败: %w", err)
	}

	gen := generator.NewGenerator(cfg)
	return NewApp(cfg, gen), nil
}

// Run 运行应用程序，执行完整的代码生成流程
// 流程包括：
// 1. 解析数据库表结构（从数据库或SQL文件）
// 2. 根据配置过滤需要处理的表
// 3. 生成Model实体类代码
// 4. 生成DTO数据传输对象代码
// 5. 生成DAO数据访问对象代码
// 6. 生成Tool工具类代码
func (a *App) Run() error {
	var schemas []model.Schema
	var err error

	// 根据配置的生成模式选择不同的解析器
	// 采用延迟初始化策略：只在需要时才创建对应的解析器
	// 这样可以避免 statement 模式下不必要的数据库连接尝试
	switch a.Config.GenerateConfig.GenerateMode {
	case "database":
		// 从数据库解析表结构
		log.Println("🚀 开始从数据库解析表结构...")

		// 动态创建 DatabaseParser，只在 database 模式下才会尝试连接数据库
		var databaseParser *parser.DatabaseParser
		databaseParser, err = parser.NewDatabaseParser(a.Config)
		if err != nil {
			return fmt.Errorf("初始化数据库解析器失败: %w", err)
		}

		schemas, err = databaseParser.Parse()
		if err != nil {
			return err
		}
		log.Printf("✅ 数据库解析完成，共获取到 %d 个表", len(schemas))

		// 根据配置文件中的表名过滤规则，筛选需要生成代码的表
		schemas = databaseParser.FilterTables(schemas)

	case "statement":
		// 从SQL文件解析表结构
		log.Println("🚀 开始从SQL文件解析表结构...")

		// 动态创建 StatementParser，不需要数据库连接，只解析 SQL 文件
		var statementParser *parser.StatementParser
		statementParser, err = parser.NewStatementParser(a.Config)
		if err != nil {
			return fmt.Errorf("初始化SQL文件解析器失败: %w", err)
		}

		schemas, err = statementParser.Parse()
		if err != nil {
			return err
		}
		log.Printf("✅ SQL文件解析完成，共获取到 %d 个表", len(schemas))

		// 根据配置文件中的表名过滤规则，筛选需要生成代码的表
		schemas = statementParser.FilterTables(schemas)

	default:
		return fmt.Errorf("不支持的生成模式: %s，请使用 'database' 或 'statement'", a.Config.GenerateConfig.GenerateMode)
	}

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

var defaultConfigPaths = []string{
	"./application.yml",
	"./assert/application.yml",
	"/Applications/jen/application.yml",
	"/Applications/jen/assert/application.yml",
}

func main() {

	configPath := flag.StringP("config", "c", "./application.yml", "配置文件路径")

	flag.Parse()

	app, err := InitializeApp(*configPath)
	if err != nil {
		for _, path := range defaultConfigPaths {
			app, err = InitializeApp(path)
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	if err = app.Run(); err != nil {
		log.Fatalf("运行应用失败: %v", err)
	}
}
