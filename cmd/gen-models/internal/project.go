package internal

import (
	"gorm.io/gen"
)

func ProjectModels(g *gen.Generator) []interface{} {
	project := g.GenerateModel("projects")
	return []interface{}{
		project,
	}
}
