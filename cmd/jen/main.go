package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/constant"
	"github.com/LingoJack/model_infrax/internal/infra"
	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/parser"
	flag "github.com/spf13/pflag"
)

const Version = "v2.0.2"

func main() {
	flag.Usage = printUsage
	showVersion := flag.BoolP("version", "v", false, "显示版本号")
	configPath := flag.StringP("config", "c", "", "指定配置文件路径")
	flag.Parse()

	// 处理 jen init 子命令
	if len(flag.Args()) > 0 && flag.Args()[0] == "init" {
		if err := generateConfigTemplate(); err != nil {
			logger.Errorf("[main] 初始化配置失败: %v", err)
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
		printVersion()
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
