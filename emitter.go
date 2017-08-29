package events

import (
	"sync"
)

// EmitterOption defines option for Emitter
type EmitterOption struct {
	apply func(*Emitter)
}

// WithEventStategy sets delivery strategy for provided event
func WithEventStategy(event string, strategy DispatchStrategy) EmitterOption {
	return EmitterOption{func(emitter *Emitter) {
		if dispatcher, exists := emitter.dispatchers[event]; exists {
			dispatcher.strategy = strategy
			return
		}

		emitter.dispatchers[event] = NewDispatcher(strategy)
	}}
}

// WithDefaultStrategy sets default delivery strategy for event emitter
func WithDefaultStrategy(strategy DispatchStrategy) EmitterOption {
	return EmitterOption{func(emitter *Emitter) {
		emitter.strategy = strategy
	}}
}

// NewEmitter creates new event emitter
func NewEmitter(options ...EmitterOption) *Emitter {
	emitter := new(Emitter)
	emitter.strategy = Broadcast
	emitter.dispatchers = make(map[string]*Dispatcher)

	for _, option := range options {
		option.apply(emitter)
	}

	return emitter
}

// Emitter
type Emitter struct {
	guard       sync.Mutex
	strategy    DispatchStrategy
	dispatchers map[string]*Dispatcher
}

// On subscribes listeners to provided event and return emitter
// usefull for chain subscriptions
func (emitter *Emitter) On(event string, handlers ...Listener) *Emitter {
	emitter.AddEventListeners(event, handlers...)
	return emitter
}

// AddEventListeners subscribes listeners to provided event
func (emitter *Emitter) AddEventListeners(event string, handlers ...Listener) {
	emitter.guard.Lock()
	defer emitter.guard.Unlock()

	if _, exists := emitter.dispatchers[event]; !exists {
		emitter.dispatchers[event] = NewDispatcher(emitter.strategy)
	}

	emitter.dispatchers[event].AddSubscribers(handlers)
}

// AddEventListener subscribes listeners to provided events
func (emitter *Emitter) AddEventListener(handler Listener, events ...string) {
	emitter.guard.Lock()
	defer emitter.guard.Unlock()

	for _, event := range events {
		if _, exists := emitter.dispatchers[event]; !exists {
			emitter.dispatchers[event] = NewDispatcher(emitter.strategy)
		}

		emitter.dispatchers[event].AddSubscriber(handler)
	}
}

// RemoveEventListeners unsubscribe all listeners from provided event
func (emitter *Emitter) RemoveEventListeners(event string) {
	emitter.guard.Lock()
	defer emitter.guard.Unlock()

	delete(emitter.dispatchers, event)
}

// RemoveEventListener unsubscribe provided listener from all events
func (emitter *Emitter) RemoveEventListener(handler Listener) {
	emitter.guard.Lock()
	defer emitter.guard.Unlock()

	for _, dispatcher := range emitter.dispatchers {
		dispatcher.RemoveSubscriber(handler)
	}
}

// Fire start delivering event to listeners
func (emitter *Emitter) Fire(data interface{}) {
	event := New(data)
	if dispatcher, ok := emitter.dispatchers[event.Key]; ok {
		dispatcher.Dispatch(event)
	}
}
