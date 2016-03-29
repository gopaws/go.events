package dispatcher

import (
	"gopkg.in/ADone/go.events.v1"
)

func ParallelBroadcastFactory() events.Dispatcher {
	return &ParallelBroadcastDispatcher{make([]events.Listener, 0)}
}

type ParallelBroadcastDispatcher struct {
	Subscribers []events.Listener
}

func (dispatcher *ParallelBroadcastDispatcher) AddSubscribers(subscribers ...events.Listener) {
	dispatcher.Subscribers = append(dispatcher.Subscribers, subscribers...)
}

func (dispatcher *ParallelBroadcastDispatcher) Dispatch(event events.Event) {
	for _, subscriber := range dispatcher.Subscribers {
		go subscriber.Handle(event)
	}
}
