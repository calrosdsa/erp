package uom

type UOMPaths struct {
	Base                    string
}

func NewPluginPaths(base string) UOMPaths {
	return UOMPaths{
		Base:                    base,
	}
}
