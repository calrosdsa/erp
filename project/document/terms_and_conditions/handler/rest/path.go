package terms_and_conditions_rest

type Paths struct {
	Base   string
	Detail string
	UpdateStatus string 
}

func NewPaths(base string) Paths {
	return Paths{
		Base:   base,
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
	}
}
