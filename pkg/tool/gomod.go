package tool

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadModulePath 从当前工作目录的 go.mod 文件中读取模块路径
// 返回:
//   - string: 模块路径，如 "github.com/foo/myapp"
//   - error: 读取失败时返回错误
func ReadModulePath() (string, error) {
	f, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("打开 go.mod 失败: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if mod, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(mod), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取 go.mod 失败: %w", err)
	}
	return "", fmt.Errorf("go.mod 中未找到 module 声明")
}
