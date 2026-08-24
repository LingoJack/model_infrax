package generator

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/constant"
)

//go:embed template/*
var fs embed.FS

const (
	templatePathPrefix = "template/"
	templatePathSuffix = ".tpl"
)

// getFramework 获取配置的框架类型
// 必须在 conf.InitWithPath 之后调用
func getFramework() string {
	return conf.ValueStr("generate_option.use_framework")
}

// frameworkPrefix 根据框架配置返回模板目录前缀（template/itea-go 或 template/gorm）
func frameworkPrefix() string {
	framework := getFramework()
	if framework != constant.GenerateOptionFrameworkGorm && framework != constant.GenerateOptionFrameworkIteaGo {
		panic(fmt.Sprintf("[frameworkPrefix] 不支持的框架: %s，请使用 '%s' 或 '%s'",
			framework,
			constant.GenerateOptionFrameworkGorm,
			constant.GenerateOptionFrameworkIteaGo))
	}
	return filepath.Join(templatePathPrefix, framework)
}

func voTemplatePath() string {
	return filepath.Join(frameworkPrefix(), "vo"+templatePathSuffix)
}

func daoTemplatePath() string {
	return filepath.Join(frameworkPrefix(), "dao"+templatePathSuffix)
}

func dtoTemplatePath() string {
	return filepath.Join(frameworkPrefix(), "dto"+templatePathSuffix)
}

func poTemplatePath() string {
	return filepath.Join(frameworkPrefix(), "po"+templatePathSuffix)
}

func toolTemplateDir() string {
	return filepath.Join(templatePathPrefix, "tool")
}

func toolTemplatePath(templateName string) string {
	return filepath.Join(toolTemplateDir(), templateName)
}
