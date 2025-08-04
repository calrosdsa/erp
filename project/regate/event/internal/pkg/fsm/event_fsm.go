package event_fsm

import (
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/pkg/fsm"
)

type eventFsm struct{}

func NewEventFsm() fsm.FsmState{
	return &eventFsm{}
}

func (f *eventFsm) NextState(state string,events []int32,args ...interface{})(string,error) {
	eventsState := make([]proto.EventState, len(events))
	for i, event := range events {
		eventsState[i] = proto.EventState(event)
	}
	stateMachine := fsm.New()
	draftState := stateMachine.NewState(proto.State_DRAFT.String())
	enabledState := stateMachine.NewState(proto.State_ENABLED.String())
	cancellledState := stateMachine.NewState(proto.State_CANCELLED.String())
	completedState := stateMachine.NewState(proto.State_COMPLETED.String())

	submit := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_SUBMIT_EVENT))
	cancel := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_CANCEL_EVENT))
	complete := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_COMPLETED_EVENT))
	enabled := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_ENABLED_EVENT))
	

	stateMachine.LinkStates(draftState, enabledState, submit)
	stateMachine.LinkStates(enabledState, cancellledState, cancel)
	stateMachine.LinkStates(enabledState, completedState, complete)
	stateMachine.LinkStates(cancellledState, enabledState, enabled)
	stateMachine.LinkStates(completedState, enabledState, enabled)

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