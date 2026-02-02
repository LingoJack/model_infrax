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

// ReplaceNodeIDToZero 将雪花 ID 的节点 ID 替换为 0
//
// Snowflake ID 结构：
//   - 1 位符号位（始终为 0）
//   - 41 位时间戳（毫秒级）
//   - 10 位节点 ID
//   - 12 位序列号
//
// 参数:
//   - id: 原始雪花 ID（int64 格式）
//
// 返回:
//   - int64: 节点 ID 被替换为 0 后的新 ID
//
// 示例:
//
//	newID := ReplaceNodeIDToZero(1234567890123456789)
//	// 返回的 newID 中 node ID 部分为 0
func ReplaceNodeIDToZero(id int64) int64 {
	// Snowflake ID 布局:
	// | 1 bit (符号) | 41 bits (时间戳) | 10 bits (节点ID) | 12 bits (序列号) |
	//
	// 节点 ID 位于第 12-21 位（从右往左数，0-indexed）
	// 需要将这 10 位清零

	// 创建掩码：将第 12-21 位设为 0，其他位设为 1
	// 0x3FF 是 10 个 1（节点 ID 的位数）
	// 左移 12 位后取反，得到节点 ID 位置为 0 的掩码
	nodeIDMask := int64(0x3FF) << 12 // 节点 ID 的掩码（第 12-21 位为 1）
	clearMask := ^nodeIDMask          // 取反得到清除掩码（第 12-21 位为 0）

	// 使用 AND 运算清除节点 ID 位
	return id & clearMask
}

// ReplaceNodeIDToZeroString 将雪花 ID 字符串的节点 ID 替换为 0
//
// 参数:
//   - idStr: 原始雪花 ID（字符串格式）
//
// 返回:
//   - string: 节点 ID 被替换为 0 后的新 ID 字符串
//   - error: 解析失败时返回错误
//
// 示例:
//
//	newIDStr, err := ReplaceNodeIDToZeroString("1234567890123456789")
//	if err != nil {
//	    log.Fatal(err)
//	}
func ReplaceNodeIDToZeroString(idStr string) (string, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("[ReplaceNodeIDToZeroString] 解析ID失败 idStr=%s: %w", idStr, err)
	}
	newID := ReplaceNodeIDToZero(id)
	return strconv.FormatInt(newID, 10), nil
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

// ==================== LRU 节点实例池 ====================

// lruNode LRU 双向链表节点
type lruNode struct {
	nodeID    int64
	snowflake *Snowflake
	prev      *lruNode
	next      *lruNode
}

// SnowflakePool 雪花算法节点实例池（基于 LRU 缓存）
//
// 特性:
//   - 支持多节点 ID 的雪花算法实例管理
//   - 使用 LRU 策略自动淘汰最少使用的实例
//   - 线程安全
//   - 自动创建和缓存实例
type SnowflakePool struct {
	capacity int                    // 最大缓存容量
	cache    map[int64]*lruNode     // nodeID -> lruNode 映射
	head     *lruNode               // LRU 链表头（最近使用）
	tail     *lruNode               // LRU 链表尾（最久未使用）
	mu       sync.RWMutex           // 读写锁
}

// NewSnowflakePool 创建一个新的雪花算法节点实例池
//
// 参数:
//   - capacity: 最大缓存容量，建议设置为 10-100
//
// 返回:
//   - *SnowflakePool: 实例池对象
//
// 示例:
//
//	pool := NewSnowflakePool(50)
//	id := pool.NextID(123) // 使用节点 123 生成 ID
func NewSnowflakePool(capacity int) *SnowflakePool {
	if capacity <= 0 {
		capacity = 10 // 默认容量
	}

	return &SnowflakePool{
		capacity: capacity,
		cache:    make(map[int64]*lruNode),
		head:     nil,
		tail:     nil,
	}
}

// moveToHead 将节点移动到链表头部（标记为最近使用）
func (p *SnowflakePool) moveToHead(node *lruNode) {
	if node == p.head {
		return
	}

	// 从当前位置移除
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node == p.tail {
		p.tail = node.prev
	}

	// 移动到头部
	node.prev = nil
	node.next = p.head
	if p.head != nil {
		p.head.prev = node
	}
	p.head = node

	// 如果链表为空，设置尾部
	if p.tail == nil {
		p.tail = node
	}
}

// addToHead 添加新节点到链表头部
func (p *SnowflakePool) addToHead(node *lruNode) {
	node.prev = nil
	node.next = p.head

	if p.head != nil {
		p.head.prev = node
	}
	p.head = node

	if p.tail == nil {
		p.tail = node
	}
}

// removeTail 移除链表尾部节点（最久未使用）
func (p *SnowflakePool) removeTail() *lruNode {
	if p.tail == nil {
		return nil
	}

	node := p.tail
	if p.tail.prev != nil {
		p.tail.prev.next = nil
		p.tail = p.tail.prev
	} else {
		p.head = nil
		p.tail = nil
	}

	return node
}

// GetOrCreate 获取或创建指定节点 ID 的雪花算法实例
//
// 参数:
//   - nodeID: 节点ID，取值范围 0-1023
//
// 返回:
//   - *Snowflake: 雪花算法实例
//   - error: 创建失败时返回错误
//
// 特性:
//   - 如果实例已存在，直接返回并更新 LRU 顺序
//   - 如果实例不存在，自动创建并缓存
//   - 当缓存满时，自动淘汰最久未使用的实例
//   - 线程安全
func (p *SnowflakePool) GetOrCreate(nodeID int64) (*Snowflake, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 检查缓存是否存在
	if node, exists := p.cache[nodeID]; exists {
		// 移动到头部（标记为最近使用）
		p.moveToHead(node)
		return node.snowflake, nil
	}

	// 创建新实例
	sf, err := NewSnowflake(nodeID)
	if err != nil {
		return nil, fmt.Errorf("[SnowflakePool.GetOrCreate] 创建实例失败 nodeID=%d: %w", nodeID, err)
	}

	// 创建新节点
	newNode := &lruNode{
		nodeID:    nodeID,
		snowflake: sf,
	}

	// 添加到缓存
	p.cache[nodeID] = newNode
	p.addToHead(newNode)

	// 检查容量，如果超出则淘汰最久未使用的
	if len(p.cache) > p.capacity {
		tail := p.removeTail()
		if tail != nil {
			delete(p.cache, tail.nodeID)
		}
	}

	return sf, nil
}

// NextID 使用指定节点 ID 生成唯一 ID（int64 格式）
//
// 参数:
//   - nodeID: 节点ID，取值范围 0-1023
//
// 返回:
//   - int64: 唯一的 64 位整数 ID
//   - error: 生成失败时返回错误
//
// 示例:
//
//	pool := NewSnowflakePool(50)
//	id, err := pool.NextID(123)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (p *SnowflakePool) NextID(nodeID int64) (int64, error) {
	sf, err := p.GetOrCreate(nodeID)
	if err != nil {
		return 0, err
	}
	return sf.NextID(), nil
}

// NextIDString 使用指定节点 ID 生成唯一 ID（字符串格式）
//
// 参数:
//   - nodeID: 节点ID，取值范围 0-1023
//
// 返回:
//   - string: 唯一 ID 的字符串表示
//   - error: 生成失败时返回错误
//
// 示例:
//
//	pool := NewSnowflakePool(50)
//	idStr, err := pool.NextIDString(123)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (p *SnowflakePool) NextIDString(nodeID int64) (string, error) {
	sf, err := p.GetOrCreate(nodeID)
	if err != nil {
		return "", err
	}
	return sf.NextIDString(), nil
}

// Size 返回当前缓存的实例数量
func (p *SnowflakePool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.cache)
}

// Clear 清空所有缓存的实例
func (p *SnowflakePool) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cache = make(map[int64]*lruNode)
	p.head = nil
	p.tail = nil
}

// 全局默认实例池
var (
	defaultPool     *SnowflakePool
	poolOnce        sync.Once
	defaultCapacity = 50 // 默认容量
)

// ensurePoolInitialized 确保全局默认实例池已初始化
func ensurePoolInitialized() {
	poolOnce.Do(func() {
		defaultPool = NewSnowflakePool(defaultCapacity)
	})
}

// SetPoolCapacity 设置全局默认实例池的容量（可选）
//
// 参数:
//   - capacity: 最大缓存容量
//
// 注意:
//   - 必须在首次调用 SnowflakeIDWithNode 等方法之前调用
//   - 如果已经初始化过，此方法不会重复初始化
//
// 示例:
//
//	tool.SetPoolCapacity(100)
func SetPoolCapacity(capacity int) {
	poolOnce.Do(func() {
		defaultPool = NewSnowflakePool(capacity)
	})
}

// SnowflakeIDWithNode 使用指定节点 ID 生成唯一 ID（int64 格式）
//
// 参数:
//   - nodeID: 节点ID，取值范围 0-1023
//
// 返回:
//   - int64: 唯一的 64 位整数 ID
//
// 特性:
//   - 开箱即用，自动管理实例池
//   - 线程安全
//   - 自动缓存和复用实例
//
// 示例:
//
//	id := tool.SnowflakeIDWithNode(123)
//	fmt.Println(id) // 输出: 1234567890123456789
func SnowflakeIDWithNode(nodeID int64) int64 {
	ensurePoolInitialized()
	id, err := defaultPool.NextID(nodeID)
	if err != nil {
		panic(fmt.Sprintf("[SnowflakeIDWithNode] 生成ID失败 nodeID=%d: %v", nodeID, err))
	}
	return id
}

// SnowflakeIDStringWithNode 使用指定节点 ID 生成唯一 ID（字符串格式）
//
// 参数:
//   - nodeID: 节点ID，取值范围 0-1023
//
// 返回:
//   - string: 唯一 ID 的字符串表示
//
// 特性:
//   - 开箱即用，自动管理实例池
//   - 线程安全
//   - 自动缓存和复用实例
//   - 字符串格式便于在 JSON 中传输，避免精度丢失
//
// 示例:
//
//	idStr := tool.SnowflakeIDStringWithNode(123)
//	fmt.Println(idStr) // 输出: "1234567890123456789"
func SnowflakeIDStringWithNode(nodeID int64) string {
	ensurePoolInitialized()
	idStr, err := defaultPool.NextIDString(nodeID)
	if err != nil {
		panic(fmt.Sprintf("[SnowflakeIDStringWithNode] 生成ID失败 nodeID=%d: %v", nodeID, err))
	}
	return idStr
}