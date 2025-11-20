# ✅ 目录结构优化 - 操作清单

## 📋 需要手动执行的操作

### 第一步：重命名目录 ⭐ 最重要

```bash
cd /Users/jacklingo/dev/model_infrax
mv assert assets
```

### 第二步：验证修改

```bash
# 1. 检查是否还有遗漏的 assert 引用
grep -r "assert" . --exclude-dir=".git" --exclude="*.md" --exclude-dir="assets"

# 2. 运行测试
go test ./...

# 3. 构建项目
go build ./cmd/jen

# 4. 尝试运行
./jen -c ./assets/application.yml
```

### 第三步：提交更改

```bash
# 查看所有修改
git status

# 添加所有新文件和修改
git add .

# 提交
git commit -m "feat: 优化目录结构

- 重命名 assert 为 assets（修正拼写错误）
- 新增 examples 目录，包含4个使用示例
- 新增 MIGRATION.md 迁移指南
- 新增 QUICKSTART.md 快速开始指南
- 新增 .gitignore 文件
- 更新所有相关文件的路径引用
- 更新 README.md 项目结构说明"

# 推送到远程（如果需要）
git push origin main
```

## 📊 已完成的修改

### ✅ 新增文件（9个）

1. ✅ `examples/README.md` - 示例说明文档
2. ✅ `examples/basic/main.go` - 基础使用示例
3. ✅ `examples/database/main.go` - 数据库模式示例
4. ✅ `examples/statement/main.go` - SQL文件模式示例
5. ✅ `examples/advanced/main.go` - 高级用法示例
6. ✅ `MIGRATION.md` - 迁移指南
7. ✅ `OPTIMIZATION_SUMMARY.md` - 优化总结
8. ✅ `QUICKSTART.md` - 快速开始指南
9. ✅ `.gitignore` - Git忽略文件

### ✅ 修改文件（6个）

1. ✅ `cmd/jen/main.go` - 更新默认配置路径
2. ✅ `assets/application.yml` - 更新SQL文件路径
3. ✅ `build.sh` - 更新构建脚本路径
4. ✅ `README.md` - 更新项目结构和资源链接
5. ✅ `parser/database_parser_test.go` - 更新所有测试路径
6. ✅ `parser/statement_parser_test.go` - 更新所有测试路径

### ⏳ 待执行操作（1个）

1. ⏳ 重命名目录：`assert/` → `assets/`

## 🎯 验证清单

完成上述操作后，请逐项检查：

- [ ] `assert` 目录已成功重命名为 `assets`
- [ ] 运行 `grep -r "assert"` 没有发现遗漏的引用（除了文档）
- [ ] 所有测试通过：`go test ./...`
- [ ] 项目可以正常构建：`go build ./cmd/jen`
- [ ] 命令行工具可以正常运行：`./jen -c ./assets/application.yml`
- [ ] 所有示例代码可以正常编译
- [ ] Git 状态正常，没有意外的修改
- [ ] 文档链接都能正常访问

## 📁 优化后的目录结构

```
model_infrax/
├── cmd/jen/                      # 命令行工具
├── config/                       # 配置管理
├── parser/                       # 解析器
├── generator/                    # 代码生成器
├── model/                        # 数据模型
├── pkg/app/                      # 公共包
├── tool/                         # 工具函数
├── examples/                     # 使用示例 ✨ 新增
│   ├── README.md
│   ├── basic/
│   ├── database/
│   ├── statement/
│   └── advanced/
├── assets/                       # 资源文件 ✨ 重命名
│   ├── application.yml
│   ├── schema.sql
│   ├── install.sh
│   ├── jcode
│   └── prompt/
├── api.go
├── wire.go
├── wire_gen.go
├── build.sh
├── go.mod
├── go.sum
├── .gitignore                    # ✨ 新增
├── README.md
├── README_API.md
├── MIGRATION.md                  # ✨ 新增
├── OPTIMIZATION_SUMMARY.md       # ✨ 新增
├── QUICKSTART.md                 # ✨ 新增
└── TODO.md                       # 本文件
```

## 💡 优化亮点

### 1. 修正拼写错误
- ❌ `assert` (断言) → ✅ `assets` (资源)
- 更符合语义，避免误解

### 2. 完善示例代码
- 4个完整的使用示例
- 涵盖所有常见场景
- 详细的注释说明

### 3. 完善文档体系
- 迁移指南 - 帮助用户升级
- 快速开始 - 快速上手
- 优化总结 - 了解变更

### 4. 规范项目结构
- 添加 .gitignore
- 清晰的目录组织
- 便于维护和扩展

## 🚀 下一步建议

完成当前优化后，可以考虑：

1. **创建 `internal` 目录**
   - 将内部包移到 `internal/` 下
   - 防止外部项目导入内部实现

2. **创建 `scripts` 目录**
   - 将 `build.sh` 移到 `scripts/` 下
   - 添加更多构建和部署脚本

3. **创建 `docs` 目录**
   - 将文档集中管理
   - 添加更多技术文档

4. **移动 `tool` 到 `pkg/tool`**
   - 如果需要对外暴露
   - 更符合 Go 项目规范

## 🆘 遇到问题？

### 问题1: 重命名后测试失败

**解决方案**：
```bash
# 检查是否有硬编码的路径
grep -r "/assert/" .
grep -r "assert/" .
```

### 问题2: Git 显示大量删除和新增

**解决方案**：
```bash
# Git 应该能自动识别重命名
# 如果没有，可以手动告诉 Git
git add -A
git status  # 应该显示 renamed: assert/... -> assets/...
```

### 问题3: 构建失败

**解决方案**：
```bash
# 清理缓存重新构建
go clean -cache
go mod tidy
go build ./cmd/jen
```

## 📞 完成后

完成所有操作后：

1. ✅ 删除本文件（TODO.md）
2. ✅ 查看 [QUICKSTART.md](./QUICKSTART.md) 快速开始
3. ✅ 查看 [examples/](./examples/) 了解使用方法
4. ✅ 开始使用优化后的项目！

---

**祝你使用愉快！** 🎉