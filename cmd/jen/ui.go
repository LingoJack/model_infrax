package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/LingoJack/model_infrax/internal/webui"
)

// runUI 启动 Web UI 配置界面（配置已在 main 中加载）
func runUI(port int) error {
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	openBrowser(url)
	return webui.Start(port)
}

// openBrowser 尝试用系统默认浏览器打开 URL（失败仅忽略，用户可手动访问）
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
