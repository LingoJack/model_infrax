package tool

import (
	"encoding/json"
)

// DeepCopyByJson 深度拷贝对象
// 使用 JSON 序列化/反序列化实现深度拷贝，支持泛型类型
// 注意：该方法依赖 JSON 编解码，因此：
//  1. 源对象必须可以被 JSON 序列化（字段需要导出或有 json tag）
//  2. 不会拷贝未导出的字段
//  3. 不会拷贝函数、channel 等不可序列化的类型
//  4. 时间类型会被转换为 RFC3339 格式字符串后再解析
//
// 参数:
//   - src: 源对象，需要被拷贝的对象
//
// 返回:
//   - T: 拷贝后的新对象，类型与源对象相同
//   - error: 如果拷贝过程中发生错误（序列化或反序列化失败），返回错误信息
//
// 示例:
//
//	type User struct {
//	    Name string
//	    Age  int
//	    Tags []string
//	}
//	original := User{Name: "Alice", Age: 30, Tags: []string{"admin", "user"}}
//	copied, err := DeepCopy(original)
//	if err != nil {
//	    // 处理错误
//	}
//	// copied 是 User 类型，无需类型转换
//	copied.Name = "Bob" // 修改拷贝不影响原对象
func DeepCopyByJson[T any](src T) (T, error) {
	var dst T

	// 序列化源对象为 JSON
	data, err := json.Marshal(src)
	if err != nil {
		return dst, err
	}

	// 反序列化 JSON 到目标对象
	err = json.Unmarshal(data, &dst)
	if err != nil {
		return dst, err
	}

	return dst, nil
}

// DeepCopyToByJson 深度拷贝对象到不同类型
// 使用 JSON 序列化/反序列化实现类型转换和深度拷贝
// 适用于需要在不同但结构相似的类型之间转换的场景
//
// 注意：
//  1. 源类型和目标类型的字段名需要匹配（或通过 json tag 匹配）
//  2. 目标类型中不存在的字段会被忽略
//  3. 源类型中不存在但目标类型存在的字段会使用零值
//
// 参数:
//   - src: 源对象，需要被拷贝的对象
//
// 返回:
//   - T: 目标类型的新对象
//   - error: 如果拷贝过程中发生错误（序列化或反序列化失败），返回错误信息
//
// 示例:
//
//	type UserDTO struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	type UserEntity struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	    ID   int    `json:"id"` // 这个字段在 UserDTO 中不存在，会使用零值
//	}
//	dto := UserDTO{Name: "Alice", Age: 30}
//	entity, err := DeepCopyTo[UserEntity](dto)
//	if err != nil {
//	    // 处理错误
//	}
//	// entity 是 UserEntity 类型，无需类型转换
func DeepCopyToByJson[T any](src any) (T, error) {
	var dst T

	// 序列化源对象为 JSON
	data, err := json.Marshal(src)
	if err != nil {
		return dst, err
	}

	// 反序列化 JSON 到目标对象
	err = json.Unmarshal(data, &dst)
	if err != nil {
		return dst, err
	}

	return dst, nil
}
