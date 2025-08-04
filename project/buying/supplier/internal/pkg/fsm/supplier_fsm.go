package supplier_fsm

import (
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/pkg/fsm"
)

type supplierFsm struct{}

func NewSupplierFsm() fsm.FsmState{
	return &supplierFsm{}
}
func (f *supplierFsm) NextState(state string,events []int32,args ...interface{})(string,error) {
	eventsState := make([]proto.EventState, len(events))
	for i, event := range events {
		eventsState[i] = proto.EventState(event)
	}
	stateMachine := fsm.New()
	enabledState := stateMachine.NewState(proto.State_ENABLED.String())
	disabledState := stateMachine.NewState(proto.State_DISABLED.String())

	enabled := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_ENABLED_EVENT))
	disabled := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_DISABLED_EVENT))

	stateMachine.LinkStates(enabledState, disabledState, disabled)
	stateMachine.LinkStates(disabledState, enabledState, enabled)

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
// func (f *supplierFsm) NextState(state string,events []int32,args ...interface{})(string,error) {
// 	eventsState := make([]proto.EventState, len(events))
// 	for i, event := range events {
// 		eventsState[i] = proto.EventState(event)
// 	}
// 	stateMachine := fsm.New()
// 	draftState := stateMachine.NewState(proto.State_DRAFT.String())
// 	submittedState := stateMachine.NewState(proto.State_SUBMITTED.String())
// 	cancellledState := stateMachine.NewState(proto.State_CANCELLED.String())

// 	submit := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_SUBMIT_EVENT))
// 	cancel := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_CANCEL_EVENT))

// 	stateMachine.LinkStates(draftState, submittedState, submit)
// 	stateMachine.LinkStates(submittedState, cancellledState, cancel)

// 	err := stateMachine.SetInialState(state)
// 	if err != nil {
// 		return "", err
// 	}

// 	stateMachine.Compute(eventsState, true)
// 	if currentState, ok := stateMachine.PresentState.Value.(string); ok {
// 		return currentState, nil
// 	} else {
// 		return "", domain.FAIL_TYPE_ASSERTION
// 	}
// }