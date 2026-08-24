package webui

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/LingoJack/model_infrax/internal/app"
	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/generator"
	"github.com/LingoJack/model_infrax/internal/logger"
)

// configField 描述一个可编辑的配置键
type configField struct {
	Key   string   `json:"key"`
	Label string   `json:"label"`
	Type  string   `json:"type"`  // text | password | number | bool | select | list
	Opts  []string `json:"opts"`  // select 的选项
	Group string   `json:"group"` // 所属分区
	Hint  string   `json:"hint"`  // 补充说明
}

// configFields 全部可编辑配置键（按展示顺序）
var configFields = []configField{
	{Key: "generate_config.generate_mode", Label: "生成模式", Type: "select", Opts: []string{"statement", "database"}, Group: "生成模式", Hint: "statement=解析 SQL 文件, database=连库解析"},
	{Key: "generate_config.sql_file_path", Label: "SQL 文件路径", Type: "text", Group: "生成模式"},
	{Key: "generate_config.host", Label: "数据库 Host", Type: "text", Group: "数据库连接"},
	{Key: "generate_config.port", Label: "数据库端口", Type: "number", Group: "数据库连接"},
	{Key: "generate_config.username", Label: "数据库用户名", Type: "text", Group: "数据库连接"},
	{Key: "generate_config.password", Label: "数据库密码", Type: "password", Group: "数据库连接"},
	{Key: "generate_config.database_name", Label: "数据库名", Type: "text", Group: "数据库连接"},
	{Key: "generate_config.all_tables", Label: "处理全部表", Type: "bool", Group: "表范围", Hint: "关闭后仅处理下方表名列表"},
	{Key: "generate_config.table_names", Label: "表名列表", Type: "list", Group: "表范围"},
	{Key: "generate_option.use_framework", Label: "使用框架", Type: "select", Opts: []string{"itea-go", "gorm"}, Group: "输出与框架"},
	{Key: "generate_option.output_path", Label: "输出路径", Type: "text", Group: "输出与框架"},
	{Key: "generate_option.gorm_tag_style", Label: "gorm tag 风格", Type: "select", Opts: []string{"full", "standard", "minimal"}, Group: "输出与框架", Hint: "full=全量属性 / standard=column+type+主键 / minimal=仅 column+主键"},
	{Key: "generate_option.comment_style", Label: "注释风格", Type: "select", Opts: []string{"full", "brief", "none"}, Group: "输出与框架", Hint: "full=详细 godoc / brief=一句话 / none=无注释"},
	{Key: "generate_option.package.po", Label: "PO 包路径", Type: "text", Group: "包路径"},
	{Key: "generate_option.package.dto", Label: "DTO 包路径", Type: "text", Group: "包路径"},
	{Key: "generate_option.package.vo", Label: "VO 包路径", Type: "text", Group: "包路径"},
	{Key: "generate_option.package.dao", Label: "DAO 包路径", Type: "text", Group: "包路径"},
	{Key: "generate_option.package.tool", Label: "Tool 包路径", Type: "text", Group: "包路径"},
}

// Start 启动 Web UI HTTP 服务（阻塞运行）
func Start(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", pageHandler)
	mux.HandleFunc("/api/config", configHandler)
	mux.HandleFunc("/api/generate", generateHandler)
	mux.HandleFunc("/api/files", filesHandler)
	mux.HandleFunc("/api/file", fileContentHandler)

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听 %s 失败: %w", addr, err)
	}
	url := fmt.Sprintf("http://%s", ln.Addr().String())
	logger.ColorPrintf(logger.ColorHiGreen, "jen Web UI 已启动: %s\n", url)
	logger.ColorPrintf(logger.ColorWhite, "按 Ctrl+C 停止服务\n")

	srv := &http.Server{Handler: mux}
	return srv.Serve(ln)
}

// pageHandler 服务内嵌的 React 前端静态资源（dist 目录）
// 处理 "/"（index.html）与 "/assets/*"（js/css），未命中路径由 FileServer 返回 404
// 禁缓存，防止浏览器保留旧版本页面
func pageHandler(w http.ResponseWriter, r *http.Request) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		http.Error(w, "前端资源缺失", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	http.FileServer(http.FS(sub)).ServeHTTP(w, r)
}

// configHandler GET 返回配置值+元数据；POST 保存（可选同时触发生成）
func configHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, buildConfigSnapshot())

	case http.MethodPost:
		var req struct {
			Values map[string]any `json:"values"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "请求体解析失败: "+err.Error())
			return
		}
		if err := saveConfig(req.Values); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, map[string]any{"ok": true})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// generateHandler POST 触发一次代码生成（使用当前已保存配置）
func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	logger.ColorPrintf(logger.ColorHiCyan, "\n[ui] 开始生成...\n")
	if err := app.Run(truncatedCtx()); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]any{"ok": true})
}

// truncatedCtx 生成流程上下文（限制时长防止 UI 请求无限挂起）
func truncatedCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	_ = cancel // 请求结束由 server 生命周期管理，无需显式取消
	return ctx
}

// buildConfigSnapshot 汇总当前配置值与页面元数据
func buildConfigSnapshot() map[string]any {
	values := map[string]any{}
	for _, f := range configFields {
		switch f.Type {
		case "bool":
			values[f.Key] = conf.ValueBool(f.Key)
		case "number":
			values[f.Key] = conf.ValueInt(f.Key)
		case "list":
			v, _ := conf.ValueStrSlice(f.Key)
			if v == nil {
				v = []string{}
			}
			values[f.Key] = v
		default:
			values[f.Key] = conf.ValueStr(f.Key)
		}
	}

	return map[string]any{
		"config_path": conf.ActivateConfigPath(),
		"values":      values,
		"fields":      configFields,
		"dao_methods": generator.AllDaoMethods,
	}
}

// saveConfig 校验并写回配置文件，然后重载
func saveConfig(values map[string]any) error {
	if err := validateValues(values); err != nil {
		return err
	}

	path := conf.ActivateConfigPath()
	if path == "" {
		return fmt.Errorf("配置文件未加载")
	}

	normalized := normalizeValues(values)
	if err := PatchYAMLFile(path, normalized); err != nil {
		return fmt.Errorf("写回配置文件失败: %w", err)
	}

	// 重新加载（沿用当前生效的配置文件路径，UI 会话内后续生成使用新值）
	if err := conf.InitWithPath(path); err != nil {
		return fmt.Errorf("重新加载配置失败: %w", err)
	}
	logger.ColorPrintf(logger.ColorHiGreen, "[ui] 配置已保存并生效: %s\n", path)
	return nil
}

// validateValues 校验枚举类取值
func validateValues(values map[string]any) error {
	for key, val := range values {
		switch key {
		case "generate_config.generate_mode":
			if s, ok := val.(string); ok && s != "statement" && s != "database" {
				return fmt.Errorf("%s 非法值: %s", key, s)
			}
		case "generate_option.use_framework":
			if s, ok := val.(string); ok && s != "itea-go" && s != "gorm" {
				return fmt.Errorf("%s 非法值: %s", key, s)
			}
		case "generate_option.gorm_tag_style":
			if s, ok := val.(string); ok && !strIn(s, generator.ValidGormTagStyles) {
				return fmt.Errorf("%s 非法值: %s，可选: %v", key, s, generator.ValidGormTagStyles)
			}
		case "generate_option.comment_style":
			if s, ok := val.(string); ok && !strIn(s, generator.ValidCommentStyles) {
				return fmt.Errorf("%s 非法值: %s，可选: %v", key, s, generator.ValidCommentStyles)
			}
		case "generate_option.dao_methods":
			list, err := toStringSlice(val)
			if err != nil {
				return fmt.Errorf("%s 类型错误: %w", key, err)
			}
			valid := map[string]bool{}
			for _, m := range generator.AllDaoMethods {
				valid[m.ID] = true
			}
			for _, m := range list {
				if !valid[m] {
					return fmt.Errorf("%s 含非法方法名: %s", key, m)
				}
			}
		}
	}
	return nil
}

// normalizeValues 将前端 JSON 值规整为 PatchYAMLFile 支持的类型
func normalizeValues(values map[string]any) map[string]any {
	out := map[string]any{}
	for key, val := range values {
		switch key {
		case "generate_option.dao_methods":
			list, _ := toStringSlice(val)
			out[key] = list
		case "generate_config.port":
			if f, ok := val.(float64); ok {
				out[key] = int(f)
			} else {
				out[key] = val
			}
		default:
			out[key] = val
		}
	}
	return out
}

// toStringSlice 兼容 JSON 数组与逗号分隔字符串
func toStringSlice(val any) ([]string, error) {
	switch v := val.(type) {
	case []string:
		return v, nil
	case []any:
		res := make([]string, 0, len(v))
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("列表元素不是字符串: %v", item)
			}
			res = append(res, s)
		}
		return res, nil
	case string:
		var res []string
		for _, item := range strings.Split(v, ",") {
			if item = strings.TrimSpace(item); item != "" {
				res = append(res, item)
			}
		}
		return res, nil
	default:
		return nil, fmt.Errorf("期望数组, 实际 %T", val)
	}
}

func strIn(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": false, "error": msg})
}
