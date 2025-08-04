package invoice_fsm

import (
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/pkg/fsm"
)

type orderFsm struct{}

func NewInvoiceFsm() fsm.FsmState{
	return &orderFsm{}
}

func (f *orderFsm) NextState(state string,events []int32,args ...interface{})(string,error) {
	eventsState := make([]proto.EventState, len(events))
	for i, event := range events {
		eventsState[i] = proto.EventState(event)
	}
	stateMachine := fsm.New()
	draftState := stateMachine.NewState(proto.State_DRAFT.String())
	unPaidState := stateMachine.NewState(proto.State_UNPAID.String())
	paidState := stateMachine.NewState(proto.State_PAID.String())
	cancellledState := stateMachine.NewState(proto.State_CANCELLED.String())
	partiallyPaid := stateMachine.NewState(proto.State_PARTIALLY_PAID.String())

	submit := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_SUBMIT_EVENT))
	cancel := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_CANCEL_EVENT))

	stateMachine.LinkStates(draftState, unPaidState, submit)
	stateMachine.LinkStates(unPaidState, cancellledState, cancel)
	stateMachine.LinkStates(paidState,cancellledState,cancel)
	stateMachine.LinkStates(partiallyPaid,cancellledState,cancel)
	// stateMachine.LinkStates(pendingState, cancellledState, cancel)

	err := stateMachine.SetInialState(state)
	if err != nil {
		return "", err
	}

	stateMachine.Compute(eventsState, true)
	if currentState, ok := stateMachine.PresentState.Value.(string); ok {
		return currentState, nil
	} else {
		return "", domain.FAIL_TYPE_ASSERTION
	}
}