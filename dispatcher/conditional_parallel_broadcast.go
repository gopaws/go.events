package dispatcher

import (
	"gopkg.in/ADone/go.events.v1"
)

func ConditionalParallelBroadcastFactory() events.Dispatcher {
	return &ConditionalParallelBroadcastDispatcher{make([]events.Listener, 0)}
}

type ConditionalParallelBroadcastDispatcher struct {
	Subscribers []events.Listener
}

func (dispatcher *ConditionalParallelBroadcastDispatcher) AddSubscribers(subscribers ...events.Listener) {
	dispatcher.Subscribers = append(dispatcher.Subscribers, subscribers...)
}

func (dispatcher *ConditionalParallelBroadcastDispatcher) Dispatch(event events.Event) {
	if _, ok := event.Context["_sync"]; ok {
		delete(event.Context, "_sync")
		for _, subscriber := range dispatcher.Subscribers {
			subscriber.Handle(event)
		}
	} else {
		for _, subscriber := range dispatcher.Subscribers {
			go subscriber.Handle(event)
		}
	}
}
