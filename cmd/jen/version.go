package main

import (
	"fmt"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/constant"
	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/tool"
	"github.com/fatih/color"
)

// printUsage 打印帮助信息
func printUsage() {
	logger.ColorPrintf(logger.ColorHiGreen, "Model Infrax %s — Go 代码生成器 CLI\n\n", Version)
	logger.ColorPrintf(logger.ColorHiCyan, "用法:\n")
	logger.ColorPrintf(logger.ColorWhite, "  jen init          初始化 .model_infrax 配置目录（已有配置会询问是否覆盖）\n")
	logger.ColorPrintf(logger.ColorWhite, "  jen ui            启动本地 Web UI 配置界面（默认 http://127.0.0.1:8899）\n")
	logger.ColorPrintf(logger.ColorWhite, "  jen               加载 .model_infrax/config.yml 生成代码\n")
	logger.ColorPrintf(logger.ColorWhite, "  jen -c <path>     指定配置文件路径生成代码\n")
	logger.ColorPrintf(logger.ColorWhite, "  jen -v            显示版本信息\n")
	logger.ColorPrintf(logger.ColorWhite, "  jen -h            显示帮助信息\n")
	fmt.Println()
	logger.ColorPrintf(logger.ColorHiCyan, "生成项覆盖 flag（可选，优先级高于配置文件）:\n")
	logger.ColorPrintf(logger.ColorWhite, "  --gorm-tag <full|standard|minimal>    gorm tag 风格\n")
	logger.ColorPrintf(logger.ColorWhite, "  --comment <full|brief|none>           注释风格\n")
	logger.ColorPrintf(logger.ColorWhite, "  --methods <a,b,...>                   仅生成指定方法（如 SelectList,Insert,SelectByPk）\n")
	logger.ColorPrintf(logger.ColorWhite, "  --without-methods <a,b,...>           排除指定方法，其余全部生成\n")
	logger.ColorPrintf(logger.ColorWhite, "  示例: jen --gorm-tag=minimal --comment=brief --methods=Insert,SelectByPk\n")
	fmt.Println()
	logger.ColorPrintf(logger.ColorHiCyan, "快速上手:\n")
	logger.ColorPrintf(logger.ColorWhite, "  1. jen init                            初始化配置\n")
	logger.ColorPrintf(logger.ColorWhite, "  2. 编辑 .model_infrax/schema.sql       编写建表语句\n")
	logger.ColorPrintf(logger.ColorWhite, "  3. jen                                 生成代码到 target/jen/\n")
	fmt.Println()
	logger.ColorPrintf(logger.ColorHiCyan, "更多信息:\n")
	logger.ColorPrintf(logger.ColorWhite, "  https://github.com/LingoJack/model_infrax\n")
}

// printVersion 打印版本信息（必须在 conf.InitWithPath 之后调用）
func printVersion() {
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
}
