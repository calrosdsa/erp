package internal

import (
	"gorm.io/gen"
)

func PianoModels(g *gen.Generator) []interface{} {
	pianoForm := g.GenerateModelAs("piano_form","PianoForm")
	return []interface{}{
		pianoForm,
	}
}

