package model

import (
	"encoding/json"
	"fmt"
)

type IndexName string

type Schema struct {
	Name        string
	Columns     []Column
	Comment     string
	PrimaryKey  Index
	UniqueIndex []Index
	Indexes     []Index
}

func (t Schema) Json() string {
	byts, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}

func (t Schema) JsonIndent() string {
	byts, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(byts)
}
