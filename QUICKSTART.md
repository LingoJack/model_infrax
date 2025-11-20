# 🚀 快速开始指南

## 📦 安装

```bash
# 方式1: 安装命令行工具
go install github.com/LingoJack/model_infrax/cmd/jen@latest

# 方式2: 作为库使用
go get github.com/LingoJack/model_infrax@latest
```

## 💻 使用方式

### 1️⃣ 命令行模式

```bash
# 使用配置文件
jen -c ./application.yml

# 使用默认配置
jen
```

### 2️⃣ 编程模式 - 数据库

```go
import "github.com/LingoJack/model_infrax"

model_infrax.NewBuilder().
    DatabaseMode("localhost", 3306, "mydb", "root", "password").
    AllTables().
    OutputPath("./output").
    BuildAndGenerate()
```

### 3️⃣ 编程模式 - SQL文件

```go
model_infrax.NewBuilder().
    StatementMode("./schema.sql").
    Tables("t_user", "t_order").
    OutputPath("./output").
    IgnoreTableNamePrefix(true).
    BuildAndGenerate()
```

### 4️⃣ 配置文件模式

```go
model_infrax.GenerateFromConfig("./application.yml")
```

## 🎯 常用配置

```go
builder := model_infrax.NewBuilder().
    // 数据源配置（二选一）
    DatabaseMode("host", port, "db", "user", "pass").  // 从数据库
    // StatementMode("./schema.sql").                  // 从SQL文件
    
    // 表选择（二选一）
    AllTables().                                        // 所有表
    // Tables("t_user", "t_order").                    // 指定表
    
    // 输出配置
    OutputPath("./output").                             // 输出路径
    
    // 生成选项
    IgnoreTableNamePrefix(true).                       // 去掉表名前缀
    CrudOnlyIdx(true).                                 // 只为索引生成CRUD
    ModelAllInOneFile(false, "").                      // 每表一个文件
    UseFramework("itea-go").                           // 使用框架
    
    // 包名配置
    Packages("po", "dto", "vo", "dao", "tool")         // 批量设置
```

## 📚 示例代码

查看 `examples/` 目录获取更多示例：

- `examples/basic/` - 基础使用
- `examples/database/` - 数据库模式
- `examples/statement/` - SQL文件模式
- `examples/advanced/` - 高级用法

## 🔧 配置文件示例

```yaml
generate_config:
  generate_mode: database  # 或 statement
  
  # database 模式
  host: localhost
  port: 3306
  database_name: mydb
  username: root
  password: password
  
  # statement 模式
  # sql_file_path: ./schema.sql
  
  # 表选择
  all_tables: false
  table_names:
    - t_user
    - t_order

generate_option:
  output_path: ./output
  ignore_table_name_prefix: true
  crud_only_idx: false
  all_model_in_one_file: false
  use_framework: itea-go
  
  package_name:
    po_package: model/entity
    dto_package: model/query
    vo_package: model/view
    dao_package: dao
    tool_package: tool
```

## 🆘 常见问题

### Q: 如何只生成指定的表？

```go
.Tables("t_user", "t_order", "t_product")
```

### Q: 如何去掉表名前缀？

```go
.IgnoreTableNamePrefix(true)  // t_user -> User
```

### Q: 如何合并到一个文件？

```go
.ModelAllInOneFile(true, "models.go")
```

### Q: 如何使用环境变量？

```go
host := os.Getenv("DB_HOST")
port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
// ...
```

## 📖 完整文档

- [README.md](./README.md) - 完整文档
- [README_API.md](./README_API.md) - API参考
- [MIGRATION.md](./MIGRATION.md) - 迁移指南
- [examples/README.md](./examples/README.md) - 示例说明

## 🎉 快速测试

```bash
# 1. 克隆项目
git clone https://github.com/LingoJack/model_infrax.git
cd model_infrax

# 2. 运行示例
cd examples/basic
go run main.go

# 3. 查看生成的代码
ls -la ./output
```

---

**提示**: 更多详细信息请查看完整文档 📚