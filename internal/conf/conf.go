package conf

import (
	"fmt"
	"os"
	"sync"

	"github.com/LingoJack/model_infrax/internal/constant"
	"github.com/LingoJack/model_infrax/internal/logger"
	"github.com/LingoJack/model_infrax/internal/tool"
	"gopkg.in/yaml.v3"
)

var (
	config     *appConfig
	configLock sync.RWMutex

	// override 保存 CLI flag 等来源的运行时覆盖值，读取优先级高于配置文件
	override     = map[string]any{}
	overrideLock sync.RWMutex
)

// DefaultConfigPaths 默认配置文件查找路径（按顺序探测）
var DefaultConfigPaths = []string{
	"./.model_infrax/config.yml",
}

// ActivateConfigPath 返回当前生效的配置文件路径（未加载时为空）
func ActivateConfigPath() string {
	return ValueStr(constant.ActivateConfigPathKey)
}

// InitWithPath 初始化配置文件
// customPath 非空时优先使用，否则按 DefaultConfigPaths 顺序探测
// 可重复调用（如 UI 保存配置后重载），每次整体替换已加载的配置
func InitWithPath(customPath string) error {
	var configPath string

	if customPath != "" {
		if !tool.IsValidFilePath(customPath) {
			return fmt.Errorf("指定的配置文件不存在: %s", customPath)
		}
		configPath = customPath
	} else {
		for _, p := range DefaultConfigPaths {
			if tool.IsValidFilePath(p) {
				configPath = p
				break
			}
		}
	}

	if configPath == "" {
		return fmt.Errorf("未找到有效的配置文件，请使用 -c 参数指定配置文件路径")
	}

	logger.Infof("[InitWithPath] 加载配置文件: %s", configPath)

	if err := load(configPath); err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	// 将实际生效的配置文件路径写入配置，供外部读取
	if len(ValueStr(constant.ActivateConfigPathKey)) > 0 {
		return fmt.Errorf("配置文件异常，非法key: %s", constant.ActivateConfigPathKey)
	}
	configLock.Lock()
	config.data[constant.ActivateConfigPathKey] = configPath
	configLock.Unlock()

	return nil
}

// SetOverride 设置运行时覆盖值（如 CLI flag）
// 覆盖值优先于配置文件被 Value* 系列读取到
func SetOverride(key string, val any) {
	overrideLock.Lock()
	defer overrideLock.Unlock()
	override[key] = val
}

// load 加载配置文件并整体替换当前配置（支持重载）
func load(file string) error {
	c, err := newConfig(file)
	if err != nil {
		return err
	}
	configLock.Lock()
	config = c
	configLock.Unlock()
	return nil
}

// appConfig 保存 YAML 解析后的配置数据
type appConfig struct {
	data map[string]interface{}
}

// newConfig 读取并解析 YAML 配置文件
func newConfig(file string) (*appConfig, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析为 map 以支持动态键访问
	var rawData map[string]interface{}
	if err = yaml.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &appConfig{data: rawData}, nil
}

// get 通过点分隔的键逐层访问嵌套 map
// 返回 (值, 是否存在)；任何一层不存在或不是 map 时返回 (nil, false)
func (c *appConfig) get(key string) (interface{}, bool) {
	if c.data == nil {
		return nil, false
	}

	keys := splitKey(key)
	var current interface{} = c.data
	for _, k := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current, ok = m[k]
		if !ok {
			return nil, false
		}
	}

	return current, true
}
