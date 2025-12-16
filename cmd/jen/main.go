package main

import (
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

// Version 当前版本号
const Version = "1.1.0"

// defaultConfigPaths 默认配置文件路径列表
// 按照优先级顺序查找配置文件，找到第一个可用的就使用
var defaultConfigPaths = []string{
	"./application.yml",                                 // 当前目录下的配置文件
	"./assets/application.yml",                          // assets目录下的配置文件
	"/Applications/model_infrax/application.yml",        // 系统安装目录下的配置文件
	"/Applications/model_infrax/assets/application.yml", // 系统安装目录assets子目录下的配置文件
}

// main 主函数，程序入口点
// 采用优先级自动降级策略，按以下顺序尝试：
//  1. 如果用户指定了 --config 参数，使用指定的配置文件（最高优先级，用户意图优先）
//  2. 按默认路径列表查找并使用第一个可用的配置文件
//
// 支持的命令行参数：
//
//	-c, --config: 指定配置文件路径（可选）
//	-v, --version: 显示版本号
//
// 使用示例：
//
//	jen                                    # 自动选择最合适的方式
//	jen -c ./my-config.yml                 # 强制使用指定的配置文件（最高优先级）
//	jen --config /path/to/config.yml       # 使用长格式参数
//	jen -v                                 # 显示版本号
//	jen --version                          # 显示版本号（长格式）
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

	// 尝试默认配置文件路径
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

	// 初始化应用实例
	appInstance, err := InitializeApp(configPath)
	if err != nil {
		return err
	}

	// 运行应用
	if err = appInstance.Run(); err != nil {
		return err
	}

	return nil
}
