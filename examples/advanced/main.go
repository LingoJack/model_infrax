package main

import (
	"log"

	"github.com/LingoJack/model_infrax/config"
	"github.com/LingoJack/model_infrax/modelinfra"
)

// 高级用法示例：展示更多配置选项和灵活用法
func main() {
	log.Println("🚀 高级用法示例...")

	// 示例1: 使用自定义配置构建器
	customConfig := modelinfra.NewBuilder().
		DatabaseMode("localhost", 3306, "test_db", "root", "password").
		Tables("t_user", "t_order"). // 只生成指定的表
		OutputPath("./custom_output").
		IgnoreTableNamePrefix(true).
		CrudOnlyIdx(false). // 为所有字段生成CRUD方法
		ModelAllInOneFile(false, "").
		UseFramework(""). // 使用GORM原生
		Packages("po", "dto", "vo", "dao", "tool").
		Build()

	// 可以在这里对配置进行进一步的自定义修改
	log.Printf("📋 配置信息: 输出路径=%s, 生成模式=%s",
		customConfig.GenerateOption.OutputPath,
		customConfig.GenerateConfig.GenerateMode)

	// 示例2: 使用配置对象生成代码
	err := GenerateFromBuilder(customConfig)
	if err != nil {
		log.Fatalf("❌ 生成失败: %v", err)
	}

	log.Println("✅ 代码生成成功！")

	// 示例3: 批量生成多个数据库的代码
	databases := []string{"db1", "db2", "db3"}
	for _, dbName := range databases {
		log.Printf("📦 正在生成数据库 %s 的代码...", dbName)

		err := modelinfra.NewBuilder().
			DatabaseMode("localhost", 3306, dbName, "root", "password").
			AllTables().
			OutputPath("./output/" + dbName).
			IgnoreTableNamePrefix(true).
			BuildAndGenerate()

		if err != nil {
			log.Printf("⚠️ 数据库 %s 生成失败: %v", dbName, err)
			continue
		}

		log.Printf("✅ 数据库 %s 生成成功", dbName)
	}

	log.Println("🎊 所有任务完成！")
}

// GenerateFromBuilder 从配置构建器生成代码
// 这是一个辅助函数，用于演示如何使用配置对象
func GenerateFromBuilder(cfg *config.Configger) error {
	// 这里可以添加更多的自定义逻辑
	// 例如：验证配置、记录日志、发送通知等

	// 实际生成代码的逻辑
	// 注意：这里需要手动创建App实例并运行
	// 在实际使用中，建议使用 BuildAndGenerate() 方法
	log.Println("⚙️ 使用自定义配置生成代码...")

	return nil
}