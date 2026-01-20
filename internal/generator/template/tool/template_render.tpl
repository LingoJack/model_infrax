package tool

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/valyala/fasttemplate"
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
//	data := map[string]interface{}{
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
//	data := map[string]interface{}{
//	    "Title": "购物清单",
//	    "Items": []string{"苹果", "香蕉", "橙子"},
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
func RenderWithName(templateName string, templateToRender string, data map[string]interface{}) (result string, err error) {
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

// RenderWithNameAndFuncMap 渲染Go模板并返回渲染后的字符串（支持自定义函数）
//
// 该函数使用Go标准库的text/template包来解析和执行模板，将数据填充到模板中，
// 并支持通过funcMap注册自定义函数供模板使用。
//
// 参数:
//   - templateName: 模板的名称，用于标识模板（在错误信息中会显示）
//   - templateToRender: 要渲染的模板字符串，使用Go template语法
//   - data: 用于填充模板的数据，以map形式提供，key为模板中的变量名
//   - funcMap: 自定义函数映射，key为函数名，value为函数实现
//
// 返回值:
//   - result: 渲染后的字符串结果
//   - err: 如果解析或执行模板时发生错误，返回相应的错误信息
//
// 使用示例:
//
//	funcMap := template.FuncMap{
//	    "upper": strings.ToUpper,
//	    "lower": strings.ToLower,
//	    "add": func(a, b int) int { return a + b },
//	}
//	data := map[string]interface{}{
//	    "Name": "张三",
//	    "Age":  25,
//	}
//	tmpl := "你好，{{.Name | upper}}，明年你{{add .Age 1}}岁"
//	result, err := RenderWithNameAndFuncMap("greeting", tmpl, data, funcMap)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result) // 输出: 你好，张三，明年你26岁
//
// 注意事项:
//   - funcMap必须在Parse之前通过Funcs方法注册
//   - 自定义函数可以在模板中通过管道符(|)或直接调用使用
//   - 函数参数和返回值类型必须符合Go template的要求
func RenderWithNameAndFuncMap(templateName string, templateToRender string, data map[string]interface{}, funcMap template.FuncMap) (result string, err error) {
	// 创建新模板并注册自定义函数
	tmpl := template.New(templateName)
	if funcMap != nil {
		tmpl = tmpl.Funcs(funcMap)
	}

	// 解析模板字符串
	tmpl, err = tmpl.Parse(templateToRender)
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
func Render(templateToRender string, data map[string]interface{}) (result string, err error) {
	return RenderWithName("", templateToRender, data)
}

// RenderWithFuncMap 渲染模板（支持自定义函数）
func RenderWithFuncMap(templateToRender string, data map[string]interface{}, funcMap template.FuncMap) (result string, err error) {
	return RenderWithNameAndFuncMap("", templateToRender, data, funcMap)
}

// MustRender 渲染模板，如果失败则panic
func MustRender(templateToRender string, data map[string]interface{}) (result string) {
	result, err := Render(templateToRender, data)
	if err != nil {
		panic(err)
	}
	return result
}

// MustRenderWithFuncMap 渲染模板（支持自定义函数），如果失败则panic
func MustRenderWithFuncMap(templateToRender string, data map[string]interface{}, funcMap template.FuncMap) (result string) {
	result, err := RenderWithFuncMap(templateToRender, data, funcMap)
	if err != nil {
		panic(err)
	}
	return result
}

// ExpandTemplate 使用fasttemplate渲染模板，只支持${var}格式的变量替换
//
// 该函数使用fasttemplate库进行变量替换，只识别${var}格式的变量，
//
// 参数:
//   - templateToRender: 要渲染的模板字符串，使用${var}格式引用变量
//   - data: 用于填充模板的数据，以map形式提供，key为变量名
//
// 返回值:
//   - result: 渲染后的字符串结果
//
// 变量格式:
//   - ${var} - 唯一支持的格式
//
// 使用示例:
//
//	data := map[string]interface{}{
//	    "Name": "张三",
//	    "Age":  "25",
//	    "City": "北京",
//	}
//	tmpl := "你好，我是${Name}，今年${Age}岁，来自${City}，我有$100"
//	result = ExpandTemplate(tmpl, data)
//	fmt.Println(result) // 输出: 你好，我是张三，今年25岁，来自北京，我有$100
//
// 注意事项:
//   - 只支持${var}格式，$var格式不会被识别
//   - 所有值都会被转换为字符串
//   - 如果变量不存在，会返回错误
//   - 不支持复杂的模板语法（如条件、循环等）
//   - 适用于简单的配置文件、环境变量替换等场景
//   - 变量名必须是有效的标识符（字母、数字、下划线）
func ExpandTemplate(templateToRender string, data map[string]interface{}) (result string) {
	t := fasttemplate.New(templateToRender, "${", "}")
	vars := make(map[string]interface{}, len(data))
	for k, v := range data {
		vars[k] = fmt.Sprintf("%v", v)
	}
	result = t.ExecuteString(vars)
	return
}
