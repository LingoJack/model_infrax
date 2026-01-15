package main

import (
	"context"
	"fmt"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/constant"
	"github.com/LingoJack/model_infrax/internal/generator"
	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/internal/parser/database_parser"
	"github.com/LingoJack/model_infrax/internal/parser/statement_parser"
	flag "github.com/spf13/pflag"
)

var (
	generateMode = conf.ValueStr("generate_config.generate_mode")
)

const Version = "1.2.0"

func main() {
	showVersion := flag.BoolP("version", "v", false, "显示版本号")
	flag.Parse()

	if *showVersion {
		logger.Infof("版本号: %s", Version)
		return
	}

	if err := GenerateCode(context.Background()); err != nil {
		logger.Infof("生成代码失败: %v", err)
		return
	}
}

func GenerateCode(ctx context.Context) (err error) {
	if err = process(ctx); err != nil {
		logger.Errorf("生成代码失败: %v", err)
		return
	}
	return
}

func process(ctx context.Context) (err error) {
	var schemas []model.Schema
	logger.Infof("生成模式: %s", generateMode)
	switch generateMode {
	case constant.GenerateModeDatabase:
		schemas, err = database_parser.Parse(ctx)
		if err != nil {
			logger.Errorf("解析数据库失败: %v", err)
			return
		}
		schemas = database_parser.Filter(schemas)
	case constant.GenerateModeStatement:
		schemas, err = statement_parser.Parse()
		if err != nil {
			logger.Errorf("解析语句失败: %v", err)
			return
		}
		schemas = statement_parser.Filter(schemas)
	default:
		err = fmt.Errorf("不支持的生成模式: %s，请使用 '%s' 或 '%s'", generateMode, constant.GenerateModeDatabase, constant.GenerateModeStatement)
		logger.Errorf("生成代码失败: %v", err)
		return
	}

	logger.Infof("找到 %d 张表", len(schemas))

	if len(schemas) == 0 {
		logger.Infof("未找到任何表")
		return nil
	}

	logger.Infof("生成 model 代码")
	err = generator.GenerateModelOneByOne(schemas)
	if err != nil {
		logger.Errorf("生成Model代码失败: %v", err)
		err = fmt.Errorf("生成Model代码失败: %w", err)
		return
	}

	logger.Infof("生成 dto 代码")
	err = generator.GenerateDtoOneByOne(schemas)
	if err != nil {
		logger.Errorf("生成DTO代码失败: %v", err)
		err = fmt.Errorf("生成DTO代码失败: %w", err)
		return
	}

	logger.Infof("生成 vo 代码")
	err = generator.GenerateVoOneByOne(schemas)
	if err != nil {
		logger.Errorf("生成VO代码失败: %v", err)
		err = fmt.Errorf("生成VO代码失败: %w", err)
		return
	}

	logger.Infof("生成 dao 代码")
	err = generator.GenerateDaoOneByOne(schemas)
	if err != nil {
		err = fmt.Errorf("生成DAO代码失败: %w", err)
		logger.Errorf("生成DAO代码失败: %v", err)
		return
	}

	logger.Infof("生成 tool 代码")
	err = generator.GenerateAllTools()
	if err != nil {
		err = fmt.Errorf("生成工具代码失败: %w", err)
		logger.Errorf("生成工具代码失败: %v", err)
		return
	}

	return nil
}
