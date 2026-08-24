package conf

import (
	"errors"
	"fmt"
	"strings"
)

// splitKey 将点分隔的配置键拆分为路径片段
func splitKey(key string) []string {
	return strings.Split(key, ".")
}

// value 读取指定键的原始配置值（线程安全）
// 返回 (值, 是否存在)
func rawValue(key string) (any, bool) {
	configLock.RLock()
	defer configLock.RUnlock()

	if config == nil {
		return nil, false
	}
	return config.get(key)
}

// ValueStr 获取字符串配置值
// 配置不存在返回空字符串；非字符串类型使用 %v 转换
func ValueStr(key string) string {
	val, exists := rawValue(key)
	if !exists || val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

// ValueInt 获取整数配置值
// 支持 int/int64/float64 自动转换（float64 直接截断小数）
// 配置不存在或类型不匹配时返回 0
func ValueInt(key string) int {
	val, exists := rawValue(key)
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

// ValueBool 获取布尔配置值
// 只支持 bool 类型（不做字符串转换）；配置不存在或类型不匹配时返回 false
func ValueBool(key string) bool {
	val, exists := rawValue(key)
	if !exists || val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

// ValueStrSlice 获取字符串切片配置值
// 兼容 YAML 解析出的 []interface{}，非字符串元素使用 %v 转换
// 配置不存在或不是切片时返回错误
func ValueStrSlice(key string) (res []string, err error) {
	val, exists := rawValue(key)
	if !exists || val == nil {
		return nil, errors.New("配置不存在")
	}

	// 先尝试直接断言为 []string
	if slice, ok := val.([]string); ok {
		return slice, nil
	}

	// YAML 解析后通常是 []interface{}，需要转换
	slice, ok := val.([]interface{})
	if !ok {
		return nil, errors.New("配置不是切片")
	}

	res = make([]string, 0, len(slice))
	for _, item := range slice {
		if str, ok := item.(string); ok {
			res = append(res, str)
		} else {
			res = append(res, fmt.Sprintf("%v", item))
		}
	}

	return
}
