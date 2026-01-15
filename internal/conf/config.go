package conf

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	config *Config
)

func Load(file string) (err error) {
	if config != nil {
		err = fmt.Errorf("配置已加载")
		return
	}
	c, err := NewConfig(file)
	if err != nil {
		return
	}
	config = c
	return
}

// Value 通过点分隔的键获取配置值
// 例如: Value("app.product_id") 返回配置值的原始类型
func Value(key string) (value any) {
	if config == nil {
		return nil
	}
	return config.get(key)
}

// ValueStr 通过点分隔的键获取字符串类型的配置值
// 例如: ValueStr("app.product_id") 返回 "N212XOYZ86"
func ValueStr(key string) (value string) {
	val := Value(key)
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

// ValueInt 通过点分隔的键获取整数类型的配置值
// 例如: ValueInt("app.port") 返回 8080
func ValueInt(key string) (value int) {
	val := Value(key)
	if val == nil {
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

// ValueBool 通过点分隔的键获取布尔类型的配置值
// 例如: ValueBool("app.debug") 返回 true 或 false
func ValueBool(key string) (value bool) {
	val := Value(key)
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

// ValueFloat 通过点分隔的键获取浮点数类型的配置值
// 例如: ValueFloat("app.timeout") 返回 3.14
func ValueFloat(key string) (value float64) {
	val := Value(key)
	if val == nil {
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

// ValueType 获取配置值，赋值到 dst 指针
// 例如: var port int; ValueType("app.port", &port)
func ValueType(key string, dst any) error {
	val := Value(key)
	if val == nil {
		return fmt.Errorf("配置键 %s 不存在", key)
	}

	// 这里可以使用反射或类型断言来赋值
	// 简单实现：直接使用类型断言
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

// Config 配置结构，使用 map 存储以支持灵活的键值访问
type Config struct {
	data map[string]interface{}
}

// NewConfig 从文件加载配置
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
// 例如: get("app.product_id") 会依次访问 data["app"]["product_id"]
func (c *Config) get(key string) interface{} {
	if c.data == nil {
		return nil
	}

	// 分割键，支持多层嵌套访问
	keys := strings.Split(key, ".")
	var current interface{} = c.data

	// 逐层访问嵌套的 map
	for _, k := range keys {
		// 确保当前值是 map 类型
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}

		// 获取下一层的值
		current, ok = m[k]
		if !ok {
			return nil
		}
	}

	return current
}
