package config

import (
	"fmt"
	"model_infrax/tool"
	"os"

	"gopkg.in/yaml.v3"
)

type Configger struct {
	GenerateConfig GenerateConfig `yaml:"generate_config"`
	GenerateOption GenerateOption `yaml:"generate_option"`
}

type GenerateConfig struct {
	GenerateMode string   `yaml:"generate_mode"` // 生成模式
	DatabaseName string   `yaml:"database_name"` // 数据库名称
	Host         string   `yaml:"host"`          // 数据库主机地址
	Port         int      `yaml:"port"`          // 数据库端口
	URLTemplate  string   `yaml:"url_template"`  // URL模板
	Username     string   `yaml:"username"`      // 数据库用户名
	Password     string   `yaml:"password"`      // 数据库密码
	AllTables    bool     `yaml:"all_tables"`
	TableNames   []string `yaml:"table_names"`
}

type GenerateOption struct {
	OutputPath            string        `yaml:"output_path"`              // 输出路径
	IgnoreTableNamePrefix bool          `yaml:"ignore_table_name_prefix"` // 是否忽略表名前缀
	CrudOnlyIdx           bool          `yaml:"crud_only_idx"`            // 是否仅CRUD索引
	Package               PackageConfig `yaml:"package"`                  // 包配置
}

type PackageConfig struct {
	Model  string `yaml:"model"`
	DTO    string `yaml:"dto"`
	VO     string `yaml:"vo"`
	Tool   string `yaml:"tool"`
	Mapper string `yaml:"mapper"`
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

	return &config, nil
}
