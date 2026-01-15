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

var (
	framework = conf.ValueStr("generate_option.use_framework")
)

const (
	templatePathPrefix = "template/"
	templatePathSuffix = ".tpl"
)

func FrameworkPrefix() string {
	if framework != constant.GenerateOptionFrameworkGorm && framework != constant.GenerateOptionFrameworkIteaGo {
		panic(fmt.Sprintf("不支持的框架: %s", framework))
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
	return filepath.Join(FrameworkPrefix(), "tool")
}

func ToolTemplatePath(templateName string) string {
	return filepath.Join(ToolTemplateDir(), templateName)
}
