package party_rest

type PartyPaths struct {
	Base              string
	TypeUsers         string
	PartyByReferences string
	References        string
	ReferencesType    string
	PartyReferences   string
	Connections       string
}

func NewPartyPaths(base string) PartyPaths {
	return PartyPaths{
		Base:              base,
		TypeUsers:         base + "/type/users",
		PartyByReferences: base + "/parties-by-references/{party_type}",
		References:        base + "/references",
		ReferencesType:    base + "/references/type",
		Connections:       base + "/connections/{id}",
	}
}
