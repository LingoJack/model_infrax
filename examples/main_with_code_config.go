//go:build example
// +build example

package main

import (
	"log"
	"model_infrax/config"
)

// 这个文件展示了如何使用 Go 代码配置来替代 YAML 配置文件
// 类似于 Wire 的使用方式，但更加灵活

func main() {
	// 方式1: 使用 Builder 直接创建应用
	app, err := NewAppFromBuilder(
		config.NewBuilder().
			DatabaseMode("localhost", 3306, "mydb", "root", "password").
			Tables("users", "orders", "products").
			OutputPath("./output").
			IgnoreTableNamePrefix(true),
	)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("运行失败: %v", err)
	}
}

// 方式2: 分步构建配置
func Example_StepByStep() {
	// 创建构建器
	builder := config.NewBuilder()

	// 配置数据库连接
	builder.DatabaseMode("localhost", 3306, "mydb", "root", "password")

	// 配置要生成的表
	builder.Tables("users", "orders")

	// 配置输出选项
	builder.OutputPath("./output")
	builder.IgnoreTableNamePrefix(true)
	builder.CrudOnlyIdx(true)

	// 配置包名
	builder.Packages("model", "dto", "vo", "dao", "utils")

	// 创建应用
	app, err := NewAppFromBuilder(builder)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}

	// 运行生成
	if err := app.Run(); err != nil {
		log.Fatalf("运行失败: %v", err)
	}
}

// 方式3: 使用条件逻辑（这是 Go 代码配置相比 YAML 的优势）
func Example_WithLogic() {
	builder := config.NewBuilder()

	// 根据环境决定配置
	isDev := true // 可以从环境变量读取

	if isDev {
		// 开发环境：从本地数据库生成
		builder.DatabaseMode("localhost", 3306, "dev_db", "root", "password")
	} else {
		// 生产环境：从SQL文件生成
		builder.StatementMode("./production_schema.sql")
	}

	// 动态决定要生成的表
	tables := []string{"users"}
	if needOrders := true; needOrders {
		tables = append(tables, "orders")
	}
	builder.Tables(tables...)

	// 创建并运行
	app, err := NewAppFromBuilder(builder)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("运行失败: %v", err)
	}
}