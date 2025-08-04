package rest_group

type GroupPaths struct {
	Base string
	Type string
	Detail string
	Decendents string 
}

func NewGroupPaths(base string) GroupPaths {
	return GroupPaths{
		Base: base,
		Type: base + "/{party}",
		Detail: base + "/detail/{id}",
		Decendents: base + "/descendents/{id}",
	}
}
