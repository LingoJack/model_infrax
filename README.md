# Model Infrax

一个强大的 Go 代码生成工具，支持从数据库或 SQL 文件生成 Model、DAO、DTO、VO 等代码。

## ✨ 特性

- 🚀 **多种使用方式**: 支持命令行工具和编程式 API
- 🎯 **类型安全**: 使用 Go 代码配置，编译时检查
- 🔧 **灵活配置**: 支持 YAML 配置文件和 Builder 模式
- 📦 **易于分发**: 可通过 `go install` 安装
- 🎨 **框架支持**: 支持多种框架模板（如 itea-go）
- 🔌 **依赖注入**: 使用 Wire 进行依赖注入

## 📦 安装方式

### 方式 1: 安装命令行工具（推荐）

如果你只是想使用命令行工具生成代码，这是最简单的方式：

```bash
# 从 GitHub 安装（需要先发布到 GitHub）
go install github.com/LingoJack/model_infrax/cmd/jen@latest

# 或者从本地安装（开发阶段）
cd /path/to/model_infrax
go install ./cmd/model_infrax
```

安装完成后，可以直接使用 `jen` 命令：

```bash
# 使用配置文件生成代码
jen -c ./application.yml

# 或者使用默认配置文件路径
jen
```

### 方式 2: 作为 Go 库使用

如果你想在自己的 Go 项目中使用这个工具，可以作为库导入：

```bash
# 在你的项目目录下执行
go get github.com/LingoJack/model_infrax@latest

# 或者指定版本
go get github.com/LingoJack/model_infrax@v1.0.0
```

### 方式 3: 克隆源码本地使用

```bash
# 克隆仓库
git clone https://github.com/LingoJack/model_infrax.git
cd model_infrax

# 安装依赖
go mod download

# 直接运行
go run main.go -c ./application.yml

# 或者构建二进制文件
go build -o model_infrax main.go
./model_infrax -c ./application.yml
```

## 🚀 快速开始

### 1. 命令行模式（使用 YAML 配置）

创建配置文件 `application.yml`：

```yaml
generate:
  generate_mode: database
  host: localhost
  port: 3306
  database_name: mydb
  username: root
  password: password
  all_tables: false
  table_names:
    - users
    - orders

option:
  output_path: ./generated
  ignore_table_name_prefix: true
  crud_only_idx: false
```

运行生成：

```bash
model_infrax -c ./application.yml
```

### 2. 编程模式（在 Go 代码中使用）

在你的 Go 项目中创建文件 `generate.go`：

```go
package main

import (
    "log"
    "github.com/LingoJack/model_infrax"
)

func main() {
    // 使用 Builder 模式配置
    err := model_infrax.Generate(
        model_infrax.NewBuilder().
            DatabaseMode("localhost", 3306, "mydb", "root", "password").
            Tables("users", "orders", "products").
            OutputPath("./generated").
            IgnoreTableNamePrefix(true),
    )
    
    if err != nil {
        log.Fatalf("生成失败: %v", err)
    }
    
    log.Println("✅ 代码生成成功！")
}
```

运行：

```bash
go run generate.go
```

## 📖 详细使用示例

### 示例 1: 从数据库生成所有表

```go
package main

import (
    "log"
    "github.com/LingoJack/model_infrax"
)

func main() {
    err := model_infrax.Generate(
        model_infrax.NewBuilder().
            DatabaseMode("localhost", 3306, "mydb", "root", "password").
            AllTables().  // 生成所有表
            OutputPath("./generated"),
    )
    
    if err != nil {
        log.Fatal(err)
    }
}
```

### 示例 2: 从 SQL 文件生成

```go
err := model_infrax.Generate(
    model_infrax.NewBuilder().
        StatementMode("~/schema.sql").  // 从 SQL 文件生成
        AllTables().
        OutputPath("./generated"),
)
```

### 示例 3: 完整配置示例

```go
err := model_infrax.Generate(
    model_infrax.NewBuilder().
        // 数据库配置
        DatabaseMode("localhost", 3306, "mydb", "root", "password").
        
        // 指定要生成的表
        Tables("t_user", "t_order", "t_product").
        
        // 输出配置
        OutputPath("./output").
        
        // 生成选项
        IgnoreTableNamePrefix(true).   // 去掉表名前缀 t_
        CrudOnlyIdx(true).             // 只为索引字段生成 CRUD
        ModelAllInOneFile(false, "").  // 每个表一个文件
        
        // 自定义包名
        Packages("entity", "dto", "vo", "dao", "util").
        
        // 使用框架模板
        UseFramework("itea-go"),
)
```

### 示例 4: 批量生成多个数据库

```go
package main

import (
    "log"
    "github.com/LingoJack/model_infrax"
)

func main() {
    databases := []struct {
        name   string
        tables []string
    }{
        {"user_db", []string{"users", "profiles"}},
        {"order_db", []string{"orders", "order_items"}},
        {"product_db", []string{"products", "categories"}},
    }
    
    for _, db := range databases {
        log.Printf("🚀 生成数据库 %s...", db.name)
        
        err := model_infrax.Generate(
            model_infrax.NewBuilder().
                DatabaseMode("localhost", 3306, db.name, "root", "password").
                Tables(db.tables...).
                OutputPath("./generated/" + db.name),
        )
        
        if err != nil {
            log.Printf("❌ 失败: %v", err)
            continue
        }
        
        log.Printf("✅ 成功")
    }
}
```

### 示例 5: 在 Web 服务中使用

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/LingoJack/model_infrax"
)

type GenerateRequest struct {
    Host     string   `json:"host"`
    Port     int      `json:"port"`
    Database string   `json:"database"`
    Username string   `json:"username"`
    Password string   `json:"password"`
    Tables   []string `json:"tables"`
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
    var req GenerateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    builder := model_infrax.NewBuilder().
        DatabaseMode(req.Host, req.Port, req.Database, req.Username, req.Password).
        OutputPath("./generated")
    
    if len(req.Tables) > 0 {
        builder.Tables(req.Tables...)
    } else {
        builder.AllTables()
    }
    
    err := model_infrax.Generate(builder)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
    })
}

func main() {
    http.HandleFunc("/generate", handleGenerate)
    log.Println("服务启动在 :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## 🔧 配置选项说明

### Builder API 完整列表

| 方法 | 说明 | 示例 |
|------|------|------|
| `DatabaseMode(host, port, db, user, pass)` | 从数据库生成 | `.DatabaseMode("localhost", 3306, "mydb", "root", "pwd")` |
| `StatementMode(sqlFile)` | 从 SQL 文件生成 | `.StatementMode("~/schema.sql")` |
| `AllTables()` | 生成所有表 | `.AllTables()` |
| `Tables(names...)` | 指定表名 | `.Tables("users", "orders")` |
| `OutputPath(path)` | 输出路径 | `.OutputPath("./generated")` |
| `IgnoreTableNamePrefix(bool)` | 忽略表名前缀 | `.IgnoreTableNamePrefix(true)` |
| `CrudOnlyIdx(bool)` | 只为索引生成 CRUD | `.CrudOnlyIdx(true)` |
| `ModelAllInOneFile(bool, name)` | 合并到一个文件 | `.ModelAllInOneFile(true, "models.go")` |
| `UseFramework(name)` | 使用框架模板 | `.UseFramework("itea-go")` |
| `Packages(po, dto, vo, dao, tool)` | 批量设置包名 | `.Packages("entity", "dto", "vo", "dao", "util")` |
| `PoPackage(name)` | 设置 PO 包名 | `.PoPackage("entity")` |
| `DtoPackage(name)` | 设置 DTO 包名 | `.DtoPackage("dto")` |
| `VoPackage(name)` | 设置 VO 包名 | `.VoPackage("vo")` |
| `DaoPackage(name)` | 设置 DAO 包名 | `.DaoPackage("dao")` |
| `ToolPackage(name)` | 设置 Tool 包名 | `.ToolPackage("util")` |

## 🌍 环境变量配置

可以使用环境变量来管理敏感信息：

```go
package main

import (
    "log"
    "os"
    "strconv"
    "github.com/LingoJack/model_infrax"
)

func main() {
    // 从环境变量读取配置
    host := os.Getenv("DB_HOST")
    port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
    dbName := os.Getenv("DB_NAME")
    user := os.Getenv("DB_USER")
    pass := os.Getenv("DB_PASS")
    
    err := model_infrax.Generate(
        model_infrax.NewBuilder().
            DatabaseMode(host, port, dbName, user, pass).
            AllTables().
            OutputPath("./generated"),
    )
    
    if err != nil {
        log.Fatal(err)
    }
}
```

使用时设置环境变量：

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=mydb
export DB_USER=root
export DB_PASS=password

go run generate.go
```

## 🆚 对比：命令行 vs 编程式

| 特性 | 命令行模式 | 编程式模式 |
|------|-----------|-----------|
| 配置方式 | YAML 文件 | Go 代码 |
| 类型安全 | ❌ | ✅ |
| IDE 支持 | ❌ | ✅ |
| 动态配置 | ❌ | ✅ |
| 适用场景 | 独立使用 | 集成到应用 |

## 🏗️ 项目结构

```
model_infrax/
├── cmd/
│   └── model_infrax/      # 命令行工具入口
│       └── main.go
├── config/                # 配置模块
│   ├── config.go
│   └── builder.go         # Builder 模式配置
├── parser/                # 解析器模块
├── generator/             # 代码生成器模块
├── examples/              # 使用示例
│   └── programmatic_usage.go
├── api.go                 # 对外暴露的 API
├── main.go                # 原始入口（保留兼容）
├── wire.go                # Wire 依赖注入配置
├── README_API.md          # API 文档
└── README.md              # 本文件
```

## 🔍 使用场景

### 1. 微服务开发
为多个微服务批量生成数据访问层代码

### 2. CI/CD 集成
在构建流程中自动生成代码

### 3. Web 服务
提供代码生成 API 服务

### 4. 开发工具
集成到 IDE 插件或开发工具中

## 🐛 常见问题

### Q1: 如何安装到全局？

```bash
# 方式 1: 使用 go install
go install github.com/LingoJack/model_infrax/cmd/model_infrax@latest

# 方式 2: 手动构建并移动
go build -o model_infrax ./cmd/model_infrax
sudo mv model_infrax /usr/local/bin/
```

### Q2: 如何指定 Go 版本？

在 `go.mod` 中已经指定了 Go 1.25.1，确保你的 Go 版本 >= 1.25.1：

```bash
go version
```

### Q3: 如何更新到最新版本？

```bash
# 更新命令行工具
go install github.com/LingoJack/model_infrax/cmd/model_infrax@latest

# 更新库依赖
go get -u github.com/LingoJack/model_infrax@latest
go mod tidy
```

### Q4: 如何在 CI/CD 中使用？

在 `.gitlab-ci.yml` 或 `.github/workflows/generate.yml` 中：

```yaml
generate:
  stage: build
  script:
    - go install github.com/LingoJack/model_infrax/cmd/model_infrax@latest
    - model_infrax -c ./application.yml
  artifacts:
    paths:
      - generated/
```

### Q5: 如何处理私有仓库？

如果你的项目在私有仓库，需要配置 Git 凭证：

```bash
# 配置 Git 使用 SSH
git config --global url."git@github.com:".insteadOf "https://github.com/"

# 或者使用 GOPRIVATE
export GOPRIVATE=github.com/LingoJack/*
```

## 📚 更多资源

- **[API 文档](./README_API.md)** - 完整的 API 参考
- **[示例代码](./examples/)** - 更多使用示例
- **[配置示例](./assert/application.yml)** - YAML 配置示例

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🙏 致谢

- [Wire](https://github.com/google/wire) - 依赖注入框架
- [GORM](https://gorm.io/) - ORM 框架