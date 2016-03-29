package dispatcher

import (
	"gopkg.in/ADone/go.events.v1"

	"testing"
)

type TestHandler struct {
	Target bool
}

func (handler *TestHandler) Handle(event events.Event) {
	handler.Target = true
}

func TestBroadcastFactory(t *testing.T) {
	dispatcher := BroadcastFactory()

	if len(dispatcher.(*BroadcastDispatcher).Subscribers) != 0 {
		t.Fail()
	}
}

func TestAddSubscribers(t *testing.T) {
	dispatcher := BroadcastFactory()
	handler := new(TestHandler)
	dispatcher.AddSubscribers(handler)
	if len(dispatcher.(*BroadcastDispatcher).Subscribers) != 1 {
		t.Log("Fail add subscriber")
		t.Fail()
	}
}

func TestDispatch(t *testing.T) {
	dispatcher := BroadcastFactory()
	handler := new(TestHandler)
	dispatcher.AddSubscribers(handler)

	dispatcher.Dispatch(events.New("test"))

	if !handler.Target {
		t.Log("Fail dispatch event")
		t.Fail()
	}
}
