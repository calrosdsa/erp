package tac_rest

type TacPaths struct {
	Base string
}

func NewTacPaths(base string)TacPaths{
	return TacPaths{
		Base: base,
	}
}