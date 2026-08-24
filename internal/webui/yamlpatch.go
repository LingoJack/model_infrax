package webui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// PatchYAMLFile 按点分键更新 YAML 文件，保留原有注释与键顺序
// sets 中的值支持 string / int / bool / []string；键不存在时追加到对应层级
func PatchYAMLFile(path string, sets map[string]any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return fmt.Errorf("配置文件为空或格式异常")
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return fmt.Errorf("配置文件根节点不是 mapping")
	}

	// 排序保证写回稳定（按插入时的确定性顺序无要求，map 迭代顺序不影响单个键的正确性）
	for key, val := range sets {
		if err := setPath(root, key, val); err != nil {
			return err
		}
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	return os.WriteFile(path, out, 0644)
}

// setPath 在 mapping node 上按点分路径设置值
func setPath(m *yaml.Node, key string, val any) error {
	parts := splitKey(key)
	cur := m
	for _, part := range parts[:len(parts)-1] {
		next := mappingGet(cur, part)
		if next == nil {
			// 中间层级不存在：创建空 mapping
			created := &yaml.Node{Kind: yaml.MappingNode}
			mappingSet(cur, part, created)
			cur = created
			continue
		}
		if next.Kind != yaml.MappingNode {
			return fmt.Errorf("配置键 %s 的中间节点 %s 不是 mapping", key, part)
		}
		cur = next
	}

	leaf := parts[len(parts)-1]
	valueNode, err := toNode(val)
	if err != nil {
		return fmt.Errorf("配置键 %s 的值非法: %w", key, err)
	}
	mappingSet(cur, leaf, valueNode)
	return nil
}

// mappingGet 在 mapping node 中按 key 查找 value 子节点
func mappingGet(m *yaml.Node, key string) *yaml.Node {
	if m.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}

// mappingSet 设置 mapping node 的 key（存在则替换 value，不存在则追加）
func mappingSet(m *yaml.Node, key string, value *yaml.Node) {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			m.Content[i+1] = value
			return
		}
	}
	m.Content = append(m.Content, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: key,
	}, value)
}

// toNode 将 Go 值转换为 yaml.Node
func toNode(val any) (*yaml.Node, error) {
	switch v := val.(type) {
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v}, nil
	case bool:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(v)}, nil
	case int:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.Itoa(v)}, nil
	case float64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: strconv.FormatFloat(v, 'f', -1, 64)}, nil
	case []string:
		seq := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range v {
			seq.Content = append(seq.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: item})
		}
		return seq, nil
	case nil:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}, nil
	default:
		// 兜底：走标准序列化
		raw, err := yaml.Marshal(val)
		if err != nil {
			return nil, err
		}
		var n yaml.Node
		if err := yaml.Unmarshal(raw, &n); err != nil {
			return nil, err
		}
		if len(n.Content) == 0 {
			return nil, fmt.Errorf("无法序列化值: %v", val)
		}
		return n.Content[0], nil
	}
}

func splitKey(key string) []string {
	return strings.Split(key, ".")
}
