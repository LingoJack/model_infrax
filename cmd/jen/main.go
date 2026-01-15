package main

import (
	"log"
	"os"

	"github.com/LingoJack/model_infrax/internal/tool"
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

	path := ""
	if *configPath != "" {
		log.Printf("📋 使用用户指定的配置文件: %s", *configPath)
		path = tool.DeStringPtr(configPath)
	} else {
		for _, p := range defaultConfigPaths {
			if tool.IsValidFilePath(p) {
				path = p
				break
			}
		}
	}

	log.Printf("当前生效配置文件(%s)", path)
	if len(path) == 0 {
		log.Println("未找到有效的配置文件")
		return
	}
	if err := StartGenerateWithConf(path); err != nil {
		log.Fatalf("使用配置文件 %s 失败: %v", path, err)
	}

	os.Exit(1)
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
