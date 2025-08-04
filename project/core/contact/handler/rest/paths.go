package contact_rest

type ContactPaths struct {
	Base       string
	Detail     string
	References string
}

func NewContactPaths(base string) ContactPaths {
	return ContactPaths{
		Base:   base,
		Detail: base + "/detail/{id}",
	}
}
