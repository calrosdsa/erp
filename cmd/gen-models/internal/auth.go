package internal

import (
	"gorm.io/gen"
)

func AuthModels(g *gen.Generator) []interface{} {
	profile := g.GenerateModel("profiles")
	role := g.GenerateModel("roles")
	session := g.GenerateModel("sessions")
	return []interface{}{
		profile,
		role,
		session,
	}
}
