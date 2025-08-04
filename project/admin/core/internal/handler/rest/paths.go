package rest_core

type CorePaths struct {
	Base string 
	Detail string
	Action string
}

func NewCorePaths(base string) CorePaths{
	return CorePaths{
		Base: base,
		Detail: base + "/detail/{id}",
		Action:base + "/action",
	}
}