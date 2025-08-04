package user_rest

type UserPaths  struct {
	Base string 
}

func NewUserPaths(base string) UserPaths{
	return UserPaths{
		Base: base,
	}
}