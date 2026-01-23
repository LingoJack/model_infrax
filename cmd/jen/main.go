package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LingoJack/model_infrax/assets"
	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/constant"
	"github.com/LingoJack/model_infrax/internal/generator"
	"github.com/LingoJack/model_infrax/internal/infra/db_infra"
	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/model"
	"github.com/LingoJack/model_infrax/internal/parser/database_parser"
	"github.com/LingoJack/model_infrax/internal/parser/statement_parser"
	"github.com/LingoJack/model_infrax/pkg/tool"
	"github.com/fatih/color"
	"github.com/samber/lo"
	flag "github.com/spf13/pflag"
)

const Version = "v1.2.0"

// initModules 初始化各个模块
// 必须在 conf.InitWithPath 之后调用
func initModules() error {
	generateMode := conf.ValueStr("generate_config.generate_mode")
	logger.Infof("[initModules] 生成模式: %s", generateMode)

	switch generateMode {
	case constant.GenerateModeDatabase:
		// 数据库模式：初始化数据库连接和 database_parser
		if err := db_infra.Init(); err != nil {
			return fmt.Errorf("初始化数据库连接失败: %w", err)
		}
		if err := database_parser.Init(); err != nil {
			return fmt.Errorf("初始化 database_parser 失败: %w", err)
		}
	case constant.GenerateModeStatement:
		// SQL 语句模式：初始化 statement_parser
		if err := statement_parser.Init(); err != nil {
			return fmt.Errorf("初始化 statement_parser 失败: %w", err)
		}
	default:
		return fmt.Errorf("不支持的生成模式: %s，请使用 '%s' 或 '%s'", generateMode, constant.GenerateModeDatabase, constant.GenerateModeStatement)
	}

	logger.Infof("[initModules] 模块初始化完成")
	return nil
}

func main() {
	showVersion := flag.BoolP("version", "v", false, "显示版本号")
	initConfig := flag.Bool("init", false, "在当前目录生成配置文件模板 model_infrax.yml")
	configPath := flag.StringP("config", "c", "", "指定配置文件路径")
	flag.Parse()

	// 处理 --init 命令
	if *initConfig {
		if err := generateConfigTemplate(); err != nil {
			logger.Errorf("[main] 生成配置文件模板失败: %v", err)
			os.Exit(1)
		}
		return
	}

	// 初始化配置（init 命令除外，version 命令也需要加载配置以显示当前生效的配置文件）
	if err := conf.InitWithPath(*configPath); err != nil {
		logger.Errorf("[main] 初始化配置失败: %v", err)
		os.Exit(1)
	}

	// 处理 --version 命令
	if *showVersion {
		logger.ColorPrintf(logger.ColorHiGreen, "═══════════════════════════════════════\n")
		logger.ColorPrintf(logger.ColorHiGreen, "        Model Infrax 代码生成器\n")
		logger.ColorPrintf(logger.ColorHiGreen, "═══════════════════════════════════════\n")
		fmt.Println()

		infos := []struct {
			key   string
			value string
			color color.Attribute
		}{
			{"版本", Version, logger.ColorHiGreen},
			{"github", "https://github.com/LingoJack/model_infrax", logger.ColorHiBlue},
			{"email", "3065225677@qq.com", logger.ColorHiYellow},
			{"作者", "达不溜勾勾", logger.ColorWhite},
			{"配置文件路径(优先级逆序)", tool.JsonifyIndent(conf.DefaultConfigPaths), logger.ColorHiCyan},
		}

		for _, info := range infos {
			logger.ColorPrintf(logger.ColorHiWhite, "%s: ", info.key)
			logger.ColorPrintf(info.color, "%s\n", info.value)
		}
		logger.ColorPrintf(logger.ColorGreen, "● 当前生效配置文件(%s)\n", conf.ValueStr(constant.ActivateConfigPathKey))

		fmt.Println()
		logger.ColorPrintf(logger.ColorHiGreen, "═══════════════════════════════════════\n")

		return
	}

	// 初始化各个模块（必须在配置加载后，且非 version 命令）
	if err := initModules(); err != nil {
		logger.Errorf("[main] 初始化模块失败: %v", err)
		os.Exit(1)
	}

	if err := generate(context.Background()); err != nil {
		logger.Errorf("[main] 生成代码失败: %v", err)
		os.Exit(1)
	}
}

func generate(ctx context.Context) (err error) {
	generateMode := conf.ValueStr("generate_config.generate_mode")
	var schemas []model.Schema
	logger.Infof("[generate] 生成模式: %s", generateMode)

	switch generateMode {
	case constant.GenerateModeDatabase:
		schemas, err = database_parser.Parse(ctx)
		if err != nil {
			logger.Errorf("[generate] 解析数据库失败: %v", err)
			err = fmt.Errorf("解析数据库失败: %w", err)
			return
		}
		schemas = database_parser.Filter(schemas)
	case constant.GenerateModeStatement:
		var statements []string
		statements, err = statement_parser.SqlStatements()
		if err != nil {
			logger.Errorf("[generate] 读取SQL文件失败: %v", err)
			return
		}
		schemas, err = statement_parser.Parse(statements)
		if err != nil {
			logger.Errorf("[generate] 解析语句失败: %v", err)
			return
		}
		logger.Infof("[generate] 过滤前，找到 %d 张表: %v", len(schemas), lo.Map(schemas, func(schema model.Schema, index int) string { return schema.Name }))
		schemas = statement_parser.Filter(schemas)
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

	logger.Infof("[generate] ====== 生成 model 代码 ======")
	err = generator.GenerateModels(schemas)
	if err != nil {
		logger.Errorf("[generate] 生成Model代码失败: %v", err)
		err = fmt.Errorf("生成Model代码失败: %w", err)
		return
	}

	logger.Infof("[generate] ====== 生成 dto 代码 ======")
	err = generator.GenerateDtos(schemas)
	if err != nil {
		logger.Errorf("[generate] 生成DTO代码失败: %v", err)
		err = fmt.Errorf("生成DTO代码失败: %w", err)
		return
	}

	logger.Infof("[generate] ====== 生成 vo 代码 ======")
	err = generator.GenerateVos(schemas)
	if err != nil {
		logger.Errorf("[generate] 生成VO代码失败: %v", err)
		err = fmt.Errorf("生成VO代码失败: %w", err)
		return
	}

	logger.Infof("[generate] ====== 生成 dao 代码 ======")
	err = generator.GenerateDaos(schemas)
	if err != nil {
		err = fmt.Errorf("生成DAO代码失败: %w", err)
		logger.Errorf("[generate] 生成DAO代码失败: %v", err)
		return
	}

	logger.Infof("[generate] ====== 生成 tool 代码 ======")
	err = generator.GenerateTools()
	if err != nil {
		err = fmt.Errorf("生成工具代码失败: %w", err)
		logger.Errorf("[generate] 生成工具代码失败: %v", err)
		return
	}

	return nil
}

// generateConfigTemplate 在当前目录生成配置文件模板
func generateConfigTemplate() error {
	const configFileName = "model_infrax.yml"

	logger.Infof("[generateConfigTemplate] 开始生成配置文件模板...")

	// 检查文件是否已存在
	if _, err := os.Stat(configFileName); err == nil {
		logger.Errorf("[generateConfigTemplate] 配置文件已存在: %s", configFileName)
		return fmt.Errorf("配置文件已存在: %s，请先删除或重命名现有文件", configFileName)
	}

	// 写入文件
	if err := os.WriteFile(configFileName, []byte(assets.ApplicationYml), 0644); err != nil {
		logger.Errorf("[generateConfigTemplate] 写入配置文件失败: %v, 文件路径: %s", err, configFileName)
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	logger.ColorPrintf(logger.ColorHiGreen, "✓ 配置文件模板生成成功: %s\n", configFileName)
	logger.ColorPrintf(logger.ColorHiCyan, "\n使用说明:\n")
	logger.ColorPrintf(logger.ColorWhite, "1. 编辑 %s 文件，配置数据库连接或SQL文件路径\n", configFileName)
	logger.ColorPrintf(logger.ColorWhite, "2. 运行 jen 命令开始生成代码\n")
	logger.ColorPrintf(logger.ColorWhite, "3. 使用 jen -c <path> 可以指定其他配置文件路径\n")

	return nil
}
