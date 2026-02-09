package tool

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

func ChineseToEnglish(chinese string) (res string) {
	// 默认
	a := pinyin.NewArgs()
	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}

	py := pinyin.Pinyin(chinese, a)
	for _, item := range py {
		res += strings.Join(item, "_")
	}

	return res
}
