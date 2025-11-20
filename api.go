package main

import (
	"fmt"

	"github.com/LingoJack/model_infrax/config"
	"github.com/LingoJack/model_infrax/generator"
	"github.com/LingoJack/model_infrax/pkg/app"
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
	appInstance, err := app.NewAppFromBuilder(builder)
	if err != nil {
		return fmt.Errorf("初始化应用失败: %w", err)
	}

	return appInstance.Run()
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
	appInstance, err := initializeAppForAPI(configPath)
	if err != nil {
		return fmt.Errorf("初始化应用失败: %w", err)
	}

	return appInstance.Run()
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

// initializeAppForAPI 为API提供应用初始化功能
// 重新实现原来wire生成的InitializeApp函数的逻辑，但使用pkg/app包
//
// 参数:
//   - configPath: 配置文件的路径
//
// 返回:
//   - *app.App: 初始化完成的应用实例
//   - error: 初始化过程中的错误，nil表示成功
func initializeAppForAPI(configPath string) (*app.App, error) {
	// 加载配置文件
	configger, err := config.NewConfigger(configPath)
	if err != nil {
		return nil, err
	}

	// 创建代码生成器实例
	generatorGenerator := generator.NewGenerator(configger)

	// 创建应用实例，使用pkg/app包中的NewApp函数
	appInstance := app.NewApp(configger, generatorGenerator)
	
	return appInstance, nil
}
