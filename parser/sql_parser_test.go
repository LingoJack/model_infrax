package parser

import (
	"fmt"
	"os"
)

// loadSQL 加载SQL文件内容
func loadSQL() string {
	byts, err := os.ReadFile("/Users/jacklingo/dev/model_infrax/assert/database.sql")
	if err != nil {
		panic(fmt.Sprintf("读取SQL文件失败: %v", err))
	}
	return string(byts)
}
