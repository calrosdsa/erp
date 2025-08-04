package charges_template_rest

type Paths struct {
	Base   string
	Detail string
}

func NewPaths(base string) Paths {
	return Paths{
		Base:   base,
		Detail: base + "/detail/{id}",
	}
}
