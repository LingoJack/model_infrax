package tool

import (
	"encoding/json"
	"fmt"
)

func JsonifyIndent(v interface{}) string {
	byts, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}
