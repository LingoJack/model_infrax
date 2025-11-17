package parser

import (
	"log"
	"model_infrax/config"
	"model_infrax/tool"
	"testing"
)

func TestStatementParser_Parse(t *testing.T) {
	sqlFilePath := "/Users/jacklingo/dev/model_infrax/assert/database.sql"

	configger, err := config.NewConfigger("/Users/jacklingo/dev/model_infrax/assert/application.yml")
	if err != nil {
		panic(err)
	}

	parser, err := NewStatementParser(configger, sqlFilePath)
	if err != nil {
		panic(err)
	}

	schemas, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	log.Printf("schemas: %s", tool.JsonifyIndent(schemas))
}
