# Model Infrax - GORM 模型与 DAO 生成器

这是一个基于数据库表结构自动生成完整分层架构代码的工具，包括 GORM 模型、DAO 层、DTO/PO 结构体和工具函数。

## 🌟 功能特性

- 🚀 **自动从数据库读取表结构**：支持 MySQL 数据库连接和表结构解析
- 📝 **生成符合 GORM 规范的 Go 结构体**：包含完整的标签和注释
- 🏗️ **完整的分层架构生成**：
  - **PO (Persistent Object)**：数据库实体模型
  - **DTO (Data Transfer Object)**：查询和传输对象
  - **DAO (Data Access Object)**：完整的数据访问层
- 🎯 **智能索引支持**：根据索引类型自动生成对应的方法
  - 主键索引：生成单条记录查询方法
  - 唯一索引：生成单条记录查询方法  
  - 普通索引：生成列表查询方法
- 🔧 **灵活的配置选项**：支持多种自定义配置
- 🛠️ **工具函数生成**：自动生成指针操作和字符串处理工具
- 📦 **使用 Wire 进行依赖注入**：现代化的依赖管理

## 解析结构示例
```json
{
    "Name": "t_llm_history",
    "Columns": [
        {
            "ColumnName": "id",
            "Collate": "",
            "Comment": "主键ID",
            "Type": "bigint unsigned",
            "Default": null,
            "IsAutoIncrement": true,
            "IsNullable": false,
            "IsIndexed": true,
            "IsUnique": true,
            "IsPrimaryKey": true
        },
        {
            "ColumnName": "model",
            "Collate": "utf8mb4_unicode_ci",
            "Comment": "模型名称",
            "Type": "varchar(128)",
            "Default": "",
            "IsAutoIncrement": false,
            "IsNullable": false,
            "IsIndexed": true,
            "IsUnique": false,
            "IsPrimaryKey": false
        },
        {
            "ColumnName": "input",
            "Collate": "utf8mb4_unicode_ci",
            "Comment": "输入内容",
            "Type": "text",
            "Default": null,
            "IsAutoIncrement": false,
            "IsNullable": false,
            "IsIndexed": false,
            "IsUnique": false,
            "IsPrimaryKey": false
        },
        {
            "ColumnName": "output",
            "Collate": "utf8mb4_unicode_ci",
            "Comment": "输出内容",
            "Type": "text",
            "Default": null,
            "IsAutoIncrement": false,
            "IsNullable": false,
            "IsIndexed": false,
            "IsUnique": false,
            "IsPrimaryKey": false
        },
        {
            "ColumnName": "createTime",
            "Collate": "",
            "Comment": "创建时间",
            "Type": "datetime",
            "Default": "CURRENT_TIMESTAMP",
            "IsAutoIncrement": false,
            "IsNullable": false,
            "IsIndexed": true,
            "IsUnique": false,
            "IsPrimaryKey": false
        },
        {
            "ColumnName": "updateTime",
            "Collate": "",
            "Comment": "更新时间",
            "Type": "datetime",
            "Default": "CURRENT_TIMESTAMP",
            "IsAutoIncrement": false,
            "IsNullable": false,
            "IsIndexed": false,
            "IsUnique": false,
            "IsPrimaryKey": false
        }
    ],
    "Comment": "LLM历史记录表",
    "PrimaryKey": {
        "IndexName": "PRIMARY",
        "Columns": [
            {
                "ColumnName": "id",
                "Collate": "",
                "Comment": "主键ID",
                "Type": "bigint unsigned",
                "Default": null,
                "IsAutoIncrement": true,
                "IsNullable": false,
                "IsIndexed": false,
                "IsUnique": false,
                "IsPrimaryKey": false
            }
        ]
    },
    "UniqueIndex": [
        {
            "IndexName": "PRIMARY",
            "Columns": [
                {
                    "ColumnName": "id",
                    "Collate": "",
                    "Comment": "主键ID",
                    "Type": "bigint unsigned",
                    "Default": null,
                    "IsAutoIncrement": true,
                    "IsNullable": false,
                    "IsIndexed": false,
                    "IsUnique": false,
                    "IsPrimaryKey": false
                }
            ]
        }
    ],
    "Indexes": [
        {
            "IndexName": "idx_model_createTime",
            "Columns": [
                {
                    "ColumnName": "model",
                    "Collate": "utf8mb4_unicode_ci",
                    "Comment": "模型名称",
                    "Type": "varchar(128)",
                    "Default": "",
                    "IsAutoIncrement": false,
                    "IsNullable": false,
                    "IsIndexed": false,
                    "IsUnique": false,
                    "IsPrimaryKey": false
                },
                {
                    "ColumnName": "createTime",
                    "Collate": "",
                    "Comment": "创建时间",
                    "Type": "datetime",
                    "Default": "CURRENT_TIMESTAMP",
                    "IsAutoIncrement": false,
                    "IsNullable": false,
                    "IsIndexed": false,
                    "IsUnique": false,
                    "IsPrimaryKey": false
                }
            ]
        },
        {
            "IndexName": "PRIMARY",
            "Columns": [
                {
                    "ColumnName": "id",
                    "Collate": "",
                    "Comment": "主键ID",
                    "Type": "bigint unsigned",
                    "Default": null,
                    "IsAutoIncrement": true,
                    "IsNullable": false,
                    "IsIndexed": false,
                    "IsUnique": false,
                    "IsPrimaryKey": false
                }
            ]
        }
    ]
}
```

## 🚀 快速开始

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
  all_tables: true
  table_names:
    - t_artifact

generate_option:
  # 输出路径
  output_path: ~/dev/model_infrax/output

  # 是否将所有模型放在一个文件中
  all_model_in_one_file: false

  # 所有模型放在一个文件中的文件名
  all_model_in_one_file_name: model.go

  # 只为有索引的字段生成 infrax 方法
  crud_only_idx: false

  # go 的 package 映射
  package_name:
    po_package: model/entity      # PO 层包名
    dto_package: model/query      # DTO 层包名
    vo_package: model/view        # VO 层包名（预留）
    dao_package: dao              # DAO 层包名
    tool_package: tool            # 工具函数包名
```

### 2. 运行生成器

```bash
# 生成 Wire 依赖注入代码
go generate ./...

# 运行程序
go run .
```

### 3. 查看生成的代码

生成的代码将按以下结构组织：

```
output/
├── dao/                    # DAO 层：数据访问对象
│   ├── t_artifact_dao.go
│   └── t_feedback_task_dao.go
├── model/                  # 模型层
│   ├── entity/            # PO：持久化对象
│   │   ├── t_artifact.go
│   │   └── t_feedback_task.go
│   ├── query/             # DTO：数据传输对象
│   │   ├── t_artifact_dto.go
│   │   └── t_feedback_task_dto.go
│   └── view/              # VO：视图对象（预留）
└── tool/                  # 工具函数
    ├── ptr.go             # 指针操作工具
    └── str.go             # 字符串处理工具
```

## 📁 项目结构

```
model_infrax/
├── assert/                      # 配置和模板文件
│   ├── application.yml         # 配置文件
│   ├── database.sql            # 测试数据库脚本
│   └── template/               # 代码模板
│       ├── dao.template        # DAO 层模板
│       ├── dto.template        # DTO 结构体模板
│       ├── po.template         # PO 结构体模板
│       └── tools/              # 工具函数模板
│           ├── ptr.template    # 指针工具模板
│           └── str.template    # 字符串工具模板
├── config/                     # 配置管理
├── generator/                  # 代码生成器
│   ├── generator.go           # 生成器主逻辑
│   ├── template_func.go       # 模板函数
│   └── template_func_test.go  # 单元测试
├── model/                      # 数据模型定义
├── parser/                     # SQL 解析器
├── tool/                       # 工具函数
├── main.go                     # 程序入口
├── wire.go                     # Wire 依赖注入配置
└── wire_gen.go                 # Wire 生成的代码
```

## 🎨 模板系统

### 模板文件说明

- **`po.template`**：生成 GORM 实体模型，包含完整的标签和注释
- **`dto.template`**：生成查询和传输用的 DTO 结构体
- **`dao.template`**：生成完整的数据访问层，包含 CRUD 操作
- **`tools/ptr.template`**：生成指针操作工具函数
- **`tools/str.template`**：生成字符串处理工具函数

### 可用的模板函数

- `ToPascalCase`: 将字符串转换为 PascalCase（大驼峰）
  - 例如: `t_artifact` -> `TArtifact`
  - 例如: `artifactId` -> `ArtifactID`

- `GetGoType`: 根据列信息返回对应的 Go 类型
  - 自动识别 ID、时间、整数等类型
  - 支持可空类型（指针类型）

- `GetMySQLType`: 根据列信息返回对应的 MySQL 类型
  - 自动推断合适的数据库类型

## 💻 生成的代码示例

### PO 实体模型

```go
package entity

import (
    "time"
)

// TArtifact 任务执行流程中生成的中间产物表
type TArtifact struct {
    ID           uint64    `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement;comment:主键ID;not null" json:"id"`
    ArtifactID   string    `gorm:"column:artifactId;type:varchar(128);comment:产物ID;not null" json:"artifactId"`
    ArtifactName string    `gorm:"column:artifactName;type:varchar(128);comment:产物名称;not null" json:"artifactName"`
    SessionID    string    `gorm:"column:sessionId;type:varchar(128);comment:所属的会话;not null" json:"sessionId"`
    Content      *string   `gorm:"column:content;type:text;comment:内容" json:"content"`
    CreateTime   time.Time `gorm:"column:createTime;type:datetime;comment:创建时间;not null" json:"createTime"`
    UpdateTime   time.Time `gorm:"column:updateTime;type:datetime;comment:更新时间;not null" json:"updateTime"`
}

// TableName 返回表名
func (t *TArtifact) TableName() string {
    return "t_artifact"
}
```

### DAO 数据访问层

```go
package dao

import (
    "context"
    "your_project/model/entity"
    "your_project/model/query"
    "gorm.io/gorm"
)

// TArtifactDAO 数据访问对象
type TArtifactDAO struct {
    db *gorm.DB
}

// NewTArtifactDAO 创建 DAO 实例
func NewTArtifactDAO(db *gorm.DB) *TArtifactDAO {
    return &TArtifactDAO{db: db}
}

// SelectById 根据 ID 查询单条记录
func (dao *TArtifactDAO) SelectById(ctx context.Context, id uint64) (*entity.TArtifact, error) {
    var result entity.TArtifact
    err := dao.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
    return &result, err
}

// SelectList 根据条件查询列表
func (dao *TArtifactDAO) SelectList(ctx context.Context, dto *query.TArtifactDTO, options ...*query.TArtifactQueryOptions) ([]*entity.TArtifact, error) {
    db := dao.db.WithContext(ctx).Model(&entity.TArtifact{})
    
    // 构建 WHERE 条件
    if dto != nil {
        if dto.ArtifactID != nil {
            db = db.Where("artifactId = ?", *dto.ArtifactID)
        }
        // ... 更多条件
    }
    
    var results []*entity.TArtifact
    err := db.Find(&results).Error
    return results, err
}

// Insert 插入记录
func (dao *TArtifactDAO) Insert(ctx context.Context, po *entity.TArtifact) error {
    return dao.db.WithContext(ctx).Create(po).Error
}

// Update 更新记录
func (dao *TArtifactDAO) Update(ctx context.Context, po *entity.TArtifact) error {
    return dao.db.WithContext(ctx).Save(po).Error
}

// Delete 删除记录
func (dao *TArtifactDAO) Delete(ctx context.Context, id uint64) error {
    return dao.db.WithContext(ctx).Delete(&entity.TArtifact{}, id).Error
}
```

### DTO 查询对象

```go
package query

// TArtifactDTO 查询传输对象
type TArtifactDTO struct {
    ID           *uint64  `json:"id,omitempty"`           // 主键ID
    ArtifactID   *string  `json:"artifactId,omitempty"`   // 产物ID
    ArtifactName *string  `json:"artifactName,omitempty"` // 产物名称
    SessionID    *string  `json:"sessionId,omitempty"`    // 所属的会话
    Content      *string  `json:"content,omitempty"`      // 内容
    // ... 更多字段
}

// TArtifactQueryOptions 查询选项
type TArtifactQueryOptions struct {
    OrderBy   string `json:"orderBy,omitempty"`   // 排序字段
    PageSize  int    `json:"pageSize,omitempty"`  // 页面大小
    PageIndex int    `json:"pageIndex,omitempty"` // 页面索引
}
```

## 🔧 配置说明

### generate_config 配置项

| 字段 | 类型 | 说明 |
|------|------|------|
| `generate_mode` | string | 生成模式，目前支持 `database` |
| `database_name` | string | 数据库名称 |
| `host` | string | 数据库主机地址 |
| `port` | int | 数据库端口 |
| `username` | string | 数据库用户名 |
| `password` | string | 数据库密码 |
| `all_tables` | bool | 是否处理所有表 |
| `table_names` | []string | 指定要处理的表名列表 |

### generate_option 配置项

| 字段 | 类型 | 说明 |
|------|------|------|
| `output_path` | string | 输出路径 |
| `all_model_in_one_file` | bool | 是否将所有模型放在一个文件中 |
| `all_model_in_one_file_name` | string | 合并文件时的文件名 |
| `crud_only_idx` | bool | 是否只为有索引的字段生成 CRUD 方法 |
| `package_name` | object | 各层包名映射 |

## 🧪 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./generator/...

# 运行测试并显示覆盖率
go test -cover ./...
```

## 🔨 开发指南

### 添加新的模板函数

1. 在 `generator/template_func.go` 中添加函数
2. 在 `generator/generator.go` 的 `FuncMap` 中注册函数
3. 在模板中使用新函数

### 自定义模板

你可以修改 `assert/template/` 目录下的模板文件来自定义生成的代码格式：

- **修改 PO 结构体**：编辑 `po.template`
- **修改 DTO 结构体**：编辑 `dto.template`
- **修改 DAO 方法**：编辑 `dao.template`
- **修改工具函数**：编辑 `tools/` 下的模板

### 扩展新功能

1. 在 `parser/` 中添加新的解析逻辑
2. 在 `generator/` 中添加生成逻辑
3. 创建新的模板文件
4. 更新配置文件结构

## 📚 核心特性说明

### 智能索引处理

- **主键索引**：自动生成 `SelectById`、`UpdateById`、`DeleteById` 等方法
- **唯一索引**：为每个唯一索引生成 `SelectByXXX` 方法
- **普通索引**：为每个普通索引生成 `SelectListByXXX` 方法

### 零值覆盖处理

生成的 DAO 方法会自动处理零值覆盖问题：

```go
// 使用 DTO 更新时，只有非 nil 字段会被更新
func (dao *TArtifactDAO) UpdateByDTO(ctx context.Context, id uint64, dto *query.TArtifactDTO) error {
    updates := make(map[string]interface{})
    
    if dto.ArtifactName != nil {
        updates["artifactName"] = *dto.ArtifactName
    }
    // nil 值不会被包含在 updates 中，避免零值覆盖
    
    return dao.db.WithContext(ctx).Model(&entity.TArtifact{}).Where("id = ?", id).Updates(updates).Error
}
```

### 类型安全

- 所有方法都使用强类型参数
- 自动处理可空字段的指针类型
- 提供类型转换和验证

## 📦 依赖项

- [GORM](https://gorm.io/) - ORM 库
- [Wire](https://github.com/google/wire) - 依赖注入
- [lo](https://github.com/samber/lo) - 函数式编程工具
- [yaml.v3](https://github.com/go-yaml/yaml) - YAML 解析

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

MIT License