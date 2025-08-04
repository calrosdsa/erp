package internal

import (
	"gorm.io/gen"
)

func CompanyModels(g *gen.Generator) []interface{} {
	company := g.GenerateModel("companies")
	companyEntities := g.GenerateModel("company_entities")
	companyDefaults := g.GenerateModel("company_defaults")
	return []interface{}{
		company,
		companyEntities,
		companyDefaults,
	}
}
