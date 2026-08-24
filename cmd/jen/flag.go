package main

import (
	"fmt"
	"strings"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/generator"
	flag "github.com/spf13/pflag"
)

// CLI 覆盖 flag 定义（与 generate_option 配置键一一对应）
var (
	gormTagFlag = flag.String("gorm-tag", "", "gorm tag 风格: full|standard|minimal（覆盖配置文件）")
	commentFlag = flag.String("comment", "", "注释风格: full|brief|none（覆盖配置文件）")
	methodsFlag = flag.StringSlice("methods", nil, "仅生成指定方法，逗号分隔（覆盖配置文件），如: SelectList,Insert,SelectByPk")
	withoutFlag = flag.StringSlice("without-methods", nil, "排除指定方法，逗号分隔，其余全部生成；不可与 --methods 同用")
)

var (
	validGormTagStyles = []string{"full", "standard", "minimal"}
	validCommentStyles = []string{"full", "brief", "none"}
)

// applyFlagOverrides 将显式给出的 CLI flag 写入 conf 覆盖层
// 必须在 conf.InitWithPath 之后调用
func applyFlagOverrides() error {
	if *gormTagFlag != "" {
		if !contains(validGormTagStyles, *gormTagFlag) {
			return fmt.Errorf("--gorm-tag 非法值: %s，可选: %s", *gormTagFlag, strings.Join(validGormTagStyles, "|"))
		}
		conf.SetOverride("generate_option.gorm_tag_style", *gormTagFlag)
	}

	if *commentFlag != "" {
		if !contains(validCommentStyles, *commentFlag) {
			return fmt.Errorf("--comment 非法值: %s，可选: %s", *commentFlag, strings.Join(validCommentStyles, "|"))
		}
		conf.SetOverride("generate_option.comment_style", *commentFlag)
	}

	if len(*methodsFlag) > 0 && len(*withoutFlag) > 0 {
		return fmt.Errorf("--methods 与 --without-methods 不可同时使用")
	}

	if len(*methodsFlag) > 0 {
		// 校验方法名合法性（防止拼写错误静默失效）
		for _, m := range *methodsFlag {
			if err := validateMethodName(m); err != nil {
				return err
			}
		}
		conf.SetOverride("generate_option.dao_methods", normalizeMethods(*methodsFlag))
	}

	if len(*withoutFlag) > 0 {
		exclude := map[string]bool{}
		for _, m := range *withoutFlag {
			if err := validateMethodName(m); err != nil {
				return err
			}
			exclude[strings.TrimSpace(m)] = true
		}
		var kept []string
		for _, meta := range generator.AllDaoMethodIDs() {
			if !exclude[meta] {
				kept = append(kept, meta)
			}
		}
		conf.SetOverride("generate_option.dao_methods", kept)
	}

	return nil
}

// validateMethodName 校验方法名是否在支持清单中
func validateMethodName(name string) error {
	name = strings.TrimSpace(name)
	for _, id := range generator.AllDaoMethodIDs() {
		if id == name {
			return nil
		}
	}
	return fmt.Errorf("不支持的 DAO 方法名: %s，可用值见 --help 或 jen ui", name)
}

// normalizeMethods 去空白后聚合方法名列表
func normalizeMethods(methods []string) []string {
	res := make([]string, 0, len(methods))
	for _, m := range methods {
		m = strings.TrimSpace(m)
		if m != "" {
			res = append(res, m)
		}
	}
	return res
}

func contains(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}
