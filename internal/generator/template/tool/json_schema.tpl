package tool

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/jsonschema"
)

// StructToJSONSchema 将Go结构体转换为JSON Schema字符串
//
// 该函数使用invopop/jsonschema库通过反射机制将Go结构体转换为符合JSON Schema Draft 2020-12标准的Schema定义。
//
// 参数:
//   - structInstance: 要转换的结构体实例（必须传入指针或实例）
//
// 返回值:
//   - schemaJSON: JSON Schema的字符串表示（格式化后的JSON）
//   - err: 如果转换或序列化时发生错误，返回相应的错误信息
//
// 使用示例:
//
//	type User struct {
//	    Name  string `json:"name" jsonschema:"title=用户名,description=用户的姓名,minLength=1,maxLength=50"`
//	    Age   int    `json:"age" jsonschema:"title=年龄,description=用户的年龄,minimum=0,maximum=150"`
//	    Email string `json:"email" jsonschema:"title=邮箱,description=用户的电子邮箱,format=email"`
//	}
//	schema, err := StructToJSONSchema(&User{})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(schema)
//
// 输出示例:
//
//	{
//	  "$schema": "https://json-schema.org/draft/2020-12/schema",
//	  "type": "object",
//	  "properties": {
//	    "name": {
//	      "type": "string",
//	      "title": "用户名",
//	      "description": "用户的姓名",
//	      "minLength": 1,
//	      "maxLength": 50
//	    },
//	    "age": {
//	      "type": "integer",
//	      "title": "年龄",
//	      "description": "用户的年龄",
//	      "minimum": 0,
//	      "maximum": 150
//	    },
//	    "email": {
//	      "type": "string",
//	      "title": "邮箱",
//	      "description": "用户的电子邮箱",
//	      "format": "email"
//	    }
//	  }
//	}
//
// 支持的jsonschema标签:
//   - title: 字段标题
//   - description: 字段描述
//   - required: 必填字段（在结构体级别使用）
//   - enum: 枚举值（例如: enum=active|inactive|pending）
//   - minLength/maxLength: 字符串长度限制
//   - minimum/maximum: 数值范围限制
//   - format: 格式验证（如email、uri、date-time等）
//   - pattern: 正则表达式验证
//   - default: 默认值
//
// 注意事项:
//   - 建议传入结构体指针以获得更准确的Schema
//   - 使用json标签定义字段名称，使用jsonschema标签定义验证规则
//   - 支持嵌套结构体、切片、映射等复杂类型
//   - 生成的Schema符合JSON Schema Draft 2020-12标准
func StructToJSONSchema(structInstance interface{}) (schemaJSON string, err error) {
	// 创建默认的Reflector实例
	reflector := &jsonschema.Reflector{}

	// 通过反射生成JSON Schema
	schema := reflector.Reflect(structInstance)

	// 将Schema对象序列化为格式化的JSON字符串
	schemaBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		err = errors.New(fmt.Sprintf("json schema marshal error: %v", err))
		return
	}

	schemaJSON = string(schemaBytes)
	return
}

// StructToJSONSchemaWithOptions 将Go结构体转换为JSON Schema字符串（支持自定义选项）
//
// 该函数提供了更灵活的配置选项，允许自定义Reflector的行为。
//
// 参数:
//   - structInstance: 要转换的结构体实例
//   - options: 自定义配置函数，用于配置Reflector的行为
//
// 返回值:
//   - schemaJSON: JSON Schema的字符串表示
//   - err: 如果转换或序列化时发生错误，返回相应的错误信息
//
// 使用示例:
//
//	type Product struct {
//	    ID    int    `json:"id"`
//	    Name  string `json:"name"`
//	    Price float64 `json:"price"`
//	}
//
//	// 自定义配置：不生成$id字段，展开所有引用
//	schema, err := StructToJSONSchemaWithOptions(&Product{}, func(r *jsonschema.Reflector) {
//	    r.Anonymous = true              // 不生成$id字段
//	    r.ExpandedStruct = true         // 展开嵌套结构体而不是使用$ref
//	    r.DoNotReference = true         // 不使用$ref引用
//	    r.AllowAdditionalProperties = false  // 不允许额外属性
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(schema)
//
// 常用配置选项:
//   - Anonymous: 设置为true时不生成$id字段
//   - ExpandedStruct: 设置为true时展开嵌套结构体
//   - DoNotReference: 设置为true时不使用$ref引用
//   - AllowAdditionalProperties: 是否允许额外属性
//   - RequiredFromJSONSchemaTags: 从jsonschema标签读取required信息
//
// 注意事项:
//   - 配置函数可以为nil，此时使用默认配置
//   - 不同的配置会影响生成的Schema结构和大小
func StructToJSONSchemaWithOptions(structInstance interface{}, options func(*jsonschema.Reflector)) (schemaJSON string, err error) {
	// 创建Reflector实例
	reflector := &jsonschema.Reflector{}

	// 应用自定义配置
	if options != nil {
		options(reflector)
	}

	// 通过反射生成JSON Schema
	schema := reflector.Reflect(structInstance)

	// 将Schema对象序列化为格式化的JSON字符串
	schemaBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		err = errors.New(fmt.Sprintf("json schema marshal error: %v", err))
		return
	}

	schemaJSON = string(schemaBytes)
	return
}

// StructToJSONSchemaObject 将Go结构体转换为JSON Schema对象
//
// 该函数返回jsonschema.Schema对象而不是JSON字符串，适用于需要进一步处理Schema的场景。
//
// 参数:
//   - structInstance: 要转换的结构体实例
//
// 返回值:
//   - schema: jsonschema.Schema对象，可以进行进一步的操作和修改
//
// 使用示例:
//
//	type Config struct {
//	    Host string `json:"host"`
//	    Port int    `json:"port"`
//	}
//
//	schema := StructToJSONSchemaObject(&Config{})
//
//	// 可以进一步修改Schema对象
//	schema.Title = "服务器配置"
//	schema.Description = "服务器的配置信息"
//
//	// 序列化为JSON
//	schemaJSON, _ := json.MarshalIndent(schema, "", "  ")
//	fmt.Println(string(schemaJSON))
//
// 注意事项:
//   - 返回的Schema对象可以被修改和扩展
//   - 适用于需要动态调整Schema的场景
//   - 需要手动序列化为JSON字符串
func StructToJSONSchemaObject(structInstance interface{}) *jsonschema.Schema {
	reflector := &jsonschema.Reflector{}
	return reflector.Reflect(structInstance)
}

// StructToJSONSchemaCompact 将Go结构体转换为紧凑格式的JSON Schema字符串
//
// 该函数生成不带缩进的紧凑JSON字符串，适用于网络传输或存储场景。
//
// 参数:
//   - structInstance: 要转换的结构体实例
//
// 返回值:
//   - schemaJSON: 紧凑格式的JSON Schema字符串
//   - err: 如果转换或序列化时发生错误，返回相应的错误信息
//
// 使用示例:
//
//	type Message struct {
//	    Content string `json:"content"`
//	    Sender  string `json:"sender"`
//	}
//
//	schema, err := StructToJSONSchemaCompact(&Message{})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(schema)
//	// 输出: {"$schema":"https://json-schema.org/draft/2020-12/schema","type":"object",...}
//
// 注意事项:
//   - 生成的JSON字符串不包含换行和缩进，体积更小
//   - 适用于API传输、配置文件存储等场景
//   - 可读性较差，调试时建议使用StructToJSONSchema
func StructToJSONSchemaCompact(structInstance interface{}) (schemaJSON string, err error) {
	reflector := &jsonschema.Reflector{}
	schema := reflector.Reflect(structInstance)

	// 使用Marshal而不是MarshalIndent生成紧凑格式
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		err = errors.New(fmt.Sprintf("json schema marshal error: %v", err))
		return
	}

	schemaJSON = string(schemaBytes)
	return
}

// MustStructToJSONSchema 将Go结构体转换为JSON Schema字符串，如果失败则panic
//
// 该函数是StructToJSONSchema的panic版本，适用于确定不会出错的场景。
//
// 参数:
//   - structInstance: 要转换的结构体实例
//
// 返回值:
//   - schemaJSON: JSON Schema的字符串表示
//
// 使用示例:
//
//	type Settings struct {
//	    Theme    string `json:"theme"`
//	    Language string `json:"language"`
//	}
//
//	// 在初始化阶段使用，如果失败会panic
//	schema := MustStructToJSONSchema(&Settings{})
//	fmt.Println(schema)
//
// 注意事项:
//   - 仅在确定结构体定义正确且不会出错时使用
//   - 适用于程序初始化阶段的配置加载
//   - 如果转换失败会导致程序panic
func MustStructToJSONSchema(structInstance interface{}) (schemaJSON string) {
	schemaJSON, err := StructToJSONSchema(structInstance)
	if err != nil {
		panic(err)
	}
	return schemaJSON
}

// ValidateJSONAgainstStruct 验证JSON数据是否符合结构体定义的Schema
//
// 该函数先将结构体转换为JSON Schema，然后验证给定的JSON数据是否符合该Schema。
// 注意：此函数仅生成Schema，实际验证需要配合JSON Schema验证库使用。
//
// 参数:
//   - structInstance: 结构体实例，用于生成Schema
//   - jsonData: 要验证的JSON字符串
//
// 返回值:
//   - valid: 是否有效（当前实现始终返回true，需要配合验证库）
//   - schema: 生成的JSON Schema字符串
//   - err: 如果生成Schema时发生错误，返回相应的错误信息
//
// 使用示例:
//
//	type User struct {
//	    Name string `json:"name" jsonschema:"required"`
//	    Age  int    `json:"age" jsonschema:"minimum=0"`
//	}
//
//	jsonData := `{"name": "张三", "age": 25}`
//	valid, schema, err := ValidateJSONAgainstStruct(&User{}, jsonData)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Valid: %v\nSchema: %s\n", valid, schema)
//
// 注意事项:
//   - 当前实现仅生成Schema，不执行实际验证
//   - 如需完整验证功能，建议配合github.com/xeipuuv/gojsonschema等验证库使用
//   - 返回的schema可用于前端验证或文档生成
func ValidateJSONAgainstStruct(structInstance interface{}, jsonData string) (valid bool, schema string, err error) {
	// 生成JSON Schema
	schema, err = StructToJSONSchema(structInstance)
	if err != nil {
		return
	}

	// 注意：这里仅生成Schema，实际验证需要使用专门的JSON Schema验证库
	// 例如: github.com/xeipuuv/gojsonschema
	// 这里返回true仅作为示例，实际使用时应该实现真正的验证逻辑
	valid = true

	return
}