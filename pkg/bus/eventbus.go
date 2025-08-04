package bus

import (
	"context"
	"fmt"
	// "regexp"
	"time"

	"github.com/google/uuid"
)

const (
	// CtxKeyTxID tx id context key
	CtxKeyTxID = ctxKey(116)

	// CtxKeySource source context key
	CtxKeySource = ctxKey(117)

	// Version syncs with package version
	Version = "3.0.3"

	empty = ""
)
type uuidGenerator struct{}

func NewGenerator() IDGenerator{
	return &uuidGenerator{}
}
func (g *uuidGenerator)Generate() string {
	return uuid.NewString()
}


// NewBus inits a new bus
func NewBus() (Bus) {
	generatorId := NewGenerator()
	return &bus{
		idgen:    generatorId.Generate,
		topics:   make(map[string][]Handler),
		handlers: make(map[string]Handler),
	}
}

// WithID returns an option to set event's id field
func WithID(id string) EventOption {
	return func(e Event) Event {
		e.ID = id
		return e
	}
}

// WithTxID returns an option to set event's txID field
func WithTxID(txID string) EventOption {
	return func(e Event) Event {
		e.TxID = txID
		return e
	}
}

// WithSource returns an option to set event's source field
func WithSource(source string) EventOption {
	return func(e Event) Event {
		e.Source = source
		return e
	}
}

// WithOccurredAt returns an option to set event's occurredAt field
func WithOccurredAt(time time.Time) EventOption {
	return func(e Event) Event {
		e.OccurredAt = time
		return e
	}
}



// Emit inits a new event and delivers to the interested in handlers with
// sync safety
func (b *bus) Emit(ctx context.Context, topic string, data interface{}) error {
	b.mutex.RLock()
	handlers, ok := b.topics[topic]
	b.mutex.RUnlock()

	if !ok {
		return fmt.Errorf("bus: topic(%s) not found", topic)
	}

	source, _ := ctx.Value(CtxKeySource).(string)
	txID, _ := ctx.Value(CtxKeyTxID).(string)
	if txID == empty {
		txID = b.idgen()
		ctx = context.WithValue(ctx, CtxKeyTxID, txID)
	}

	e := Event{
		ID:         b.idgen(),
		Topic:      topic,
		Data:       data,
		OccurredAt: time.Now(),
		TxID:       txID,
		Source:     source,
	}
	fmt.Println("HANDLERS",len(handlers))
	for _, h := range handlers {
		fmt.Println("HANDLER",h.Matcher,h.key)
		err := h.Handle(ctx, e)
		if h.AbortOnError {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// EmitWithOpts inits a new event and delivers to the interested in handlers
// with sync safety and options
func (b *bus) EmitWithOpts(ctx context.Context, topic string, data interface{}, opts ...EventOption) error {
	b.mutex.RLock()
	handlers, ok := b.topics[topic]
	b.mutex.RUnlock()

	if !ok {
		return fmt.Errorf("bus: topic(%s) not found", topic)
	}

	e := Event{Topic: topic, Data: data}
	for _, o := range opts {
		e = o(e)
	}

	if e.TxID == empty {
		e.TxID = b.idgen()
	}
	if e.ID == empty {
		e.ID = b.idgen()
	}
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}

	for _, h := range handlers {
		h.Handle(ctx, e)
	}

	return nil
}

// Topics lists the all registered topics
func (b *bus) Topics() []string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	topics, index := make([]string, len(b.topics)), 0

	for topic := range b.topics {
		topics[index] = topic
		index++
	}
	return topics
}

// RegisterTopics registers topics and fullfills handlers
func (b *bus) RegisterTopics(topics ...string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, n := range topics {
		b.registerTopic(n)
	}
}

// DeregisterTopics deletes topic
func (b *bus) DeregisterTopics(topics ...string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, n := range topics {
		b.deregisterTopic(n)
	}
}

// TopicHandlerKeys returns all handlers for the topic
func (b *bus) TopicHandlerKeys(topic string) []string {
	b.mutex.RLock()
	handlers := b.topics[topic]
	b.mutex.RUnlock()

	keys := make([]string, len(handlers))

	for i, h := range handlers {
		keys[i] = h.key
	}

	return keys
}

// HandlerKeys returns list of registered handler keys
func (b *bus) HandlerKeys() []string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	keys, index := make([]string, len(b.handlers)), 0

	for k := range b.handlers {
		keys[index] = k
		index++
	}
	return keys
}

// HandlerTopicSubscriptions returns all topic subscriptions of the handler
func (b *bus) HandlerTopicSubscriptions(handlerKey string) []string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.handlerTopicSubscriptions(handlerKey)
}

// RegisterHandler re/register the handler to the registry
func (b *bus) RegisterHandler(key string, h Handler) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	h.key = key
	b.registerHandler(h)
}

// DeregisterHandler deletes handler from the registry
func (b *bus) DeregisterHandler(key string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.deregisterHandler(key)
}

// Generate is an implementation of IDGenerator for bus.Next fn type
func (n Next) Generate() string {
	return n()
}

func (b *bus) registerHandler(h Handler) {
	// b.deregisterHandler(h.key)
	b.handlers[h.key] = h
	for _, t := range b.handlerTopicSubscriptions(h.key) {
		b.registerTopicHandler(t, h)
	}
}

func (b *bus) deregisterHandler(handlerKey string) {
	if _, ok := b.handlers[handlerKey]; ok {
		for _, t := range b.handlerTopicSubscriptions(handlerKey) {
			b.deregisterTopicHandler(t, handlerKey)
		}
		delete(b.handlers, handlerKey)
	}
}

func (b *bus) registerTopicHandler(topic string, h Handler) {
	b.topics[topic] = append(b.topics[topic], h)
}

func (b *bus) deregisterTopicHandler(topic, handlerKey string) {
	l := len(b.topics[topic])
	for i, h := range b.topics[topic] {
		if h.key == handlerKey {
			b.topics[topic][i] = b.topics[topic][l-1]
			b.topics[topic] = b.topics[topic][:l-1]
			break
		}
	}
}

func (b *bus) registerTopic(topic string) {
	if _, ok := b.topics[topic]; ok {
		return
	}

	b.topics[topic] = b.buildHandlers(topic)
}

func (b *bus) deregisterTopic(topic string) {
	delete(b.topics, topic)
}

func (b *bus) buildHandlers(topic string) []Handler {
	handlers := make([]Handler, 0)
	for _, h := range b.handlers {
		if h.Matcher == topic {
			handlers = append(handlers, h)
		}
		// if matched, _ := regexp.MatchString(h.Matcher, topic); matched {
		// 	handlers = append(handlers, h)
		// }
	}
	return handlers
}

func (b *bus) handlerTopicSubscriptions(handlerKey string) []string {
	var subscriptions []string
	h, ok := b.handlers[handlerKey]
	if !ok {
		return subscriptions
	}

	for topic := range b.topics {
		if h.Matcher == topic {
			subscriptions = append(subscriptions, topic)
		}
		// if matched, _ := regexp.MatchString(h.Matcher, topic); matched {
		// 	subscriptions = append(subscriptions, topic)
		// }
	}
	return subscriptions
}
