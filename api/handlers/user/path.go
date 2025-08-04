package user 
type ProfilePaths struct  {
	Base string
	Detail string 
	Me string 
}

func NewProfilePaths(base  string) ProfilePaths {
	return ProfilePaths{
		Base: base,
		Detail: base + "/detail/${id}",
		Me:base + "/me",
	}
}