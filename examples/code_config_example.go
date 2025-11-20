//go:build example
// +build example

package main

import (
	"log"
	"model_infrax/config"
	"model_infrax/generator"
)

// 示例1: 从数据库生成代码
func Example_DatabaseMode() {
	// 使用 Builder 模式配置生成器
	// 这种方式类似 Wire，但更加灵活和类型安全
	cfg := config.NewBuilder().
		DatabaseMode("localhost", 3306, "mydb", "root", "password").
		Tables("users", "orders", "products"). // 指定要生成的表
		OutputPath("./output").                // 输出路径
		IgnoreTableNamePrefix(true).           // 忽略表名前缀
		CrudOnlyIdx(true).                     // 只为索引字段生成CRUD
		Packages("model", "dto", "vo", "dao", "utils"). // 自定义包名
		MustBuild() // 构建配置，如果有错误会panic

	// 创建生成器
	gen := generator.NewGenerator(cfg)

	// 创建应用并运行
	app := &App{
		Config:    cfg,
		Generator: gen,
	}

	if err := app.Run(); err != nil {
		log.Fatalf("生成失败: %v", err)
	}
}

// 示例2: 从SQL文件生成代码
func Example_StatementMode() {
	cfg := config.NewBuilder().
		StatementMode("./schema.sql"). // 从SQL文件解析
		AllTables().                   // 生成所有表
		OutputPath("~/projects/myapp/generated"). // 支持 ~ 符号
		ModelAllInOneFile(true, "models.go").     // 所有Model放在一个文件
		UseFramework("itea-go").                  // 使用特定框架模板
		MustBuild()

	gen := generator.NewGenerator(cfg)
	app := &App{
		Config:    cfg,
		Generator: gen,
	}

	if err := app.Run(); err != nil {
		log.Fatalf("生成失败: %v", err)
	}
}

// 示例3: 高级配置 - 使用条件逻辑
func Example_AdvancedConfig() {
	builder := config.NewBuilder()

	// 根据环境变量决定使用哪种模式
	// 这是 Go 代码配置相比 YAML 的优势：可以使用编程逻辑
	useDatabase := true // 实际可以从环境变量读取

	if useDatabase {
		builder.DatabaseMode("localhost", 3306, "mydb", "root", "password")
	} else {
		builder.StatementMode("./schema.sql")
	}

	// 动态决定要生成的表
	tables := []string{"users", "orders"}
	if needProducts := true; needProducts {
		tables = append(tables, "products")
	}
	builder.Tables(tables...)

	// 构建配置
	cfg, err := builder.
		OutputPath("./output").
		IgnoreTableNamePrefix(true).
		Build()

	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}

	gen := generator.NewGenerator(cfg)
	app := &App{
		Config:    cfg,
		Generator: gen,
	}

	if err := app.Run(); err != nil {
		log.Fatalf("生成失败: %v", err)
	}
}

// 示例4: 最小配置
func Example_MinimalConfig() {
	// 最简单的配置方式
	cfg := config.NewBuilder().
		DatabaseMode("localhost", 3306, "mydb", "root", "password").
		AllTables().
		MustBuild()

	gen := generator.NewGenerator(cfg)
	app := &App{
		Config:    cfg,
		Generator: gen,
	}

	if err := app.Run(); err != nil {
		log.Fatalf("生成失败: %v", err)
	}
}

// 示例5: 自定义包名
func Example_CustomPackages() {
	cfg := config.NewBuilder().
		StatementMode("./schema.sql").
		AllTables().
		OutputPath("./output").
		// 方式1: 一次性设置所有包名
		Packages("entity", "request", "response", "repository", "helper").
		MustBuild()

	// 或者使用方式2: 单独设置每个包名
	cfg2 := config.NewBuilder().
		StatementMode("./schema.sql").
		AllTables().
		OutputPath("./output").
		PoPackage("entity").
		DtoPackage("request").
		VoPackage("response").
		DaoPackage("repository").
		ToolPackage("helper").
		MustBuild()

	_ = cfg2 // 避免未使用变量警告

	gen := generator.NewGenerator(cfg)
	app := &App{
		Config:    cfg,
		Generator: gen,
	}

	if err := app.Run(); err != nil {
		log.Fatalf("生成失败: %v", err)
	}
}

// 示例6: 在测试中使用
func Example_InTest() {
	// 在测试中可以快速创建配置
	cfg := config.NewBuilder().
		StatementMode("./testdata/schema.sql").
		Tables("test_users").
		OutputPath("./testdata/output").
		MustBuild()

	gen := generator.NewGenerator(cfg)
	app := &App{
		Config:    cfg,
		Generator: gen,
	}

	if err := app.Run(); err != nil {
		log.Fatalf("测试生成失败: %v", err)
	}
}