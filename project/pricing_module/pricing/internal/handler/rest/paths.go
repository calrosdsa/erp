package pricing_rest

type Paths struct {
	Base   string
	Detail string
	UpdateStatus string 
	GeneratePo string 
	GenerateQuotation string 
}

func NewPaths(base string) Paths {
	return Paths{
		Base:   base,
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
		GeneratePo: base + "/generate-po",
		GenerateQuotation: base + "/generate-quotation",
	}
}
