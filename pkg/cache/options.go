package cache

type options struct {
	KeyEntity string 
	ID string 
	Revalidate bool
	Data interface{}
	TypeAssertion func(d interface{})
}

var Options options

type Option func(o *options) 

func (*options) WithKeyEntity(keyEntity string) Option {
	return func(o *options) {
		o.KeyEntity = keyEntity
	}
}
func (*options)WithTypeAssertion(f func(d interface{})) Option {
	return func(o *options) {
		o.TypeAssertion = f
	}
}

func (*options) WithID(id string) Option {
	return func(o *options) {
		o.ID = id
	}
}

func (*options) WithRevalidate(revalidate bool)Option {
	return func(o *options) {
		o.Revalidate = revalidate
	}
}

func (*options) WithData(data interface{}) Option {
	return func(o *options) {
		o.Data = data
	}
}


func (o *options) Apply(opts ...Option) options {
	r := options{}
	for _,option := range opts {
		option(&r)
	}
	return r
}
