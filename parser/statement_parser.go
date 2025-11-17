package parser

import (
	"model_infrax/config"
	"model_infrax/model"
)

type StatementParser struct {
	configger *config.Configger
}

func NewStatementParser(cfg *config.Configger) (*StatementParser, error) {
	return &StatementParser{
		configger: cfg,
	}, nil
}

func (p *StatementParser) Parse() (schemas []model.Schema, err error) {
	return nil, nil
}

func (p *StatementParser) FilterTables(schemas []model.Schema) (filtered []model.Schema) {
	return
}
