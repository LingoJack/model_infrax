package app

import (
	"fmt"
	"log"

	"github.com/LingoJack/model_infrax/config"
	"github.com/LingoJack/model_infrax/generator"
	"github.com/LingoJack/model_infrax/model"
	"github.com/LingoJack/model_infrax/parser"
)

// App 应用程序核心结构体
// 负责整个代码生成流程的协调和执行
// 包含配置管理和代码生成器两个核心组件
type App struct {
	// Config 配置管理器，负责加载和管理应用配置
	Config    *config.Configger
	// Generator 代码生成器，负责生成各种类型的代码文件
	Generator *generator.Generator
}

// NewApp 创建应用实例
// 使用依赖注入模式，将配置和生成器注入到App中
// 注意：DatabaseParser 和 StatementParser 不再作为依赖注入，而是在 Run 方法中根据模式动态创建
// 这样可以避免 statement 模式下不必要的数据库连接，提升性能
//
// 参数:
//   - cfg: 配置管理器实例
//   - g: 代码生成器实例
//
// 返回:
//   - *App: 应用实例指针
func NewApp(cfg *config.Configger, g *generator.Generator) *App {
	return &App{
		Config:    cfg,
		Generator: g,
	}
}

// NewAppFromBuilder 从配置构建器创建应用实例
// 这是使用 Go 代码配置的便捷方法，提供链式调用的配置方式
// 适用于需要在代码中直接配置应用场景，避免手动创建配置文件的繁琐
//
// 使用示例:
//
//	app, err := NewAppFromBuilder(
//	    config.NewBuilder().
//	        DatabaseMode("localhost", 3306, "mydb", "root", "password").
//	        AllTables().
//	        OutputPath("./output"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	app.Run()
//
// 参数:
//   - builder: 配置构建器实例，用于构建应用配置
//
// 返回:
//   - *App: 应用实例指针
//   - error: 构建过程中的错误，nil表示成功
func NewAppFromBuilder(builder *config.ConfiggerBuilder) (*App, error) {
	// 构建配置对象
	cfg, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("构建配置失败: %w", err)
	}

	// 基于配置创建代码生成器
	gen := generator.NewGenerator(cfg)
	
	// 创建并返回应用实例
	return NewApp(cfg, gen), nil
}

// Run 运行应用程序，执行完整的代码生成流程
// 这是应用的核心方法，负责协调整个代码生成过程
//
// 生成流程包括：
// 1. 根据配置模式（database/statement）选择合适的解析器
// 2. 解析数据库表结构或SQL文件
// 3. 根据配置过滤需要处理的表
// 4. 生成Model实体类代码
// 5. 生成DTO数据传输对象代码
// 6. 生成DAO数据访问对象代码
// 7. 生成Tool工具类代码
//
// 返回:
//   - error: 执行过程中的错误，nil表示成功完成
func (a *App) Run() error {
	var schemas []model.Schema
	var err error

	// 根据配置的生成模式选择不同的解析器
	// 采用延迟初始化策略：只在需要时才创建对应的解析器
	// 这样可以避免 statement 模式下不必要的数据库连接尝试，提升启动速度
	switch a.Config.GenerateConfig.GenerateMode {
	case "database":
		// 从数据库解析表结构模式
		// 适用于可以直接连接数据库的场景，能够获取到最准确的表结构信息
		log.Println("🚀 开始从数据库解析表结构...")

		// 动态创建 DatabaseParser，只在 database 模式下才会尝试连接数据库
		// 这种设计避免了在不需要数据库连接时进行连接尝试
		var databaseParser *parser.DatabaseParser
		databaseParser, err = parser.NewDatabaseParser(a.Config)
		if err != nil {
			return fmt.Errorf("初始化数据库解析器失败: %w", err)
		}

		// 解析数据库表结构
		schemas, err = databaseParser.Parse()
		if err != nil {
			return err
		}
		log.Printf("✅ 数据库解析完成，共获取到 %d 个表", len(schemas))

		// 根据配置文件中的表名过滤规则，筛选需要生成代码的表
		// 支持白名单、黑名单、正则表达式等多种过滤方式
		schemas = databaseParser.FilterTables(schemas)

	case "statement":
		// 从SQL文件解析表结构模式
		// 适用于无法直接连接数据库，但有SQL DDL文件的场景
		log.Println("🚀 开始从SQL文件解析表结构...")

		// 动态创建 StatementParser，不需要数据库连接，只解析 SQL 文件
		// 这种模式适用于离线环境或者CI/CD流程
		var statementParser *parser.StatementParser
		statementParser, err = parser.NewStatementParser(a.Config)
		if err != nil {
			return fmt.Errorf("初始化SQL文件解析器失败: %w", err)
		}

		// 解析SQL文件中的表结构定义
		schemas, err = statementParser.Parse()
		if err != nil {
			return err
		}
		log.Printf("✅ SQL文件解析完成，共获取到 %d 个表", len(schemas))

		// 根据配置文件中的表名过滤规则，筛选需要生成代码的表
		schemas = statementParser.FilterTables(schemas)

	default:
		// 不支持的生成模式，返回明确的错误信息
		return fmt.Errorf("不支持的生成模式: %s，请使用 'database' 或 'statement'", a.Config.GenerateConfig.GenerateMode)
	}

	// 输出过滤后的表数量，方便用户了解处理范围
	log.Printf("🔍 过滤后需要处理的表数量: %d", len(schemas))

	// 检查是否有表需要处理，如果没有则提前退出
	if len(schemas) == 0 {
		log.Println("⚠️ 没有找到需要处理的表，请检查配置文件中的表过滤规则")
		return nil
	}

	// 开始生成Model代码
	// Model是数据实体类，用于表示数据库表结构
	log.Println("🏗️ 开始生成 Model 代码...")

	// 根据配置决定是生成到一个文件还是分别生成到多个文件
	// 支持两种模式：
	// 1. 所有Model生成到一个文件（适合小项目）
	// 2. 每个Model生成到独立文件（适合大项目，便于维护）
	if a.Config.GenerateOption.ModelAllInOneFile {
		err = a.Generator.GenerateModel(schemas, a.Config.GenerateOption.ModelAllInOneFileName)
	} else {
		err = a.Generator.GenerateModelOneByOne(schemas)
	}
	if err != nil {
		return fmt.Errorf("生成Model代码失败: %w", err)
	}

	log.Println("✅ Model 代码生成完成")

	// 开始生成DTO代码
	// DTO（Data Transfer Object）是数据传输对象，用于API接口的数据交换
	// 包含请求和响应的数据结构，以及数据验证和转换逻辑
	log.Println("📝 开始生成 DTO 代码...")

	// 生成DTO数据传输对象代码，每个表生成对应的DTO结构
	err = a.Generator.GenerateDTOOneByOne(schemas)
	if err != nil {
		return fmt.Errorf("生成DTO代码失败: %w", err)
	}

	log.Println("✅ DTO 代码生成完成")

	// 开始生成DAO代码
	// DAO（Data Access Object）是数据访问对象，提供数据库操作方法
	// 包含增删改查等基本操作，以及复杂查询方法
	log.Println("🗄️ 开始生成 DAO 代码...")

	// 生成DAO数据访问对象代码，每个表生成对应的DAO结构
	err = a.Generator.GenerateDAOOneByOne(schemas)
	if err != nil {
		return fmt.Errorf("生成DAO代码失败: %w", err)
	}

	log.Println("✅ DAO 代码生成完成")

	// 开始生成Tool工具类代码
	// Tool提供通用的辅助功能，如字符串处理、数据转换、验证等
	// 这些工具类可以在整个项目中复用
	log.Println("🛠️ 开始生成 Tool 工具代码...")

	// 生成Tool工具类代码，包括各种通用的辅助方法
	err = a.Generator.GenerateAllTools()
	if err != nil {
		return fmt.Errorf("生成Tool代码失败: %w", err)
	}

	log.Println("🎉 所有代码生成完成！")
	log.Printf("📊 生成统计: %d个表 -> Model + DTO + DAO + Tools", len(schemas))

	return nil
}