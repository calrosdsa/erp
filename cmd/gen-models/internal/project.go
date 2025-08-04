package internal

import (
	"gorm.io/gen"
)

func ProjectModels(g *gen.Generator) []interface{} {
	project := g.GenerateModel("projects")
	task := g.GenerateModel("tasks")
	return []interface{}{
		project,
		task,
	}
}
