package internal

import (
	"gorm.io/gen"
)

func BuyingModels(g *gen.Generator) []interface{} {
	supplier := g.GenerateModel("suppliers")
	return []interface{}{
		supplier,
	}
}
