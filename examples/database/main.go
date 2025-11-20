package main

import (
	"log"

	"github.com/LingoJack/model_infrax"
)

// 数据库模式示例：从数据库读取表结构生成代码
// 适合已有数据库，需要生成对应的Go代码的场景
func main() {
	log.Println("🚀 开始从数据库生成代码...")

	// 使用 Builder 模式配置并生成代码
	// 这种方式不需要配置文件，所有配置都在代码中完成
	err := model_infrax.NewBuilder().
		// 配置数据库连接信息
		DatabaseMode("localhost", 3306, "test_db", "root", "password").
		// 生成所有表
		AllTables().
		// 配置输出路径
		OutputPath("./output").
		// 忽略表名前缀（如 t_user -> User）
		IgnoreTableNamePrefix(true).
		// 只为有索引的字段生成CRUD方法
		CrudOnlyIdx(true).
		// 配置包名
		Packages("model/entity", "model/query", "model/view", "dao", "tool").
		// 构建配置并生成代码
		BuildAndGenerate()

	if err != nil {
		log.Fatalf("❌ 生成失败: %v", err)
	}

	log.Println("✅ 代码生成成功！")
}