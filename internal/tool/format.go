package tool

import (
	"go/format"
	"log"
)

func FormatGoCode(code string) string {
	formatted, err := format.Source([]byte(code))
	if err != nil {
		log.Printf("格式化 Go 代码失败: %v", err)
		return code
	}
	return string(formatted)
}
