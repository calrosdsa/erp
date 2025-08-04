package party

type PartyPaths struct {
	Base              string
	TypeUsers         string
	PartyByReferences string
	References        string
	ReferencesType        string
	PartyReferences   string
}

func NewPartyPaths(base string) PartyPaths {
	return PartyPaths{
		Base:              base,
		TypeUsers:         base + "/type/users",
		PartyByReferences: base + "/parties-by-references/{party_type}",
		References:        base + "/references",
		ReferencesType:        base + "/references/type",
	}
}

type PartyAddressPaths struct {
	Base       string
	Detail     string
	References string
}

func NewPartyAddressPaths(base string) PartyAddressPaths {
	return PartyAddressPaths{
		Base:       base,
		Detail:     base + "/detail/{id}",
		References: base + "/references",
	}
}


type PartyContactPaths struct {
	Base       string
	Detail     string
	References string
}

func NewPartyContactPaths(base string) PartyContactPaths {
	return PartyContactPaths{
		Base:       base,
		Detail:     base + "/detail/{id}",
	}
}
