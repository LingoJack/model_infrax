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

func FrameworkPrefix() string {
	framework := getFramework()
	if framework != constant.GenerateOptionFrameworkGorm && framework != constant.GenerateOptionFrameworkIteaGo {
		panic(fmt.Sprintf("[FrameworkPrefix] 不支持的框架: %s，请使用 '%s' 或 '%s'", 
			framework, 
			constant.GenerateOptionFrameworkGorm, 
			constant.GenerateOptionFrameworkIteaGo))
	}
	return filepath.Join(templatePathPrefix, framework)
}

func VoTemplatePath() string {
	return filepath.Join(FrameworkPrefix(), "vo"+templatePathSuffix)
}

func DaoTemplatePath() string {
	return filepath.Join(FrameworkPrefix(), "dao"+templatePathSuffix)
}

func DtoTemplatePath() string {
	return filepath.Join(FrameworkPrefix(), "dto"+templatePathSuffix)
}

func PoTemplatePath() string {
	return filepath.Join(FrameworkPrefix(), "po"+templatePathSuffix)
}

func ToolTemplateDir() string {
	return filepath.Join(templatePathPrefix, "tool")
}

func ToolTemplatePath(templateName string) string {
	return filepath.Join(ToolTemplateDir(), templateName)
}
