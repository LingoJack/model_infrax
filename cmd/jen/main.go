package main

import (
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

const Version = "1.1.0"

var defaultConfigPaths = []string{
	"./application.yml",
	"./assets/application.yml",
	"/Applications/model_infrax/application.yml",
	"/Applications/model_infrax/assets/application.yml",
}

func main() {
	configPath := flag.StringP("config", "c", "", "配置文件路径")
	showVersion := flag.BoolP("version", "v", false, "显示版本号")
	flag.Parse()

	if *showVersion {
		log.Printf("jen version %s", Version)
		return
	}

	if *configPath != "" {
		log.Printf("📋 使用用户指定的配置文件: %s", *configPath)
		if err := StartGenerateWithConf(*configPath); err != nil {
			log.Fatalf("❌ 使用配置文件 %s 失败: %v", *configPath, err)
		}
		log.Println("🎊 程序执行完成")
		return
	}

	log.Println("🔍 未找到 model_infra.go，尝试使用默认配置文件...")
	for _, path := range defaultConfigPaths {
		if fileExists(path) {
			log.Printf("📁 找到配置文件: %s", path)
			if err := StartGenerateWithConf(path); err != nil {
				log.Printf("⚠️ 配置文件 %s 加载失败: %v，继续尝试下一个...", path, err)
				continue
			}
			log.Println("🎊 程序执行完成")
			return
		}
	}

	os.Exit(1)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func StartGenerateWithConf(configPath string) error {
	log.Println("🚀 开始执行代码生成...")

	appInstance, err := InitializeApp(configPath)
	if err != nil {
		return err
	}

	if err = appInstance.Run(); err != nil {
		return err
	}

	return nil
}
