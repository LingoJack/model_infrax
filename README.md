# Model Infrax - GORM 模型生成器

这是一个基于数据库表结构自动生成 GORM 模型的工具。

## 功能特性

- 🚀 自动从数据库读取表结构
- 📝 生成符合 GORM 规范的 Go 结构体
- 🎯 支持自定义模板
- 🔧 灵活的配置选项
- 📦 使用 Wire 进行依赖注入

## 快速开始

### 1. 配置数据库连接

编辑 `assert/application.yml` 文件：

```yaml
generate_config:
  generate_mode: database
  database_name: test_db
  host: localhost
  port: 3306
  username: root
  password: your_password
  all_tables: false
  table_names:
    - t_artifact

generate_option:
  output_path: ~/dev/model_infrax/output
  crud_only_idx: false
```

### 2. 运行生成器

```bash
# 生成 Wire 依赖注入代码
go generate ./...

# 运行程序
go run .
```

### 3. 查看生成的代码

生成的模型文件将保存在配置的 `output_path` 目录下。

## 项目结构

```
model_infrax/
├── assert/
│   ├── application.yml      # 配置文件
│   ├── database.sql          # 测试数据库脚本
│   └── template/
│       └── model.template    # GORM 模型模板
├── config/                   # 配置管理
├── generator/                # 代码生成器
│   ├── generator.go          # 生成器主逻辑
│   ├── template_func.go      # 模板函数
│   └── template_func_test.go # 单元测试
├── model/                    # 数据模型定义
├── parser/                   # SQL 解析器
├── tool/                     # 工具函数
├── main.go                   # 程序入口
├── wire.go                   # Wire 依赖注入配置
└── wire_gen.go               # Wire 生成的代码

```

## 模板说明

模板文件位于 `assert/template/model.template`，使用 Go 的 `text/template` 语法。

### 可用的模板函数

- `ToPascalCase`: 将字符串转换为 PascalCase（大驼峰）
  - 例如: `t_artifact` -> `TArtifact`
  - 例如: `artifactId` -> `ArtifactID`

- `GetGoType`: 根据列信息返回对应的 Go 类型
  - 自动识别 ID、时间、整数等类型
  - 支持可空类型（指针类型）

- `GetMySQLType`: 根据列信息返回对应的 MySQL 类型
  - 自动推断合适的数据库类型

### 模板示例

```go
{{- range . }}
type {{ .Name | ToPascalCase }} struct {
{{- range .Columns }}
	{{ .ColumnName | ToPascalCase }} {{ . | GetGoType }} `gorm:"column:{{ .ColumnName }};type:{{ . | GetMySQLType }};comment:{{ .Comment }}" json:"{{ .ColumnName }}"`
{{- end }}
}

func (t *{{ .Name | ToPascalCase }}) TableName() string {
	return "{{ .Name }}"
}
{{- end }}
```

## 生成的代码示例

```go
package model

import (
	"time"
)

// TArtifact 任务执行流程中生成的中间产物表
type TArtifact struct {
	ID           uint64    `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement;comment:主键ID;not null" json:"id"`
	ArtifactID   string    `gorm:"column:artifactId;type:varchar(128);comment:产物ID;not null" json:"artifactId"`
	ArtifactName string    `gorm:"column:artifactName;type:varchar(128);comment:产物名称;not null" json:"artifactName"`
	SessionID    string    `gorm:"column:sessionId;type:varchar(128);comment:所属的会话;not null" json:"sessionId"`
	Step         int       `gorm:"column:step;type:int(11);comment:大的步骤点;not null" json:"step"`
	SubStep      string    `gorm:"column:subStep;type:varchar(128);comment:小的步骤点;not null" json:"subStep"`
	Content      *string   `gorm:"column:content;type:text;comment:内容" json:"content"`
	Version      *string   `gorm:"column:version;type:varchar(128);comment:版本" json:"version"`
	CreateTime   time.Time `gorm:"column:createTime;type:datetime;comment:创建时间;not null" json:"createTime"`
	UpdateTime   time.Time `gorm:"column:updateTime;type:datetime;comment:更新时间;not null" json:"updateTime"`
}

// TableName 返回表名
func (t *TArtifact) TableName() string {
	return "t_artifact"
}
```

## 开发指南

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./generator/...

# 运行测试并显示覆盖率
go test -cover ./...
```

### 添加新的模板函数

1. 在 `generator/template_func.go` 中添加函数
2. 在 `generator/generator.go` 的 `FuncMap` 中注册函数
3. 在模板中使用新函数

### 自定义模板

你可以修改 `assert/template/model.template` 来自定义生成的代码格式。

## 依赖项

- [GORM](https://gorm.io/) - ORM 库
- [Wire](https://github.com/google/wire) - 依赖注入
- [lo](https://github.com/samber/lo) - 函数式编程工具
- [yaml.v3](https://github.com/go-yaml/yaml) - YAML 解析

## 许可证

MIT License