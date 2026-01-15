package main

import (
	"context"
	"log"

	"github.com/LingoJack/model_infrax/internal/app"
	"github.com/LingoJack/model_infrax/internal/conf"
	flag "github.com/spf13/pflag"
)

const Version = "1.1.0"

func main() {
	showVersion := flag.BoolP("version", "v", false, "显示版本号")
	flag.Parse()

	if *showVersion {
		log.Printf("jen version %s", Version)
		return
	}

	if err := GenerateCode(context.Background()); err != nil {
		log.Fatalf("使用配置文件 %s 失败: %v", conf.ValueStr("config_path"), err)
		return
	}
}

func GenerateCode(ctx context.Context) (err error) {
	if err = app.Run(ctx); err != nil {
		return err
	}

	return nil
}
