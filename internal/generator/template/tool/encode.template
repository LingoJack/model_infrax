package tool

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Jsonify(v interface{}) string {
	byts, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}

func JsonifyBytes(v interface{}) []byte {
	byts, err := json.Marshal(v)
	if err != nil {
		return []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
	return byts
}

func JsonifyIndent(v interface{}) string {
	byts, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}

func JsonifyIndentBytes(v interface{}) []byte {
	byts, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
	return byts
}

func parseLLMAnswer(data string, fmtType string) string {
	// 去除首尾空白字符
	data = strings.TrimSpace(data)

	// 如果字符串为空，直接返回
	if len(data) == 0 {
		return data
	}

	splitter := "```"
	startIdx := strings.Index(data, splitter+fmtType)
	prefixLen := 3 + len(fmtType)
	if startIdx == -1 {
		startIdx = strings.Index(data, splitter)
		prefixLen = 3
		if startIdx == -1 {
			return data
		}
	}

	// 查找```的结束位置
	endIdx := strings.LastIndex(data, splitter)
	if endIdx == -1 {
		return data
	}

	// 如果起始和结束位置相同，返回```json之后的内容
	if endIdx == startIdx {
		return data[startIdx+prefixLen:]
	}

	// 返回```json和```之间的内容
	return strings.TrimSpace(data[startIdx+prefixLen : endIdx])
}

func ParseLLMRespJson(data string) string {
	return parseLLMAnswer(data, "json")
}

func ParseLLMRespYaml(data string) string {
	return parseLLMAnswer(data, "yaml")
}
