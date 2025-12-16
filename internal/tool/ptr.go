package tool

import "time"

// Ptr 系列函数用于快速将值转换为指针
// 主要用于 Builder 模式中需要传递指针类型的场景

// StringPtr 将 string 值转换为 *string 指针
// 参数:
//   - v: 字符串值
//
// 返回:
//   - *string: 指向该字符串的指针
//
// 示例:
//   - builder.WithName(tool.StringPtr("test"))
func StringPtr(v string) *string {
	return &v
}

// IntPtr 将 int 值转换为 *int 指针
// 参数:
//   - v: 整数值
//
// 返回:
//   - *int: 指向该整数的指针
//
// 示例:
//   - builder.WithAge(tool.IntPtr(18))
func IntPtr(v int) *int {
	return &v
}

// Int8Ptr 将 int8 值转换为 *int8 指针
// 参数:
//   - v: int8 值
//
// 返回:
//   - *int8: 指向该值的指针
func Int8Ptr(v int8) *int8 {
	return &v
}

// Int16Ptr 将 int16 值转换为 *int16 指针
// 参数:
//   - v: int16 值
//
// 返回:
//   - *int16: 指向该值的指针
func Int16Ptr(v int16) *int16 {
	return &v
}

// Int32Ptr 将 int32 值转换为 *int32 指针
// 参数:
//   - v: int32 值
//
// 返回:
//   - *int32: 指向该值的指针
func Int32Ptr(v int32) *int32 {
	return &v
}

// Int64Ptr 将 int64 值转换为 *int64 指针
// 参数:
//   - v: int64 值
//
// 返回:
//   - *int64: 指向该值的指针
func Int64Ptr(v int64) *int64 {
	return &v
}

// UintPtr 将 uint 值转换为 *uint 指针
// 参数:
//   - v: uint 值
//
// 返回:
//   - *uint: 指向该值的指针
func UintPtr(v uint) *uint {
	return &v
}

// Uint8Ptr 将 uint8 值转换为 *uint8 指针
// 参数:
//   - v: uint8 值
//
// 返回:
//   - *uint8: 指向该值的指针
func Uint8Ptr(v uint8) *uint8 {
	return &v
}

// Uint16Ptr 将 uint16 值转换为 *uint16 指针
// 参数:
//   - v: uint16 值
//
// 返回:
//   - *uint16: 指向该值的指针
func Uint16Ptr(v uint16) *uint16 {
	return &v
}

// Uint32Ptr 将 uint32 值转换为 *uint32 指针
// 参数:
//   - v: uint32 值
//
// 返回:
//   - *uint32: 指向该值的指针
func Uint32Ptr(v uint32) *uint32 {
	return &v
}

// Uint64Ptr 将 uint64 值转换为 *uint64 指针
// 参数:
//   - v: uint64 值
//
// 返回:
//   - *uint64: 指向该值的指针
func Uint64Ptr(v uint64) *uint64 {
	return &v
}

// Float32Ptr 将 float32 值转换为 *float32 指针
// 参数:
//   - v: float32 值
//
// 返回:
//   - *float32: 指向该值的指针
func Float32Ptr(v float32) *float32 {
	return &v
}

// Float64Ptr 将 float64 值转换为 *float64 指针
// 参数:
//   - v: float64 值
//
// 返回:
//   - *float64: 指向该值的指针
func Float64Ptr(v float64) *float64 {
	return &v
}

// BoolPtr 将 bool 值转换为 *bool 指针
// 参数:
//   - v: bool 值
//
// 返回:
//   - *bool: 指向该值的指针
func BoolPtr(v bool) *bool {
	return &v
}

// TimePtr 将 time.Time 值转换为 *time.Time 指针
// 参数:
//   - v: time.Time 值
//
// 返回:
//   - *time.Time: 指向该时间的指针
func TimePtr(v time.Time) *time.Time {
	return &v
}

// BytePtr 将 byte 值转换为 *byte 指针
// 参数:
//   - v: byte 值
//
// 返回:
//   - *byte: 指向该值的指针
func BytePtr(v byte) *byte {
	return &v
}

// RunePtr 将 rune 值转换为 *rune 指针
// 参数:
//   - v: rune 值
//
// 返回:
//   - *rune: 指向该值的指针
func RunePtr(v rune) *rune {
	return &v
}
