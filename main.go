package main

import (
	"log"
	"model_infrax/config"
	"model_infrax/parser"
	"model_infrax/tool"
)

type App struct {
	Config *config.Configger
	Parser *parser.Parser
}

func NewApp(cfg *config.Configger, p *parser.Parser) *App {
	return &App{
		Config: cfg,
		Parser: p,
	}
}

// Run 运行应用程序
func (a *App) Run() error {
	schemas, err := a.Parser.AllTables()
	if err != nil {
		return err
	}

	log.Println(tool.JsonifyIndent(schemas))

	log.Println("================ filter ==================")

	schemas = a.Parser.FilterTables(schemas)

	log.Println(tool.JsonifyIndent(schemas))

	return nil
}

func main() {
	configPath := "./assert/application.yml"

	app, err := InitializeApp(configPath)
	if err != nil {
		log.Fatalf("初始化应用失败: %v", err)
	}

	if err = app.Run(); err != nil {
		log.Fatalf("运行应用失败: %v", err)
	}
}
