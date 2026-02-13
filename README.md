# Model Infrax

> Go 代码生成器 CLI，根据 SQL 建表语句或数据库连接，自动生成 PO / DTO / VO / DAO / Tool 等代码。

## 安装

```bash
go install github.com/LingoJack/model_infrax/cmd/jen@latest
```

## 快速上手

### 1. 初始化配置

在你的项目根目录执行：

```bash
jen --init
```

这会在当前目录创建 `.model_infrax/` 目录，包含：

| 文件 | 说明 |
|------|------|
| `.model_infrax/config.yml` | 代码生成配置文件 |
| `.model_infrax/schema.sql` | SQL 建表语句文件 |

### 2. 编写 SQL

在 `.model_infrax/schema.sql` 中编写你的建表语句：

```sql
CREATE TABLE IF NOT EXISTS `t_user`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name`       varchar(128)        NOT NULL COMMENT '用户名',
    `createTime` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户表';
```

### 3. 编辑配置（可选）

按需修改 `.model_infrax/config.yml`：

```yaml
generate_config:
  # 生成模式: database(从数据库解析) 或 statement(从SQL文件解析)
  generate_mode: statement

  # statement 模式配置（相对于项目根目录）
  sql_file_path: .model_infrax/schema.sql

  # 通用配置
  all_tables: true        # true: 生成所有表; false: 只生成 table_names 指定的表
  table_names:
    - t_example

generate_option:
  # 输出路径（相对于项目根目录）
  output_path: target/jen

  # go 的 package 映射
  package:
    po: model/entity
    dto: model/query
    vo: model/view
    dao: dao
    tool: tool

  # 使用框架, 可选: itea-go, gorm
  use_framework: itea-go
```

如需使用数据库模式，配置 `generate_mode: database` 并填写数据库连接信息：

```yaml
generate_config:
  generate_mode: database
  database_name: test_db
  host: localhost
  port: 3306
  username: root
  password: 123456
```

### 4. 生成代码

```bash
jen
```

生成的代码输出到 `target/jen/` 目录下。

也可以通过 `-c` 指定配置文件路径：

```bash
jen -c /path/to/your/config.yml
```

## 命令参考

| 命令 | 说明 |
|------|------|
| `jen --init` | 在当前目录初始化 `.model_infrax` 配置目录 |
| `jen` | 加载 `.model_infrax/config.yml` 并生成代码 |
| `jen -c <path>` | 指定配置文件路径生成代码 |
| `jen -v` | 显示版本信息 |

## 生成产物目录结构

```
target/jen/
├── model/
│   ├── entity/     # PO (Persistent Object)
│   ├── query/      # DTO (Data Transfer Object)
│   └── view/       # VO (View Object)
├── dao/            # DAO (Data Access Object)
└── tool/           # 工具代码
```