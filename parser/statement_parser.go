package parser

import (
	"fmt"
	"log"
	"model_infrax/config"
	"model_infrax/model"
	"model_infrax/tool"
	"os"
	"strings"

	"github.com/samber/lo"
)

type StatementParser struct {
	configger  *config.Configger
	statements []string
}

func NewStatementParser(cfg *config.Configger, sqlFilePath string) (*StatementParser, error) {
	path := tool.EscapeHomeDir(sqlFilePath)

	byts, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	statements := strings.Split(string(byts), ";")

	return &StatementParser{
		configger:  cfg,
		statements: statements,
	}, nil
}

func (p *StatementParser) Parse() (schemas []model.Schema, err error) {
	for _, statement := range p.statements {
		// 跳过空语句
		trimmed := strings.TrimSpace(statement)
		if trimmed == "" {
			continue
		}

		log.Printf("⌛️ parsing statement: %s", statement)
		var schema model.Schema
		schema, err = p.parseStatement(statement)
		if err != nil {
			return nil, fmt.Errorf("解析语句失败: %w", err)
		}
		schemas = append(schemas, schema)
	}
	return
}

func (p *StatementParser) FilterTables(schemas []model.Schema) (filtered []model.Schema) {
	if p.configger.GenerateConfig.AllTables {
		filtered = schemas
		return
	}
	filtered = lo.Filter(schemas, func(schema model.Schema, index int) bool {
		return lo.Contains(p.configger.GenerateConfig.TableNames, schema.Name)
	})
	return
}

func (p *StatementParser) parseStatement(statement string) (schema model.Schema, err error) {
	// TODO
	return schema, nil
}
