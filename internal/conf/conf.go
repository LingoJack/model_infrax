package conf

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/tool"
	"gopkg.in/yaml.v3"
)

var (
	config     *Config
	configLock sync.RWMutex
	loadOnce   sync.Once
)

func init() {
	var defaultConfigPaths = []string{
		"./application.yml",
		"./assets/application.yml",
		"/Applications/model_infrax/application.yml",
		"/Applications/model_infrax/assets/application.yml",
	}

	path := ""
	for _, p := range defaultConfigPaths {
		if tool.IsValidFilePath(p) {
			path = p
			break
		}
	}

	logger.Infof("当前生效配置文件(%s)", path)
	if len(path) == 0 {
		logger.Errorf("未找到有效的配置文件")
		return
	}

	err := Load(path)
	if err != nil {
		logger.Errorf("加载配置文件失败: %v", err)
		return
	}

	config.data["config_path"] = path

	logger.Infof("当前生效配置文件(%s), 读取到内容:\n%s", path, tool.JsonifyIndent(config.data))
}

// Load 加载配置文件（只能加载一次）
// 使用 sync.Once 确保配置只加载一次，即使多次调用也只会执行一次
//
// 参数:
//   - file: 配置文件路径，支持相对路径和绝对路径
//
// 返回:
//   - error: 加载失败时返回错误，成功返回 nil
//
// 使用示例:
//
//	if err := tool.Load("config.yaml"); err != nil {
//	    log.Fatalf("[Load]  加载配置失败: %v", err)
//	}
//
// 注意事项:
//   - 该方法只能成功加载一次，重复调用会被忽略（不会返回错误）
//   - 如果需要重新加载配置，请使用 Reload 方法
//   - 建议在程序启动时调用，确保配置在使用前已加载
func Load(file string) (err error) {
	loadOnce.Do(func() {
		if config != nil {
			err = fmt.Errorf("配置已加载")
			return
		}
		var c *Config
		c, err = NewConfig(file)
		if err != nil {
			return
		}
		configLock.Lock()
		config = c
		configLock.Unlock()
		return
	})
	return
}

// Reload 重新加载配置文件（用于测试或热加载场景）
// 与 Load 不同，Reload 可以多次调用，每次都会重新读取配置文件
//
// 参数:
//   - file: 配置文件路径，支持相对路径和绝对路径
//
// 返回:
//   - error: 加载失败时返回错误，成功返回 nil
//
// 使用示例:
//
//	// 测试场景下切换配置
//	if err := tool.Reload("config_test.yaml"); err != nil {
//	    log.Printf("[Reload]  重新加载配置失败: %v", err)
//	}
//
//	// 热加载场景（监听配置文件变化）
//	watcher.OnChange(func() {
//	    if err := tool.Reload("config.yaml"); err != nil {
//	        log.Printf("[Reload]  热加载配置失败: %v", err)
//	    }
//	})
//
// 注意事项:
//   - 该方法会完全替换当前配置，不是增量更新
//   - 重新加载期间会加写锁，可能会短暂阻塞其他读取操作
//   - 建议在生产环境谨慎使用，避免配置不一致导致的问题
func Reload(file string) error {
	c, err := NewConfig(file)
	if err != nil {
		logger.Errorf("[Reload]  重新加载配置失败: %v", err)
		return err
	}
	configLock.Lock()
	config = c
	configLock.Unlock()
	return nil
}

// Value 通过点分隔的键获取配置值
// 支持嵌套键访问，返回原始类型的值和是否存在的标志
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问，如 "app.database.host"
//
// 返回:
//   - value: 配置值，保持原始类型（string, int, bool, map, slice 等）
//   - exists: 配置键是否存在，true 表示存在，false 表示不存在
//
// 使用示例:
//
//	// 基本使用
//	if val, exists := tool.Value("app.port"); exists {
//	    fmt.Printf("端口配置: %v\n", val)
//	} else {
//	    fmt.Println("端口配置不存在")
//	}
//
//	// 嵌套访问
//	if val, exists := tool.Value("database.mysql.host"); exists {
//	    fmt.Printf("数据库主机: %v\n", val)
//	}
//
//	// 区分"不存在"和"值为 nil"
//	val, exists := tool.Value("optional.feature")
//	if !exists {
//	    fmt.Println("配置项不存在")
//	} else if val == nil {
//	    fmt.Println("配置项存在但值为 nil")
//	}
//
// 注意事项:
//   - 返回的 value 是 interface{} 类型，需要根据实际情况进行类型断言
//   - 建议使用类型安全的方法如 ValueStr、ValueInt 等
//   - 该方法是线程安全的，使用读锁保护
func Value(key string) (value any, exists bool) {
	configLock.RLock()
	defer configLock.RUnlock()

	if config == nil {
		return nil, false
	}
	return config.get(key)
}

// ValueWithDefault 获取配置值，如果不存在则返回默认值
// 这是一个便捷方法，避免了手动检查配置是否存在
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//   - defaultValue: 默认值，当配置不存在时返回
//
// 返回:
//   - any: 配置值或默认值
//
// 使用示例:
//
//	// 获取端口，不存在时使用默认值 8080
//	port := tool.ValueWithDefault("app.port", 8080)
//
//	// 获取超时时间，不存在时使用默认值 30 秒
//	timeout := tool.ValueWithDefault("app.timeout", 30)
//
//	// 获取配置对象，不存在时使用默认配置
//	defaultConfig := map[string] interface{}{"enabled": false}
//	config := tool.ValueWithDefault("feature.config", defaultConfig)
//
// 注意事项:
//   - 返回值类型是 interface{}，需要进行类型断言
//   - 建议使用类型安全的方法如 ValueStrWithDefault、ValueIntWithDefault 等
//   - defaultValue 的类型应该与期望的配置值类型一致
func ValueWithDefault(key string, defaultValue any) any {
	if val, exists := Value(key); exists {
		return val
	}
	return defaultValue
}

// ValueStr 通过点分隔的键获取字符串类型的配置值
// 自动进行类型转换，如果配置不存在或为 nil，返回空字符串
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//
// 返回:
//   - string: 字符串类型的配置值，不存在时返回空字符串 ""
//
// 使用示例:
//
//	// 获取产品 ID
//	productID := tool.ValueStr("app.product_id")
//	fmt.Printf("产品ID: %s\n", productID)
//
//	// 获取数据库主机
//	dbHost := tool.ValueStr("database.host")
//
//	// 获取嵌套配置
//	apiKey := tool.ValueStr("services.payment.api_key")
//
// 注意事项:
//   - 如果配置值不是字符串类型，会使用 fmt.Sprintf("%v", val) 进行转换
//   - 配置不存在时返回空字符串，无法区分"不存在"和"值为空字符串"
//   - 如需区分这两种情况，请使用 Value 方法
func ValueStr(key string) (value string) {
	val, exists := Value(key)
	if !exists || val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

// ValueStrWithDefault 获取字符串配置值，如果不存在则返回默认值
// 这是类型安全的便捷方法，推荐使用
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//   - defaultValue: 默认值，当配置不存在或为 nil 时返回
//
// 返回:
//   - string: 字符串类型的配置值或默认值
//
// 使用示例:
//
//	// 获取环境配置，默认为 "development"
//	env := tool.ValueStrWithDefault("app.env", "development")
//
//	// 获取日志级别，默认为 "info"
//	logLevel := tool.ValueStrWithDefault("log.level", "info")
//
//	// 获取 API 地址，默认为本地地址
//	apiURL := tool.ValueStrWithDefault("api.url", "http://localhost:8080")
//
// 注意事项:
//   - 如果配置值不是字符串类型，会使用 fmt.Sprintf("%v", val) 进行转换
//   - 推荐在生产环境使用此方法，确保总是有合理的默认值
func ValueStrWithDefault(key string, defaultValue string) string {
	val, exists := Value(key)
	if !exists || val == nil {
		return defaultValue
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

// ValueInt 通过点分隔的键获取整数类型的配置值
// 支持多种数值类型的自动转换（int, int64, float64）
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//
// 返回:
//   - int: 整数类型的配置值，不存在或类型不匹配时返回 0
//
// 使用示例:
//
//	// 获取服务端口
//	port := tool.ValueInt("app.port")
//	fmt.Printf("服务端口: %d\n", port)
//
//	// 获取最大连接数
//	maxConn := tool.ValueInt("database.max_connections")
//
//	// 获取超时时间（秒）
//	timeout := tool.ValueInt("app.timeout_seconds")
//
// 注意事项:
//   - 支持的类型转换: int, int64, float64（会截断小数部分）
//   - 配置不存在或类型不匹配时返回 0，无法区分"不存在"和"值为 0"
//   - 如需区分这两种情况，请使用 Value 方法
//   - float64 转 int 时会直接截断小数部分，不会四舍五入
func ValueInt(key string) (value int) {
	val, exists := Value(key)
	if !exists || val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

// ValueIntWithDefault 获取整数配置值，如果不存在则返回默认值
// 这是类型安全的便捷方法，推荐使用
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//   - defaultValue: 默认值，当配置不存在或为 nil 时返回
//
// 返回:
//   - int: 整数类型的配置值或默认值
//
// 使用示例:
//
//	// 获取服务端口，默认 8080
//	port := tool.ValueIntWithDefault("app.port", 8080)
//
//	// 获取工作线程数，默认为 CPU 核心数
//	workers := tool.ValueIntWithDefault("app.workers", runtime.NumCPU())
//
//	// 获取重试次数，默认 3 次
//	retryCount := tool.ValueIntWithDefault("api.retry_count", 3)
//
//	// 获取缓存大小（MB），默认 100MB
//	cacheSize := tool.ValueIntWithDefault("cache.size_mb", 100)
//
// 注意事项:
//   - 支持的类型转换: int, int64, float64（会截断小数部分）
//   - 类型不匹配时也会返回默认值
//   - 推荐在生产环境使用此方法，确保总是有合理的默认值
func ValueIntWithDefault(key string, defaultValue int) int {
	val, exists := Value(key)
	if !exists || val == nil {
		return defaultValue
	}
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return defaultValue
	}
}

// ValueBool 通过点分隔的键获取布尔类型的配置值
// 用于获取开关、标志等布尔配置
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//
// 返回:
//   - bool: 布尔类型的配置值，不存在或类型不匹配时返回 false
//
// 使用示例:
//
//	// 获取调试模式开关
//	debug := tool.ValueBool("app.debug")
//	if debug {
//	    log.SetLevel(log.DebugLevel)
//	}
//
//	// 获取功能开关
//	enableCache := tool.ValueBool("feature.cache_enabled")
//
//	// 获取 SSL 配置
//	useSSL := tool.ValueBool("database.ssl")
//
// 注意事项:
//   - 只支持 bool 类型，不会进行字符串到布尔的转换（如 "true" 不会转为 true）
//   - 配置不存在或类型不匹配时返回 false，无法区分"不存在"和"值为 false"
//   - 如需区分这两种情况，请使用 Value 方法或 ValueBoolWithDefault
func ValueBool(key string) (value bool) {
	val, exists := Value(key)
	if !exists || val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

// ValueBoolWithDefault 获取布尔配置值，如果不存在则返回默认值
// 这是类型安全的便捷方法，推荐使用
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//   - defaultValue: 默认值，当配置不存在或为 nil 时返回
//
// 返回:
//   - bool: 布尔类型的配置值或默认值
//
// 使用示例:
//
//	// 获取调试模式，默认关闭
//	debug := tool.ValueBoolWithDefault("app.debug", false)
//
//	// 获取自动重连开关，默认开启
//	autoReconnect := tool.ValueBoolWithDefault("database.auto_reconnect", true)
//
//	// 获取功能开关，默认启用
//	enableNewFeature := tool.ValueBoolWithDefault("feature.new_ui", true)
//
//	// 获取日志输出开关，默认开启
//	enableLog := tool.ValueBoolWithDefault("log.enabled", true)
//
// 注意事项:
//   - 只支持 bool 类型，不会进行字符串到布尔的转换
//   - 类型不匹配时也会返回默认值
//   - 推荐使用此方法，可以明确表达配置的默认行为
func ValueBoolWithDefault(key string, defaultValue bool) bool {
	val, exists := Value(key)
	if !exists || val == nil {
		return defaultValue
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return defaultValue
}

// ValueFloat 通过点分隔的键获取浮点数类型的配置值
// 支持多种数值类型的自动转换（float64, float32, int, int64）
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//
// 返回:
//   - float64: 浮点数类型的配置值，不存在或类型不匹配时返回 0.0
//
// 使用示例:
//
//	// 获取超时时间（秒，支持小数）
//	timeout := tool.ValueFloat("app.timeout")
//	fmt.Printf("超时时间: %.2f 秒\n", timeout)
//
//	// 获取价格
//	price := tool.ValueFloat("product.price")
//
//	// 获取比率
//	ratio := tool.ValueFloat("algorithm.learning_rate")
//
//	// 获取阈值
//	threshold := tool.ValueFloat("monitor.cpu_threshold")
//
// 注意事项:
//   - 支持的类型转换: float64, float32, int, int64
//   - 配置不存在或类型不匹配时返回 0.0，无法区分"不存在"和"值为 0.0"
//   - 如需区分这两种情况，请使用 Value 方法
func ValueFloat(key string) (value float64) {
	val, exists := Value(key)
	if !exists || val == nil {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

// ValueFloatWithDefault 获取浮点数配置值，如果不存在则返回默认值
// 这是类型安全的便捷方法，推荐使用
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//   - defaultValue: 默认值，当配置不存在或为 nil 时返回
//
// 返回:
//   - float64: 浮点数类型的配置值或默认值
//
// 使用示例:
//
//	// 获取超时时间，默认 30.5 秒
//	timeout := tool.ValueFloatWithDefault("app.timeout", 30.5)
//
//	// 获取 CPU 使用率阈值，默认 80%
//	cpuThreshold := tool.ValueFloatWithDefault("monitor.cpu_threshold", 80.0)
//
//	// 获取价格折扣，默认 0.9（九折）
//	discount := tool.ValueFloatWithDefault("product.discount", 0.9)
//
//	// 获取学习率，默认 0.001
//	learningRate := tool.ValueFloatWithDefault("ml.learning_rate", 0.001)
//
// 注意事项:
//   - 支持的类型转换: float64, float32, int, int64
//   - 类型不匹配时也会返回默认值
//   - 推荐在生产环境使用此方法，确保总是有合理的默认值
func ValueFloatWithDefault(key string, defaultValue float64) float64 {
	val, exists := Value(key)
	if !exists || val == nil {
		return defaultValue
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return defaultValue
	}
}

// ValueSlice 获取切片类型的配置值
// 用于获取数组、列表等配置
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//
// 返回:
//   - [] interface{}: 切片类型的配置值，不存在或类型不匹配时返回 nil
//
// 使用示例:
//
//	// 获取服务器列表
//	servers := tool.ValueSlice("app.servers")
//	for i, server := range servers {
//	    fmt.Printf("服务器 %d: %v\n", i, server)
//	}
//
//	// 获取允许的 IP 列表
//	allowedIPs := tool.ValueSlice("security.allowed_ips")
//
//	// 获取标签列表
//	tags := tool.ValueSlice("product.tags")
//
//	// 配置示例 (YAML):
//	// app:
//	//   servers:
//	//     - "192.168.1.1"
//	//     - "192.168.1.2"
//	//     - "192.168.1.3"
//
// 注意事项:
//   - 返回的是 [] interface{} 类型，元素需要进行类型断言
//   - 配置不存在或类型不匹配时返回 nil
//   - 如果需要特定类型的切片，需要手动转换
func ValueSlice(key string) []interface{} {
	val, exists := Value(key)
	if !exists || val == nil {
		return nil
	}
	if slice, ok := val.([]interface{}); ok {
		return slice
	}
	return nil
}

func ValueStrSlice(key string) (res []string, err error) {
	val, exists := Value(key)
	if !exists || val == nil {
		return nil, errors.New("配置不存在")
	}

	// 先尝试直接断言为 [] string
	if slice, ok := val.([]string); ok {
		return slice, nil
	}

	// YAML 解析后通常是 [] interface{}，需要转换
	slice, ok := val.([]interface{})
	if !ok {
		return nil, errors.New("配置不是切片")
	}

	res = make([]string, 0, len(slice))
	for _, item := range slice {
		var str string
		if str, ok = item.(string); ok {
			res = append(res, str)
		} else {
			// 如果不是字符串，使用 fmt.Sprintf 转换
			res = append(res, fmt.Sprintf("%v", item))
		}
	}

	return
}

// ValueMap 获取 map 类型的配置值
// 用于获取对象、字典等配置
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//
// 返回:
//   - map[string] interface{}: map 类型的配置值，不存在或类型不匹配时返回 nil
//
// 使用示例:
//
//	// 获取数据库配置对象
//	dbConfig := tool.ValueMap("database")
//	if dbConfig != nil {
//	    host := dbConfig["host"]
//	    port := dbConfig["port"]
//	    fmt.Printf("数据库: %v:%v\n", host, port)
//	}
//
//	// 获取 Redis 配置
//	redisConfig := tool.ValueMap("redis")
//
//	// 获取功能开关配置
//	features := tool.ValueMap("features")
//	for name, enabled := range features {
//	    fmt.Printf("功能 %s: %v\n", name, enabled)
//	}
//
//	// 配置示例 (YAML):
//	// database:
//	//   host: "localhost"
//	//   port: 3306
//	//   username: "root"
//
// 注意事项:
//   - 返回的是 map[string] interface{} 类型，值需要进行类型断言
//   - 配置不存在或类型不匹配时返回 nil
//   - 如果需要结构化的配置对象，建议使用 Unmarshal 方法
func ValueMap(key string) map[string]interface{} {
	val, exists := Value(key)
	if !exists || val == nil {
		return nil
	}
	if m, ok := val.(map[string]interface{}); ok {
		return m
	}
	return nil
}

// ValueType 获取配置值，赋值到 dst 指针
// 这是一个通用的类型赋值方法，支持基本类型的指针赋值
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//   - dst: 目标指针，支持 *string, *int, *bool, *float64
//
// 返回:
//   - error: 配置不存在、值为 nil 或类型不支持时返回错误
//
// 使用示例:
//
//	// 获取字符串配置
//	var host string
//	if err := tool.ValueType("database.host", &host); err != nil {
//	    log.Printf("[ValueType]  获取配置失败: %v", err)
//	}
//	fmt.Printf("数据库主机: %s\n", host)
//
//	// 获取整数配置
//	var port int
//	if err := tool.ValueType("database.port", &port); err != nil {
//	    log.Printf("[ValueType]  获取配置失败: %v", err)
//	}
//
//	// 获取布尔配置
//	var debug bool
//	if err := tool.ValueType("app.debug", &debug); err != nil {
//	    log.Printf("[ValueType]  获取配置失败: %v", err)
//	}
//
//	// 获取浮点数配置
//	var timeout float64
//	if err := tool.ValueType("app.timeout", &timeout); err != nil {
//	    log.Printf("[ValueType]  获取配置失败: %v", err)
//	}
//
// 注意事项:
//   - dst 必须是指针类型，否则无法赋值
//   - 目前只支持 *string, *int, *bool, *float64 四种类型
//   - 配置不存在或值为 nil 时会返回错误
//   - 对于复杂类型，建议使用 Unmarshal 方法
func ValueType(key string, dst any) error {
	val, exists := Value(key)
	if !exists {
		return fmt.Errorf("配置键 %s 不存在", key)
	}
	if val == nil {
		return fmt.Errorf("配置键 %s 的值为 nil", key)
	}

	// 使用类型断言来赋值
	switch d := dst.(type) {
	case *string:
		*d = ValueStr(key)
	case *int:
		*d = ValueInt(key)
	case *bool:
		*d = ValueBool(key)
	case *float64:
		*d = ValueFloat(key)
	default:
		return fmt.Errorf("不支持的目标类型")
	}
	return nil
}

// Unmarshal 将指定键的配置解析到结构体
// 这是推荐的配置读取方式，支持将配置映射到自定义结构体
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//   - dst: 目标结构体指针，必须是指针类型
//
// 返回:
//   - error: 配置不存在、解析失败时返回错误
//
// 使用示例:
//
//	// 定义配置结构体
//	type DatabaseConfig struct {
//	    Host     string `yaml:"host"`
//	    Port     int    `yaml:"port"`
//	    Username string `yaml:"username"`
//	    Password string `yaml:"password"`
//	    MaxConn  int    `yaml:"max_connections"`
//	}
//
//	// 解析数据库配置
//	var dbConfig DatabaseConfig
//	if err := tool.Unmarshal("database", &dbConfig); err != nil {
//	    log.Fatalf("[Unmarshal]  解析数据库配置失败: %v", err)
//	}
//	fmt.Printf("数据库: %s:%d\n", dbConfig.Host, dbConfig.Port)
//
//	// 解析应用配置
//	type AppConfig struct {
//	    Name    string `yaml:"name"`
//	    Version string `yaml:"version"`
//	    Debug   bool   `yaml:"debug"`
//	}
//	var appConfig AppConfig
//	if err := tool.Unmarshal("app", &appConfig); err != nil {
//	    log.Fatalf("[Unmarshal]  解析应用配置失败: %v", err)
//	}
//
//	// 配置示例 (YAML):
//	// database:
//	//   host: "localhost"
//	//   port: 3306
//	//   username: "root"
//	//   password: "123456"
//	//   max_connections: 100
//
// 注意事项:
//   - dst 必须是结构体指针，否则无法赋值
//   - 结构体字段需要使用 yaml tag 标注字段名
//   - 支持嵌套结构体
//   - 该方法是线程安全的，使用读锁保护
func Unmarshal(key string, dst interface{}) error {
	configLock.RLock()
	defer configLock.RUnlock()

	if config == nil {
		return fmt.Errorf("配置未加载")
	}

	val, exists := config.get(key)
	if !exists {
		return fmt.Errorf("配置键 %s 不存在", key)
	}

	// 将 val 转换为 yaml 再解析到结构体
	data, err := yaml.Marshal(val)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := yaml.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("解析配置到结构体失败: %w", err)
	}

	return nil
}

// UnmarshalAll 将整个配置解析到结构体
// 用于一次性读取所有配置到一个大的配置结构体
//
// 参数:
//   - dst: 目标结构体指针，必须是指针类型
//
// 返回:
//   - error: 配置未加载、解析失败时返回错误
//
// 使用示例:
//
//	// 定义完整的配置结构体
//	type Config struct {
//	    App struct {
//	        Name    string `yaml:"name"`
//	        Version string `yaml:"version"`
//	        Port    int    `yaml:"port"`
//	        Debug   bool   `yaml:"debug"`
//	    } `yaml:"app"`
//	    Database struct {
//	        Host     string `yaml:"host"`
//	        Port     int    `yaml:"port"`
//	        Username string `yaml:"username"`
//	        Password string `yaml:"password"`
//	    } `yaml:"database"`
//	    Redis struct {
//	        Host string `yaml:"host"`
//	        Port int    `yaml:"port"`
//	    } `yaml:"redis"`
//	}
//
//	// 解析所有配置
//	var config Config
//	if err := tool.UnmarshalAll(&config); err != nil {
//	    log.Fatalf("[UnmarshalAll]  解析配置失败: %v", err)
//	}
//	fmt.Printf("应用: %s v%s\n", config.App.Name, config.App.Version)
//	fmt.Printf("数据库: %s:%d\n", config.Database.Host, config.Database.Port)
//
//	// 配置示例 (YAML):
//	// app:
//	//   name: "MyApp"
//	//   version: "1.0.0"
//	//   port: 8080
//	//   debug: false
//	// database:
//	//   host: "localhost"
//	//   port: 3306
//
// 注意事项:
//   - dst 必须是结构体指针，否则无法赋值
//   - 结构体字段需要使用 yaml tag 标注字段名
//   - 支持嵌套结构体
//   - 推荐在程序启动时使用，一次性加载所有配置
//   - 该方法是线程安全的，使用读锁保护
func UnmarshalAll(dst interface{}) error {
	configLock.RLock()
	defer configLock.RUnlock()

	if config == nil {
		return fmt.Errorf("配置未加载")
	}

	data, err := yaml.Marshal(config.data)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := yaml.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("解析配置到结构体失败: %w", err)
	}

	return nil
}

// Exists 检查配置键是否存在
// 用于在使用配置前检查配置是否存在，避免使用默认值
//
// 参数:
//   - key: 配置键，支持点分隔的嵌套访问
//
// 返回:
//   - bool: 配置键存在返回 true，不存在返回 false
//
// 使用示例:
//
//	// 检查配置是否存在
//	if tool.Exists("app.port") {
//	    port := tool.ValueInt("app.port")
//	    fmt.Printf("使用配置的端口: %d\n", port)
//	} else {
//	    fmt.Println("端口配置不存在，使用默认值")
//	}
//
//	// 检查可选配置
//	if tool.Exists("feature.new_ui") {
//	    enabled := tool.ValueBool("feature.new_ui")
//	    if enabled {
//	        // 启用新 UI
//	    }
//	}
//
//	// 条件性配置加载
//	if tool.Exists("cache.redis") {
//	    // 使用 Redis 缓存
//	    redisConfig := tool.ValueMap("cache.redis")
//	    // ...
//	} else {
//	    // 使用内存缓存
//	}
//
// 注意事项:
//   - 该方法只检查键是否存在，不检查值是否为 nil
//   - 如果配置值为 nil，该方法仍然返回 true
//   - 该方法是线程安全的
func Exists(key string) bool {
	_, exists := Value(key)
	return exists
}

// RequireKeys 验证必需的配置键是否存在
// 用于程序启动时验证必需的配置项，确保配置完整性
//
// 参数:
//   - keys: 必需的配置键列表，支持可变参数
//
// 返回:
//   - error: 如果有配置缺失，返回包含所有缺失配置的错误；全部存在返回 nil
//
// 使用示例:
//
//	// 验证必需的配置项
//	if err := tool.RequireKeys(
//	    "app.name",
//	    "app.port",
//	    "database.host",
//	    "database.port",
//	    "database.username",
//	); err != nil {
//	    log.Fatalf("[RequireKeys]  配置验证失败: %v", err)
//	}
//
//	// 分组验证
//	// 验证应用配置
//	if err := tool.RequireKeys("app.name", "app.version"); err != nil {
//	    log.Fatalf("[RequireKeys]  应用配置不完整: %v", err)
//	}
//
//	// 验证数据库配置
//	if err := tool.RequireKeys("database.host", "database.port"); err != nil {
//	    log.Fatalf("[RequireKeys]  数据库配置不完整: %v", err)
//	}
//
//	// 条件性验证
//	if tool.ValueBool("feature.cache_enabled") {
//	    if err := tool.RequireKeys("cache.redis.host", "cache.redis.port"); err != nil {
//	        log.Fatalf("[RequireKeys]  缓存配置不完整: %v", err)
//	    }
//	}
//
// 注意事项:
//   - 建议在程序启动时调用，尽早发现配置问题
//   - 错误信息会列出所有缺失的配置项
//   - 该方法是线程安全的
func RequireKeys(keys ...string) error {
	var missing []string
	for _, key := range keys {
		if !Exists(key) {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("缺少必需的配置项: %s", strings.Join(missing, ", "))
	}
	return nil
}

// AllKeys 返回所有配置键（扁平化）
// 用于调试、配置导出或配置检查
//
// 返回:
//   - [] string: 所有配置键的列表，使用点分隔表示嵌套关系
//
// 使用示例:
//
//	// 打印所有配置键
//	keys := tool.AllKeys()
//	fmt.Println("所有配置项:")
//	for _, key := range keys {
//	    fmt.Printf("  - %s\n", key)
//	}
//
//	// 检查配置完整性
//	keys := tool.AllKeys()
//	fmt.Printf("共有 %d 个配置项\n", len(keys))
//
//	// 导出配置清单
//	keys := tool.AllKeys()
//	for _, key := range keys {
//	    val, _ := tool.Value(key)
//	    fmt.Printf("%s = %v\n", key, val)
//	}
//
//	// 查找特定前缀的配置
//	keys := tool.AllKeys()
//	for _, key := range keys {
//	    if strings.HasPrefix(key, "database.") {
//	        fmt.Printf("数据库配置: %s\n", key)
//	    }
//	}
//
// 注意事项:
//   - 返回的键是扁平化的，嵌套配置使用点分隔
//   - 只返回叶子节点的键，不包含中间节点
//   - 配置未加载时返回 nil
//   - 该方法是线程安全的，使用读锁保护
func AllKeys() []string {
	configLock.RLock()
	defer configLock.RUnlock()

	if config == nil {
		return nil
	}

	return config.allKeys("", config.data)
}

// Config 配置结构，使用 map 存储以支持灵活的键值访问
// 内部使用 map[string] interface{} 存储配置数据，支持动态键访问和嵌套结构
//
// 字段说明:
//   - data: 存储配置数据的 map，键为字符串，值为任意类型
//
// 注意事项:
//   - 该结构体不应直接实例化，应使用 NewConfig 或 Load 方法
//   - 配置数据是从 YAML 文件解析而来
//   - 支持嵌套的 map 结构，可以表示复杂的配置层次
type Config struct {
	data map[string]interface{}
}

// NewConfig 从文件加载配置
// 读取 YAML 配置文件并解析为 Config 对象
//
// 参数:
//   - file: 配置文件路径，支持相对路径和绝对路径
//
// 返回:
//   - *Config: 配置对象指针
//   - error: 读取或解析失败时返回错误
//
// 使用示例:
//
//	// 加载配置文件
//	config, err := tool.NewConfig("config.yaml")
//	if err != nil {
//	    log.Fatalf("[NewConfig]  加载配置失败: %v", err)
//	}
//
//	// 获取配置值
//	val, exists := config.get("app.port")
//	if exists {
//	    fmt.Printf("端口: %v\n", val)
//	}
//
//	// 通常不需要直接调用此方法，使用 Load 方法即可
//	if err := tool.Load("config.yaml"); err != nil {
//	    log.Fatalf("[Load]  加载配置失败: %v", err)
//	}
//
// 配置文件示例 (config.yaml):
//
//	app:
//	  name: "MyApp"
//	  version: "1.0.0"
//	  port: 8080
//	  debug: false
//	database:
//	  host: "localhost"
//	  port: 3306
//	  username: "root"
//	  password: "123456"
//	redis:
//	  host: "localhost"
//	  port: 6379
//
// 注意事项:
//   - 配置文件必须是有效的 YAML 格式
//   - 文件不存在或格式错误时会返回详细的错误信息
//   - 该方法不会缓存配置，每次调用都会重新读取文件
//   - 通常应该使用 Load 方法而不是直接调用此方法
func NewConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析为 map 以支持动态键访问
	var rawData map[string]interface{}
	if err = yaml.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	cfg := &Config{
		data: rawData,
	}

	return cfg, nil
}

// get 通过点分隔的键获取值
// 支持多层嵌套访问，逐层解析配置路径
//
// 参数:
//   - key: 配置键，使用点分隔表示嵌套关系，如 "app.database.host"
//
// 返回:
//   - interface{}: 配置值，保持原始类型
//   - bool: 配置键是否存在，true 表示存在，false 表示不存在
//
// 工作原理:
//  1. 将键按点分隔符拆分为多个部分
//  2. 从根 map 开始，逐层访问嵌套的 map
//  3. 如果任何一层不存在或类型不是 map，返回 (nil, false)
//  4. 成功访问到最后一层时，返回 (value, true)
//
// 使用示例:
//
//	// 简单键访问
//	val, exists := config.get("port")
//	// 访问: data["port"]
//
//	// 嵌套键访问
//	val, exists := config.get("app.database.host")
//	// 访问: data["app"] ["database"] ["host"]
//
//	// 深层嵌套访问
//	val, exists := config.get("services.payment.api.endpoint")
//	// 访问: data["services"] ["payment"] ["api"] ["endpoint"]
//
// 注意事项:
//   - 返回 (value, exists) 可以区分"不存在"和"值为 nil"
//   - 如果中间任何一层不是 map 类型，会返回 (nil, false)
//   - 该方法不会修改配置数据，是只读操作
//   - 键名区分大小写
func (c *Config) get(key string) (interface{}, bool) {
	if c.data == nil {
		return nil, false
	}

	// 分割键，支持多层嵌套访问
	keys := strings.Split(key, ".")
	var current interface{} = c.data

	// 逐层访问嵌套的 map
	for _, k := range keys {
		// 确保当前值是 map 类型
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}

		// 获取下一层的值
		current, ok = m[k]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

// allKeys 递归获取所有配置键
// 遍历配置树，返回所有叶子节点的完整路径
//
// 参数:
//   - prefix: 当前路径前缀，用于构建完整的键路径
//   - data: 当前层级的配置数据
//
// 返回:
//   - [] string: 所有配置键的列表，使用点分隔表示嵌套关系
//
// 工作原理:
//  1. 遍历当前层级的所有键值对
//  2. 如果值是 map 类型，递归处理（中间节点）
//  3. 如果值不是 map 类型，添加到结果列表（叶子节点）
//  4. 使用 prefix 构建完整的键路径
//
// 使用示例:
//
//	// 配置结构:
//	// app:
//	//   name: "MyApp"
//	//   database:
//	//     host: "localhost"
//	//     port: 3306
//	//
//	// 返回结果:
//	// ["app.name", "app.database.host", "app.database.port"]
//
//	keys := config.allKeys("", config.data)
//	for _, key := range keys {
//	    fmt.Println(key)
//	}
//	// 输出:
//	// app.name
//	// app.database.host
//	// app.database.port
//
// 注意事项:
//   - 只返回叶子节点的键，不包含中间节点
//   - 返回的键顺序不保证，因为 map 是无序的
//   - 该方法是递归实现，配置层级过深可能影响性能
//   - 空 map 不会产生任何键
func (c *Config) allKeys(prefix string, data map[string]interface{}) []string {
	var keys []string
	for k, v := range data {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		if m, ok := v.(map[string]interface{}); ok {
			// 递归处理嵌套 map
			keys = append(keys, c.allKeys(fullKey, m)...)
		} else {
			keys = append(keys, fullKey)
		}
	}
	return keys
}
