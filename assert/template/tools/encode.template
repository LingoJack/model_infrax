package tool

import (
	"encoding/json"
	"fmt"
)

func Jsonify(v interface{}) string {
	byts, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}

func JsonifyIndent(v interface{}) string {
	byts, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}
