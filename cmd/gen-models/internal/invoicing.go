package internal

import (
	"gorm.io/gen"
)

func InvoicingModels(g *gen.Generator) []interface{} {
	purchaseRecord := g.GenerateModel("purchase_records")
	saleRecord := g.GenerateModel("sales_records")
	return []interface{}{
		purchaseRecord,
		saleRecord,
	}
}
