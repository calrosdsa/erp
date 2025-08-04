package internal

import (
	"gorm.io/gen"
)

func PartModels(g *gen.Generator) []interface{} {
	contacts := g.GenerateModel("contacts")
	return []interface{}{
		contacts,
	}
}
