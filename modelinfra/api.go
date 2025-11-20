package modelinfra

import (
	"fmt"

	"github.com/LingoJack/model_infrax/config"
	"github.com/LingoJack/model_infrax/pkg/app"
)

// Generate 是对外暴露的主要 API，用于执行代码生成
// 这个函数可以被其他 Go 项目导入使用
//
// 使用示例:
//
//	import "github.com/LingoJack/model_infrax/modelinfra"
//
//	func main() {
//	    err := modelinfra.Generate(
//	        modelinfra.NewBuilder().
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
// 使用 Wire 进行依赖注入，提供更清晰和可维护的代码结构
//
// 使用示例:
//
//	import "github.com/LingoJack/model_infrax/modelinfra"
//
//	func main() {
//	    err := modelinfra.GenerateFromConfig("./application.yml")
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	}
func GenerateFromConfig(configPath string) error {
	// 使用 Wire 生成的 InitializeAppForAPI 函数进行依赖注入
	appInstance, err := InitializeAppForAPI(configPath)
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
//	builder := modelinfra.NewBuilder().
//	    DatabaseMode("localhost", 3306, "mydb", "root", "password").
//	    AllTables().
//	    OutputPath("./output")
func NewBuilder() *config.ConfiggerBuilder {
	return config.NewBuilder()
}

// initializeAppForAPI 函数已由 Wire 自动生成
// Wire 会在当前目录的 wire_gen.go 中生成 InitializeAppForAPI 函数
// 该函数负责处理API模式下的所有依赖注入和初始化逻辑
// 不再需要手动实现此函数