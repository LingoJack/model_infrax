package webui

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/LingoJack/model_infrax/internal/conf"
)

// filesHandler GET 返回输出目录下的产物文件列表（相对路径）
func filesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	outDir := conf.ValueStr("generate_option.output_path")
	if outDir == "" {
		writeJSON(w, map[string]any{"output_path": "", "files": []string{}})
		return
	}

	var files []string
	_ = filepath.WalkDir(outDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(outDir, path)
		if err != nil {
			return nil
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	sort.Strings(files)
	if files == nil {
		files = []string{}
	}
	writeJSON(w, map[string]any{"output_path": outDir, "files": files})
}

// fileContentHandler GET 返回单个产物文件内容
// 参数 path 为相对输出目录的路径，拒绝目录穿越
func fileContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	rel := r.URL.Query().Get("path")
	outDir := conf.ValueStr("generate_option.output_path")
	if outDir == "" || rel == "" {
		writeJSONError(w, http.StatusBadRequest, "path 为空或输出目录未配置")
		return
	}

	// 防目录穿越：拼接后必须仍落在输出目录内
	clean := filepath.Clean(rel)
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		writeJSONError(w, http.StatusBadRequest, "非法文件路径")
		return
	}
	full := filepath.Join(outDir, clean)
	data, err := os.ReadFile(full)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("读取文件失败: %v", err))
		return
	}
	writeJSON(w, map[string]any{"path": filepath.ToSlash(clean), "content": string(data)})
}
