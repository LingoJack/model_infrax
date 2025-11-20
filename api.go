package main

import (
	"fmt"
	"model_infrax/config"
)

// Generate 是对外暴露的主要 API，用于执行代码生成
// 这个函数可以被其他 Go 项目导入使用
//
// 使用示例:
//
//	import "github.com/LingoJack/model_infrax"
//
//	func main() {
//	    err := model_infrax.Generate(
//	        model_infrax.NewBuilder().
//	            DatabaseMode("localhost", 3306, "mydb", "root", "password").
//	            Tables("users", "orders").
//	            OutputPath("./generated"),
//	    )
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	}
func Generate(builder *config.ConfiggerBuilder) error {
	app, err := NewAppFromBuilder(builder)
	if err != nil {
		return fmt.Errorf("初始化应用失败: %w", err)
	}

	return app.Run()
}

// GenerateFromConfig 从配置文件路径生成代码
// 这是另一种使用方式，适合需要使用 YAML 配置文件的场景
//
// 使用示例:
//
//	import "github.com/LingoJack/model_infrax"
//
//	func main() {
//	    err := model_infrax.GenerateFromConfig("./application.yml")
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	}
func GenerateFromConfig(configPath string) error {
	app, err := InitializeApp(configPath)
	if err != nil {
		return fmt.Errorf("初始化应用失败: %w", err)
	}

	return app.Run()
}

// NewBuilder 创建一个新的配置构建器
// 这是对 config.NewBuilder 的封装，提供更简洁的 API
//
// 使用示例:
//
//	builder := model_infrax.NewBuilder().
//	    DatabaseMode("localhost", 3306, "mydb", "root", "password").
//	    AllTables().
//	    OutputPath("./output")
func NewBuilder() *config.ConfiggerBuilder {
	return config.NewBuilder()
}
