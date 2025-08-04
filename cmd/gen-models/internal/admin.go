package internal

import (
	"gorm.io/gen"
)

func AdminModels(g *gen.Generator) []interface{} {
	roleTemplates := g.GenerateModel("role_templates")
	return []interface{}{
		roleTemplates,
	}
}
