package bus

import (
	"context"
	"sync"
	"time"
)

type (
	// bus is a message bus
	bus struct {
		mutex    sync.RWMutex
		idgen    Next
		topics   map[string][]Handler
		handlers map[string]Handler
	}

	Bus interface {
		Emit(ctx context.Context, topic string, data interface{}) error
		EmitWithOpts(ctx context.Context, topic string, data interface{}, opts ...EventOption) error
		Topics() []string 
		RegisterTopics(topics ...string)
		DeregisterTopics(topics ...string) 
		TopicHandlerKeys(topic string) []string
		HandlerKeys() []string 
		HandlerTopicSubscriptions(handlerKey string) []string
		RegisterHandler(key string, h Handler) 
		DeregisterHandler(key string)
	}

	// Next is a sequential unique id generator func type
	Next func() string

	// IDGenerator is a sequential unique id generator interface
	IDGenerator interface {
		Generate() string
	}

	// Event is data structure for any logs
	Event struct {
		ID         string      // identifier
		TxID       string      // transaction identifier
		Topic      string      // topic name
		Source     string      // source of the event
		OccurredAt time.Time   // creation time in nanoseconds
		Data       interface{} // actual event data
	}

	// Handler is a receiver for event reference with the given regex pattern
	Handler struct {
		key string

		// handler func to process events
		Handle func(ctx context.Context, e Event) error

		AbortOnError bool
		//handler func to process events returning error

		// topic matcher as regex pattern
		Matcher string
	}

	// EventOption is a function type to mutate event fields
	EventOption = func(Event) Event

	ctxKey int8
)
