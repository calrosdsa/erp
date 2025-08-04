package module_rest

type Paths struct {
	Base   string
	Module string
	Detail string
	UpdateStatus string 
	SearchEntities string
	ModuleSection string 
}

func NewPaths(base string) Paths {
	return Paths{
		Base:   base,
		Module: base + "/{id}",
		Detail: base + "/detail/{id}",
		UpdateStatus: base + "/update-status",
		SearchEntities: base + "/search-entities",
		ModuleSection: base + "/module-section",
	}
}
