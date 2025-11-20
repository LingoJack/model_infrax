//go:build wireinject
// +build wireinject

package main

import (
	"github.com/LingoJack/model_infrax/config"
	"github.com/LingoJack/model_infrax/generator"
	"github.com/LingoJack/model_infrax/pkg/app"
	"github.com/google/wire"
)

// provideConfiggerForAPI 提供API使用的配置管理器
func provideConfiggerForAPI(configPath string) (*config.Configger, error) {
	return config.NewConfigger(configPath)
}

// provideGeneratorForAPI 提供API使用的代码生成器
func provideGeneratorForAPI(cfg *config.Configger) *generator.Generator {
	return generator.NewGenerator(cfg)
}

// provideAppForAPI 提供API使用的应用实例
func provideAppForAPI(cfg *config.Configger, gen *generator.Generator) *app.App {
	return app.NewApp(cfg, gen)
}

// InitializeAppForAPI 使用 Wire 初始化API应用
// 这个函数专门为 api.go 中的 GenerateFromConfig 函数提供服务
//
// 参数:
//   - configPath: 配置文件的路径
//
// 返回:
//   - *app.App: 初始化完成的应用实例
//   - error: 初始化过程中的错误，nil表示成功
func InitializeAppForAPI(configPath string) (*app.App, error) {
	wire.Build(
		provideConfiggerForAPI,
		provideGeneratorForAPI,
		provideAppForAPI,
	)
	return nil, nil // 由 Wire 生成的代码会替换这个返回值
}
