package main

import (
	"log"
	"model_infrax/config"
	"model_infrax/generator"
	"model_infrax/parser"
	"model_infrax/tool"
)

type App struct {
	Config    *config.Configger
	Parser    *parser.Parser
	Generator *generator.Generator
}

func NewApp(cfg *config.Configger, p *parser.Parser, g *generator.Generator) *App {
	return &App{
		Config:    cfg,
		Parser:    p,
		Generator: g,
	}
}

// Run 运行应用程序
func (a *App) Run() error {
	log.Println("================ 开始解析数据库 ==================")

	schemas, err := a.Parser.AllTables()
	if err != nil {
		return err
	}

	log.Println("================ 解析结果 ==================")

	log.Println(tool.JsonifyIndent(schemas))

	log.Println("================ 过滤后的表 ==================")

	// 根据配置过滤表
	schemas = a.Parser.FilterTables(schemas)
	log.Println(tool.JsonifyIndent(schemas))

	log.Println("================ 开始生成 Model 代码 ==================")

	// 生成 Model 代码
	if a.Config.GenerateOption.ModelAllInOneFile {
		err = a.Generator.GenerateModel(schemas, a.Config.GenerateOption.ModelAllInOneFileName)
	} else {
		err = a.Generator.GenerateModelOneByOne(schemas)
	}
	if err != nil {
		return err
	}

	log.Println("================ Model 代码生成完成 ==================")

	log.Println("================ 开始生成 DTO 代码 ==================")

	// 生成 DTO 代码
	if a.Config.GenerateOption.ModelAllInOneFile {
		err = a.Generator.GenerateDTO(schemas, a.Config.GenerateOption.ModelAllInOneFileName)
	} else {
		err = a.Generator.GenerateDTOOneByOne(schemas)
	}
	if err != nil {
		return err
	}

	log.Println("================ DTO 代码生成完成 ==================")

	log.Println("================ 开始生成 DAO 代码 ==================")

	// 生成 DAO 代码
	if a.Config.GenerateOption.ModelAllInOneFile {
		err = a.Generator.GenerateDAO(schemas, a.Config.GenerateOption.ModelAllInOneFileName)
	} else {
		err = a.Generator.GenerateDAOOneByOne(schemas)
	}
	if err != nil {
		return err
	}

	log.Println("================ DAO 代码生成完成 ==================")

	log.Println("================ 开始生成 Tool 代码 ==================")

	// 生成 Tool 工具代码
	err = a.Generator.GenerateAllTools()
	if err != nil {
		return err
	}

	log.Println("================ Tool 代码生成完成 ==================")

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
