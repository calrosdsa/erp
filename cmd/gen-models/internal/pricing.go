package internal

import (
	"gorm.io/gen"
)

func PricingModels(g *gen.Generator) []interface{} {
	pricingLine := g.GenerateModel("pricing_line_items")
	pricingCharge := g.GenerateModel("pricing_charges")
	pricings := g.GenerateModel("pricings")
	return []interface{}{
		pricingLine,
		pricingCharge,
		pricings,
	}
}
