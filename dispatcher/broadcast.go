package dispatcher

import (
	"gopkg.in/ADone/go.events.v1"
)

func BroadcastFactory() events.Dispatcher {
	return &BroadcastDispatcher{make([]events.Listener, 0)}
}

type BroadcastDispatcher struct {
	Subscribers []events.Listener
}

func (dispatcher *BroadcastDispatcher) AddSubscribers(subscribers ...events.Listener) {
	dispatcher.Subscribers = append(dispatcher.Subscribers, subscribers...)
}

func (dispatcher *BroadcastDispatcher) Dispatch(event events.Event) {
	for _, subscriber := range dispatcher.Subscribers {
		subscriber.Handle(event)
	}
}
