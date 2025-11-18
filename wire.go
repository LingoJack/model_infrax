//go:build wireinject
// +build wireinject

package main

import (
	"model_infrax/config"
	"model_infrax/generator"
	"model_infrax/parser"

	"github.com/google/wire"
)

// InitializeApp 初始化应用程序，Wire会自动生成依赖注入代码
// 注意：解析器（DatabaseParser/StatementParser）不再在此处注入
// 而是在 App.Run() 方法中根据配置的 GenerateMode 动态创建
// 这样可以避免 statement 模式下不必要的数据库连接
func InitializeApp(configPath string) (*App, error) {
	wire.Build(
		// 配置层
		config.NewConfigger,
		// 生成器层
		generator.NewGenerator,
		// 应用层
		NewApp,
	)
	return &App{}, nil
}

// InitializeDatabaseParser 单独初始化Parser，用于需要独立使用Parser的场景
func InitializeDatabaseParser(configPath string) (*parser.DatabaseParser, error) {
	wire.Build(
		// 配置层
		config.NewConfigger,
		// 解析器层
		parser.NewDatabaseParser,
	)
	return &parser.DatabaseParser{}, nil
}

func InitializeStatementParser(configPath string) (*parser.StatementParser, error) {
	wire.Build(config.NewConfigger, parser.NewStatementParser)
	return &parser.StatementParser{}, nil
}

func InitializeGenerator(configPath string) (*generator.Generator, error) {
	wire.Build(config.NewConfigger, generator.NewGenerator)
	return &generator.Generator{}, nil
}
