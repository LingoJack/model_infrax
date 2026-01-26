package tool

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/bwmarrin/snowflake"
)

// Snowflake 雪花算法封装
// 提供分布式唯一 ID 生成能力，基于 Twitter Snowflake 算法
// ID 结构：1位符号位 + 41位时间戳 + 10位节点ID + 12位序列号
type Snowflake struct {
	node *snowflake.Node
}

// NewSnowflake 创建一个新的雪花算法实例
//
// 参数:
//   - nodeID: 节点ID，取值范围 0-1023，用于区分不同的服务节点
//
// 返回:
//   - *Snowflake: 雪花算法实例
//   - error: 创建失败时返回错误（如 nodeID 超出范围）
//
// 示例:
//
//	sf, err := NewSnowflake(1)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	id := sf.NextID()
func NewSnowflake(nodeID int64) (*Snowflake, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("[NewSnowflake] 创建雪花算法节点失败 nodeID=%d: %w", nodeID, err)
	}

	return &Snowflake{
		node: node,
	}, nil
}

// NextID 生成下一个唯一ID（int64 格式）
//
// 返回:
//   - int64: 唯一的 64 位整数 ID
//
// 注意:
//   - 此方法线程安全
//   - 同一毫秒内最多生成 4096 个 ID
func (s *Snowflake) NextID() int64 {
	return s.node.Generate().Int64()
}

// NextIDString 生成下一个唯一ID（字符串格式）
//
// 返回:
//   - string: 唯一 ID 的字符串表示
//
// 注意:
//   - 此方法线程安全
//   - 字符串格式便于在 JSON 中传输，避免精度丢失
func (s *Snowflake) NextIDString() string {
	return s.node.Generate().String()
}

// ParseID 解析 ID 字符串为 snowflake.ID 对象
//
// 参数:
//   - idStr: ID 字符串
//
// 返回:
//   - snowflake.ID: 解析后的 ID 对象，可用于提取时间戳等信息
//   - error: 解析失败时返回错误
//
// 示例:
//
//	id, err := ParseID("1234567890123456789")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	timestamp := id.Time() // 获取 ID 生成时间
func ParseID(idStr string) (snowflake.ID, error) {
	id, err := snowflake.ParseString(idStr)
	if err != nil {
		return 0, fmt.Errorf("[ParseID] 解析ID失败 idStr=%s: %w", idStr, err)
	}
	return id, nil
}

// 全局默认实例相关变量
var (
	defaultSnowflake *Snowflake // 全局默认雪花算法实例
	once             sync.Once  // 确保只初始化一次
	initErr          error      // 初始化错误
)

// getDefaultNodeID 获取默认的节点 ID
//
// 优先级:
//  1. 从环境变量 SNOWFLAKE_NODE_ID 读取（范围 0-1023）
//  2. 使用默认值 1
//
// 返回:
//   - int64: 节点 ID
//
// 注意:
//   - 如果环境变量值无效（非数字或超出范围），将使用默认值 1
func getDefaultNodeID() int64 {
	if nodeIDStr := os.Getenv("SNOWFLAKE_NODE_ID"); nodeIDStr != "" {
		if nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64); err == nil {
			if nodeID >= 0 && nodeID <= 1023 {
				return nodeID
			}
		}
	}
	return 1 // 默认使用节点 1
}

// ensureInitialized 确保全局默认实例已初始化
//
// 特性:
//   - 自动初始化，使用 getDefaultNodeID() 获取节点 ID
//   - 线程安全，使用 sync.Once 确保只初始化一次
//   - 初始化失败时会打印错误日志，但不会 panic
//
// 注意:
//   - 此方法会在首次调用 SnowflakeID 等方法时自动执行
//   - 初始化失败后，后续调用会返回错误或 panic
func ensureInitialized() {
	once.Do(func() {
		nodeID := getDefaultNodeID()
		defaultSnowflake, initErr = NewSnowflake(nodeID)
		if initErr != nil {
			// 初始化失败时记录错误，但不 panic
			fmt.Printf("[ensureInitialized] 自动初始化雪花算法失败 nodeID=%d: %v\n", nodeID, initErr)
		}
	})
}

// SetNodeID 设置全局默认实例的节点 ID（可选）
//
// 参数:
//   - nodeID: 节点ID，取值范围 0-1023
//
// 返回:
//   - error: 设置失败时返回错误
//
// 注意:
//   - 必须在首次调用 SnowflakeID 等方法之前调用
//   - 如果已经初始化过，此方法不会重复初始化
//   - 如果不调用此方法，将使用环境变量或默认值 1
//
// 示例:
//
//	// 在程序启动时设置节点 ID
//	if err := tool.SetNodeID(123); err != nil {
//	    log.Fatal(err)
//	}
func SetNodeID(nodeID int64) error {
	once.Do(func() {
		defaultSnowflake, initErr = NewSnowflake(nodeID)
	})
	return initErr
}

// SnowflakeID 生成唯一 ID（int64 格式）
//
// 返回:
//   - int64: 唯一的 64 位整数 ID
//
// 特性:
//   - 开箱即用，无需手动初始化
//   - 首次调用时自动初始化全局实例
//   - 线程安全
//
// 示例:
//
//	id := tool.SnowflakeID()
//	fmt.Println(id) // 输出: 1234567890123456789
func SnowflakeID() int64 {
	ensureInitialized()
	if defaultSnowflake == nil {
		panic(fmt.Sprintf("[SnowflakeID] 默认雪花算法实例初始化失败: %v", initErr))
	}
	return defaultSnowflake.NextID()
}

// SnowflakeIDString 生成唯一 ID（字符串格式）
//
// 返回:
//   - string: 唯一 ID 的字符串表示
//
// 特性:
//   - 开箱即用，无需手动初始化
//   - 首次调用时自动初始化全局实例
//   - 线程安全
//   - 字符串格式便于在 JSON 中传输，避免精度丢失
//
// 示例:
//
//	idStr := tool.SnowflakeIDString()
//	fmt.Println(idStr) // 输出: "1234567890123456789"
func SnowflakeIDString() string {
	ensureInitialized()
	if defaultSnowflake == nil {
		panic(fmt.Sprintf("[SnowflakeIDString] 默认雪花算法实例初始化失败: %v", initErr))
	}
	return defaultSnowflake.NextIDString()
}
