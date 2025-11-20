package config

import (
	"fmt"
	"os"

	"github.com/LingoJack/model_infrax/tool"

	"gopkg.in/yaml.v3"
)

type Configger struct {
	GenerateConfig GenerateConfig `yaml:"generate_config"`
	GenerateOption GenerateOption `yaml:"generate_option"`
}

type GenerateConfig struct {
	GenerateMode string `yaml:"generate_mode"` // 生成模式: database(从数据库解析) 或 statement(从SQL文件解析)

	// database 模式配置
	DatabaseName string `yaml:"database_name"` // 数据库名称
	Host         string `yaml:"host"`          // 数据库主机地址
	Port         int    `yaml:"port"`          // 数据库端口
	URLTemplate  string `yaml:"url_template"`  // URL模板
	Username     string `yaml:"username"`      // 数据库用户名
	Password     string `yaml:"password"`      // 数据库密码

	// statement 模式配置
	SqlFilePath string `yaml:"sql_file_path"` // SQL文件路径

	// 通用配置
	AllTables  bool     `yaml:"all_tables"`  // 是否生成所有表
	TableNames []string `yaml:"table_names"` // 表名列表
}

type GenerateOption struct {
	OutputPath            string        `yaml:"output_path"`              // 输出路径
	IgnoreTableNamePrefix bool          `yaml:"ignore_table_name_prefix"` // 是否忽略表名前缀
	CrudOnlyIdx           bool          `yaml:"crud_only_idx"`            // 是否仅CRUD索引
	Package               PackageConfig `yaml:"package_name"`             // 包配置
	ModelAllInOneFile     bool          `yaml:"all_model_in_one_file"`    // 是否将所有模型放在一个文件中
	ModelAllInOneFileName string        `yaml:"all_model_in_one_file_name"`
	UseFramework          string        `yaml:"use_framework"`
}

type PackageConfig struct {
	PoPackage   string `yaml:"po_package"`
	DtoPackage  string `yaml:"dto_package"`
	VoPackage   string `yaml:"vo_package"`
	DaoPackage  string `yaml:"dao_package"`
	ToolPackage string `yaml:"tool_package"`
}

func NewConfigger(configPath string) (*Configger, error) {
	path := tool.EscapeHomeDir(configPath)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Configger
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %w", err)
	}

	// 展开输出路径中的 ~ 符号
	config.GenerateOption.OutputPath = tool.EscapeHomeDir(config.GenerateOption.OutputPath)

	// 展开SQL文件路径中的 ~ 符号
	if config.GenerateConfig.SqlFilePath != "" {
		config.GenerateConfig.SqlFilePath = tool.EscapeHomeDir(config.GenerateConfig.SqlFilePath)
	}

	return &config, nil
}
