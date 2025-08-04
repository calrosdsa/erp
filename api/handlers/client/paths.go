package client 

type ClientPaths struct {
	Base string
}

func NewClientPaths(
	base string,
) ClientPaths{
	return ClientPaths{
		Base: base,
	}

}