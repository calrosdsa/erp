package order_fsm

import (
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/pkg/fsm"
)

type orderFsm struct{}

func NewOrderFsm() fsm.FsmState {
	return &orderFsm{}
}

func (f *orderFsm) NextState(state string, events []int32, args ...interface{}) (string, error) {
	eventsState := make([]proto.EventState, len(events))
	for i, event := range events {
		eventsState[i] = proto.EventState(event)
	}
	var (
		partyType string
	)
	if len(args) > 0 {
		partyType = args[0].(string)
	}

	stateMachine := fsm.New()
	draftState := stateMachine.NewState(proto.State_DRAFT.String())
	toReceiveAndBill := stateMachine.NewState(proto.State_TO_RECEIVE_AND_BILL.String())
	toReceive := stateMachine.NewState(proto.State_TO_RECEIVE.String())
	toDeliver := stateMachine.NewState(proto.State_TO_DELIVER.String())
	toBill := stateMachine.NewState(proto.State_TO_BILL.String())
	toDeliverAndBill := stateMachine.NewState(proto.State_TO_DELIVER_AND_BILL.String())
	cancellledState := stateMachine.NewState(proto.State_CANCELLED.String())

	submit := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_SUBMIT_EVENT))
	cancel := fsm.NewRule(fsm.Operator("eq"), fsm.Event(proto.EventState_CANCEL_EVENT))
	switch partyType {
	case proto.PartyType_purchaseOrder.String():
		stateMachine.LinkStates(draftState, toReceiveAndBill, submit)
		stateMachine.LinkStates(toReceive, cancellledState, cancel)
	case proto.PartyType_saleOrder.String():
		stateMachine.LinkStates(draftState, toDeliverAndBill, submit)
		stateMachine.LinkStates(toDeliver, cancellledState, cancel)

	}

	stateMachine.LinkStates(toBill, cancellledState, cancel)
	stateMachine.LinkStates(toReceiveAndBill, cancellledState, cancel)
	stateMachine.LinkStates(toDeliverAndBill, cancellledState, cancel)

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
