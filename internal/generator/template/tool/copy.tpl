package tool

import (
	"github.com/jinzhu/copier"
)

// Copy 深度拷贝工具函数，用于在不同结构体之间安全地复制数据。
//
// 核心特性：
//   - 深度拷贝：指针、切片、Map、嵌套结构体都会生成新对象，修改拷贝结果不会影响原数据；
//   - 字段按名称匹配：源字段多于目标字段时，只拷贝同名且类型兼容的字段；
//   - 默认忽略空值：源字段为零值时，不会覆盖目标字段（等价于 IgnoreEmpty=true）；
//   - 类型转换：支持基本类型之间的自动转换（如 int 到 int64）。
//
// 使用约定：
//   - from: 源数据，可以是值或指针；
//   - to:   目标数据，必须是指针（包括切片、Map 等引用类型也需要传 &to）；
//   - 返回值：拷贝失败时返回 error，成功时返回 nil。
//
// 使用示例：
//
//	type User struct {
//	    Name  string
//	    Age   int
//	    Email string
//	}
//
//	type UserDTO struct {
//	    Name  string
//	    Age   int
//	}
//
//	// 示例 1：结构体拷贝
//	user := User{Name: "张三", Age: 25, Email: "zhangsan@example.com"}
//	var dto UserDTO
//	if err := Copy(user, &dto); err != nil {
//	    // 处理错误
//	}
//	// dto.Name = "张三", dto.Age = 25 (Email 字段不存在于 UserDTO，被忽略)
//
//	// 示例 2：切片拷贝
//	users := []User{{Name: "张三", Age: 25}, {Name: "李四", Age: 30}}
//	var dtos []UserDTO
//	if err := Copy(users, &dtos); err != nil {
//	    // 处理错误
//	}
//	// dtos 包含两个元素的深拷贝
//
//	// 示例 3：指针拷贝
//	userPtr := &User{Name: "王五", Age: 28}
//	var dtoPtr *UserDTO
//	if err := Copy(userPtr, &dtoPtr); err != nil {
//	    // 处理错误
//	}
//	// dtoPtr 指向新分配的 UserDTO 对象
//
// 注意事项：
//   - 目标参数必须传指针，否则拷贝不会生效；
//   - 字段名称必须完全匹配（区分大小写）；
//   - 源字段为零值时不会覆盖目标字段的现有值；
//   - 不支持私有字段（小写字母开头）的拷贝；
//   - 嵌套结构体会递归深拷贝。
//
// 常见场景：
//   - DTO 转换：将数据库实体转换为 API 响应对象；
//   - 数据隔离：避免直接修改原始数据；
//   - 字段过滤：只复制需要的字段到目标结构体。
func Copy(from any, to any) (err error) {
	return copier.CopyWithOption(to, from, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
}

// CopyWithOption 带选项的深度拷贝函数，提供更灵活的拷贝控制。
//
// 与 Copy 的区别：
//   - Copy 使用固定配置（IgnoreEmpty=true, DeepCopy=true）；
//   - CopyWithOption 允许自定义拷贝行为，适用于特殊场景。
//
// 参数说明：
//   - from: 源数据（可以是值或指针）；
//   - to:   目标数据（必须是指针）；
//   - opt:  copier.Option，用于控制拷贝细节。
//
// 可用选项（copier.Option）：
//   - IgnoreEmpty (bool):   是否忽略源字段的零值，默认 false；
//     * true:  源字段为零值时不覆盖目标字段；
//     * false: 源字段为零值时会覆盖目标字段。
//   - DeepCopy (bool):      是否深度拷贝，默认 false；
//     * true:  递归拷贝嵌套结构、切片、Map 等；
//     * false: 只拷贝第一层，引用类型共享底层数据。
//   - CaseSensitive (bool): 字段名是否区分大小写，默认 true；
//   - Converters ([]TypeConverter): 自定义类型转换器。
//
// 使用示例：
//
//	type Source struct {
//	    Name  string
//	    Age   int
//	    Score float64
//	}
//
//	type Target struct {
//	    Name  string
//	    Age   int
//	    Score float64
//	}
//
//	// 示例 1：强制覆盖空值
//	source := Source{Name: "", Age: 0, Score: 85.5}
//	target := Target{Name: "原始名称", Age: 20, Score: 90.0}
//	err := CopyWithOption(source, &target, copier.Option{
//	    IgnoreEmpty: false, // 允许空值覆盖
//	    DeepCopy:    true,
//	})
//	// target.Name = "", target.Age = 0, target.Score = 85.5 (空值被覆盖)
//
//	// 示例 2：浅拷贝（共享引用）
//	type Data struct {
//	    Items []string
//	}
//	src := Data{Items: []string{"a", "b"}}
//	var dst Data
//	err := CopyWithOption(src, &dst, copier.Option{
//	    DeepCopy: false, // 浅拷贝
//	})
//	// dst.Items 和 src.Items 指向同一个底层数组
//
//	// 示例 3：忽略大小写匹配
//	err := CopyWithOption(source, &target, copier.Option{
//	    CaseSensitive: false, // 字段名忽略大小写
//	    DeepCopy:      true,
//	})
//
// 常见场景：
//   - 需要覆盖空值：例如清空某些字段（IgnoreEmpty=false）；
//   - 性能优化：浅拷贝大型切片或 Map（DeepCopy=false）；
//   - 字段名不规范：源和目标字段大小写不一致（CaseSensitive=false）；
//   - 自定义转换：需要特殊的类型转换逻辑（Converters）。
//
// 注意事项：
//   - IgnoreEmpty=false 时，零值会覆盖目标字段，可能导致数据丢失；
//   - DeepCopy=false 时，修改拷贝结果会影响原数据（共享引用）；
//   - 自定义选项需要充分理解业务场景，避免意外行为。
func CopyWithOption(from any, to any, opt copier.Option) (err error) {
	return copier.CopyWithOption(to, from, opt)
}
