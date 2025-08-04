package fsm

import (
	"erp/gen/proto"
	"errors"
	"fmt"
	"strconv"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/multi"
)

type State struct {
	Id    int64
	Value interface{}
}

type Link struct {
	Id    int64
	T, F  graph.Node
	Rules map[Operator]Event
}

func (n State) ID() int64 {
	return n.Id
}


func (l Link) From() graph.Node {
	return l.F
}

func (l Link) To() graph.Node {
	return l.T
}

func (l Link) ID() int64 {
	return l.Id
}

func (l Link) ReversedLine() graph.Line {
	return Link{F: l.T, T: l.F}
}

func (n State) String() string {
	switch n.Value.(type) {
	case int:
		return strconv.Itoa(n.Value.(int))
	case float32:
		return fmt.Sprintf("%f", n.Value.(float32))
	case float64:
		return fmt.Sprintf("%f", n.Value.(float64))
	case bool:
		return strconv.FormatBool(n.Value.(bool))
	case string:
		return n.Value.(string)
	default:
		return ""
	}
}

type Event proto.EventState
type Operator string

var NodeIDCntr = 0
var LineIdCntr = 1

type StateMachine struct {
	PresentState State
	g            *multi.DirectedGraph
}

func New() *StateMachine {
	s := &StateMachine{}
	s.g = multi.NewDirectedGraph()

	return s
}

func (s *StateMachine) SetInialState(initStateValue interface{}) error {
	nodes := s.g.Nodes() // Get all nodes in the graph

    for nodes.Next() { // Iterate through all nodes
        state := nodes.Node().(State) // Assert the node to type State
        if state.Value == initStateValue { // Check if the state's value matches the given initStateValue
            s.PresentState = state // Set the PresentState to the found state
            return nil // Return nil to indicate success
        }
    }

    // If no matching state is found, return an error
    return errors.New("state with the specified value not found")
}

func (s *StateMachine) Init(initStateValue interface{}) State {
	s.PresentState = State{Id: int64(NodeIDCntr), Value: initStateValue}
	s.g.AddNode(s.PresentState)
	NodeIDCntr++
	return s.PresentState
}

func (s *StateMachine) NewState(stateValue interface{}) State {
	state := State{Id: int64(NodeIDCntr), Value: stateValue}
	s.g.AddNode(state)
	NodeIDCntr++
	return state
}

func (s *StateMachine) LinkStates(s1, s2 State, rule map[Operator]Event) {
	s.g.SetLine(Link{F: s1, T: s2, Id: int64(LineIdCntr), Rules: rule})
	LineIdCntr++
}

func NewRule(triggerConditionOperator Operator, comparisonValue Event) map[Operator]Event {
	return map[Operator]Event{triggerConditionOperator: comparisonValue}
}


func (s *StateMachine) FireEvent(e Event) error {
	presentNode := s.PresentState

	it := s.g.From(presentNode.Id)
	for it.Next() {
		n := s.g.Node(it.Node().ID()).(State)
		line := graph.LinesOf(s.g.Lines(presentNode.Id, n.Id))[0].(Link) // There can be one defined path between two distinct states

		for key, val := range line.Rules {
			k := string(key)
			switch k {
			case "eq":
				if val == e {
					s.PresentState = n
					return nil
				}
			default:
				fmt.Printf("Sorry, the comparison operator '%s' is not supported\n", k)
				return errors.New("UNSUPPORTED_COMPARISON_OPERATOR")
			}
		}
	}
	return nil
}

func (s *StateMachine) Compute(events []proto.EventState, printState bool) State {
	for _, e := range events {
		s.FireEvent(Event(e))
		if printState {
			fmt.Printf("%s\n", s.PresentState.String())
		}
	}
	return s.PresentState
}