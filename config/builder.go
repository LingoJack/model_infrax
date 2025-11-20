package config

import (
	"fmt"

	"github.com/LingoJack/model_infrax/tool"
)

// ConfiggerBuilder 配置构建器，提供链式API用于通过Go代码配置生成器
// 使用方式类似 Wire，但更加灵活和类型安全
//
// 示例:
//
//	cfg := config.NewBuilder().
//	    DatabaseMode("localhost", 3306, "mydb", "root", "password").
//	    OutputPath("./output").
//	    Tables("users", "orders").
//	    Build()
type ConfiggerBuilder struct {
	config *Configger
}

// NewBuilder 创建一个新的配置构建器
// 返回一个带有默认值的构建器实例
func NewBuilder() *ConfiggerBuilder {
	return &ConfiggerBuilder{
		config: &Configger{
			GenerateConfig: GenerateConfig{
				GenerateMode: "database", // 默认从数据库生成
				URLTemplate:  "mysql://%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				AllTables:    false,
				TableNames:   []string{},
			},
			GenerateOption: GenerateOption{
				OutputPath:            "./output",
				IgnoreTableNamePrefix: false,
				CrudOnlyIdx:           false,
				ModelAllInOneFile:     false,
				ModelAllInOneFileName: "model.go",
				UseFramework:          "",
				Package: PackageConfig{
					PoPackage:   "po",
					DtoPackage:  "dto",
					VoPackage:   "vo",
					DaoPackage:  "dao",
					ToolPackage: "tool",
				},
			},
		},
	}
}

// DatabaseMode 配置从数据库生成模式
// host: 数据库主机地址
// port: 数据库端口
// dbName: 数据库名称
// username: 数据库用户名
// password: 数据库密码
func (b *ConfiggerBuilder) DatabaseMode(host string, port int, dbName, username, password string) *ConfiggerBuilder {
	b.config.GenerateConfig.GenerateMode = "database"
	b.config.GenerateConfig.Host = host
	b.config.GenerateConfig.Port = port
	b.config.GenerateConfig.DatabaseName = dbName
	b.config.GenerateConfig.Username = username
	b.config.GenerateConfig.Password = password
	return b
}

// StatementMode 配置从SQL文件生成模式
// sqlFilePath: SQL文件路径，支持 ~ 符号表示用户目录
func (b *ConfiggerBuilder) StatementMode(sqlFilePath string) *ConfiggerBuilder {
	b.config.GenerateConfig.GenerateMode = "statement"
	b.config.GenerateConfig.SqlFilePath = tool.EscapeHomeDir(sqlFilePath)
	return b
}

// URLTemplate 自定义数据库连接URL模板
// template: URL模板字符串，例如: "mysql://%s:%s@tcp(%s:%d)/%s?charset=utf8mb4"
func (b *ConfiggerBuilder) URLTemplate(template string) *ConfiggerBuilder {
	b.config.GenerateConfig.URLTemplate = template
	return b
}

// AllTables 配置生成所有表
// 如果设置为true，将忽略 Tables() 方法设置的表名列表
func (b *ConfiggerBuilder) AllTables() *ConfiggerBuilder {
	b.config.GenerateConfig.AllTables = true
	b.config.GenerateConfig.TableNames = []string{}
	return b
}

// Tables 配置需要生成的表名列表
// tableNames: 表名列表，可以传入多个表名
func (b *ConfiggerBuilder) Tables(tableNames ...string) *ConfiggerBuilder {
	b.config.GenerateConfig.AllTables = false
	b.config.GenerateConfig.TableNames = tableNames
	return b
}

// OutputPath 配置代码输出路径
// path: 输出路径，支持 ~ 符号表示用户目录
func (b *ConfiggerBuilder) OutputPath(path string) *ConfiggerBuilder {
	b.config.GenerateOption.OutputPath = tool.EscapeHomeDir(path)
	return b
}

// IgnoreTableNamePrefix 配置是否忽略表名前缀
// 如果设置为true，生成的类名将去除表名前缀（如 t_user -> User）
func (b *ConfiggerBuilder) IgnoreTableNamePrefix(ignore bool) *ConfiggerBuilder {
	b.config.GenerateOption.IgnoreTableNamePrefix = ignore
	return b
}

// CrudOnlyIdx 配置是否仅对索引字段生成CRUD方法
// 如果设置为true，只会为有索引的字段生成查询方法
func (b *ConfiggerBuilder) CrudOnlyIdx(onlyIdx bool) *ConfiggerBuilder {
	b.config.GenerateOption.CrudOnlyIdx = onlyIdx
	return b
}

// ModelAllInOneFile 配置是否将所有Model生成到一个文件
// allInOne: 是否合并到一个文件
// fileName: 文件名（当allInOne为true时有效）
func (b *ConfiggerBuilder) ModelAllInOneFile(allInOne bool, fileName string) *ConfiggerBuilder {
	b.config.GenerateOption.ModelAllInOneFile = allInOne
	if fileName != "" {
		b.config.GenerateOption.ModelAllInOneFileName = fileName
	}
	return b
}

// UseFramework 配置使用的框架
// framework: 框架名称，例如 "itea-go"
func (b *ConfiggerBuilder) UseFramework(framework string) *ConfiggerBuilder {
	b.config.GenerateOption.UseFramework = framework
	return b
}

// Packages 配置生成代码的包名
// po: PO（持久化对象）包名
// dto: DTO（数据传输对象）包名
// vo: VO（视图对象）包名
// dao: DAO（数据访问对象）包名
// tool: Tool（工具类）包名
func (b *ConfiggerBuilder) Packages(po, dto, vo, dao, tool string) *ConfiggerBuilder {
	if po != "" {
		b.config.GenerateOption.Package.PoPackage = po
	}
	if dto != "" {
		b.config.GenerateOption.Package.DtoPackage = dto
	}
	if vo != "" {
		b.config.GenerateOption.Package.VoPackage = vo
	}
	if dao != "" {
		b.config.GenerateOption.Package.DaoPackage = dao
	}
	if tool != "" {
		b.config.GenerateOption.Package.ToolPackage = tool
	}
	return b
}

// PoPackage 配置PO包名
func (b *ConfiggerBuilder) PoPackage(pkg string) *ConfiggerBuilder {
	b.config.GenerateOption.Package.PoPackage = pkg
	return b
}

// DtoPackage 配置DTO包名
func (b *ConfiggerBuilder) DtoPackage(pkg string) *ConfiggerBuilder {
	b.config.GenerateOption.Package.DtoPackage = pkg
	return b
}

// VoPackage 配置VO包名
func (b *ConfiggerBuilder) VoPackage(pkg string) *ConfiggerBuilder {
	b.config.GenerateOption.Package.VoPackage = pkg
	return b
}

// DaoPackage 配置DAO包名
func (b *ConfiggerBuilder) DaoPackage(pkg string) *ConfiggerBuilder {
	b.config.GenerateOption.Package.DaoPackage = pkg
	return b
}

// ToolPackage 配置Tool包名
func (b *ConfiggerBuilder) ToolPackage(pkg string) *ConfiggerBuilder {
	b.config.GenerateOption.Package.ToolPackage = pkg
	return b
}

// Build 构建最终的配置对象
// 返回构建好的 Configger 实例和可能的错误
func (b *ConfiggerBuilder) Build() (*Configger, error) {
	// 验证配置的有效性
	if err := b.validate(); err != nil {
		return nil, err
	}
	return b.config, nil
}

// validate 验证配置的有效性
func (b *ConfiggerBuilder) validate() error {
	cfg := b.config

	// 验证生成模式
	if cfg.GenerateConfig.GenerateMode != "database" && cfg.GenerateConfig.GenerateMode != "statement" {
		return fmt.Errorf("无效的生成模式: %s，必须是 'database' 或 'statement'", cfg.GenerateConfig.GenerateMode)
	}

	// 验证数据库模式的必需参数
	if cfg.GenerateConfig.GenerateMode == "database" {
		if cfg.GenerateConfig.Host == "" {
			return fmt.Errorf("数据库模式下必须指定 Host")
		}
		if cfg.GenerateConfig.Port == 0 {
			return fmt.Errorf("数据库模式下必须指定 Port")
		}
		if cfg.GenerateConfig.DatabaseName == "" {
			return fmt.Errorf("数据库模式下必须指定 DatabaseName")
		}
		if cfg.GenerateConfig.Username == "" {
			return fmt.Errorf("数据库模式下必须指定 Username")
		}
	}

	// 验证SQL文件模式的必需参数
	if cfg.GenerateConfig.GenerateMode == "statement" {
		if cfg.GenerateConfig.SqlFilePath == "" {
			return fmt.Errorf("SQL文件模式下必须指定 SqlFilePath")
		}
	}

	// 验证表名配置
	if !cfg.GenerateConfig.AllTables && len(cfg.GenerateConfig.TableNames) == 0 {
		return fmt.Errorf("必须指定要生成的表名或使用 AllTables()")
	}

	// 验证输出路径
	if cfg.GenerateOption.OutputPath == "" {
		return fmt.Errorf("必须指定输出路径")
	}

	return nil
}

// MustBuild 构建配置对象，如果出错则panic
// 适用于确定配置正确的场景
func (b *ConfiggerBuilder) MustBuild() *Configger {
	cfg, err := b.Build()
	if err != nil {
		panic(fmt.Sprintf("构建配置失败: %v", err))
	}
	return cfg
}
