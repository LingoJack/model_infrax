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
func InitializeApp(configPath string) (*App, error) {
	wire.Build(
		// 配置层
		config.NewConfigger,
		// 解析器层
		parser.NewParser,
		// 生成器层
		generator.NewGenerator,
		// 应用层
		NewApp,
	)
	return &App{}, nil
}

// InitializeParser 单独初始化Parser，用于需要独立使用Parser的场景
func InitializeParser(configPath string) (*parser.Parser, error) {
	wire.Build(
		// 配置层
		config.NewConfigger,
		// 解析器层
		parser.NewParser,
	)
	return &parser.Parser{}, nil
}