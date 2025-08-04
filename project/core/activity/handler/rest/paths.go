package activity_rest

type ActivityPaths struct {
	Comment string 
	ByPartyID string
	Base string 
}

func NewActivityPaths(base string) ActivityPaths {
	return ActivityPaths{
		Comment: base + "/comment",
		ByPartyID: base + "/{id}",
		Base: base,
	}
}