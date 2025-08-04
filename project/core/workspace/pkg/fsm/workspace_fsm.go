package workspace_fsm

import (
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/pkg/fsm"
)

type workspaceFsm struct{}

func NewFsm() fsm.FsmState{
	return &workspaceFsm{}
}

func (f *workspaceFsm) NextState(state string,events []int32,args ...interface{})(string,error) {
	eventsState := make([]proto.EventState, len(events))
	for i, event := range events {
		eventsState[i] = proto.EventState(event)
	}
	stateMachine := fsm.New()
	draftState := stateMachine.NewState(proto.State_DRAFT.String())
	enableState := stateMachine.NewState(proto.State_ENABLED.String())
	disabledState := stateMachine.NewState(proto.State_DISABLED.String())


	submit := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_SUBMIT_EVENT))
	disabled := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_DISABLED_EVENT))
	enabled := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_ENABLED_EVENT))

	stateMachine.LinkStates(draftState, enableState, submit)
	stateMachine.LinkStates(enableState, disabledState, disabled)
	stateMachine.LinkStates(disabledState,enableState,enabled)

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