package generator

import "embed"

//go:embed template/*.template
//go:embed template/itea-go/*.template
//go:embed template/tools/*.template
var templateFS embed.FS

const templatePathPrefix = "template/"
