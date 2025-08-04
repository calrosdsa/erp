package pianoform_rest

type PianoFormPaths struct {
	Base string
	Detail string
	Export string
}

func newPianoPaths(base string) PianoFormPaths {
	return PianoFormPaths{
		Base:base,
		Detail:base + "/{id}",
		Export: base + "/export",
	}
}