package parser

import (
	"log"
	"model_infrax/config"
	"model_infrax/model"
)

type StatementParser struct {
	configger  *config.Configger
	statements []string
}

func NewStatementParser(cfg *config.Configger, statements []string) (*StatementParser, error) {
	return &StatementParser{
		configger:  cfg,
		statements: statements,
	}, nil
}

func (p *StatementParser) Parse() (schemas []model.Schema, err error) {
	for _, statement := range p.statements {
		log.Printf("Parsing statement: %s", statement)
	}
	return
}

func (p *StatementParser) FilterTables(schemas []model.Schema) (filtered []model.Schema) {
	return
}
