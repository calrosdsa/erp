package booking_fsm

import (
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/pkg/fsm"
)

type bookingFsm struct{}

func NewBookingFsm() fsm.FsmState {
	return &bookingFsm{}
}

func (f *bookingFsm) NextState(state string, events []int32, args ...interface{}) (string, error) {
	eventsState := make([]proto.EventState, len(events))
	for i, event := range events {
		eventsState[i] = proto.EventState(event)
	}
	stateMachine := fsm.New()
	unPaidStatus := stateMachine.NewState(proto.State_UNPAID.String())
	partiallyPaidStatus := stateMachine.NewState(proto.State_PARTIALLY_PAID.String())
	complentedStatus := stateMachine.NewState(proto.State_COMPLETED.String())
	cancellledStatus := stateMachine.NewState(proto.State_CANCELLED.String())
	deletedStatus := stateMachine.NewState(proto.State_DELETED.String())

	completed := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_COMPLETED_EVENT))
	cancel := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_CANCEL_EVENT))
	deleted := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_DELETED_EVENT))

	stateMachine.LinkStates(unPaidStatus, complentedStatus, completed)
	stateMachine.LinkStates(partiallyPaidStatus, complentedStatus, completed)

	stateMachine.LinkStates(unPaidStatus, cancellledStatus, cancel)
	stateMachine.LinkStates(partiallyPaidStatus, cancellledStatus, cancel)
	stateMachine.LinkStates(complentedStatus, cancellledStatus, cancel)
	stateMachine.LinkStates(cancellledStatus, deletedStatus, deleted)
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
