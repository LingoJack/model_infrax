package main

import (
	"log"

	"github.com/LingoJack/model_infrax"
)

// SQL文件模式示例：从SQL文件解析表结构生成代码
// 适合没有数据库连接，但有SQL建表语句的场景
// 这种方式不需要连接数据库，速度更快
func main() {
	log.Println("🚀 开始从SQL文件生成代码...")

	// 使用 Builder 模式配置并生成代码
	err := model_infrax.NewBuilder().
		// 配置SQL文件模式
		StatementMode("./schema.sql").
		// 指定要生成的表（如果不指定则生成所有表）
		Tables("t_user", "t_order", "t_product").
		// 配置输出路径
		OutputPath("./output").
		// 忽略表名前缀
		IgnoreTableNamePrefix(true).
		// 将所有Model生成到一个文件中
		ModelAllInOneFile(true, "models.go").
		// 使用框架（如 itea-go）
		UseFramework("itea-go").
		// 配置包名
		Packages("model/entity", "model/query", "model/view", "dao", "tool").
		// 构建配置并生成代码
		BuildAndGenerate()

	if err != nil {
		log.Fatalf("❌ 生成失败: %v", err)
	}

	log.Println("✅ 代码生成成功！")
}