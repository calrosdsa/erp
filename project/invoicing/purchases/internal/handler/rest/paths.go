package purchase_record_rest

type Paths struct {
	Base   string
	Detail string
	UpdateStatus string 
	Export string
	Document string 
}

func NewPaths(base string) Paths {
	return Paths{
		Base:   base,
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
		Export: base + "/export",
		Document: base + "/export/document",
	}
}
