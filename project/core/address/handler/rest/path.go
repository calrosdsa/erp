package address_rest


type PartyAddressPaths struct {
	Base       string
	Detail     string
	References string
}

func NewPaths(base string) PartyAddressPaths {
	return PartyAddressPaths{
		Base:       base,
		Detail:     base + "/detail/{id}",
		References: base + "/references",
	}
}