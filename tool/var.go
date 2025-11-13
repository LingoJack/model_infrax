package tool

import "fmt"

func Stringify(v any) string {
	// 空字符串直接返回
	if v == "" {
		return ""
	}

	switch val := v.(type) {
	// 字符串类型
	case *string:
		if val == nil {
			return ""
		}
		return *val
	case string:
		return val

	// 有符号整数指针类型
	case *int:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%d", *val)
	case *int8:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%d", *val)
	case *int16:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%d", *val)
	case *int32:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%d", *val)
	case *int64:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%d", *val)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)

	// 无符号整数指针类型
	case *uint:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%d", *val)
	case *uint8:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%d", *val)
	case *uint16:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%d", *val)
	case *uint32:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%d", *val)
	case *uint64:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%d", *val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)

	// 浮点数指针类型
	case *float32:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%f", *val)
	case *float64:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%f", *val)
	case float32, float64:
		return fmt.Sprintf("%f", val)

	// 布尔指针类型
	case *bool:
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%t", *val)
	case bool:
		return fmt.Sprintf("%t", val)

	// 其他类型使用 JSON 格式化
	default:
		return JsonifyIndent(val)
	}
}
