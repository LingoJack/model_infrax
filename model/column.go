package model

import (
	"encoding/json"
	"fmt"
)

type Column struct {
	ColumnName      string
	Collate         string
	Comment         string
	Type            string
	IsAutoIncrement bool
	IsNullable      bool
	IsIndexed       bool
	IsUnique        bool
	IsPrimaryKey    bool
}

func (f Column) Json() string {
	byts, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}

func (f Column) JsonIndent() string {
	byts, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}
