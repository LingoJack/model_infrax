package main

import (
	"context"
	"os"

	"github.com/LingoJack/model_infrax/internal/app"
	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/logger"
	flag "github.com/spf13/pflag"
)

const Version = "v2.1.0"

func main() {
	flag.Usage = printUsage
	showVersion := flag.BoolP("version", "v", false, "显示版本号")
	configPath := flag.StringP("config", "c", "", "指定配置文件路径")
	uiPort := flag.Int("port", 8899, "Web UI 监听端口（jen ui 子命令使用）")
	flag.Parse()

	// 处理 init 子命令（无需加载配置）
	if len(flag.Args()) > 0 && flag.Args()[0] == "init" {
		if err := generateConfigTemplate(); err != nil {
			logger.Errorf("[main] 初始化配置失败: %v", err)
			os.Exit(1)
		}
		return
	}

	// 初始化配置（version/ui/generate 均需要）
	if err := conf.InitWithPath(*configPath); err != nil {
		logger.Errorf("[main] 初始化配置失败: %v", err)
		os.Exit(1)
	}

	// 处理 --version 命令
	if *showVersion {
		printVersion()
		return
	}

	// 处理 ui 子命令（加载配置后启动 Web UI）
	if len(flag.Args()) > 0 && flag.Args()[0] == "ui" {
		if err := runUI(*uiPort); err != nil {
			logger.Errorf("[main] UI 退出: %v", err)
			os.Exit(1)
		}
		return
	}

	// 应用 CLI flag 覆盖（必须在配置加载后、生成前）
	if err := applyFlagOverrides(); err != nil {
		logger.Errorf("[main] 参数错误: %v", err)
		os.Exit(1)
	}

	// 初始化各个模块并生成
	if err := app.Run(context.Background()); err != nil {
		logger.Errorf("[main] 生成代码失败: %v", err)
		os.Exit(1)
	}
}
