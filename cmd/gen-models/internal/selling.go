package internal

import (
	"gorm.io/gen"
)

func SellingModels(g *gen.Generator) []interface{} {
	customers := g.GenerateModel("customers")
	return []interface{}{
		customers,
	}
}
