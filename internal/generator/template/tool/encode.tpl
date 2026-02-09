package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// marshalNoEscapeHTML 序列化时不转义 HTML 特殊字符（如 &、<、>），避免 URL 等内容被破坏
func marshalNoEscapeHTML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	// Encode 会在末尾追加换行符，需要去掉
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b, nil
}

// marshalIndentNoEscapeHTML 带缩进的序列化，不转义 HTML 特殊字符
func marshalIndentNoEscapeHTML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b, nil
}

func Jsonify(v interface{}) string {
	byts, err := marshalNoEscapeHTML(v)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}

func JsonifyBytes(v interface{}) []byte {
	byts, err := marshalNoEscapeHTML(v)
	if err != nil {
		return []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
	return byts
}

func JsonifyIndent(v interface{}) string {
	byts, err := marshalIndentNoEscapeHTML(v)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}

func JsonifyIndentBytes(v interface{}) []byte {
	byts, err := marshalIndentNoEscapeHTML(v)
	if err != nil {
		return []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
	return byts
}