package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

// Version 当前版本号
const Version = "1.0.6"

// defaultConfigPaths 默认配置文件路径列表
// 按照优先级顺序查找配置文件，找到第一个可用的就使用
var defaultConfigPaths = []string{
	"./application.yml",                        // 当前目录下的配置文件
	"./assets/application.yml",                 // assets目录下的配置文件
	"/Applications/jen/application.yml",        // 系统安装目录下的配置文件
	"/Applications/jen/assets/application.yml", // 系统安装目录assets子目录下的配置文件
}

// defaultGoFile 默认要执行的 Go 文件
const defaultGoFile = "model_infra.go"

// main 主函数，程序入口点
// 采用优先级自动降级策略，按以下顺序尝试：
//  1. 如果用户指定了 --config 参数，使用指定的配置文件（最高优先级，用户意图优先）
//  2. 如果当前目录存在 model_infra.go 文件，直接执行它
//  3. 按默认路径列表查找并使用第一个可用的配置文件
//  4. 如果以上都失败，提示用户并退出
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
	// 定义命令行参数
	configPath := flag.StringP("config", "c", "", "配置文件路径")
	showVersion := flag.BoolP("version", "v", false, "显示版本号")
	flag.Parse()

	// 优先级 0: 显示版本号（最高优先级）
	if *showVersion {
		log.Printf("jen version %s", Version)
		return
	}

	// 优先级 1: 用户指定的配置文件（最高优先级，用户意图优先）
	if *configPath != "" {
		log.Printf("📋 使用用户指定的配置文件: %s", *configPath)
		if err := runWithConfig(*configPath); err != nil {
			log.Fatalf("❌ 使用配置文件 %s 失败: %v", *configPath, err)
		}
		log.Println("🎊 程序执行完成")
		return
	}

	// 优先级 2: 检查是否存在 model_infra.go 文件
	if fileExists(defaultGoFile) {
		log.Printf("🎯 检测到 %s 文件，直接执行...", defaultGoFile)
		if err := runGoFile(defaultGoFile); err != nil {
			log.Fatalf("❌ 执行 %s 失败: %v", defaultGoFile, err)
		}
		log.Println("🎊 程序执行完成")
		return
	}

	// 优先级 3: 尝试默认配置文件路径
	log.Println("🔍 未找到 model_infra.go，尝试使用默认配置文件...")
	for _, path := range defaultConfigPaths {
		if fileExists(path) {
			log.Printf("📁 找到配置文件: %s", path)
			if err := runWithConfig(path); err != nil {
				log.Printf("⚠️ 配置文件 %s 加载失败: %v，继续尝试下一个...", path, err)
				continue
			}
			log.Println("🎊 程序执行完成")
			return
		}
	}

	// 优先级 4: 所有方式都失败，提示用户
	log.Println("❌ 无法找到可用的配置或代码文件")
	log.Println("")
	log.Println("💡 请选择以下任一方式：")
	log.Println("")
	log.Println("   方式 1: 创建 model_infra.go 文件（推荐用于编程式控制）")
	log.Println("   -------------------------------------------------------")
	log.Println("   在当前目录创建 model_infra.go，示例：")
	log.Println("")
	log.Println("   package main")
	log.Println("")
	log.Println("   import (")
	log.Println("       \"log\"")
	log.Println("       \"github.com/LingoJack/model_infrax\"")
	log.Println("   )")
	log.Println("")
	log.Println("   func main() {")
	log.Println("       err := model_infrax.Generate(")
	log.Println("           model_infrax.NewBuilder().")
	log.Println("               DatabaseMode(\"localhost\", 3306, \"mydb\", \"root\", \"pass\").")
	log.Println("               AllTables().")
	log.Println("               OutputPath(\"./output\").")
	log.Println("               BuildAndGenerate(),")
	log.Println("       )")
	log.Println("       if err != nil {")
	log.Println("           log.Fatal(err)")
	log.Println("       }")
	log.Println("   }")
	log.Println("")
	log.Println("   方式 2: 使用配置文件（推荐用于声明式配置）")
	log.Println("   -------------------------------------------------------")
	log.Println("   创建 application.yml 配置文件，或使用 --config 参数指定")
	log.Println("   示例: jen --config ./my-config.yml")
	log.Println("")
	log.Printf("   默认配置文件查找路径: %v", defaultConfigPaths)
	log.Println("")
	os.Exit(1)
}

// fileExists 检查文件是否存在
// 参数:
//
//	filename: 要检查的文件路径
//
// 返回:
//
//	bool: 文件存在返回 true，否则返回 false
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// runGoFile 执行指定的 Go 文件
// 使用 go run 命令执行文件，并将输出重定向到当前进程的标准输出/错误输出
// 参数:
//
//	filename: 要执行的 Go 文件路径
//
// 返回:
//
//	error: 执行过程中的错误，nil 表示成功
func runGoFile(filename string) error {
	// 获取文件的绝对路径
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	log.Printf("📂 执行文件: %s", absPath)

	// 创建 go run 命令
	cmd := exec.Command("go", "run", absPath)

	// 将命令的输出重定向到当前进程的标准输出和错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// 执行命令
	return cmd.Run()
}

// runWithConfig 使用配置文件运行应用
// 参数:
//
//	configPath: 配置文件路径
//
// 返回:
//
//	error: 执行过程中的错误，nil 表示成功
func runWithConfig(configPath string) error {
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
