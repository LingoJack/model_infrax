package parser

import (
	"github.com/LingoJack/model_infrax/internal/model"
)

type Parser interface {
	Parse() (schemas []model.Schema, err error)
	FilterTables(schemas []model.Schema) (filtered []model.Schema)
}
