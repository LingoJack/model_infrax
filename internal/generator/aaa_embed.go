package generator

import (
	"embed"
	"fmt"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/constant"
)

//go:embed template/*
var fs embed.FS

const templatePathPrefix = "template/"

func FrameworkPrefix() string {
	framework := conf.ValueStr("generate_option.use_framework")
	if framework != constant.GenerateOptionFrameworkGorm && framework != constant.GenerateOptionFrameworkIteaGo {
		panic(fmt.Sprintf("不支持的框架: %s", framework))
	}
	return templatePathPrefix + framework + "/"
}
