package webui

import (
	"embed"
)

// dist 为 React 前端构建产物（internal/webui/frontend 执行 npm run build 生成）
//
//go:embed all:dist
var distFS embed.FS
