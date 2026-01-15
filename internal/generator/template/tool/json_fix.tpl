package tool

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/kaptinlin/jsonrepair"
)

// JsonRepair 修复格式错误的 JSON 字符串。
// 该函数会尝试最小化工作量，首先进行有效性检查和区域隔离，
// 仅在必要时才使用 jsonrepair 进行强力修复。
// 如果有错误，则返回原始字符串
//
// 参数：
//   - input: 待修复的 JSON 字符串
//
// 返回值：
//   - 修复后的 JSON 字符串
func JsonRepair(input string) (result string, err error) {
	// 去除首尾空白字符
	result = strings.TrimSpace(input)

	// 快速路径：如果已经是有效的 JSON，直接返回
	if strings.HasPrefix(result, "{") && strings.HasSuffix(result, "}") && json.Valid([]byte(result)) {
		return
	}

	// 如果存在 JSON 对象区域，则隔离该区域；如果仅对象部分有效，则剥离噪音
	i := strings.IndexByte(result, '{')
	j := strings.LastIndexByte(result, '}')
	if i >= 0 && j >= i {
		// 提取从第一个 '{' 到最后一个 '}' 的子串
		sub := result[i : j+1]
		result = sub
		// 如果提取的子串是有效的 JSON，直接返回
		if json.Valid([]byte(result)) {
			return
		}
		// 否则继续处理这个子串
	}

	// 移除常见的大模型生成的标记符号
	result = strings.TrimPrefix(result, "<|FunctionCallBegin|>")
	result = strings.TrimSuffix(result, "<|FunctionCallEnd|>")
	result = strings.TrimPrefix(result, "<think>")

	// 如果此时已经是有效的 JSON，直接返回
	if json.Valid([]byte(result)) {
		return
	}

	// 启发式修复：如果缺少开头或结尾的大括号，尝试补全
	if !strings.HasPrefix(result, "{") && strings.HasSuffix(result, "}") {
		// 缺少开头的 '{'
		result = "{" + result
	} else if strings.HasPrefix(result, "{") && !strings.HasSuffix(result, "}") {
		// 缺少结尾的 '}'
		result = result + "}"
	}

	// 仅在无效时尝试强力修复
	result, err = jsonrepair.JSONRepair(result)
	if err != nil {
		// 如果修复失败，返回原始字符串
		result = input
		err = errors.New(fmt.Sprintf("json repair failed, err: %s", err.Error()))
		return
	}
	return
}
