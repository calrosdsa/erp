package payment_terms_t_rest

type Paths struct {
	Base   string
	Detail string
	UpdateStatus string 
	Greet string 
}

func NewPaths(base string) Paths {
	return Paths{
		Base:   base,
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
		Greet: base + "/greet",
	}
}
