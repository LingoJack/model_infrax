package tool

import (
	"fmt"
	"os"
	"path/filepath"
)

func EscapeHomeDir(path string) string {
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		path = filepath.Join(homeDir, path[1:])
	}
	return path
}

func IsValidFilePath(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return info.IsDir()
}

func MustReadText(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// WriteFileWithDir 写入文件，如果目录不存在则自动创建
// 参数:
//   - path: 文件路径
//   - data: 要写入的数据
//   - perm: 文件权限（通常用 0644）
//
// 返回:
//   - error: 写入过程中的错误
//
// 示例:
//
//	err := WriteFileWithDir("/path/to/file.txt", [] byte("hello"), 0644)
func WriteFileWithDir(path string, data []byte, perm os.FileMode) error {
	path = EscapeHomeDir(path)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
	}

	// 写入文件
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("写入文件 %s 失败: %w", path, err)
	}

	return nil
}
