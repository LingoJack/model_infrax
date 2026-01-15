package app

import (
	"context"
	"fmt"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/constant"
	"github.com/LingoJack/model_infrax/internal/generator"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/internal/parser/database_parser"
	"github.com/LingoJack/model_infrax/internal/parser/statement_parser"
)

var (
	generateMode = conf.ValueStr("generate_config.generate_mode")
)

func Run(ctx context.Context) (err error) {
	var schemas []model.Schema
	switch generateMode {
	case constant.GenerateModeDatabase:
		schemas, err = database_parser.Parse(ctx)
		if err != nil {
			return
		}
		schemas = database_parser.Filter(schemas)
	case constant.GenerateModeStatement:
		schemas, err = statement_parser.Parse()
		if err != nil {
			return
		}
		schemas = statement_parser.Filter(schemas)
	default:
		err = fmt.Errorf("不支持的生成模式: %s，请使用 '%s' 或 '%s'", generateMode, constant.GenerateModeDatabase, constant.GenerateModeStatement)
		return
	}

	if len(schemas) == 0 {
		err = fmt.Errorf("没有找到任何表结构")
		return nil
	}

	err = generator.GenerateModelOneByOne(schemas)
	if err != nil {
		err = fmt.Errorf("生成Model代码失败: %w", err)
		return
	}

	err = generator.GenerateDtoOneByOne(schemas)
	if err != nil {
		return fmt.Errorf("生成DTO代码失败: %w", err)
	}

	err = generator.GenerateVoOneByOne(schemas)
	if err != nil {
		return fmt.Errorf("生成VO代码失败: %w", err)
	}

	err = generator.GenerateDaoOneByOne(schemas)
	if err != nil {
		return fmt.Errorf("生成DAO代码失败: %w", err)
	}

	err = generator.GenerateAllTools()
	if err != nil {
		return fmt.Errorf("生成Tool代码失败: %w", err)
	}

	return nil
}
