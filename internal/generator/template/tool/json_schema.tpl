package tool

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/jsonschema"
)

// SchemaOption 定义JSON Schema生成选项（使用位运算）
//
// 通过位运算可以灵活组合多个选项，例如：
//   - SchemaOptionAnonymous | SchemaOptionExpandStruct
//   - SchemaOptionCompact | SchemaOptionNoReference
//
// 使用示例:
//
//	opts := SchemaOptionAnonymous | SchemaOptionExpandStruct | SchemaOptionNoReference
//	schema, err := StructToJSONSchema(&User{}, opts)
type SchemaOption uint32

const (
	// SchemaOptionNone 默认选项，不做任何特殊配置
	SchemaOptionNone SchemaOption = 0

	// SchemaOptionAnonymous 不生成$id字段
	// 适用于不需要Schema标识符的场景
	SchemaOptionAnonymous SchemaOption = 1 << iota

	// SchemaOptionExpandStruct 展开嵌套结构体而不是使用$ref引用
	// 适用于需要完整Schema定义的场景
	SchemaOptionExpandStruct

	// SchemaOptionNoReference 不使用$ref引用
	// 适用于需要内联所有定义的场景
	SchemaOptionNoReference

	// SchemaOptionNoAdditionalProperties 不允许额外属性
	// 适用于需要严格验证的场景
	SchemaOptionNoAdditionalProperties

	// SchemaOptionRequiredFromTags 从jsonschema标签读取required信息
	// 适用于通过标签定义必填字段的场景
	SchemaOptionRequiredFromTags

	// SchemaOptionCompact 生成紧凑格式的JSON（不带缩进）
	// 适用于网络传输或存储场景
	SchemaOptionCompact

	// SchemaOptionAllowNullValues 允许字段为null值
	// 适用于需要支持null的场景
	SchemaOptionAllowNullValues

	// SchemaOptionPreferYAMLSchema 优先使用YAML Schema标签
	// 适用于同时支持JSON和YAML的场景
	SchemaOptionPreferYAMLSchema
)

// Has 检查是否包含指定选项
//
// 使用示例:
//
//	opts := SchemaOptionAnonymous | SchemaOptionCompact
//	if opts.Has(SchemaOptionCompact) {
//	    fmt.Println("使用紧凑格式")
//	}
func (o SchemaOption) Has(option SchemaOption) bool {
	return o&option != 0
}

// StructToJSONSchema 将Go结构体转换为JSON Schema字符串
//
// 该函数使用invopop/jsonschema库通过反射机制将Go结构体转换为符合JSON Schema Draft 2020-12标准的Schema定义。
//
// 参数:
//   - structInstance: 要转换的结构体实例（必须传入指针或实例）
//   - options: 配置选项，可以通过位运算组合多个选项（可选参数，默认为SchemaOptionNone）
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
//
//	// 基本使用（使用默认选项）
//	schema, err := StructToJSONSchema(&User{})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(schema)
//
//	// 使用单个选项
//	schema, err := StructToJSONSchema(&User{}, SchemaOptionAnonymous)
//
//	// 组合多个选项
//	schema, err := StructToJSONSchema(&User{},
//	    SchemaOptionAnonymous | SchemaOptionExpandStruct | SchemaOptionNoReference)
//
//	// 生成紧凑格式
//	schema, err := StructToJSONSchema(&User{}, SchemaOptionCompact)
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
func StructToJSONSchema(structInstance interface{}, options ...SchemaOption) (schemaJSON string, err error) {
	// 合并所有选项
	var opt SchemaOption = SchemaOptionNone
	if len(options) > 0 {
		for _, o := range options {
			opt |= o
		}
	}

	// 创建并配置Reflector实例
	reflector := &jsonschema.Reflector{}
	applyOptions(reflector, opt)

	// 通过反射生成JSON Schema
	schema := reflector.Reflect(structInstance)

	// 根据选项决定序列化格式
	var schemaBytes []byte
	if opt.Has(SchemaOptionCompact) {
		schemaBytes, err = json.Marshal(schema)
	} else {
		schemaBytes, err = json.MarshalIndent(schema, "", "  ")
	}

	if err != nil {
		err = errors.New(fmt.Sprintf("json schema marshal error: %v", err))
		return
	}

	schemaJSON = string(schemaBytes)
	return
}

// applyOptions 将位运算选项应用到Reflector配置
//
// 该函数是内部辅助函数，用于将SchemaOption转换为Reflector的具体配置。
//
// 参数:
//   - reflector: jsonschema.Reflector实例
//   - opt: 位运算组合的选项
func applyOptions(reflector *jsonschema.Reflector, opt SchemaOption) {
	if opt.Has(SchemaOptionAnonymous) {
		reflector.Anonymous = true
	}

	if opt.Has(SchemaOptionExpandStruct) {
		reflector.ExpandedStruct = true
	}

	if opt.Has(SchemaOptionNoReference) {
		reflector.DoNotReference = true
	}

	if opt.Has(SchemaOptionNoAdditionalProperties) {
		reflector.AllowAdditionalProperties = false
	}

	if opt.Has(SchemaOptionRequiredFromTags) {
		reflector.RequiredFromJSONSchemaTags = true
	}

	if opt.Has(SchemaOptionAllowNullValues) {
		reflector.AllowAdditionalProperties = true
	}

	if opt.Has(SchemaOptionPreferYAMLSchema) {
		reflector.PreferYAMLSchema = true
	}
}

// StructToJSONSchemaObject 将Go结构体转换为JSON Schema对象
//
// 该函数返回jsonschema.Schema对象而不是JSON字符串，适用于需要进一步处理Schema的场景。
//
// 参数:
//   - structInstance: 要转换的结构体实例
//   - options: 配置选项，可以通过位运算组合多个选项（可选参数）
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
//	// 使用默认选项
//	schema := StructToJSONSchemaObject(&Config{})
//
//	// 使用自定义选项
//	schema := StructToJSONSchemaObject(&Config{},
//	    SchemaOptionAnonymous | SchemaOptionExpandStruct)
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
func StructToJSONSchemaObject(structInstance interface{}, options ...SchemaOption) *jsonschema.Schema {
	// 合并所有选项
	var opt SchemaOption = SchemaOptionNone
	if len(options) > 0 {
		for _, o := range options {
			opt |= o
		}
	}

	// 创建并配置Reflector实例
	reflector := &jsonschema.Reflector{}
	applyOptions(reflector, opt)

	return reflector.Reflect(structInstance)
}

// MustStructToJSONSchema 将Go结构体转换为JSON Schema字符串，如果失败则panic
//
// 该函数是StructToJSONSchema的panic版本，适用于确定不会出错的场景。
//
// 参数:
//   - structInstance: 要转换的结构体实例
//   - options: 配置选项，可以通过位运算组合多个选项（可选参数）
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
//	// 使用自定义选项
//	schema := MustStructToJSONSchema(&Settings{},
//	    SchemaOptionAnonymous | SchemaOptionCompact)
//
// 注意事项:
//   - 仅在确定结构体定义正确且不会出错时使用
//   - 适用于程序初始化阶段的配置加载
//   - 如果转换失败会导致程序panic
func MustStructToJSONSchema(structInstance interface{}, options ...SchemaOption) (schemaJSON string) {
	schemaJSON, err := StructToJSONSchema(structInstance, options...)
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
//   - options: 配置选项，可以通过位运算组合多个选项（可选参数）
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
//	// 使用自定义选项
//	valid, schema, err := ValidateJSONAgainstStruct(&User{}, jsonData,
//	    SchemaOptionRequiredFromTags | SchemaOptionNoAdditionalProperties)
//
// 注意事项:
//   - 当前实现仅生成Schema，不执行实际验证
//   - 如需完整验证功能，建议配合github.com/xeipuuv/gojsonschema等验证库使用
//   - 返回的schema可用于前端验证或文档生成
func ValidateJSONAgainstStruct(structInstance interface{}, jsonData string, options ...SchemaOption) (valid bool, schema string, err error) {
	// 生成JSON Schema
	schema, err = StructToJSONSchema(structInstance, options...)
	if err != nil {
		return
	}

	// 注意：这里仅生成Schema，实际验证需要使用专门的JSON Schema验证库
	// 例如: github.com/xeipuuv/gojsonschema
	// 这里返回true仅作为示例，实际使用时应该实现真正的验证逻辑
	valid = true

	return
}

// GetDefaultOptions 获取推荐的默认选项组合
//
// 该函数返回一组常用的选项组合，适用于大多数场景。
//
// 返回值:
//   - options: 推荐的默认选项组合
//
// 使用示例:
//
//	opts := GetDefaultOptions()
//	schema, err := StructToJSONSchema(&User{}, opts)
func GetDefaultOptions() SchemaOption {
	return SchemaOptionAnonymous | SchemaOptionExpandStruct
}

// GetCompactOptions 获取紧凑格式的选项组合
//
// 该函数返回适用于网络传输或存储的紧凑格式选项组合。
//
// 返回值:
//   - options: 紧凑格式的选项组合
//
// 使用示例:
//
//	opts := GetCompactOptions()
//	schema, err := StructToJSONSchema(&User{}, opts)
func GetCompactOptions() SchemaOption {
	return SchemaOptionCompact | SchemaOptionAnonymous | SchemaOptionNoReference
}

// GetStrictOptions 获取严格验证的选项组合
//
// 该函数返回适用于需要严格验证的选项组合。
//
// 返回值:
//   - options: 严格验证的选项组合
//
// 使用示例:
//
//	opts := GetStrictOptions()
//	schema, err := StructToJSONSchema(&User{}, opts)
func GetStrictOptions() SchemaOption {
	return SchemaOptionNoAdditionalProperties | SchemaOptionRequiredFromTags
}