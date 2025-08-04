package fsm

type FsmState interface {
	NextState(state string,events []int32,args ...interface{})(string,error)
}