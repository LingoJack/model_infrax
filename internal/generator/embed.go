package generator

import "embed"

//go:embed template/gorm/*.template
//go:embed template/itea-go/*.template
//go:embed template/tools/*.template
var templateFS embed.FS

const templatePathPrefix = "template/"

//go:embed template/gorm/po.template
var GormTemplateModel string

//go:embed template/gorm/dao.template
var GormTemplateDao string

//go:embed template/gorm/dto.template
var GormTemplateDto string

//go:embed template/gorm/vo.template
var GormTemplateVo string
