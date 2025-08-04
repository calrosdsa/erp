package internal

import (
	"gorm.io/gen"
)

func CrmModels(g *gen.Generator) []interface{} {
	deal := g.GenerateModel("deals")
	dealParticipant := g.GenerateModel("deal_participants")

	return []interface{}{
		deal,
		dealParticipant,
	}
}
