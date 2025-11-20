# Model Infrax - Go 数据库代码生成工具

[![Go Version](https://img.shields.io/badge/Go-1.25.1+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/LingoJack/model_infrax)

Model Infrax 是一个强大的 Go 语言数据库代码生成工具，能够从数据库结构或 SQL 文件自动生成完整的模型层代码，包括 Entity、DTO、DAO 等文件。支持多种使用模式和框架适配。

## ✨ 特性

- 🚀 **多种使用方式**：支持命令行工具和 Go 库两种使用方式
- 📦 **多种生成模式**：支持从数据库连接或 SQL 文件生成代码
- 📋 **完整代码结构**：自动生成 Entity、DTO、VO、DAO 和工具类
- 🎯 **灵活配置**：支持 YAML 配置文件和 Builder 模式 API
- 🔧 **框架适配**：支持原生 GORM 和 itea-go 框架
- ⚡ **智能优化**：支持索引字段优化和表名前缀处理
- 🛠️ **依赖注入**：使用 Wire 进行依赖注入，代码结构清晰
- 📝 **类型安全**：完整的类型定义和错误处理
- 🔍 **智能配置查找**：自动按优先级查找配置文件

## 🚀 快速开始

### 安装

```bash
# 作为 Go 库使用
go get github.com/LingoJack/model_infrax

# 安装命令行工具
go install github.com/LingoJack/model_infrax/cmd/jen@latest
```

### 基础使用

#### 方式一：使用 Builder 模式 API

```go
package main

import (
    "log"
    "github.com/LingoJack/model_infrax"
)

func main() {
    // 从数据库生成代码
    err := model_infrax.Generate(
        model_infrax.NewBuilder().
            DatabaseMode("localhost", 3306, "mydb", "root", "password").
            AllTables().
            OutputPath("./output").
            IgnoreTableNamePrefix(true).
            UseFramework("itea-go"),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

#### 方式二：使用配置文件

创建 `application.yml` 配置文件：

```yaml
generate_config:
  generate_mode: database
  database_name: mydb
  host: localhost
  port: 3306
  username: root
  password: password
  all_tables: true

generate_option:
  output_path: ./output
  ignore_table_name_prefix: true
  use_framework: itea-go
  package_name:
    po_package: model/entity
    dto_package: model/query
    vo_package: model/view
    dao_package: dao
    tool_package: tool
```

然后使用配置文件生成：

```go
package main

import (
    "log"
    "github.com/LingoJack/model_infrax"
)

func main() {
    err := model_infrax.GenerateFromConfig("./application.yml")
    if err != nil {
        log.Fatal(err)
    }
}
```

## 📖 使用模式

### 1. 数据库模式

直接连接数据库，实时获取表结构：

```go
model_infrax.NewBuilder().
    DatabaseMode("localhost", 3306, "mydb", "root", "password").
    Tables("users", "orders").  // 指定表名
    OutputPath("./model").
    BuildAndGenerate()
```

### 2. SQL 文件模式

从 SQL 建表语句生成代码，无需数据库连接：

```go
model_infrax.NewBuilder().
    StatementMode("./schema.sql").
    AllTables().
    OutputPath("./model").
    BuildAndGenerate()
```

## ⚙️ 配置选项

### Builder API 完整配置

```go
model_infrax.NewBuilder().
    // 生成模式选择
    DatabaseMode("host", port, "db", "user", "pass").  // 数据库模式
    // StatementMode("./schema.sql").                   // SQL文件模式
    
    // 表选择
    AllTables().                                      // 所有表
    // Tables("users", "orders").                     // 指定表
    
    // 输出配置
    OutputPath("./output").                          // 输出路径
    IgnoreTableNamePrefix(true).                     // 忽略表名前缀
    CrudOnlyIdx(true).                               // 只为索引字段生成CRUD
    ModelAllInOneFile(true, "models.go").           // 合并到一个文件
    
    // 框架和包配置
    UseFramework("itea-go").                        // 使用框架
    Packages("po", "dto", "vo", "dao", "tool").      // 配置包名
    
    BuildAndGenerate()                               // 构建并生成
```

### 配置文件完整选项

```yaml
generate_config:
  # 生成模式: database 或 statement
  generate_mode: database
  
  # database 模式配置
  database_name: mydb
  host: localhost
  port: 3306
  username: root
  password: password
  
  # statement 模式配置
  sql_file_path: ./schema.sql
  
  # 通用配置
  all_tables: false
  table_names:
    - users
    - orders

generate_option:
  # 输出配置
  output_path: ./output
  ignore_table_name_prefix: false
  crud_only_idx: false
  all_model_in_one_file: false
  all_model_in_one_file_name: model.go
  
  # 框架配置
  use_framework: ""  # 留空为原生GORM，支持 "itea-go"
  
  # 包名配置
  package_name:
    po_package: model/entity
    dto_package: model/query
    vo_package: model/view
    dao_package: dao
    tool_package: tool
```

## 📁 生成的代码结构

```
output/
├── model/
│   ├── entity/           # 数据库实体 (PO)
│   │   └── user.go
│   ├── query/            # 查询对象 (DTO)
│   │   └── user_dto.go
│   └── view/             # 视图对象 (VO)
│       └── user_vo.go
├── dao/                  # 数据访问层
│   └── user_dao.go
└── tool/                 # 工具类
    ├── copy.go           # 对象复制工具
    ├── encode.go         # 编码工具
    ├── ptr.go           # 指针工具
    └── str.go           # 字符串工具
```

## 🎯 支持的框架

### 原生 GORM
生成标准的 GORM 模型和查询方法：

```go
// 生成的实体示例
type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"size:100"`
    Age  int
}

// 生成的查询方法
func (d *UserDAO) FindByID(id uint) (*entity.User, error) {
    var user entity.User
    err := d.db.First(&user, id).Error
    return &user, err
}
```

### itea-go 框架
生成适配 itea-go 框架的代码，包含特定的注解和工具方法。

## 📚 示例项目

查看 [`examples/`](examples/) 目录获取更多使用示例：

- [基础示例](examples/basic/) - 使用配置文件
- [数据库模式](examples/database/) - 从数据库生成
- [SQL文件模式](examples/statement/) - 从SQL文件生成
- [高级用法](examples/advanced/) - 更多配置选项

## 🔧 高级用法

### 批量生成多服务代码

```go
services := map[string][]string{
    "user_service":    {"t_user", "t_role"},
    "order_service":   {"t_order", "t_order_item"},
    "product_service": {"t_product", "t_category"},
}

for service, tables := range services {
    model_infrax.NewBuilder().
        DatabaseMode("localhost", 3306, "mydb", "root", "password").
        Tables(tables...).
        OutputPath("./services/" + service + "/model").
        IgnoreTableNamePrefix(true).
        BuildAndGenerate()
}
```

### 自定义数据库连接模板

```go
model_infrax.NewBuilder().
    DatabaseMode("localhost", 3306, "mydb", "root", "password").
    URLTemplate("mysql://%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local").
    BuildAndGenerate()
```

## 🛠️ 开发

### 环境要求

- Go 1.25.1+
- MySQL 5.7+ (使用 database 模式时)

### 构建项目

```bash
# 构建 API 库
go build .

# 构建命令行工具
go build -o jen ./cmd/jen

# 运行测试
go test ./...

# 生成 Wire 依赖注入代码
go generate ./...

# 安装命令行工具到本地
go install ./cmd/jen
```

### 项目结构

```
model_infrax/
├── api.go              # 对外 API 接口
├── cmd/                # 命令行工具
│   └── jen/           # jen 命令行工具
│       ├── main.go     # 主入口文件
│       ├── wire.go     # Wire 依赖注入配置
│       └── wire_gen.go # Wire 自动生成的代码
├── config/             # 配置管理
├── examples/           # 使用示例
├── generator/          # 代码生成器
├── model/              # 数据模型
├── parser/             # 数据库解析器
├── pkg/                # 应用核心
├── tool/               # 工具类
└── assets/             # 资源文件
```

## 🔧 命令行工具 (jen)

`jen` 是 Model Infrax 的命令行工具，提供了便捷的命令行接口。

### 安装命令行工具

```bash
go install github.com/LingoJack/model_infrax/cmd/jen@latest
```

### 命令行参数

```bash
jen [flags]

Flags:
  -c, --config string   配置文件路径 (默认: "./application.yml")
  -h, --help           显示帮助信息
```

### 使用示例

```bash
# 使用默认配置文件
jen

# 使用自定义配置文件
jen -c ./config/my-app.yml

# 使用绝对路径
jen --config /etc/jen/config.yml
```

### 配置文件优先级

如果不指定配置文件路径，工具会按以下顺序查找配置文件：
1. `./application.yml`
2. `./assets/application.yml` 
3. `/Applications/jen/application.yml`
4. `/Applications/jen/assets/application.yml`

找到第一个可用配置文件后就会使用它。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Wire](https://github.com/google/wire) - 依赖注入
- [GORM](https://gorm.io/) - ORM 框架
- [TiDB Parser](https://github.com/pingcap/tidb) - SQL 解析器

## 📞 联系方式

- 作者: LingoJack
- 项目地址: [https://github.com/LingoJack/model_infrax](https://github.com/LingoJack/model_infrax)

---

⭐ 如果这个项目对你有帮助，请给个 Star！