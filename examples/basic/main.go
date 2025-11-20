package main

import (
	"log"

	"github.com/LingoJack/model_infrax"
	"github.com/LingoJack/model_infrax/config"
)

func main() {
	log.Println("🚀 开始使用配置文件生成代码...")

	err := model_infrax.GenerateFromConfig("./application.yml")
	if err != nil {
		log.Fatalf("❌ 生成失败: %v", err)
	}

	log.Println("✅ 代码生成成功！")
}