package sales_record_rest

type Paths struct {
	Base   string
	Detail string
	UpdateStatus string 
	Export string 
}

func NewPaths(base string) Paths {
	return Paths{
		Base:   base,
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
		Export: base + "/export",
	}
}
