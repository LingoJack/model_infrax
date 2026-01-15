package generator

import "embed"

//go:embed template/gorm/*.template
//go:embed template/itea-go/*.template
//go:embed template/tools/*.template
var fs embed.FS

const templatePathPrefix = "template/"
