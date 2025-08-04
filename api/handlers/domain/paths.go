package domain

type PluginPaths struct {
	Plugins string
}

func NewConfigPaths(base string) PluginPaths {
	return PluginPaths{
		Plugins: base,
	}
}

type DomainPaths struct {
	Currency string
}

func NewDomainPaths(base string) DomainPaths {
	return DomainPaths{
		Currency: base + "/currency",
	}
}
