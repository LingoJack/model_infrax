package tool

import (
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
