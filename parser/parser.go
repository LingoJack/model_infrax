package parser

import "model_infrax/model"

type Parser interface {
	Parse() (schemas []model.Schema, err error)
	FilterTables(schemas []model.Schema) (filtered []model.Schema)
}
