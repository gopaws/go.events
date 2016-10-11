package emitter

import (
	"gopkg.in/ADone/go.meta.v1"

	"gopkg.in/ADone/go.events.v1"
	"gopkg.in/ADone/go.events.v1/dispatcher"
)

var DefaultDispatcherFactory = dispatcher.BroadcastFactory

func New(factory ...events.DispatcherFactory) *Emitter {
	emitter := new(Emitter)
	emitter.Dispatchers = make(map[string]events.Dispatcher)
	if len(factory) > 0 {
		emitter.DispatcherFactory = factory[0]
	} else {
		emitter.DispatcherFactory = DefaultDispatcherFactory
	}
	return emitter
}

type Emitter struct {
	DispatcherFactory events.DispatcherFactory
	Dispatchers       map[string]events.Dispatcher
}

func (emitter Emitter) On(event string, handlers ...events.Listener) events.Emitter {
	emitter.AddEventListener(event, handlers...)
	return emitter
}

func (emitter Emitter) AddEventListener(event string, handlers ...events.Listener) {
	if _, exists := emitter.Dispatchers[event]; !exists {
		emitter.Dispatchers[event] = emitter.DispatcherFactory()
	}
	emitter.Dispatchers[event].AddSubscribers(handlers...)
}

func (emitter Emitter) RemoveEventListeners(event string) {
	delete(emitter.Dispatchers, event)
}

func (emitter Emitter) Fire(e interface{}, context ...meta.Map) {
	var event events.Event

	switch e := e.(type) {
	case string:
		event = events.New(e)
	case events.Event:
		event = e
	}

	if len(context) > 0 {
		event.Context = event.Context.Merge(context[0])
	}

	if dispatcher, ok := emitter.Dispatchers[event.Key]; ok {
		dispatcher.Dispatch(event)
	}
}
