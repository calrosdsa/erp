package handlersquare

type SquarePath struct {
	Base string 
	Catalog string 
	Test string
	Object string 
	PaymentWeebhook string
}

func NewSquarePath(base string)SquarePath{
	return SquarePath{
		Base: base,
		Catalog:base + "/{uuid}",
		Object:base + "/{uuid}/{object_id}",
		PaymentWeebhook: "/payment/webhook",
		Test: base + "/test",
	}
}


type SquareSubscriptionPaths struct {
	Base string
	Cancel string
}

func NewSquareSubscriptionPaths(base string)SquareSubscriptionPaths{
	return SquareSubscriptionPaths{
		Base: base,
		Cancel: base + "/cancel",
	}
}