package main

import (
	"context"
	"log"

	"github.com/LingoJack/model_infrax/internal/conf"
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
	var err error
	configPath := flag.StringP("config", "c", "", "配置文件路径")
	showVersion := flag.BoolP("version", "v", false, "显示版本号")
	flag.Parse()

	if *showVersion {
		log.Printf("jen version %s", Version)
		return
	}

	path := ""
	if *configPath != "" {
		log.Printf("使用用户指定的配置文件: %s", *configPath)
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

	err = conf.Load(path)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
		return
	}

	if err = GenerateCode(context.Background(), path); err != nil {
		log.Fatalf("使用配置文件 %s 失败: %v", path, err)
		return
	}
}

func GenerateCode(ctx context.Context, path string) error {
	appInstance, err := InitializeApp(path)
	if err != nil {
		return err
	}

	if err = appInstance.Run(ctx); err != nil {
		return err
	}

	return nil
}
