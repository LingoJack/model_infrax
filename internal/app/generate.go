package app

import (
	"context"
	"fmt"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/constant"
	"github.com/LingoJack/model_infrax/internal/generator"
	"github.com/LingoJack/model_infrax/internal/infra"
	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/internal/parser"
	"github.com/samber/lo"
)

// initModules 初始化各个模块
// 必须在 conf.InitWithPath 之后调用
func initModules() error {
	generateMode := conf.ValueStr("generate_config.generate_mode")
	logger.Infof("[initModules] 生成模式: %s", generateMode)

	switch generateMode {
	case constant.GenerateModeDatabase:
		// 数据库模式：初始化数据库连接和 database 解析器
		if err := infra.InitDB(); err != nil {
			return fmt.Errorf("初始化数据库连接失败: %w", err)
		}
		if err := parser.InitDatabase(); err != nil {
			return fmt.Errorf("初始化 database 解析器失败: %w", err)
		}
	case constant.GenerateModeStatement:
		// SQL 语句模式：初始化 statement 解析器
		if err := parser.InitStatement(); err != nil {
			return fmt.Errorf("初始化 statement 解析器失败: %w", err)
		}
	default:
		return fmt.Errorf("不支持的生成模式: %s，请使用 '%s' 或 '%s'", generateMode, constant.GenerateModeDatabase, constant.GenerateModeStatement)
	}

	logger.Infof("[initModules] 模块初始化完成")
	return nil
}

// Run 执行完整生成流程：初始化模块 → 解析表结构 → 依次生成 PO/DTO/VO/DAO/Tool
// 必须在 conf.InitWithPath（及可选的 flag 覆盖）之后调用
func Run(ctx context.Context) error {
	if err := initModules(); err != nil {
		return err
	}
	return generate(ctx)
}

// generate 主流程：解析表结构后依次生成 PO/DTO/VO/DAO/Tool 代码
func generate(ctx context.Context) (err error) {
	generateMode := conf.ValueStr("generate_config.generate_mode")
	var schemas []model.Schema
	logger.Infof("[generate] 生成模式: %s", generateMode)

	switch generateMode {
	case constant.GenerateModeDatabase:
		schemas, err = parser.ParseDatabase(ctx)
		if err != nil {
			logger.Errorf("[generate] 解析数据库失败: %v", err)
			err = fmt.Errorf("解析数据库失败: %w", err)
			return
		}
		schemas = parser.FilterTables(schemas)
	case constant.GenerateModeStatement:
		var statements []string
		statements, err = parser.SqlStatements()
		if err != nil {
			logger.Errorf("[generate] 读取SQL文件失败: %v", err)
			return
		}
		schemas, err = parser.ParseStatements(statements)
		if err != nil {
			logger.Errorf("[generate] 解析语句失败: %v", err)
			return
		}
		logger.Infof("[generate] 过滤前，找到 %d 张表: %v", len(schemas), lo.Map(schemas, func(schema model.Schema, index int) string { return schema.Name }))
		schemas = parser.FilterTables(schemas)
	default:
		err = fmt.Errorf("不支持的生成模式: %s，请使用 '%s' 或 '%s'", generateMode, constant.GenerateModeDatabase, constant.GenerateModeStatement)
		logger.Errorf("[generate] 生成代码失败: %v", err)
		return
	}

	logger.Infof("[generate] 过滤后，找到 %d 张表: %v", len(schemas), lo.Map(schemas, func(schema model.Schema, index int) string { return schema.Name }))

	if len(schemas) == 0 {
		logger.Infof("[generate] 未找到任何表")
		err = fmt.Errorf("未找到任何表")
		return
	}

	generators := []struct {
		label string
		run   func([]model.Schema) error
	}{
		{"po", generator.GeneratePos},
		{"dto", generator.GenerateDtos},
		{"vo", generator.GenerateVos},
		{"dao", generator.GenerateDaos},
	}

	for _, g := range generators {
		logger.Infof("[generate] ====== 生成 %s 代码 ======", g.label)
		if err = g.run(schemas); err != nil {
			logger.Errorf("[generate] 生成 %s 代码失败: %v", g.label, err)
			err = fmt.Errorf("生成 %s 代码失败: %w", g.label, err)
			return
		}
	}

	logger.Infof("[generate] ====== 生成 tool 代码 ======")
	if err = generator.GenerateTools(); err != nil {
		logger.Errorf("[generate] 生成工具代码失败: %v", err)
		err = fmt.Errorf("生成工具代码失败: %w", err)
		return
	}

	return nil
}
