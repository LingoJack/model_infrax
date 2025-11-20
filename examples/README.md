# Examples - 使用示例

本目录包含了 `model_infrax` 的各种使用示例，帮助你快速上手。

## 📁 目录结构

```
examples/
├── basic/          # 基础示例 - 使用配置文件
├── database/       # 数据库模式 - 从数据库生成代码
├── statement/      # SQL文件模式 - 从SQL文件生成代码
└── advanced/       # 高级用法 - 更多配置选项
```

## 🚀 快速开始

### 1. 基础示例 - 使用配置文件

最简单的使用方式，适合快速开始：

```bash
cd examples/basic
go run main.go
```

**特点：**
- ✅ 配置简单，只需一个YAML文件
- ✅ 适合团队协作，配置文件可以版本控制
- ✅ 支持所有配置选项

### 2. 数据库模式 - 从数据库生成代码

从现有数据库读取表结构生成代码：

```bash
cd examples/database
go run main.go
```

**特点：**
- ✅ 直接连接数据库，实时获取最新表结构
- ✅ 支持生成所有表或指定表
- ✅ 适合已有数据库的项目

**注意：** 需要修改代码中的数据库连接信息。

### 3. SQL文件模式 - 从SQL文件生成代码

从SQL建表语句生成代码，无需数据库连接：

```bash
cd examples/statement
go run main.go
```

**特点：**
- ✅ 不需要数据库连接，速度更快
- ✅ 适合开发初期，数据库还未搭建的场景
- ✅ 可以从设计文档直接生成代码

### 4. 高级用法 - 更多配置选项

展示更多高级配置和灵活用法：

```bash
cd examples/advanced
go run main.go
```

**特点：**
- ✅ 自定义配置构建器
- ✅ 批量生成多个数据库
- ✅ 更灵活的配置方式

## 📝 配置说明

### Builder 模式 API

所有示例都使用 Builder 模式进行配置，支持链式调用：

```go
model_infrax.NewBuilder().
    DatabaseMode("host", port, "dbname", "user", "pass").  // 数据库模式
    // 或
    StatementMode("./schema.sql").                         // SQL文件模式
    
    AllTables().                                           // 生成所有表
    // 或
    Tables("t_user", "t_order").                          // 指定表
    
    OutputPath("./output").                                // 输出路径
    IgnoreTableNamePrefix(true).                          // 忽略表名前缀
    CrudOnlyIdx(true).                                    // 只为索引字段生成CRUD
    ModelAllInOneFile(true, "models.go").                 // 合并到一个文件
    UseFramework("itea-go").                              // 使用框架
    Packages("po", "dto", "vo", "dao", "tool").           // 配置包名
    BuildAndGenerate()                                     // 构建并生成
```

### 配置文件方式

也可以使用YAML配置文件：

```go
model_infrax.GenerateFromConfig("./application.yml")
```

配置文件示例请参考：[../assets/application.yml](../assets/application.yml)

## 🔧 常见场景

### 场景1: 新项目，从SQL设计文档生成代码

```go
model_infrax.NewBuilder().
    StatementMode("./design/schema.sql").
    AllTables().
    OutputPath("./internal/model").
    IgnoreTableNamePrefix(true).
    UseFramework("itea-go").
    BuildAndGenerate()
```

### 场景2: 已有数据库，生成指定表的代码

```go
model_infrax.NewBuilder().
    DatabaseMode("localhost", 3306, "mydb", "root", "password").
    Tables("t_user", "t_order", "t_product").
    OutputPath("./model").
    IgnoreTableNamePrefix(true).
    CrudOnlyIdx(true).
    BuildAndGenerate()
```

### 场景3: 微服务项目，批量生成多个服务的代码

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

## 📚 更多资源

- **[主文档](../README.md)** - 完整的使用文档
- **[API文档](../README_API.md)** - API参考
- **[配置示例](../assets/application.yml)** - YAML配置示例

## 💡 提示

1. **数据库连接**: 使用数据库模式时，确保数据库可访问
2. **输出路径**: 建议使用相对路径，支持 `~` 表示用户目录
3. **表名前缀**: 如果表名有统一前缀（如 `t_`），建议开启 `IgnoreTableNamePrefix`
4. **框架选择**: 
   - 不指定框架：生成GORM原生代码
   - `itea-go`：生成适配itea-go框架的代码

## 🤝 贡献

欢迎提交更多使用示例！