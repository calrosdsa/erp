package plugin

type PluginPaths struct {
	Base                    string
	Detail                  string
}

func NewPluginPaths(base string) PluginPaths {
	return PluginPaths{
		Base:                    base,
		Detail:                  base + "/{plugin}",
	}
}
