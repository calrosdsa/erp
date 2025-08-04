package payment_terms_rest

type Paths struct {
	Base   string
	Detail string
	UpdateStatus string 
	Lines string 
}

func NewPaths(base string) Paths {
	return Paths{
		Base:   base,
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
		Lines: base + "/{id}/lines",
	}
}
