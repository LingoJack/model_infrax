package assets

import _ "embed"

//go:embed config.yml
var DefaultConfigYml string

//go:embed schema.sql
var DefaultSchemaSql string
