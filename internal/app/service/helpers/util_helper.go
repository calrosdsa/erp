package helpers 

type Util interface {
	Min(a,b int32) int32
	Abs(x int32) int32
}

type util struct {

}

func NewUtil()Util {
	return &util{}
}

func (r *util) Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func (r *util) Min(x, y int32) int32 {
	if x < y {
		return x
	} else {
		return y
	}
}