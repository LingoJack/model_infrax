package tool

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

// RenderWithName 渲染Go模板并返回渲染后的字符串
//
// 该函数使用Go标准库的text/template包来解析和执行模板，将数据填充到模板中。
//
// 参数:
//   - templateName: 模板的名称，用于标识模板（在错误信息中会显示）
//   - templateToRender: 要渲染的模板字符串，使用Go template语法
//   - data: 用于填充模板的数据，以map形式提供，key为模板中的变量名
//
// 返回值:
//   - result: 渲染后的字符串结果
//   - err: 如果解析或执行模板时发生错误，返回相应的错误信息
//
// 模板语法示例:
//   - {{.VariableName}} - 访问数据中的变量
//   - {{if .Condition}}...{{end}} - 条件判断
//   - {{range .List}}...{{end}} - 循环遍历
//
// 使用示例:
//
//	data := map[string] interface{}{
//	    "Name": "张三",
//	    "Age":  25,
//	    "City": "北京",
//	}
//	tmpl := "你好，我是{{.Name}}，今年{{.Age}}岁，来自{{.City}}"
//	result, err := Render("greeting", tmpl, data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result) // 输出: 你好，我是张三，今年25岁，来自北京
//
// 复杂示例（使用条件和循环）:
//
//	data := map[string] interface{}{
//	    "Title": "购物清单",
//	    "Items": [] string{"苹果", "香蕉", "橙子"},
//	    "Total": 3,
//	}
//	tmpl := `{{.Title}}:
//	{{range .Items}}- {{.}}
//	{{end}}共计: {{.Total}}项`
//	result, err := Render("shopping", tmpl, data)
//
// 错误处理:
//   - 当模板语法错误时，返回"template parse error: ..."
//   - 当模板执行错误时（如访问不存在的变量），返回"template execute error: ..."
//
// 注意事项:
//   - 模板中引用的变量必须在data map中存在，否则会导致执行错误
//   - 模板语法严格遵循Go text/template规范
//   - 对于复杂的模板，建议先进行充分测试
func RenderWithName(templateName string, templateToRender string, data map[string] interface{}) (result string, err error) {
	// 创建新模板并解析模板字符串
	tmpl, err := template.New(templateName).Parse(templateToRender)
	if err != nil {
		err = errors.New(fmt.Sprintf("template parse error: %v", err))
		return
	}

	// 创建缓冲区用于存储渲染结果
	var buf bytes.Buffer

	// 执行模板，将数据填充到模板中
	err = tmpl.Execute(&buf, data)
	if err != nil {
		err = errors.New(fmt.Sprintf("template execute error: %v", err))
		return
	}

	// 将缓冲区内容转换为字符串返回
	result = buf.String()
	return
}

// Render 渲染模板
func Render(templateToRender string, data map[string] interface{}) (result string, err error) {
	return RenderWithName("", templateToRender, data)
}

// MustRender 渲染模板，如果失败则panic
func MustRender(templateToRender string, data map[string] interface{}) (result string) {
	result, err := Render(templateToRender, data)
	if err != nil {
		panic(err)
	}
	return result
}
