package tool

import (
	"github.com/jinzhu/copier"
)

// Copy 深度拷贝工具函数。
//
// 特性：
//   - 深度拷贝：指针、切片、Map、嵌套结构体都会生成新对象；
//   - 字段按名称匹配：源字段多于目标字段时，只拷贝同名且类型兼容的字段；
//   - 默认忽略空值：源字段为零值时，不会覆盖目标字段（等价于 IgnoreEmpty=true）。
//
// 使用约定：
//   - from 可以是值或指针；
//   - to 必须是指针，切片、Map 等引用类型同样需要传 &to；
//   - 修改拷贝结果不会影响原数据。
func Copy(from any, to any) (err error) {
	return copier.CopyWithOption(to, from, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
}

// CopyWithOption 带选项的深度拷贝函数。
//
// 与 Copy 的区别：允许自定义拷贝行为，例如：
//   - 强制覆盖空值（IgnoreEmpty=false）；
//   - 控制是否深度拷贝（DeepCopy）。
//
// 参数说明：
//   - from: 源数据（可以是值或指针）；
//   - to:   目标数据（必须是指针）；
//   - opt:  copier.Option，用于控制拷贝细节。
func CopyWithOption(from any, to any, opt copier.Option) (err error) {
	return copier.CopyWithOption(to, from, opt)
}